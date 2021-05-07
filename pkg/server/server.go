package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andylibrian/terjang/pkg/messages"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/rakyll/statik/fs"

	// Import statik package
	_ "github.com/andylibrian/terjang/pkg/server/statik"
)

var logger *zap.SugaredLogger

func init() {
	l, err := zap.NewProduction()

	if err != nil {
		panic("Can not create logger")
	}

	logger = l.Sugar()
}

// SetLogger is ...
func SetLogger(l *zap.SugaredLogger) {
	logger = l
}

// Server is ...
type Server struct {
	upgrader            websocket.Upgrader
	workerService       *WorkerService
	notificationService *NotificationService
	httpServer          *http.Server
	loadTestState       int
}

// NewServer creates a new instance of server.
func NewServer() *Server {
	return &Server{
		upgrader:            websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		workerService:       NewWorkerService(),
		notificationService: NewNotificationService(),
		loadTestState:       messages.ServerStateNotStarted,
	}
}

// GetWorkerService returns the worker service.
func (s *Server) GetWorkerService() *WorkerService {
	return s.workerService
}

// Run listens on the specified port and serve requests.
func (s *Server) Run(addr string) error {
	router, err := s.setupRouter()

	if err != nil {
		return err
	}

	go s.runNotificationLoop()
	go s.watchWorkerStateChange()

	s.httpServer = &http.Server{Addr: addr, Handler: router}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		s.httpServer.Shutdown(ctx)
	}()

	logger.Infow("Server is listening on", "address", addr)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("Server failed to listen and serve: %w", err)
	}

	return nil
}

// Close is ...
func (s *Server) Close() error {
	if s.httpServer == nil {
		return nil
	}

	return s.httpServer.Close()
}

func (s *Server) setupRouter() (*httprouter.Router, error) {
	statikFs, err := fs.New()
	if err != nil {
		return nil, err
	}

	router := httprouter.New()

	// static files
	router.GET("/", serveStatikFile(statikFs, "/index.html"))
	router.GET("/favicon.ico", serveStatikFile(statikFs, "/favicon.ico"))
	router.Handler("GET", "/js/*filepath", http.FileServer(statikFs))
	router.Handler("GET", "/css/*filepath", http.FileServer(statikFs))

	router.GET("/cluster/join", s.acceptWorkerConn)
	router.GET("/notifications", s.acceptNotificationConn)
	router.POST("/api/v1/load_test", s.handleStartLoadTest)
	router.DELETE("/api/v1/load_test", s.handleStopLoadTest)

	router.GET("/healthz", s.handleHealthz)

	// CORS
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Headers", "*")
		}

		w.WriteHeader(204)
	})

	return router, nil
}

func serveStatikFile(fs http.FileSystem, path string) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		reader, err := fs.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()
		contents, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		}

		w.Write(contents)
	}
}

func (s *Server) acceptWorkerConn(responseWriter http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	name := ""

	names, ok := req.URL.Query()["name"]
	if ok && len(names[0]) > 0 {
		name = names[0]
	}

	conn, err := s.upgrader.Upgrade(responseWriter, req, nil)
	if err != nil {
		logger.Warnw("Failed to upgrade websocket connection", "error", err)
		// TODO: should respond?
		return
	}

	s.workerService.AddWorker(conn, name)

	logger.Infow("Worker connected", "name", name)

	defer s.stopLoadTestIfNoWorkerRemaining()
	defer logger.Infow("Worker removed", "name", name)
	defer s.workerService.RemoveWorker(conn)
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		s.workerService.GetMessageHandler().HandleMessage(conn, message)
	}
}

func (s *Server) stopLoadTestIfNoWorkerRemaining() {
	if len(s.workerService.workers) == 0 {
		s.loadTestState = messages.ServerStateStopped
	}
}

