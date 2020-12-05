package worker

import (
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Worker struct {
	conn                 *websocket.Conn
	connWriteLock        *sync.Mutex
	messageHandler       MessageHandler
	connectRetryInterval time.Duration
	isConnectedCh        chan struct{}
}

type MessageHandler interface {
	HandleMessage(message []byte)
}

type defaultMessageHandler struct {
}

func NewWorker() *Worker {
	return &Worker{
		messageHandler:       &defaultMessageHandler{},
		connectRetryInterval: 5 * time.Second,
		isConnectedCh:        make(chan struct{}),
	}
}

func (w *Worker) Run() {
	serverURL := url.URL{Scheme: "ws", Host: "127.0.0.1:9009", Path: "/cluster/join"}

	serverURLStr := serverURL.String()

	var conn *websocket.Conn
	var err error

	for i := 0; i < 10; i++ {
		conn, _, err = websocket.DefaultDialer.Dial(serverURLStr, nil)

		if err == nil {
			break
		}

		time.Sleep(w.connectRetryInterval)
	}

	w.conn = conn
	defer conn.Close()
	defer close(w.isConnectedCh)
	w.isConnectedCh <- struct{}{}

	for {
		_, message, err := conn.ReadMessage()
		w.messageHandler.HandleMessage(message)

		if err != nil {
			return
		}
	}
}

func (w *Worker) SetConnectRetryInterval(d time.Duration) {
	w.connectRetryInterval = d
}

func (w *Worker) SendMessageToServer(message []byte) {
	if w.conn == nil {
		// should indicate error
	} else {
		w.conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (w *Worker) GetMessageHandler() MessageHandler {
	return w.messageHandler
}

func (w *Worker) SetMessageHandler(h MessageHandler) {
	w.messageHandler = h
}

func (w *Worker) IsConnectedCh() <-chan struct{} {
	return w.isConnectedCh
}

func (h *defaultMessageHandler) HandleMessage(message []byte) {
}
