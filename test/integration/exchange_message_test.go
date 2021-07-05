package integration

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/andylibrian/terjang/pkg/messages"
	"github.com/andylibrian/terjang/pkg/server"
	"github.com/andylibrian/terjang/pkg/worker"
	"github.com/gorilla/websocket"

	"github.com/stretchr/testify/assert"
)

type serverMessageHandlerStub struct {
	handlerDelegate     server.MessageHandler
	messageCount        int
	metricsMessageCount int
	lastMetrics         *messages.WorkerLoadTestMetrics
}

func (s *serverMessageHandlerStub) HandleMessage(conn *websocket.Conn, message []byte) {
	s.messageCount++

	var envelope messages.Envelope
	json.Unmarshal(message, &envelope)

	if envelope.Kind == messages.KindWorkerLoadTestMetrics {
		s.metricsMessageCount++

		var metrics messages.WorkerLoadTestMetrics
		json.Unmarshal([]byte(envelope.Data), &metrics)

		s.lastMetrics = &metrics
	}
}

func (s *serverMessageHandlerStub) MessageCount() int {
	return s.messageCount
}

func (s *serverMessageHandlerStub) MetricsMessageCount() int {
	return s.metricsMessageCount
}

type workerMessageHandlerStub struct {
	handlerDelegate worker.MessageHandler
	messageCount    int
}

func (s *workerMessageHandlerStub) HandleMessage(message []byte) {
	s.messageCount++
}

func (s *workerMessageHandlerStub) MessageCount() int {
	return s.messageCount
}

func TestWorkerSendMessageToServer(t *testing.T) {
	server := server.NewServer()
	defaultServerMsgHandler := server.GetWorkerService().GetMessageHandler()

	serverMsgHandlerStub := serverMessageHandlerStub{handlerDelegate: defaultServerMsgHandler}
	server.GetWorkerService().SetMessageHandler(&serverMsgHandlerStub)

	go server.Run("127.0.0.1:9009")
	defer server.Close()

	worker := worker.NewWorker()
	worker.SetConnectRetryInterval(connectRetryInterval)
	// Wait for worker to be connected
	connected := make(chan struct{})
	worker.AddConnectedCallback(func() {
		connected <- struct{}{}
	})

	go worker.Run("127.0.0.1:9009")
	<-connected

	worker.SendMessageToServer([]byte("msg1"))
	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, 1, serverMsgHandlerStub.MessageCount())
}

func TestServerBroadcastMessagesToWorker(t *testing.T) {
	server := server.NewServer()
	go server.Run("127.0.0.1:9009")
	defer server.Close()

	worker1 := worker.NewWorker()
	defaultWorker1MsgHandler := worker1.GetMessageHandler()
	worker1MsgHandlerStub := workerMessageHandlerStub{handlerDelegate: defaultWorker1MsgHandler}
	worker1.SetMessageHandler(&worker1MsgHandlerStub)

	connected1 := make(chan struct{})
	worker1.AddConnectedCallback(func() {
		connected1 <- struct{}{}
	})

	worker2 := worker.NewWorker()
	defaultWorker2MsgHandler := worker2.GetMessageHandler()
	worker2MsgHandlerStub := workerMessageHandlerStub{handlerDelegate: defaultWorker2MsgHandler}
	worker2.SetMessageHandler(&worker2MsgHandlerStub)

	connected2 := make(chan struct{})
	worker2.AddConnectedCallback(func() {
		connected2 <- struct{}{}
	})

	go worker1.Run("127.0.0.1:9009")
	go worker2.Run("127.0.0.1:9009")

	<-connected1
	<-connected2

	server.GetWorkerService().BroadcastMessageToWorkers([]byte("msg1"))
	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, 1, worker1MsgHandlerStub.MessageCount())
	assert.Equal(t, 1, worker2MsgHandlerStub.MessageCount())
}
