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

		message := r.Message
		if r.Status == runner.Skip || r.Status == runner.XFail {
			message = ""
		}

		if r.Source != "" || message != "" || len(r.Warnings) > 0 {
			fmt.Fprintln(w, "  ---")

			if r.Source != "" {
				fmt.Fprintf(w, "  source: %q\n", r.Source)
			}

			if message != "" {
				fmt.Fprintf(w, "  message: %q\n", message)
			}

			if len(r.Warnings) > 0 {
				fmt.Fprintln(w, "  warnings:")

				for _, warning := range r.Warnings {
					fmt.Fprintf(w, "    - %q\n", warning)
				}
			}

			fmt.Fprintln(w, "  ...")
		}
	}

	return nil
}
