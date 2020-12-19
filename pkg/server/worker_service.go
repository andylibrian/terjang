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

type WorkerService struct {
	messageHandler MessageHandler
	workers        map[*websocket.Conn]*worker
	workersLock    sync.RWMutex
	stateUpdatedCh chan struct{}
}

type MessageHandler interface {
	HandleMessage(conn *websocket.Conn, message []byte)
}

type defaultMessageHandler struct {
	workerService *WorkerService
}

func NewWorkerService() *WorkerService {
	w := &WorkerService{
		workers:        make(map[*websocket.Conn]*worker),
		stateUpdatedCh: make(chan struct{}),
	}

	w.messageHandler = &defaultMessageHandler{workerService: w}

	return w
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

	w.workers[conn] = &worker{conn: conn}
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
