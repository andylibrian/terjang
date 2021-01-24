package integration

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/andylibrian/terjang/pkg/messages"
	"github.com/andylibrian/terjang/pkg/server"
	"github.com/andylibrian/terjang/pkg/worker"
	"github.com/stretchr/testify/assert"
)

func TestWorkerSendMetricsDuringLoadTest(t *testing.T) {
	target := targetServer{}
	go target.listenAndServe(":10080")

	server := server.NewServer()
	defaultServerMsgHandler := server.GetWorkerService().GetMessageHandler()

	serverMsgHandlerStub := serverMessageHandlerStub{handlerDelegate: defaultServerMsgHandler}
	server.GetWorkerService().SetMessageHandler(&serverMsgHandlerStub)
	go server.Run("127.0.0.1:9049")
	defer server.Close()

	worker := worker.NewWorker()
	worker.SetConnectRetryInterval(connectRetryInterval)

	// Wait for worker to be connected
	connected := make(chan struct{})
	worker.AddConnectedCallback(func() {
		connected <- struct{}{}
	})

	go worker.Run("127.0.0.1:9049")
	<-connected

	duration := 2
	rate := 10
	startLoadTestRequest := messages.StartLoadTestRequest{
		Method:   "POST",
		URL:      "http://127.0.0.1:10080/hello",
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

	assert.Greater(t, serverMsgHandlerStub.MetricsMessageCount(), 0)
	assert.Less(t, serverMsgHandlerStub.MetricsMessageCount(), 3)

	assert.Greater(t, serverMsgHandlerStub.lastMetrics.Duration.Seconds(), float64(0))
	assert.Greater(t, serverMsgHandlerStub.lastMetrics.BytesIn.Total, uint64(0))
	assert.Equal(t, serverMsgHandlerStub.lastMetrics.Success, float64(1))
}
