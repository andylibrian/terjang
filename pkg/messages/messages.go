package messages

import (
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Envelope struct {
	Kind string
	Data string
}

const KindStartLoadTestRequest = "StartLoadTestRequest"
const KindStopLoadTestRequest = "StopLoadTestRequest"
const KindWorkerLoadTestMetrics = "WorkerLoadTestMetrics"

const KindServerInfo = "ServerInfo"
const KindWorkerInfo = "WorkerInfo"
const KindWorkersInfo = "WorkersInfo"

type StartLoadTestRequest struct {
	Method   string `json:"method"`
	Url      string `json:"url"`
	Duration uint64 `json:"duration,string"`
	// Rate per worker
	Rate   uint64 `json:"rate,string"`
	Header string `json:"header"`
	Body   string `json:"body"`
}

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

type WorkerState int

const WorkerStateNotStarted = WorkerState(0)
const WorkerStateRunning = WorkerState(1)
const WorkerStateDone = WorkerState(2)
const WorkerStateStopped = WorkerState(3)

type WorkerInfo struct {
	State WorkerState `json:"state"`
}

const ServerStateNotStarted = 0
const ServerStateRunning = 1
const ServerStateDone = 2
const ServerStateStopped = 3

type ServerInfo struct {
	NumOfWorkers int
	State        string
}
