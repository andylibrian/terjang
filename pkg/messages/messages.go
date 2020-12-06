package messages

type Envelope struct {
	Kind string
	Data string
}

const KindStartLoadTestRequest = "StartLoadTestRequest"
const KindStopLoadTestRequest = "StopLoadTestRequest"

type StartLoadTestRequest struct {
	Method   string `json:"method"`
	Url      string `json:"url"`
	Duration uint64 `json:"duration,string"`
	// Rate per worker
	Rate   uint64 `json:"rate,string"`
	Header string `json:"header"`
	Body   string `json:"body"`
}
