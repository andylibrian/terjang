package messages

import (
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

/* Envelope is a struct type that has two fields;
* string constant Kind and string Data variables
* Kind values include KindStartLoadTestRequest, KindStopLoadTestRequest, KindWorkerLoadTestMetrics, KindServerInfo, KindWorkerInfo and KindWorkersInfo
************************************************/
type Envelope struct {
	Kind string `json:"kind"`
	Data string `json:"data"`
}

// KindStartLoadTestRequest is a const string flagging data for Starting Load Test Request
const KindStartLoadTestRequest = "StartLoadTestRequest"

// KindStopLoadTestRequest is a const string flagging data for Stopping Load Test Request
const KindStopLoadTestRequest = "StopLoadTestRequest"

// KindWorkerLoadTestMetrics is a const string flagging data for Worker when loading Test Metrics
const KindWorkerLoadTestMetrics = "WorkerLoadTestMetrics"

// KindServerInfo is a const string flagging data for Server Info
const KindServerInfo = "ServerInfo"

// KindWorkerInfo is a const string flagging data for Worker Info
const KindWorkerInfo = "WorkerInfo"

// KindWorkersInfo is a const string flagging data for Workers Info
const KindWorkersInfo = "WorkersInfo"

/* StartLoadTestRequest is a struct type that has 6 fields;
* string Method, string URL, uint64 Duration, uint64 Rate, string Header and string Body
***********************************************/
type StartLoadTestRequest struct {
	Method   string `json:"method"`
	URL      string `json:"url"`
	Duration uint64 `json:"duration,string"`
	// Rate per worker
	Rate   uint64 `json:"rate,string"`
	Header string `json:"header"`
	Body   string `json:"body"`
}

/* WorkerLoadTestMetrics is a struct type that has 11 fields;
* time.Duration Duration, time.Duration Wait, uint64 Requests, float64 Rate,  float64 Throughput,
* float64 Success, vegeta.LatencyMetrics Latency, vegeta.ByteMetrics BytesIn, vegeta.ByteMetrics BytesOut,
* map[string]int StatusCodes, []string Errors
***********************************************/
type WorkerLoadTestMetrics struct {
	// Duration is the duration of the attack.
	Duration time.Duration `json:"duration"`
	// Wait is the extra time waiting for responses from targets.
	Wait time.Duration `json:"wait"`
	// Requests is the total number of requests executed.
	Requests uint64 `json:"requests"`
	// Rate is the rate of sent requests per second.
	Rate float64 `json:"rate"`
	// Throughput is the rate of successful requests per second.
	Throughput float64 `json:"throughput"`
	// Success is the percentage of non-error responses.
	Success float64 `json:"success"`
	// Latencies holds computed request latency metrics.
	Latencies vegeta.LatencyMetrics `json:"latencies"`
	// BytesIn holds computed incoming byte metrics.
	BytesIn vegeta.ByteMetrics `json:"bytes_in"`
	// BytesOut holds computed outgoing byte metrics.
	BytesOut vegeta.ByteMetrics `json:"bytes_out"`
	// StatusCodes is a histogram of the responses' status codes.
	StatusCodes map[string]int `json:"status_codes"`
	// Errors is a set of unique errors returned by the targets during the attack.
	Errors []string `json:"errors"`
}

// WorkerState is a custom int type
type WorkerState int

// WorkerStateNotStarted is a const WorkerState custom type with value 0
const WorkerStateNotStarted = WorkerState(0)

// WorkerStateRunning is a const WorkerState custom type with value 1
const WorkerStateRunning = WorkerState(1)

// WorkerStateDone is a const WorkerState custom type with value 2
const WorkerStateDone = WorkerState(2)

// WorkerStateStopped is a const WorkerState custom type with value 3
const WorkerStateStopped = WorkerState(3)

// WorkerInfo is a struct type that has 1 field; WorkerState State
type WorkerInfo struct {
	State WorkerState `json:"state"`
}

// ServerStateNotStarted is a const string with value 0
const ServerStateNotStarted = 0

// ServerStateRunning is a const string with value 1
const ServerStateRunning = 1

// ServerStateDone is a const string with value 2
const ServerStateDone = 2

// ServerStateStopped is a const string with value 3
const ServerStateStopped = 3

// ServerInfo is a struct type that has 2 fields; int NumOfWorkers and string State
type ServerInfo struct {
	NumOfWorkers int    `json:"num_of_workers"`
	State        string `json:"state"`
}
