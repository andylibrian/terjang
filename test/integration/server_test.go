package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andylibrian/terjang/pkg/messages"
	"github.com/andylibrian/terjang/pkg/server"
	"github.com/andylibrian/terjang/pkg/worker"
	"github.com/stretchr/testify/assert"
)

/*********************************
	{
		"num_of_workers":0,
		"state":"NotStarted"
	}
**********************************/
type ServerStruct struct {
	Num_of_workers int
	State          string
}

func TestHandleServerInfo(t *testing.T) {

	assert := assert.New(t)

	//Mock server
	server := server.NewServer()
	go server.Run("127.0.0.1:9029")
	defer server.Close()

	//Http request GET server_info
	req := httptest.NewRequest(http.MethodGet, "/api/v1/server_info", nil)
	w := httptest.NewRecorder()

	server.HandleServerInfo(w, req, nil)

	resp := w.Result()

	if resp.StatusCode == http.StatusOK {
		var ServerInfo ServerStruct
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		json.Unmarshal([]byte(bodyString), &ServerInfo)
		fmt.Println(bodyString)
		log.Printf("Num of Workers: %d, States: %s", ServerInfo.Num_of_workers, ServerInfo.State)

		var expectedNumOfWorkers = 0
		var expectedState = "NotStarted"
		assert.Equal(expectedNumOfWorkers, ServerInfo.Num_of_workers, "The two number of workers should be the same")
		assert.Equal(expectedState, ServerInfo.State, "The two states should be the same")
	} else {
		assert.Equal(resp.StatusCode, http.StatusOK)
	}

}

type WorkersStruct struct {
	Name    string                         `json:"name"`
	Metrics messages.WorkerLoadTestMetrics `json:"metrics"`
	State   string                         `json:"state"`
}

func TestHandleWorkersInfo(t *testing.T) {

	assert := assert.New(t)

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

	worker.SetName("worker1")
	go worker.Run("127.0.0.1:9019")
	<-connected

	//Http request GET server_info
	req := httptest.NewRequest(http.MethodGet, "/api/v1/worker_info", nil)
	w := httptest.NewRecorder()

	server.HandleWorkersInfo(w, req, nil)

	resp := w.Result()

	if resp.StatusCode == http.StatusOK {

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)

		var result []WorkersStruct

		err = json.Unmarshal([]byte(bodyString), &result)
		if err != nil {
			panic(err)
		}

		var expectedName = "worker1"
		assert.Equal(expectedName, result[0].Name, "The two number of workers should be the same")

	} else {
		assert.Equal(resp.StatusCode, http.StatusOK)
	}
}
