package integration

import (
	"testing"
	"time"

	"github.com/andylibrian/terjang/pkg/server"
	"github.com/andylibrian/terjang/pkg/worker"

	"github.com/stretchr/testify/assert"
)

type serverMessageHandlerStub struct {
	handlerDelegate server.MessageHandler
	messageCount    int
}

func (s *serverMessageHandlerStub) HandleMessage(message []byte) {
	s.messageCount++
}

func (s *serverMessageHandlerStub) MessageCount() int {
	return s.messageCount
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

	go server.Run()
	defer server.Close()

	time.Sleep(1 * time.Second)

	worker := worker.NewWorker()
	go worker.Run()
	time.Sleep(1 * time.Second)

	worker.SendMessageToServer([]byte("msg1"))
	time.Sleep(1 * time.Second)

	assert.Equal(t, 1, serverMsgHandlerStub.MessageCount())
}

func TestServerBroadcastMessagesToWorker(t *testing.T) {
	server := server.NewServer()
	go server.Run()
	defer server.Close()

	time.Sleep(1 * time.Second)

	worker1 := worker.NewWorker()
	defaultWorker1MsgHandler := worker1.GetMessageHandler()
	worker1MsgHandlerStub := workerMessageHandlerStub{handlerDelegate: defaultWorker1MsgHandler}
	worker1.SetMessageHandler(&worker1MsgHandlerStub)
	go worker1.Run()

	worker2 := worker.NewWorker()
	defaultWorker2MsgHandler := worker2.GetMessageHandler()
	worker2MsgHandlerStub := workerMessageHandlerStub{handlerDelegate: defaultWorker2MsgHandler}
	worker2.SetMessageHandler(&worker2MsgHandlerStub)
	go worker2.Run()
	time.Sleep(1 * time.Second)

	server.GetWorkerService().BroadcastMessageToWorkers([]byte("msg1"))
	time.Sleep(1 * time.Second)

	assert.Equal(t, 1, worker1MsgHandlerStub.MessageCount())
	assert.Equal(t, 1, worker2MsgHandlerStub.MessageCount())
}
