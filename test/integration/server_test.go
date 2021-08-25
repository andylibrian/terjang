package integration

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andylibrian/terjang/pkg/messages"
	"github.com/andylibrian/terjang/pkg/server"
	"github.com/andylibrian/terjang/pkg/worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleServerInfo(t *testing.T) {

	server := server.NewServer()
	go server.Run("127.0.0.1:9029")
	defer server.Close()

	//Http request GET server_info
	req := httptest.NewRequest(http.MethodGet, "/api/v1/server_info", nil)
	w := httptest.NewRecorder()

	server.HandleServerInfo(w, req, nil)

	resp := w.Result()

	require.Equal(t, resp.StatusCode, http.StatusOK)

	var ServerResult messages.ServerInfo
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	json.Unmarshal(bodyBytes, &ServerResult)
	log.Printf("%s", string(bodyBytes))
	log.Printf("Num of Workers: %d, States: %s", ServerResult.NumOfWorkers, ServerResult.State)

	assert.Equal(t, 0, ServerResult.NumOfWorkers)
	assert.Equal(t, "NotStarted", ServerResult.State)

}

type WorkerStruct struct {
	Name    string                         `json:"name"`
	Metrics messages.WorkerLoadTestMetrics `json:"metrics"`
	State   string                         `json:"state"`
}

func TestHandleWorkersInfo(t *testing.T) {

	server := server.NewServer()
	go server.Run("127.0.0.1:9019")
	defer server.Close()

	worker := worker.NewWorker()
	worker.SetConnectRetryInterval(connectRetryInterval)

	// Wait for worker to be connected
	connected := make(chan struct{})
	worker.AddConnectedCallback(func() {
		connected <- struct{}{}
	})

	worker.SetName("worker1")
	go worker.Run("127.0.0.1:9019")
	<-connected

	//Http request GET server_info
	req := httptest.NewRequest(http.MethodGet, "/api/v1/worker_info", nil)
	w := httptest.NewRecorder()

	server.HandleWorkersInfo(w, req, nil)

	resp := w.Result()

	require.Equal(t, resp.StatusCode, http.StatusOK)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var result []WorkerStruct

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "worker1", result[0].Name)

}
