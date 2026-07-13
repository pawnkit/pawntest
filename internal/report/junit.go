package report

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/pawnkit/pawntest/internal/runner"
)

type junitSuite struct {
	XMLName  xml.Name    `xml:"testsuite"`
	Name     string      `xml:"name,attr,omitempty"`
	Tests    int         `xml:"tests,attr"`
	Failures int         `xml:"failures,attr"`
	Errors   int         `xml:"errors,attr"`
	Skipped  int         `xml:"skipped,attr"`
	Cases    []junitCase `xml:"testcase"`
}

type junitCase struct {
	Name      string        `xml:"name,attr"`
	Class     string        `xml:"classname,attr,omitempty"`
	File      string        `xml:"file,attr,omitempty"`
	Time      string        `xml:"time,attr"`
	Failure   *junitFailure `xml:"failure,omitempty"`
	Error     *junitFailure `xml:"error,omitempty"`
	Skipped   *junitSkipped `xml:"skipped,omitempty"`
	SystemOut string        `xml:"system-out,omitempty"`
}

type junitFailure struct {
	Message string `xml:"message,attr"`
}

type junitSkipped struct {
	Message string `xml:"message,attr,omitempty"`
}

func JUnit(w io.Writer, suite runner.Suite) error {
	js := junitSuite{Name: "pawntest", Tests: len(suite.Results)}
	for _, r := range suite.Results {
		c := junitCase{Name: r.Name, Class: r.Source, File: r.Source, Time: fmt.Sprintf("%.3f", r.Duration.Seconds())}
		if len(r.Warnings) > 0 {
			c.SystemOut = strings.Join(r.Warnings, "\n")
		}

		switch r.Status {
		case runner.Pass:
		case runner.Fail:
			js.Failures++
			c.Failure = &junitFailure{Message: r.Message}
		case runner.Error:
			js.Errors++
			c.Error = &junitFailure{Message: r.Message}
		case runner.Skip:
			js.Skipped++
			c.Skipped = &junitSkipped{Message: r.Message}
		case runner.XFail:
			js.Skipped++
			c.Skipped = &junitSkipped{Message: "expected failure: " + r.Message}
		case runner.XPass:
			js.Failures++
			c.Failure = &junitFailure{Message: r.Message}
		}

		js.Cases = append(js.Cases, c)
	}

	_, _ = io.WriteString(w, xml.Header)

	return xml.NewEncoder(w).Encode(js)
}
