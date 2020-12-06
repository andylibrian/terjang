package worker

import (
	"encoding/json"
	"net/url"
	"sync"
	"time"

	"github.com/andylibrian/terjang/pkg/messages"
	"github.com/gorilla/websocket"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Worker struct {
	conn                 *websocket.Conn
	connWriteLock        *sync.Mutex
	messageHandler       MessageHandler
	connectRetryInterval time.Duration
	isConnectedCh        chan struct{}
	attacker             *vegeta.Attacker
	metrics              vegeta.Metrics
}

type MessageHandler interface {
	HandleMessage(message []byte)
}

type defaultMessageHandler struct {
	worker *Worker
}

func NewWorker() *Worker {
	worker := &Worker{
		connectRetryInterval: 5 * time.Second,
		isConnectedCh:        make(chan struct{}),
		attacker:             vegeta.NewAttacker(),
	}

	msgHandler := &defaultMessageHandler{worker: worker}
	worker.messageHandler = msgHandler

	return worker
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
	var envelope messages.Envelope
	err := json.Unmarshal(message, &envelope)

	if err != nil {
		return
	}

	if envelope.Kind == messages.KindStartLoadTestRequest {
		var req messages.StartLoadTestRequest
		err = json.Unmarshal([]byte(envelope.Data), &req)

		if err != nil {
			return
		}

		rate := vegeta.Rate{Freq: int(req.Rate), Per: time.Second}
		duration := time.Duration(req.Duration) * time.Second
		targeter := vegeta.NewStaticTargeter(vegeta.Target{
			Method: req.Method,
			URL:    req.Url,
			Body:   []byte(req.Body),
		})

		go h.worker.startLoadTest(targeter, rate, duration, "terjang")
	} else if envelope.Kind == messages.KindStopLoadTestRequest {
		h.worker.stopLoadTest()
	}
}

func (w *Worker) startLoadTest(tr vegeta.Targeter, p vegeta.Pacer, du time.Duration, name string) {
	for res := range w.attacker.Attack(tr, p, du, "terjang") {
		w.metrics.Add(res)
	}
}

func (w *Worker) stopLoadTest() {
	w.attacker.Stop()
}
