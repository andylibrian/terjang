package integration

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/andylibrian/terjang/pkg/messages"
	"github.com/andylibrian/terjang/pkg/server"
	"github.com/andylibrian/terjang/pkg/worker"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

type stubNotificationClient struct {
	isConnectedCh chan struct{}
	messages      []messages.Envelope
}

func (s *stubNotificationClient) run() {
	serverURL := url.URL{Scheme: "ws", Host: "127.0.0.1:9009", Path: "/notifications"}
	serverURLStr := serverURL.String()

	var conn *websocket.Conn
	var err error

	for i := 0; i < 3; i++ {
		conn, _, err = websocket.DefaultDialer.Dial(serverURLStr, nil)

		if err == nil {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	if err != nil {
		fmt.Printf("Error connecting to server %s\n", err)
		return
	}

	defer conn.Close()
	defer close(s.isConnectedCh)
	s.isConnectedCh <- struct{}{}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var envelope messages.Envelope
		err = json.Unmarshal(msg, &envelope)

		if err == nil {
			s.messages = append(s.messages, envelope)
		}
	}
}

func TestServerSendServerInfoNotification(t *testing.T) {
	server := server.NewServer()
	go server.Run()
	defer server.Close()

	clientStub := stubNotificationClient{isConnectedCh: make(chan struct{})}
	go clientStub.run()

	<-clientStub.isConnectedCh

	// Wait for a notification that comes every second
	time.Sleep(1*time.Second + 100*time.Millisecond)

	lastMsg := clientStub.messages[len(clientStub.messages)-1]
	assert.Equal(t, messages.KindServerInfo, lastMsg.Kind)

	var serverInfo messages.ServerInfo
	json.Unmarshal([]byte(lastMsg.Data), &serverInfo)

	assert.Equal(t, 0, serverInfo.NumOfWorkers)
	assert.Equal(t, "NotStarted", serverInfo.State)

	worker := worker.NewWorker()
	go worker.Run()

	<-worker.IsConnectedCh()

	time.Sleep(1*time.Second + 100*time.Millisecond)

	// assert server info
	lastMsg = clientStub.messages[len(clientStub.messages)-1]
	assert.Equal(t, messages.KindServerInfo, lastMsg.Kind)

	json.Unmarshal([]byte(lastMsg.Data), &serverInfo)

	assert.Equal(t, 1, serverInfo.NumOfWorkers)
	assert.Equal(t, "NotStarted", serverInfo.State)
}
