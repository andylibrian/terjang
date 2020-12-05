package server

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	upgrader      websocket.Upgrader
	workerService *WorkerService
	httpServer    *http.Server
}

// NewServer creates a new instance of server.
func NewServer() *Server {
	return &Server{
		upgrader:      websocket.Upgrader{},
		workerService: NewWorkerService(),
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
