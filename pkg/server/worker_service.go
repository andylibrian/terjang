package server

import (
	"encoding/json"
	"sync"

	"github.com/andylibrian/terjang/pkg/messages"
	"github.com/gorilla/websocket"
)

type worker struct {
	Name     string `json:"name"`
	conn     *websocket.Conn
	Metrics  messages.WorkerLoadTestMetrics `json:"metrics"`
	state    messages.WorkerState
	StateStr string `json:"state"`
}

// WorkerService maintains a collection of workers and
// provide a function to broadcast messages to them.
type WorkerService struct {
	messageHandler MessageHandler
	workers        map[*websocket.Conn]*worker
	workersLock    sync.RWMutex
	stateUpdatedCh chan struct{}
}

// MessageHandler is the interface to handle message from a worker.
type MessageHandler interface {
	HandleMessage(conn *websocket.Conn, message []byte)
}

type defaultMessageHandler struct {
	workerService *WorkerService
}

// NewWorkerService creates a new worker service.
func NewWorkerService() *WorkerService {
	w := &WorkerService{
		workers:        make(map[*websocket.Conn]*worker),
		stateUpdatedCh: make(chan struct{}),
	}

	w.messageHandler = &defaultMessageHandler{workerService: w}

	return w
}

// GetMessageHandler returns the registered message handler.
func (w *WorkerService) GetMessageHandler() MessageHandler {
	return w.messageHandler
}

// SetMessageHandler registers a message handler to be used by WorkerService.
func (w *WorkerService) SetMessageHandler(h MessageHandler) {
	w.messageHandler = h
}

// AddWorker registers a worker.
func (w *WorkerService) AddWorker(conn *websocket.Conn, name string) {
	w.workersLock.Lock()
	defer w.workersLock.Unlock()

	w.workers[conn] = &worker{conn: conn, Name: name}
}

// RemoveWorker removes a worker from the collection.
func (w *WorkerService) RemoveWorker(conn *websocket.Conn) {
	w.workersLock.Lock()
	defer w.workersLock.Unlock()

	delete(w.workers, conn)
}

// BroadcastMessageToWorkers sends a message to the registered workers.
func (w *WorkerService) BroadcastMessageToWorkers(message []byte) {
	w.workersLock.RLock()
	defer w.workersLock.RUnlock()

	for conn := range w.workers {
		// TODO: conn should be synced
		conn.WriteMessage(websocket.TextMessage, message)
	}
}

// HandleMessage handle messages from a worker.
func (h *defaultMessageHandler) HandleMessage(conn *websocket.Conn, message []byte) {
	var envelope messages.Envelope
	err := json.Unmarshal(message, &envelope)

	if err != nil {
		return
	}

	if envelope.Kind == messages.KindWorkerInfo {
		var workerInfo messages.WorkerInfo
		json.Unmarshal([]byte(envelope.Data), &workerInfo)

		w := h.workerService.workers[conn]

		if w.state != workerInfo.State {
			w.state = workerInfo.State
			h.workerService.stateUpdatedCh <- struct{}{}
		}
	} else if envelope.Kind == messages.KindWorkerLoadTestMetrics {
		w := h.workerService.workers[conn]
		json.Unmarshal([]byte(envelope.Data), &w.Metrics)
	}
}
