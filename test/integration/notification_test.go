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

type stubWorker struct {
	Name     string                         `json:"name`
	Metrics  messages.WorkerLoadTestMetrics `json:"metrics"`
	StateStr string                         `json:"state"`
}

type stubNotificationClient struct {
	isConnectedCh   chan struct{}
	messages        []messages.Envelope
	serverInfoMsgs  []messages.Envelope
	workersInfoMsgs []messages.Envelope
}

func (s *stubNotificationClient) run(addr string) {
	serverURL := url.URL{Scheme: "ws", Host: addr, Path: "/notifications"}
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

		if envelope.Kind == messages.KindServerInfo {
			s.serverInfoMsgs = append(s.serverInfoMsgs, envelope)
		} else if envelope.Kind == messages.KindWorkersInfo {
			s.workersInfoMsgs = append(s.workersInfoMsgs, envelope)
		}
	}
}

func TestServerSendServerInfoNotification(t *testing.T) {
	server := server.NewServer()
	go server.Run("127.0.0.1:9029")
	defer server.Close()

	clientStub := stubNotificationClient{isConnectedCh: make(chan struct{})}
	go clientStub.run("127.0.0.1:9029")

	<-clientStub.isConnectedCh

	// Wait for a notification that comes every second
	time.Sleep(1*time.Second + 100*time.Millisecond)

	lastMsg := clientStub.serverInfoMsgs[len(clientStub.serverInfoMsgs)-1]
	assert.Equal(t, messages.KindServerInfo, lastMsg.Kind)

	var serverInfo messages.ServerInfo
	json.Unmarshal([]byte(lastMsg.Data), &serverInfo)

	assert.Equal(t, 0, serverInfo.NumOfWorkers)
	assert.Equal(t, "NotStarted", serverInfo.State)

	worker := worker.NewWorker()

	// Wait for worker to be connected
	connected := make(chan struct{})
	worker.AddConnectedCallback(func() {
		connected <- struct{}{}
	})

	go worker.Run("127.0.0.1:9029")
	<-connected

	time.Sleep(1*time.Second + 100*time.Millisecond)

	// assert server info
	lastMsg = clientStub.serverInfoMsgs[len(clientStub.serverInfoMsgs)-1]
	assert.Equal(t, messages.KindServerInfo, lastMsg.Kind)

	json.Unmarshal([]byte(lastMsg.Data), &serverInfo)

	assert.Equal(t, 1, serverInfo.NumOfWorkers)
	assert.Equal(t, "NotStarted", serverInfo.State)
}

func TestServerUpdateServerInfoNotification(t *testing.T) {
	target := targetServer{}
	go target.listenAndServe(":10090")

	server := server.NewServer()
	go server.Run("127.0.0.1:9039")
	defer server.Close()

	clientStub := stubNotificationClient{isConnectedCh: make(chan struct{})}
	go clientStub.run("127.0.0.1:9039")

	worker := worker.NewWorker()

	<-clientStub.isConnectedCh
	// Wait for worker to be connected
	connected := make(chan struct{})
	worker.AddConnectedCallback(func() {
		connected <- struct{}{}
	})

	go worker.Run("127.0.0.1:9039")
	<-connected

	duration := 2
	rate := 10
	startLoadTestRequest := messages.StartLoadTestRequest{
		Method:   "POST",
		URL:      "http://127.0.0.1:10090/hello",
		Duration: uint64(duration),
		Rate:     uint64(rate),
	}

	server.StartLoadTest(&startLoadTestRequest)

	// During load test
	time.Sleep(1 * time.Second)
	time.Sleep(100 * time.Millisecond)

	lastMsg := clientStub.serverInfoMsgs[len(clientStub.serverInfoMsgs)-1]
	assert.Equal(t, messages.KindServerInfo, lastMsg.Kind)

	var serverInfo messages.ServerInfo
	json.Unmarshal([]byte(lastMsg.Data), &serverInfo)

	assert.Equal(t, "Running", serverInfo.State)

	// After load test
	time.Sleep(3 * time.Second)
	time.Sleep(100 * time.Millisecond)

	lastMsg = clientStub.serverInfoMsgs[len(clientStub.serverInfoMsgs)-1]
	assert.Equal(t, messages.KindServerInfo, lastMsg.Kind)

	json.Unmarshal([]byte(lastMsg.Data), &serverInfo)

	assert.Equal(t, "Done", serverInfo.State)

	lastWorkersInfo := clientStub.workersInfoMsgs[len(clientStub.workersInfoMsgs)-1]

	var workersInfo []stubWorker
	json.Unmarshal([]byte(lastWorkersInfo.Data), &workersInfo)
	assert.Equal(t, uint64(duration*rate), workersInfo[0].Metrics.Requests)
	assert.Equal(t, float64(1), workersInfo[0].Metrics.Success)
}
