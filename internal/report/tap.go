package report

import (
	"fmt"
	"io"

	"github.com/pawnkit/pawntest/internal/runner"
)

func TAP(w io.Writer, suite runner.Suite) error {
	fmt.Fprintf(w, "1..%d\n", len(suite.Results))
	for i, r := range suite.Results {
		status := "ok"
		if r.Status == runner.Fail || r.Status == runner.Error || r.Status == runner.XPass {
			status = "not ok"
		}
		directive := ""
		if r.Status == runner.Skip {
			directive = " # SKIP"
			if r.Message != "" {
				directive += " " + r.Message
			}
		}
		if r.Status == runner.XFail {
			directive = " # TODO expected failure"
		}
		fmt.Fprintf(w, "%s %d - %s%s\n", status, i+1, r.Name, directive)
		if r.Message != "" && r.Status != runner.Skip && r.Status != runner.XFail {
			fmt.Fprintf(w, "  ---\n  message: %q\n  ...\n", r.Message)
		}
	}
	return nil
}
