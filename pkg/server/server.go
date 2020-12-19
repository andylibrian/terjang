package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/andylibrian/terjang/pkg/messages"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

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
func (s *Server) Run() error {
	router, err := s.setupRouter()

	if err != nil {
		return err
	}

	go s.runNotificationLoop()
	go s.watchWorkerStateChange()

	s.httpServer = &http.Server{Addr: "0.0.0.0:9009", Handler: router}
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		// TODO: log err
	}

	return nil
}

func (s *Server) Close() error {
	if s.httpServer == nil {
		return nil
	}

	return s.httpServer.Close()
}

func (s *Server) setupRouter() (*httprouter.Router, error) {
	router := httprouter.New()

	router.GET("/cluster/join", s.acceptWorkerConn)
	router.GET("/notifications", s.acceptNotificationConn)

	// CORS
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Headers", "*")
		}

		w.WriteHeader(204)
	})

	return router, nil
}

func (s *Server) acceptWorkerConn(responseWriter http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	conn, err := s.upgrader.Upgrade(responseWriter, req, nil)
	if err != nil {
		// TODO: should respond? should probably log
		return
	}

	s.workerService.AddWorker(conn)

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

func (s *Server) acceptNotificationConn(responseWriter http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	conn, err := s.upgrader.Upgrade(responseWriter, req, nil)
	if err != nil {
		// TODO: should respond? should probably log
		return
	}

	s.notificationService.AddSubscriber(conn)

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

func (s *Server) StartLoadTest(r *messages.StartLoadTestRequest) {
	req, _ := json.Marshal(r)
	envelope, _ := json.Marshal(messages.Envelope{Kind: messages.KindStartLoadTestRequest, Data: string(req)})

	s.loadTestState = messages.ServerStateRunning
	s.GetWorkerService().BroadcastMessageToWorkers(envelope)
}

func (s *Server) watchWorkerStateChange() {
	for {
		<-s.workerService.stateUpdatedCh
		s.loadTestState = s.summarizeWorkerStates()
	}
}

func (s *Server) summarizeWorkerStates() int {
	var serverState int = s.loadTestState

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
