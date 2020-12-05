package server

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WorkerService struct {
	messageHandler MessageHandler
	workers        map[*websocket.Conn]struct{}
	workersLock    sync.RWMutex
}

type MessageHandler interface {
	HandleMessage(message []byte)
}

type defaultMessageHandler struct {
}

func NewWorkerService() *WorkerService {
	return &WorkerService{
		messageHandler: &defaultMessageHandler{},
		workers:        make(map[*websocket.Conn]struct{}),
	}
}

func (w *WorkerService) GetMessageHandler() MessageHandler {
	return w.messageHandler
}

func (w *WorkerService) SetMessageHandler(h MessageHandler) {
	w.messageHandler = h
}

func (w *WorkerService) AddWorker(conn *websocket.Conn) {
	w.workersLock.Lock()
	defer w.workersLock.Unlock()

	w.workers[conn] = struct{}{}
}

func (w *WorkerService) RemoveWorker(conn *websocket.Conn) {
	w.workersLock.Lock()
	defer w.workersLock.Unlock()

	delete(w.workers, conn)
}

func (w *WorkerService) BroadcastMessageToWorkers(message []byte) {
	w.workersLock.RLock()
	defer w.workersLock.RUnlock()

	for conn := range w.workers {
		// TODO: conn should be synced
		conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (h *defaultMessageHandler) HandleMessage(message []byte) {

}
