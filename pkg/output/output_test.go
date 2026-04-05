package output

import (
	"bytes"
	"testing"
)

func TestNewPrinter(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		want   string
	}{
		{"json format", FormatJSON, "*output.JSONPrinter"},
		{"table format", FormatTable, "*output.TablePrinter"},
		{"simple format", FormatSimple, "*output.SimplePrinter"},
		{"default format", Format("unknown"), "*output.SimplePrinter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{Format: tt.format}
			printer := NewPrinter(opts)
			if printer == nil {
				t.Error("NewPrinter returned nil")
			}
		})
	}
}

func TestJSONPrinter_PrintIssues(t *testing.T) {
	printer := &JSONPrinter{opts: &Options{}}

	tests := []struct {
		name    string
		issues  interface{}
		wantErr bool
	}{
		{
			name: "empty list",
			issues: []map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "single issue",
			issues: []map[string]interface{}{
				{"number": 1, "state": "open", "title": "Test Issue"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printer.PrintIssues(&buf, tt.issues)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintIssues() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTablePrinter_PrintIssues(t *testing.T) {
	printer := &TablePrinter{opts: &Options{}}

	issues := []map[string]interface{}{
		{"number": 123, "state": "open", "title": "Test Issue", "author": "user1", "labels": []string{"bug"}},
	}

	var buf bytes.Buffer
	err := printer.PrintIssues(&buf, issues)
	if err != nil {
		t.Errorf("PrintIssues() error = %v", err)
	}

	if buf.Len() == 0 {
		t.Error("PrintIssues() produced no output")
	}
}

func TestSimplePrinter_PrintIssues(t *testing.T) {
	printer := &SimplePrinter{opts: &Options{}}

	issues := []map[string]interface{}{
		{"number": 123, "state": "open", "title": "Test Issue"},
	}

	var buf bytes.Buffer
	err := printer.PrintIssues(&buf, issues)
	if err != nil {
		t.Errorf("PrintIssues() error = %v", err)
	}

	if buf.Len() == 0 {
		t.Error("PrintIssues() produced no output")
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input string
		max   int
		want  string
	}{
		{"short", 10, "short"},
		{"this is a long string", 10, "this is..."},
		{"exact", 5, "exact"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := truncate(tt.input, tt.max)
			if got != tt.want {
				t.Errorf("truncate() = %v, want %v", got, tt.want)
			}
		})
	}
}