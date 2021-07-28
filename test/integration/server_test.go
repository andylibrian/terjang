package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andylibrian/terjang/pkg/server"
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

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got status %d but wanted %d", resp.StatusCode, http.StatusOK)
	}
}
func TestHandleWorkersInfo(t *testing.T) {

	//Mock server
	server := server.NewServer()
	go server.Run("127.0.0.1:9029")
	defer server.Close()

	//Http request GET server_info
	req := httptest.NewRequest(http.MethodGet, "/api/v1/worker_info", nil)
	w := httptest.NewRecorder()

	server.HandleWorkersInfo(w, req, nil)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got status %d but wanted %d", resp.StatusCode, http.StatusOK)
	}
}
