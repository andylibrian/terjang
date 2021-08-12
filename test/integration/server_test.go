package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

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

/*********************************
	{
		"name":"worker1",
		"state":""
	}
**********************************/
type Latencies struct {
	Total string
	Mean  string
	/**********
	50th string
	90th string
	95th string
	99th string
	**********/
	Max string
	Min string
}
type Bytes struct {
	Total string
	Mean  string
}

type Metrics struct {
	Duration     string
	Requests     string
	Rate         string
	Throughput   string
	Success      string
	Latencies    Latencies
	Bytes_in     Bytes
	Bytes_out    Bytes
	Status_codes string
	Errors       string
}
type WorkersStruct struct {
	Name    string
	Metrics Metrics
	State   string
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

	//assert.Equal(t, resp.StatusCode, http.StatusOK)

	if resp.StatusCode == http.StatusOK {
		//var WorkerInfo WorkersStruct

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)

		var result WorkersStruct

		//ref. https://play.golang.org/p/zOUMUNH4w9
		json.Unmarshal([]byte(bodyString), &result)

		fmt.Println(bodyString)
		fmt.Println(result)

		//var expectedName = "worker1"
		//assert.Equal(expectedName, WorkerInfo.name, "The two number of workers should be the same")

	} else {
		assert.Equal(resp.StatusCode, http.StatusOK)
	}
}
