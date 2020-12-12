package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/andylibrian/terjang/pkg/messages"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

const LoadTestNotStarted = 0
const LoadTestRunning = 1
const LoadTestDone = 2
const LoadTestStopped = 3

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
		upgrader:            websocket.Upgrader{},
		workerService:       NewWorkerService(),
		notificationService: NewNotificationService(),
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

		s.workerService.GetMessageHandler().HandleMessage(message)
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
		serverInfo := messages.ServerInfo{NumOfWorkers: len(s.workerService.workers), State: loadTestStateToString(s.loadTestState)}
		serverInfoMsg, _ := json.Marshal(serverInfo)

		envelope := messages.Envelope{Kind: messages.KindServerInfo, Data: string(serverInfoMsg)}
		envelopeMsg, _ := json.Marshal(envelope)
		s.notificationService.BroadcastMessageToSubscribers([]byte(envelopeMsg))

		time.Sleep(1 * time.Second)
	}
}

func loadTestStateToString(s int) string {
	switch s {
	case LoadTestNotStarted:
		return "NotStarted"
	case LoadTestRunning:
		return "Running"
	case LoadTestDone:
		return "Done"
	case LoadTestStopped:
		return "Stopped"
	}

	return ""
}
