package integration

import (
	"io"
	"io/ioutil"
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

	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ping")
	}

	//Http request GET server_info
	req := httptest.NewRequest(http.MethodGet, "/api/v1/server_info", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got status %d but wanted %d", resp.StatusCode, http.StatusOK)
	}

	t.Log(body)
}
