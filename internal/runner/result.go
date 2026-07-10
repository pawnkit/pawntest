package runner

import "time"

type Status string

const (
	Pass  Status = "pass"
	Fail  Status = "fail"
	Skip  Status = "skip"
	Error Status = "error"
	XFail Status = "xfail"
	XPass Status = "xpass"
)

type Result struct {
	Name     string        `json:"name"`
	Source   string        `json:"source,omitempty"`
	File     string        `json:"file"`
	Line     int           `json:"line,omitempty"`
	Status   Status        `json:"status"`
	Message  string        `json:"message,omitempty"`
	Duration time.Duration `json:"duration"`
}

func mergePhase(current Result, phase string, next Result) Result {
	if next.Status == Pass || next.Status == "" {
		return current
	}

	if current.Status == Pass || current.Status == "" {
		if phase != "test" && next.Message != "" {
			next.Message = phase + ": " + next.Message
		}

		return next
	}

	message := next.Message
	if message == "" {
		message = string(next.Status)
	}

	if current.Message == "" {
		current.Message = phase + ": " + message
	} else {
		current.Message += "\n" + phase + ": " + message
	}

	if statusSeverity(next.Status) > statusSeverity(current.Status) {
		current.Status = next.Status
	}

	return current
}

func statusSeverity(status Status) int {
	switch status {
	case Pass:
		return 0
	case Error:
		return 5
	case XPass:
		return 4
	case Fail:
		return 3
	case Skip:
		return 2
	case XFail:
		return 1
	}

	return 0
}
