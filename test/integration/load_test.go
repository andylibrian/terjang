package integration

import (
	"encoding/json"
	"io/ioutil"
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
	counter  uint32
	lastReq  *http.Request
	lastBody []byte
}

func (t *targetServer) helloHandler(w http.ResponseWriter, req *http.Request) {
	t.lastBody, _ = ioutil.ReadAll(req.Body)
	atomic.AddUint32(&t.counter, 1)
	t.lastReq = req

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello"))
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

	duration := 1
	rate := 10
	startLoadTestRequest := messages.StartLoadTestRequest{
		Method:   "POST",
		Url:      "http://127.0.0.1:10080/hello",
		Duration: uint64(duration),
		Rate:     uint64(rate),
		Header:   "X-load-test: MyLoadTest\nX-Foo: Bar",
		Body:     "thebody",
	}

	req, _ := json.Marshal(startLoadTestRequest)
	envelope, _ := json.Marshal(messages.Envelope{Kind: messages.KindStartLoadTestRequest, Data: string(req)})
	server.GetWorkerService().BroadcastMessageToWorkers(envelope)

	// Wait for the load test to complete.
	time.Sleep(time.Duration(duration) * time.Second)
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, rate*duration, int(target.counter))
	assert.Equal(t, "POST", target.lastReq.Method)
	assert.Equal(t, "thebody", string(target.lastBody))
	assert.Equal(t, "MyLoadTest", target.lastReq.Header.Get("X-Load-Test"))
	assert.Equal(t, "Bar", target.lastReq.Header.Get("X-Foo"))
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

	duration := 2
	rate := 10

	startLoadTestRequest := messages.StartLoadTestRequest{
		Method:   "GET",
		Url:      "http://127.0.0.1:10081/hello",
		Duration: uint64(duration),
		Rate:     uint64(rate),
	}

	req, _ := json.Marshal(startLoadTestRequest)
	envelope, _ := json.Marshal(messages.Envelope{Kind: messages.KindStartLoadTestRequest, Data: string(req)})
	server.GetWorkerService().BroadcastMessageToWorkers(envelope)

	time.Sleep(500 * time.Millisecond)

	envelope, _ = json.Marshal(messages.Envelope{Kind: messages.KindStopLoadTestRequest})
	server.GetWorkerService().BroadcastMessageToWorkers(envelope)

	// Sleep for the load test duration if it wouldn't be stopped.
	time.Sleep(time.Duration(duration) * time.Second)
	time.Sleep(200 * time.Millisecond)

	// Expect incomplete, but not zero
	assert.Less(t, int(target.counter), duration*rate)
	assert.Greater(t, int(target.counter), 0)
}
