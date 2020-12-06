package integration

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/andylibrian/terjang/pkg/messages"
	"github.com/andylibrian/terjang/pkg/server"
	"github.com/andylibrian/terjang/pkg/worker"
	"github.com/stretchr/testify/assert"
)

type targetServer struct {
	counter uint32
}

func (t *targetServer) helloHandler(w http.ResponseWriter, req *http.Request) {
	atomic.AddUint32(&t.counter, 1)

	w.WriteHeader(http.StatusOK)
}

func (t *targetServer) listenAndServe(addr string) {
	handler := http.NewServeMux()
	handler.HandleFunc("/hello", t.helloHandler)
	target := &http.Server{Addr: addr, Handler: handler}
	target.ListenAndServe()
}

func TestStartLoadTest(t *testing.T) {
	target := targetServer{}
	go target.listenAndServe(":10080")

	server := server.NewServer()
	go server.Run()
	defer server.Close()

	worker := worker.NewWorker()
	worker.SetConnectRetryInterval(connectRetryInterval)
	go worker.Run()

	<-worker.IsConnectedCh()

	duration := uint64(1)
	startLoadTestRequest := messages.StartLoadTestRequest{
		Method:   "GET",
		Url:      "http://127.0.0.1:10080/hello",
		Duration: duration,
		Rate:     10,
	}

	req, _ := json.Marshal(startLoadTestRequest)
	envelope, _ := json.Marshal(messages.Envelope{Kind: messages.KindStartLoadTestRequest, Data: string(req)})
	server.GetWorkerService().BroadcastMessageToWorkers(envelope)

	// Wait for the load test to complete.
	time.Sleep(time.Duration(duration) * time.Second)
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 10, int(target.counter))
}

func TestStopLoadTest(t *testing.T) {
	target := targetServer{}
	go target.listenAndServe(":10081")

	server := server.NewServer()
	go server.Run()
	defer server.Close()

	worker := worker.NewWorker()
	worker.SetConnectRetryInterval(connectRetryInterval)
	go worker.Run()

	<-worker.IsConnectedCh()

	duration := uint64(2)
	rate := uint64(10)

	startLoadTestRequest := messages.StartLoadTestRequest{
		Method:   "GET",
		Url:      "http://127.0.0.1:10081/hello",
		Duration: duration,
		Rate:     rate,
	}

	req, _ := json.Marshal(startLoadTestRequest)
	envelope, _ := json.Marshal(messages.Envelope{Kind: messages.KindStartLoadTestRequest, Data: string(req)})
	server.GetWorkerService().BroadcastMessageToWorkers(envelope)

	time.Sleep(500 * time.Millisecond)

	envelope, _ = json.Marshal(messages.Envelope{Kind: messages.KindStopLoadTestRequest})
	server.GetWorkerService().BroadcastMessageToWorkers(envelope)

	// Sleep for the load test duration if it wouldn't be stopped.
	time.Sleep(3 * time.Second)

	// Expect incomplete, but not zero
	assert.Less(t, int(target.counter), int(duration*rate))
	assert.Greater(t, int(target.counter), 0)
}
