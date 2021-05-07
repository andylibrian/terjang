package messages

import (
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// Envelope is ...
type Envelope struct {
	Kind string `json:"kind"`
	Data string `json:"data"`
}

// KindStartLoadTestRequest is ...
const KindStartLoadTestRequest = "StartLoadTestRequest"

// KindStopLoadTestRequest is ...
const KindStopLoadTestRequest = "StopLoadTestRequest"

// KindWorkerLoadTestMetrics is ...
const KindWorkerLoadTestMetrics = "WorkerLoadTestMetrics"

// KindServerInfo is ...
const KindServerInfo = "ServerInfo"

// KindWorkerInfo is ...
const KindWorkerInfo = "WorkerInfo"

// KindWorkersInfo is ...
const KindWorkersInfo = "WorkersInfo"

// StartLoadTestRequest is ...
type StartLoadTestRequest struct {
	Method   string `json:"method"`
	URL      string `json:"url"`
	Duration uint64 `json:"duration,string"`
	// Rate per worker
	Rate   uint64 `json:"rate,string"`
	Header string `json:"header"`
	Body   string `json:"body"`
}

// WorkerLoadTestMetrics is ...
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

// WorkerState is ...
type WorkerState int

// WorkerStateNotStarted is ...
const WorkerStateNotStarted = WorkerState(0)

// WorkerStateRunning is ...
const WorkerStateRunning = WorkerState(1)

// WorkerStateDone is ...
const WorkerStateDone = WorkerState(2)

// WorkerStateStopped is ...
const WorkerStateStopped = WorkerState(3)

// WorkerInfo is ...
type WorkerInfo struct {
	State WorkerState `json:"state"`
}

// ServerStateNotStarted is ...
const ServerStateNotStarted = 0

// ServerStateRunning is ...
const ServerStateRunning = 1

// ServerStateDone is ...
const ServerStateDone = 2

// ServerStateStopped is ...
const ServerStateStopped = 3

// ServerInfo is ...
type ServerInfo struct {
	NumOfWorkers int    `json:"num_of_workers"`
	State        string `json:"state"`
}
