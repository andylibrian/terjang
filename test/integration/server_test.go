package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andylibrian/terjang/pkg/server"
	"github.com/andylibrian/terjang/pkg/worker"
	"github.com/stretchr/testify/assert"
)

func TestHandleServerInfo(t *testing.T) {

	//Mock server
	server := server.NewServer()
	go server.Run("127.0.0.1:9029")
	defer server.Close()

	//Http request GET server_info
	req := httptest.NewRequest(http.MethodGet, "/api/v1/server_info", nil)
	w := httptest.NewRecorder()

	server.HandleServerInfo(w, req, nil)

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusOK)
}
func TestHandleWorkersInfo(t *testing.T) {

	//Mock server
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

	go worker.Run("127.0.0.1:9019")
	<-connected

	//Http request GET server_info
	req := httptest.NewRequest(http.MethodGet, "/api/v1/worker_info", nil)
	w := httptest.NewRecorder()

	server.HandleWorkersInfo(w, req, nil)

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusOK)
}