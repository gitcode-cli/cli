package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"gitcode.com/gitcode-cli/cli/api"
)

func TestIssueListPrinterSimple(t *testing.T) {
	printer, err := NewIssueListPrinter(IssueListOptions{
		Format:     FormatSimple,
		TimeFormat: TimeFormatAbsolute,
	})
	if err != nil {
		t.Fatalf("NewIssueListPrinter() error = %v", err)
	}

	issues := []api.Issue{{
		Number:    "123",
		State:     "open",
		Title:     "Test issue",
		UpdatedAt: api.FlexibleTime{Time: time.Date(2026, 4, 7, 12, 30, 0, 0, time.UTC)},
	}}

	var buf bytes.Buffer
	if err := printer.Print(&buf, issues); err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "#123") || !strings.Contains(got, "Test issue") || !strings.Contains(got, "2026-04-07 12:30") {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestIssueListPrinterTable(t *testing.T) {
	printer, err := NewIssueListPrinter(IssueListOptions{
		Format:     FormatTable,
		TimeFormat: TimeFormatRelative,
	})
	if err != nil {
		t.Fatalf("NewIssueListPrinter() error = %v", err)
	}

	issues := []api.Issue{{
		Number:    "123",
		State:     "closed",
		Title:     "Test issue",
		UpdatedAt: api.FlexibleTime{Time: time.Now().Add(-2 * time.Hour)},
	}}

	var buf bytes.Buffer
	if err := printer.Print(&buf, issues); err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "NUMBER") || !strings.Contains(got, "#123") || !strings.Contains(got, "closed") {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestIssueListPrinterTemplate(t *testing.T) {
	printer, err := NewIssueListPrinter(IssueListOptions{
		Template: "{{range .}}{{upper .Title}} {{time .UpdatedAt}}{{end}}",
	})
	if err != nil {
		t.Fatalf("NewIssueListPrinter() error = %v", err)
	}

	issues := []api.Issue{{
		Title:     "Test issue",
		UpdatedAt: api.FlexibleTime{Time: time.Date(2026, 4, 7, 12, 30, 0, 0, time.UTC)},
	}}

	var buf bytes.Buffer
	if err := printer.Print(&buf, issues); err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "TEST ISSUE") || !strings.Contains(got, "2026-04-07 12:30") {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestIssueListPrinterInvalidTemplate(t *testing.T) {
	if _, err := NewIssueListPrinter(IssueListOptions{Template: "{{range ."}); err == nil {
		t.Fatal("expected template parse error")
	}
}
