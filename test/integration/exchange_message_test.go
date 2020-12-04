package integration

import (
	"testing"
	"time"

	"github.com/andylibrian/terjang/pkg/server"
	"github.com/andylibrian/terjang/pkg/worker"

	"github.com/stretchr/testify/assert"
)

type Stub struct {
	handlerDelegate server.MessageHandler
	messageCount    int
}

func (s *Stub) HandleMessage(message []byte) {
	s.messageCount++
}

func (s *Stub) MessageCount() int {
	return s.messageCount
}

func TestWorkerSendMessageToServer(t *testing.T) {
	server := server.NewServer()
	defaultServerMsgHandler := server.GetWorkerService().GetMessageHandler()

	serverMsgHandlerStub := Stub{handlerDelegate: defaultServerMsgHandler}
	server.GetWorkerService().SetMessageHandler(&serverMsgHandlerStub)

	go server.Run()
	time.Sleep(1 * time.Second)

	worker := worker.NewWorker()
	go worker.Run()
	time.Sleep(1 * time.Second)

	worker.SendMessageToServer("msg1")
	time.Sleep(1 * time.Second)

	assert.Equal(t, 1, serverMsgHandlerStub.MessageCount())
}
