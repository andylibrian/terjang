package messages

import (
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// Envelope is a messaging wrapper for the communication between the server, workers, and UI.
type Envelope struct {
	Kind string `json:"kind"`
	Data string `json:"data"`
}

// KindStartLoadTestRequest is a kind that indicates a request to start a load test.
const KindStartLoadTestRequest = "StartLoadTestRequest"

// KindStopLoadTestRequest is a kind that indicates a request to stop a load test.
const KindStopLoadTestRequest = "StopLoadTestRequest"

// KindWorkerLoadTestMetrics is a kind that indicates the envelope contains load test metrics from worker.
const KindWorkerLoadTestMetrics = "WorkerLoadTestMetrics"

// KindServerInfo is a kind that indicates the envelope contains server info.
const KindServerInfo = "ServerInfo"

// KindWorkerInfo is a kind that indicates the envelope contains worker info.
const KindWorkerInfo = "WorkerInfo"

// KindWorkersInfo is a kind that indicates the envelope contains workers info.
const KindWorkersInfo = "WorkersInfo"

// StartLoadTestRequest is a struct type containing the detail of a load test request.
// It is sent from server to workers. Upon receiving this, workers should start
// running the load test.
type StartLoadTestRequest struct {
	Method   string `json:"method"`
	URL      string `json:"url"`
	Duration uint64 `json:"duration,string"`
	// Rate per worker
	Rate   uint64 `json:"rate,string"`
	Header string `json:"header"`
	Body   string `json:"body"`
}

// WorkerLoadTestMetrics is a struct type containing load test metrics from a worker.
// It is sent from workers to the server.
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

// WorkerState indicates worker state
type WorkerState int

// WorkerStateNotStarted indicates that the worker is not started.
const WorkerStateNotStarted = WorkerState(0)

// WorkerStateRunning indicates that the worker is running a load test.
const WorkerStateRunning = WorkerState(1)

// WorkerStateDone indicates that the worker is done running a load test.
const WorkerStateDone = WorkerState(2)

// WorkerStateStopped indicates that the worker has stopped running a load test before completed.
const WorkerStateStopped = WorkerState(3)

// WorkerInfo is a messaging type containing worker information.
type WorkerInfo struct {
	State WorkerState `json:"state"`
}

// ServerStateNotStarted indicates that the server sees that its workers are not started.
const ServerStateNotStarted = 0

// ServerStateRunning indicates that the server sees that its workers are running a load test.
const ServerStateRunning = 1

// ServerStateDone indicates that the server sees that its workers are done running a load test.
const ServerStateDone = 2

// ServerStateStopped indicates that the server sees that its workers have stopped running a load test before completed.
const ServerStateStopped = 3

// ServerInfo is a messaging  type containing server information.
type ServerInfo struct {
	NumOfWorkers int    `json:"num_of_workers"`
	State        string `json:"state"`
}
