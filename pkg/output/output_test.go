package output

import (
	"bytes"
	"testing"
	"time"

	"gitcode.com/gitcode-cli/cli/api"
)

func TestNewPrinter(t *testing.T) {
	tests := []struct {
		name   string
		format Format
	}{
		{"json format", FormatJSON},
		{"table format", FormatTable},
		{"simple format", FormatSimple},
		{"default format", Format("unknown")},
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
			name:    "empty list",
			issues:  []map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "single issue",
			issues: []map[string]interface{}{
				{"number": 1, "state": "open", "title": "Test Issue"},
			},
			wantErr: false,
		},
		{
			name:    "api issue slice",
			issues:  []api.Issue{},
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

	tests := []struct {
		name    string
		issues  interface{}
		wantErr bool
	}{
		{
			name:    "map issues",
			issues:  []map[string]interface{}{{"number": 123, "state": "open", "title": "Test"}},
			wantErr: false,
		},
		{
			name:    "api issues",
			issues:  []api.Issue{{Number: "123", State: "open", Title: "Test"}},
			wantErr: false,
		},
		{
			name:    "api issue pointers",
			issues:  []*api.Issue{{Number: "123", State: "open", Title: "Test"}},
			wantErr: false,
		},
		{
			name:    "empty api issues",
			issues:  []api.Issue{},
			wantErr: false,
		},
		{
			name:    "invalid type",
			issues:  "invalid",
			wantErr: true,
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

func TestSimplePrinter_PrintIssues(t *testing.T) {
	printer := &SimplePrinter{opts: &Options{}}

	tests := []struct {
		name    string
		issues  interface{}
		wantErr bool
	}{
		{
			name:    "map issues",
			issues:  []map[string]interface{}{{"number": 123, "state": "open", "title": "Test"}},
			wantErr: false,
		},
		{
			name:    "api issues",
			issues:  []api.Issue{{Number: "123", State: "open", Title: "Test"}},
			wantErr: false,
		},
		{
			name:    "api issue pointers",
			issues:  []*api.Issue{{Number: "123", State: "open", Title: "Test"}},
			wantErr: false,
		},
		{
			name:    "empty issues",
			issues:  []api.Issue{},
			wantErr: false,
		},
		{
			name:    "invalid type",
			issues:  "invalid",
			wantErr: true,
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

func TestTablePrinter_PrintPRs(t *testing.T) {
	printer := &TablePrinter{opts: &Options{}}

	tests := []struct {
		name    string
		prs     interface{}
		wantErr bool
	}{
		{
			name:    "api pr pointers",
			prs:     []*api.PullRequest{{Number: 123, State: "open", Title: "Test PR"}},
			wantErr: false,
		},
		{
			name:    "map prs",
			prs:     []map[string]interface{}{{"number": 123, "state": "open", "title": "Test"}},
			wantErr: false,
		},
		{
			name:    "invalid type",
			prs:     "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printer.PrintPRs(&buf, tt.prs)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintPRs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSimplePrinter_PrintPRs(t *testing.T) {
	printer := &SimplePrinter{opts: &Options{}}

	tests := []struct {
		name    string
		prs     interface{}
		wantErr bool
	}{
		{
			name:    "api pr pointers",
			prs:     []*api.PullRequest{{Number: 123, State: "open", Title: "Test PR"}},
			wantErr: false,
		},
		{
			name:    "map prs",
			prs:     []map[string]interface{}{{"number": 123, "state": "open", "title": "Test"}},
			wantErr: false,
		},
		{
			name:    "invalid type",
			prs:     "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printer.PrintPRs(&buf, tt.prs)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintPRs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name   string
		format TimeFormat
	}{
		{"relative", TimeFormatRelative},
		{"absolute", TimeFormatAbsolute},
		{"default", TimeFormat("unknown")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTime(now, tt.format)
			if result == "" {
				t.Error("FormatTime returned empty string")
			}
		})
	}
}

func TestFormatRelativeTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		duration time.Duration
	}{
		{"just now", 0},
		{"minutes ago", 5 * time.Minute},
		{"hours ago", 2 * time.Hour},
		{"days ago", 2 * 24 * time.Hour},
		{"weeks ago", 2 * 7 * 24 * time.Hour},
		{"months ago", 60 * 24 * time.Hour},
		{"years ago", 400 * 24 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			past := now.Add(-tt.duration)
			result := formatRelativeTime(past)
			if result == "" {
				t.Error("formatRelativeTime returned empty string")
			}
		})
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
		{"", 10, ""},
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

func TestFormatAPILabels(t *testing.T) {
	tests := []struct {
		name   string
		labels []*api.Label
		want   string
	}{
		{"empty", nil, ""},
		{"single", []*api.Label{{Name: "bug"}}, "bug"},
		{"multiple", []*api.Label{{Name: "bug"}, {Name: "enhancement"}}, "bug, enhancement"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatAPILabels(tt.labels)
			if got != tt.want {
				t.Errorf("formatAPILabels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplatePrinter(t *testing.T) {
	tests := []struct {
		name       string
		template   string
		data       interface{}
		wantErr    bool
		wantOutput bool
	}{
		{
			name:       "simple template",
			template:   "{{.Title}}",
			data:       map[string]string{"Title": "Test"},
			wantErr:    false,
			wantOutput: true,
		},
		{
			name:       "invalid template",
			template:   "{{.Title",
			data:       map[string]string{"Title": "Test"},
			wantErr:    true,
			wantOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer, err := NewTemplatePrinter(tt.template, &Options{})
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			var buf bytes.Buffer
			err = printer.PrintIssues(&buf, tt.data)
			if err != nil {
				t.Errorf("PrintIssues error: %v", err)
			}

			if tt.wantOutput && buf.Len() == 0 {
				t.Error("expected output but got empty")
			}
		})
	}
}