package output

import (
	"bytes"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/api"
)

func TestPrintSimpleIncludesJobID(t *testing.T) {
	jobs := []api.WorkflowRunJob{
		{ID: "abc123", Name: "codecheck", Status: "COMPLETED", Steps: []api.WorkflowRunStep{{}, {}}},
		{ID: "def456", Name: "Antipoison", Status: "FAILED", Steps: []api.WorkflowRunStep{}},
	}

	var buf bytes.Buffer
	printer := &WorkflowJobListPrinter{opts: WorkflowJobListOptions{Format: FormatSimple}}
	if err := printer.Print(&buf, jobs); err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "abc123") {
		t.Fatalf("simple output = %q, missing job ID abc123", got)
	}
	if !strings.Contains(got, "def456") {
		t.Fatalf("simple output = %q, missing job ID def456", got)
	}
}

func TestPrintTableIncludesJobID(t *testing.T) {
	jobs := []api.WorkflowRunJob{
		{ID: "abc123", Name: "codecheck", Identifier: "check", Status: "COMPLETED", Sequence: 0, Steps: []api.WorkflowRunStep{}},
	}

	var buf bytes.Buffer
	printer := &WorkflowJobListPrinter{opts: WorkflowJobListOptions{Format: FormatTable}}
	if err := printer.Print(&buf, jobs); err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "ID") {
		t.Fatalf("table header = %q, missing ID column", got)
	}
	if !strings.Contains(got, "abc123") {
		t.Fatalf("table output = %q, missing job ID abc123", got)
	}
}

func TestPrintSimpleEmptyJobs(t *testing.T) {
	jobs := []api.WorkflowRunJob{}

	var buf bytes.Buffer
	printer := &WorkflowJobListPrinter{opts: WorkflowJobListOptions{Format: FormatSimple}}
	if err := printer.Print(&buf, jobs); err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	if buf.String() != "" {
		t.Fatalf("empty output = %q, want empty", buf.String())
	}
}

func TestPrintSimpleJobNameFallback(t *testing.T) {
	jobs := []api.WorkflowRunJob{
		{ID: "abc123", Name: "", Identifier: "fallback-name", Status: "COMPLETED", Steps: []api.WorkflowRunStep{}},
	}

	var buf bytes.Buffer
	printer := &WorkflowJobListPrinter{opts: WorkflowJobListOptions{Format: FormatSimple}}
	if err := printer.Print(&buf, jobs); err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "fallback-name") {
		t.Fatalf("simple output = %q, expected fallback to Identifier", got)
	}
}