func (s *Server) acceptNotificationConn(responseWriter http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	conn, err := s.upgrader.Upgrade(responseWriter, req, nil)
	if err != nil {
		logger.Warnw("Failed to upgrade websocket connection", "error", err)
		return
	}

	s.notificationService.AddSubscriber(conn)

	logger.Infow("Notification subscriber connected")

	defer logger.Infow("Notification subscriber removed")
	defer s.notificationService.RemoveSubscriber(conn)
	defer conn.Close()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (s *Server) runNotificationLoop() {
	for {
		// Server Info
		serverInfo := messages.ServerInfo{NumOfWorkers: len(s.workerService.workers), State: loadTestStateToString(s.loadTestState)}
		serverInfoMsg, _ := json.Marshal(serverInfo)

		envelope := messages.Envelope{Kind: messages.KindServerInfo, Data: string(serverInfoMsg)}
		envelopeMsg, _ := json.Marshal(envelope)
		s.notificationService.BroadcastMessageToSubscribers([]byte(envelopeMsg))

		// Workers Info
		var wks []*worker
		for _, v := range s.workerService.workers {
			wks = append(wks, v)
		}
		workersInfoMsg, _ := json.Marshal(wks)

		envelope = messages.Envelope{Kind: messages.KindWorkersInfo, Data: string(workersInfoMsg)}
		envelopeMsg, _ = json.Marshal(envelope)

		s.notificationService.BroadcastMessageToSubscribers([]byte(envelopeMsg))

		time.Sleep(1 * time.Second)
	}
}

// StartLoadTest is ...
func (s *Server) StartLoadTest(r *messages.StartLoadTestRequest) {
	req, _ := json.Marshal(r)
	envelope, _ := json.Marshal(messages.Envelope{Kind: messages.KindStartLoadTestRequest, Data: string(req)})

	s.loadTestState = messages.ServerStateRunning
	s.GetWorkerService().BroadcastMessageToWorkers(envelope)

	logger.Infow("Started load test", "request", r)
}

// StopLoadTest is ...
func (s *Server) StopLoadTest() {
	envelope, _ := json.Marshal(messages.Envelope{Kind: messages.KindStopLoadTestRequest})
	s.GetWorkerService().BroadcastMessageToWorkers(envelope)

	logger.Infow("Stopped load test")
}

func (s *Server) watchWorkerStateChange() {
	for {
		<-s.workerService.stateUpdatedCh
		s.loadTestState = s.summarizeWorkerStates()
	}
}

func (s *Server) summarizeWorkerStates() int {
	serverState := s.loadTestState

	states := make(map[messages.WorkerState]int)

	for _, worker := range s.workerService.workers {
		if _, ok := states[worker.state]; ok {
			states[worker.state]++
		} else {
			states[worker.state] = 1
		}
	}

	if val, ok := states[messages.WorkerStateDone]; ok && val == len(s.workerService.workers) {
		serverState = messages.ServerStateDone
	}

	if val, ok := states[messages.WorkerStateStopped]; ok && val == len(s.workerService.workers) {
		serverState = messages.ServerStateStopped
	}

	return serverState
}

func loadTestStateToString(s int) string {
	switch s {
	case messages.ServerStateNotStarted:
		return "NotStarted"
	case messages.ServerStateRunning:
		return "Running"
	case messages.ServerStateDone:
		return "Done"
	case messages.ServerStateStopped:
		return "Stopped"
	}

	return ""
}

func (s *Server) handleStartLoadTest(responseWriter http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var startLoadTestRequest messages.StartLoadTestRequest
	err := json.NewDecoder(req.Body).Decode(&startLoadTestRequest)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	s.StartLoadTest(&startLoadTestRequest)

	header := responseWriter.Header()
	header.Set("Access-Control-Allow-Origin", "*")

	responseWriter.WriteHeader(200)
	responseWriter.Write([]byte("ok"))
}

func (s *Server) handleStopLoadTest(responseWriter http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	s.StopLoadTest()

	header := responseWriter.Header()
	header.Set("Access-Control-Allow-Origin", "*")

	responseWriter.WriteHeader(204)
}

func (s *Server) handleHealthz(responseWriter http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	responseWriter.WriteHeader(200)
}
