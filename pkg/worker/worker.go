package worker

import (
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Worker struct {
	conn           *websocket.Conn
	connWriteLock  *sync.Mutex
	messageHandler MessageHandler
}

type MessageHandler interface {
	HandleMessage(message []byte)
}

type defaultMessageHandler struct {
}

func NewWorker() *Worker {
	return &Worker{
		messageHandler: &defaultMessageHandler{},
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

		time.Sleep(5 * time.Second)
	}

	w.conn = conn
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		w.messageHandler.HandleMessage(message)

		if err != nil {
			return
		}
	}
}

func (w *Worker) SendMessageToServer(message string) {
	w.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (w *Worker) GetMessageHandler() MessageHandler {
	return w.messageHandler
}

func (w *Worker) SetMessageHandler(h MessageHandler) {
	w.messageHandler = h
}

func (h *defaultMessageHandler) HandleMessage(message []byte) {
}
