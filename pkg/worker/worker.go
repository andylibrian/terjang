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
)

type Worker struct {
	conn                 *websocket.Conn
	connWriteLock        sync.Mutex
	messageHandler       MessageHandler
	connectRetryInterval time.Duration
	isConnectedCh        chan struct{}
	attacker             *vegeta.Attacker
	metrics              vegeta.Metrics
	loadTestState        messages.WorkerState
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
		// should indicate error
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

		go h.worker.startLoadTest(targeter, rate, duration, "terjang")
	} else if envelope.Kind == messages.KindStopLoadTestRequest {
		h.worker.stopLoadTest()
	}
}

func (w *Worker) startLoadTest(tr vegeta.Targeter, p vegeta.Pacer, du time.Duration, name string) {
	w.loadTestState = messages.WorkerStateRunning
	w.sendWorkerInfoToServer()

	for res := range w.attacker.Attack(tr, p, du, "terjang") {
		w.metrics.Add(res)
	}

	w.loadTestState = messages.WorkerStateDone
	w.sendWorkerInfoToServer()
}

func (w *Worker) stopLoadTest() {
	w.attacker.Stop()

	w.loadTestState = messages.WorkerStateStopped
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
