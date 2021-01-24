package worker

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/andylibrian/terjang/pkg/messages"
	"github.com/gorilla/websocket"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	l, err := zap.NewProduction()

	if err != nil {
		panic("Can not create logger")
	}

	logger = l.Sugar()
}

func SetLogger(l *zap.SugaredLogger) {
	logger = l
}

type Worker struct {
	name                 string
	conn                 *websocket.Conn
	connWriteLock        sync.Mutex
	messageHandler       MessageHandler
	connectRetryInterval time.Duration
	attacker             *vegeta.Attacker
	metrics              vegeta.Metrics
	loadTestState        messages.WorkerState
	connectedCallbacks   []func()
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
		attacker:             vegeta.NewAttacker(),
	}

	msgHandler := &defaultMessageHandler{worker: worker}
	worker.messageHandler = msgHandler

	return worker
}

func (w *Worker) SetName(name string) {
	w.name = name
}

func (w *Worker) Run(addr string) {
	serverURL := url.URL{Scheme: "ws", Host: addr, Path: "/cluster/join", RawQuery: "name=" + w.name}

	serverURLStr := serverURL.String()

	var conn *websocket.Conn
	var err error

	for i := 0; i < 10; i++ {
		logger.Infow("Connecting to server", "address", addr)

		conn, _, err = websocket.DefaultDialer.Dial(serverURLStr, nil)

		if err == nil {
			break
		}

		time.Sleep(w.connectRetryInterval)
	}

	logger.Infow("Connected to server", "address", addr)

	w.conn = conn
	defer conn.Close()

	go func() {
		for _, callback := range w.connectedCallbacks {
			callback()
		}
	}()

	go w.LoopSendMetricsToServer()

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
		logger.Errorw("Can not send message to server because we are disconnected")
	} else {
		w.connWriteLock.Lock()
		defer w.connWriteLock.Unlock()

		w.conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (w *Worker) GetMessageHandler() MessageHandler {
	return w.messageHandler
}

func (w *Worker) SetMessageHandler(h MessageHandler) {
	w.messageHandler = h
}

func (w *Worker) AddConnectedCallback(f func()) {
	w.connectedCallbacks = append(w.connectedCallbacks, f)
}

func (h *defaultMessageHandler) HandleMessage(message []byte) {
	var envelope messages.Envelope
	err := json.Unmarshal(message, &envelope)

	logger.Debugw("Received message from server", "message", string(message))

	if err != nil {
		logger.Errorw("Failed to unmarshal a message from server", "message", string(message))
		return
	}

	if envelope.Kind == messages.KindStartLoadTestRequest {
		var req messages.StartLoadTestRequest
		err = json.Unmarshal([]byte(envelope.Data), &req)

		if err != nil {
			logger.Errorw("Failed to unmarshal a message from server", "message", string(message))
			return
		}

		header := http.Header{}
		for _, line := range strings.Split(req.Header, "\n") {
			parts := strings.Split(line, ":")

			if len(parts) != 2 {
				continue
			}

			key := parts[0]
			value := strings.TrimLeft(parts[1], " ")

			header.Add(key, value)
		}

		rate := vegeta.Rate{Freq: int(req.Rate), Per: time.Second}
		duration := time.Duration(req.Duration) * time.Second
		targeter := vegeta.NewStaticTargeter(vegeta.Target{
			Method: req.Method,
			URL:    req.Url,
			Header: header,
			Body:   []byte(req.Body),
		})

		logger.Infow("Starting load test", "request", &req)

		h.worker.resetLoadTest()
		go h.worker.startLoadTest(targeter, rate, duration, "terjang")
	} else if envelope.Kind == messages.KindStopLoadTestRequest {

		logger.Infow("Stopping load test")
		h.worker.stopLoadTest()
	}
}

func (w *Worker) resetLoadTest() {
	w.attacker = vegeta.NewAttacker()
	w.metrics = vegeta.Metrics{}
}

func (w *Worker) startLoadTest(tr vegeta.Targeter, p vegeta.Pacer, du time.Duration, name string) {
	w.loadTestState = messages.WorkerStateRunning
	w.sendWorkerInfoToServer()

	for res := range w.attacker.Attack(tr, p, du, "terjang") {
		w.metrics.Add(res)
	}

	// Preserves state if it's stopped
	if w.loadTestState != messages.WorkerStateStopped {
		w.loadTestState = messages.WorkerStateDone
	}

	w.sendWorkerInfoToServer()

	logger.Infow("Finished load test")
}

func (w *Worker) stopLoadTest() {
	w.loadTestState = messages.WorkerStateStopped
	w.attacker.Stop()
}

func (w *Worker) LoopSendMetricsToServer() {
	for {
		if w.loadTestState == messages.WorkerStateRunning || w.loadTestState == messages.WorkerStateDone {
			w.SendMetricsToServer()
		}

		time.Sleep(1 * time.Second)
	}
}

func (w *Worker) SendMetricsToServer() {
	w.metrics.Close()

	workerMetrics := messages.WorkerLoadTestMetrics{}
	workerMetrics.Duration = w.metrics.Duration
	workerMetrics.Wait = w.metrics.Wait
	workerMetrics.Rate = w.metrics.Rate
	workerMetrics.Requests = w.metrics.Requests
	workerMetrics.Success = w.metrics.Success
	workerMetrics.Throughput = w.metrics.Throughput
	workerMetrics.Latencies = w.metrics.Latencies
	workerMetrics.BytesIn = w.metrics.BytesIn
	workerMetrics.BytesOut = w.metrics.BytesOut
	workerMetrics.StatusCodes = w.metrics.StatusCodes
	workerMetrics.Errors = w.metrics.Errors

	metrics, _ := json.Marshal(workerMetrics)
	envelope := messages.Envelope{Kind: messages.KindWorkerLoadTestMetrics, Data: string(metrics)}

	msg, _ := json.Marshal(envelope)
	w.SendMessageToServer(msg)
}

func (w *Worker) sendWorkerInfoToServer() {
	workerInfo := &messages.WorkerInfo{State: w.loadTestState}
	workerInfoJson, _ := json.Marshal(workerInfo)

	envelope := &messages.Envelope{Kind: messages.KindWorkerInfo, Data: string(workerInfoJson)}
	envelopeJson, _ := json.Marshal(envelope)

	w.SendMessageToServer(envelopeJson)
}
