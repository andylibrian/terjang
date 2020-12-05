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

var counter uint32

func helloHandler(w http.ResponseWriter, req *http.Request) {
	atomic.AddUint32(&counter, 1)

	w.WriteHeader(http.StatusOK)
}

func startTargetServer() {
	handler := http.NewServeMux()
	handler.HandleFunc("/hello", helloHandler)
	target := &http.Server{Addr: ":10080", Handler: handler}
	target.ListenAndServe()
}

func TestStartLoadTest(t *testing.T) {
	go startTargetServer()

	server := server.NewServer()
	go server.Run()
	defer server.Close()

	worker := worker.NewWorker()
	worker.SetConnectRetryInterval(connectRetryInterval)
	go worker.Run()

	<-worker.IsConnectedCh()

	startLoadTestRequest := messages.StartLoadTestRequest{
		Method:   "GET",
		Url:      "http://127.0.0.1:10080/hello",
		Duration: 2,
		Rate:     10,
	}

	req, _ := json.Marshal(startLoadTestRequest)
	envelope, _ := json.Marshal(messages.Envelope{Kind: messages.KindStartLoadTestRequest, Data: string(req)})
	server.GetWorkerService().BroadcastMessageToWorkers(envelope)

	time.Sleep(3 * time.Second)

	assert.Equal(t, 20, int(counter))
}
