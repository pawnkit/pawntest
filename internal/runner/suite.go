package runner

import "time"

type Suite struct {
	Results  []Result      `json:"results"`
	Duration time.Duration `json:"duration"`
}

type Summary struct {
	Total   int `json:"total"`
	Passed  int `json:"passed"`
	Failed  int `json:"failed"`
	Skipped int `json:"skipped"`
	Errored int `json:"errored"`
	XFailed int `json:"xfailed"`
	XPassed int `json:"xpassed"`
}

func (s Suite) Summary() Summary {
	var out Summary
	for _, r := range s.Results {
		out.Total++

		switch r.Status {
		case Pass:
			out.Passed++
		case Fail:
			out.Failed++
		case Skip:
			out.Skipped++
		case Error:
			out.Errored++
		case XFail:
			out.XFailed++
		case XPass:
			out.XPassed++
		}
	}

	return out
}

func (s Suite) Failed() bool {
	for _, r := range s.Results {
		if r.Status == Fail || r.Status == Error || r.Status == XPass {
			return true
		}
	}

	return false
}
