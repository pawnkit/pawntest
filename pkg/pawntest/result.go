package pawntest

type Status string

const (
	Pass  Status = "pass"
	Fail  Status = "fail"
	Skip  Status = "skip"
	Error Status = "error"
	XFail Status = "xfail"
	XPass Status = "xpass"
)

type Public struct {
	Index int
	Name  string
}

type Result struct {
	Name     string   `json:"name"`
	Source   string   `json:"source,omitempty"`
	File     string   `json:"file"`
	Line     int      `json:"line,omitempty"`
	Status   Status   `json:"status"`
	Message  string   `json:"message,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
	Duration int64    `json:"duration_ms"`
}

type Suite struct {
	Results []Result `json:"results"`
}

func (s Suite) Failed() bool {
	for _, result := range s.Results {
		if result.Status == Fail || result.Status == Error || result.Status == XPass {
			return true
		}
	}

	return false
}
