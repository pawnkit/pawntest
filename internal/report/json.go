package report

import (
	"encoding/json"
	"io"

	"github.com/pawnkit/pawntest/internal/runner"
)

type jsonSuite struct {
	Summary runner.Summary `json:"summary"`
	Results []jsonResult   `json:"results"`
}

type jsonResult struct {
	Name       string        `json:"name"`
	Source     string        `json:"source,omitempty"`
	File       string        `json:"file,omitempty"`
	Line       int           `json:"line,omitempty"`
	Status     runner.Status `json:"status"`
	Message    string        `json:"message,omitempty"`
	Warnings   []string      `json:"warnings,omitempty"`
	DurationMS int64         `json:"duration_ms"`
}

func JSON(w io.Writer, suite runner.Suite) error {
	out := jsonSuite{Summary: suite.Summary(), Results: make([]jsonResult, 0, len(suite.Results))}
	for _, result := range suite.Results {
		out.Results = append(out.Results, jsonResult{
			Name:       result.Name,
			Source:     result.Source,
			File:       result.File,
			Line:       result.Line,
			Status:     result.Status,
			Message:    result.Message,
			Warnings:   result.Warnings,
			DurationMS: result.Duration.Milliseconds(),
		})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	return enc.Encode(out)
}
