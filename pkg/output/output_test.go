package output

import (
	"bytes"
	"fmt"
	"strings"
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
		{
			name:    "api issues with labels",
			issues:  []api.Issue{{Number: "123", State: "open", Title: "Test", Labels: []*api.Label{{Name: "bug"}}}},
			wantErr: false,
		},
		{
			name:    "api issues with user",
			issues:  []api.Issue{{Number: "123", State: "open", Title: "Test", User: &api.User{Login: "testuser"}}},
			wantErr: false,
		},
		{
			name:    "api issue pointers with labels",
			issues:  []*api.Issue{{Number: "123", State: "open", Title: "Test", Labels: []*api.Label{{Name: "enhancement"}, {Name: "priority"}}}},
			wantErr: false,
		},
		{
			name:    "api issue pointers with user",
			issues:  []*api.Issue{{Number: "123", State: "open", Title: "Test", User: &api.User{Login: "dev"}}},
			wantErr: false,
		},
		{
			name:    "map issues with all fields",
			issues:  []map[string]interface{}{{"number": 999999, "state": "closed", "title": "Very Long Title That Should Be Truncated", "labels": []string{"bug", "critical"}}},
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
			name:    "api pr pointers with user",
			prs:     []*api.PullRequest{{Number: 123, State: "open", Title: "Test PR", User: &api.User{Login: "developer"}}},
			wantErr: false,
		},
		{
			name:    "api pr pointers without user",
			prs:     []*api.PullRequest{{Number: 456, State: "closed", Title: "Another PR"}},
			wantErr: false,
		},
		{
			name:    "empty api pr pointers",
			prs:     []*api.PullRequest{},
			wantErr: false,
		},
		{
			name:    "map prs",
			prs:     []map[string]interface{}{{"number": 123, "state": "open", "title": "Test"}},
			wantErr: false,
		},
		{
			name:    "map prs with all fields",
			prs:     []map[string]interface{}{{"number": 123, "state": "open", "title": "Test", "author": "dev", "review_status": "approved"}},
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

func TestJSONPrinter_PrintPRs(t *testing.T) {
	printer := &JSONPrinter{opts: &Options{}}

	tests := []struct {
		name    string
		prs     interface{}
		wantErr bool
	}{
		{
			name:    "empty list",
			prs:     []map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    "single pr",
			prs:     []map[string]interface{}{{"number": 123, "state": "open", "title": "Test PR"}},
			wantErr: false,
		},
		{
			name:    "api pr pointers",
			prs:     []*api.PullRequest{{Number: 123, State: "open", Title: "Test PR"}},
			wantErr: false,
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

func TestJSONPrinter_PrintRepos(t *testing.T) {
	printer := &JSONPrinter{opts: &Options{}}

	tests := []struct {
		name    string
		repos   interface{}
		wantErr bool
	}{
		{
			name:    "empty list",
			repos:   []map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    "single repo",
			repos:   []map[string]interface{}{{"name": "test-repo", "visibility": "public"}},
			wantErr: false,
		},
		{
			name:    "multiple repos",
			repos:   []map[string]interface{}{{"name": "repo1"}, {"name": "repo2"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printer.PrintRepos(&buf, tt.repos)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintRepos() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && buf.Len() == 0 {
				t.Error("expected output but got empty")
			}
		})
	}
}

func TestJSONPrinter_PrintReleases(t *testing.T) {
	printer := &JSONPrinter{opts: &Options{}}

	tests := []struct {
		name     string
		releases interface{}
		wantErr  bool
	}{
		{
			name:     "empty list",
			releases: []map[string]interface{}{},
			wantErr:  false,
		},
		{
			name:     "single release",
			releases: []map[string]interface{}{{"tag": "v1.0.0", "name": "First Release"}},
			wantErr:  false,
		},
		{
			name:     "multiple releases",
			releases: []map[string]interface{}{{"tag": "v1.0.0"}, {"tag": "v2.0.0"}},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printer.PrintReleases(&buf, tt.releases)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintReleases() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && buf.Len() == 0 {
				t.Error("expected output but got empty")
			}
		})
	}
}

func TestJSONPrinter_PrintOne(t *testing.T) {
	printer := &JSONPrinter{opts: &Options{}}

	tests := []struct {
		name    string
		item    interface{}
		wantErr bool
	}{
		{
			name:    "single item",
			item:    map[string]interface{}{"number": 1, "title": "Test Issue"},
			wantErr: false,
		},
		{
			name:    "complex item",
			item:    map[string]interface{}{"id": "123", "name": "test", "value": 100},
			wantErr: false,
		},
		{
			name:    "empty item",
			item:    map[string]interface{}{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printer.PrintOne(&buf, tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintOne() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && buf.Len() == 0 {
				t.Error("expected output but got empty")
			}
		})
	}
}

func TestTablePrinter_PrintRepos(t *testing.T) {
	printer := &TablePrinter{opts: &Options{}}

	tests := []struct {
		name    string
		repos   interface{}
		wantErr bool
	}{
		{
			name:    "empty list",
			repos:   []map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    "single repo",
			repos:   []map[string]interface{}{{"name": "test-repo", "visibility": "public", "description": "A test repo", "language": "Go"}},
			wantErr: false,
		},
		{
			name:    "multiple repos",
			repos:   []map[string]interface{}{{"name": "repo1", "visibility": "private"}, {"name": "repo2", "visibility": "public"}},
			wantErr: false,
		},
		{
			name:    "invalid type",
			repos:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printer.PrintRepos(&buf, tt.repos)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintRepos() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTablePrinter_PrintReleases(t *testing.T) {
	printer := &TablePrinter{opts: &Options{}}

	tests := []struct {
		name     string
		releases interface{}
		wantErr  bool
	}{
		{
			name:     "empty list",
			releases: []map[string]interface{}{},
			wantErr:  false,
		},
		{
			name:     "single release",
			releases: []map[string]interface{}{{"tag": "v1.0.0", "name": "First Release", "type": "release", "created_at": "2026-01-01"}},
			wantErr:  false,
		},
		{
			name:     "multiple releases",
			releases: []map[string]interface{}{{"tag": "v1.0.0"}, {"tag": "v2.0.0"}},
			wantErr:  false,
		},
		{
			name:     "invalid type",
			releases: "invalid",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printer.PrintReleases(&buf, tt.releases)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintReleases() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTablePrinter_PrintOne(t *testing.T) {
	printer := &TablePrinter{opts: &Options{}}

	tests := []struct {
		name    string
		item    interface{}
		wantErr bool
	}{
		{
			name:    "single item",
			item:    map[string]interface{}{"number": 1, "title": "Test Issue"},
			wantErr: false,
		},
		{
			name:    "complex item",
			item:    map[string]interface{}{"id": "123", "name": "test", "status": "open"},
			wantErr: false,
		},
		{
			name:    "empty item",
			item:    map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    "invalid type",
			item:    "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printer.PrintOne(&buf, tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSimplePrinter_PrintRepos(t *testing.T) {
	printer := &SimplePrinter{opts: &Options{}}

	tests := []struct {
		name    string
		repos   interface{}
		wantErr bool
	}{
		{
			name:    "empty list",
			repos:   []map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    "single repo",
			repos:   []map[string]interface{}{{"name": "test-repo", "description": "A test repo"}},
			wantErr: false,
		},
		{
			name:    "multiple repos",
			repos:   []map[string]interface{}{{"name": "repo1", "description": "desc1"}, {"name": "repo2", "description": "desc2"}},
			wantErr: false,
		},
		{
			name:    "invalid type",
			repos:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printer.PrintRepos(&buf, tt.repos)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintRepos() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSimplePrinter_PrintReleases(t *testing.T) {
	printer := &SimplePrinter{opts: &Options{}}

	tests := []struct {
		name     string
		releases interface{}
		wantErr  bool
	}{
		{
			name:     "empty list",
			releases: []map[string]interface{}{},
			wantErr:  false,
		},
		{
			name:     "single release",
			releases: []map[string]interface{}{{"tag": "v1.0.0", "name": "First Release"}},
			wantErr:  false,
		},
		{
			name:     "multiple releases",
			releases: []map[string]interface{}{{"tag": "v1.0.0", "name": "rel1"}, {"tag": "v2.0.0", "name": "rel2"}},
			wantErr:  false,
		},
		{
			name:     "invalid type",
			releases: "invalid",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printer.PrintReleases(&buf, tt.releases)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintReleases() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSimplePrinter_PrintOne(t *testing.T) {
	printer := &SimplePrinter{opts: &Options{}}

	tests := []struct {
		name    string
		item    interface{}
		wantErr bool
	}{
		{
			name:    "single item",
			item:    map[string]interface{}{"number": 1, "title": "Test Issue"},
			wantErr: false,
		},
		{
			name:    "complex item",
			item:    map[string]interface{}{"id": "123", "name": "test", "value": 100},
			wantErr: false,
		},
		{
			name:    "empty item",
			item:    map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    "invalid type",
			item:    "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printer.PrintOne(&buf, tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTemplatePrinter_PrintRepos(t *testing.T) {
	template := "{{range .}}{{.name}}\n{{end}}"
	printer, err := NewTemplatePrinter(template, &Options{})
	if err != nil {
		t.Fatalf("failed to create template printer: %v", err)
	}

	tests := []struct {
		name    string
		repos   interface{}
		wantErr bool
	}{
		{
			name:    "single repo",
			repos:   []map[string]interface{}{{"name": "test-repo"}},
			wantErr: false,
		},
		{
			name:    "multiple repos",
			repos:   []map[string]interface{}{{"name": "repo1"}, {"name": "repo2"}},
			wantErr: false,
		},
		{
			name:    "empty list",
			repos:   []map[string]interface{}{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printer.PrintRepos(&buf, tt.repos)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintRepos() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTemplatePrinter_PrintReleases(t *testing.T) {
	template := "{{range .}}{{.tag}}: {{.name}}\n{{end}}"
	printer, err := NewTemplatePrinter(template, &Options{})
	if err != nil {
		t.Fatalf("failed to create template printer: %v", err)
	}

	tests := []struct {
		name     string
		releases interface{}
		wantErr  bool
	}{
		{
			name:     "single release",
			releases: []map[string]interface{}{{"tag": "v1.0.0", "name": "First Release"}},
			wantErr:  false,
		},
		{
			name:     "multiple releases",
			releases: []map[string]interface{}{{"tag": "v1.0.0", "name": "rel1"}, {"tag": "v2.0.0", "name": "rel2"}},
			wantErr:  false,
		},
		{
			name:     "empty list",
			releases: []map[string]interface{}{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printer.PrintReleases(&buf, tt.releases)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintReleases() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTemplatePrinter_PrintOne(t *testing.T) {
	template := "{{range $k, $v := .}}{{$k}}: {{$v}}\n{{end}}"
	printer, err := NewTemplatePrinter(template, &Options{})
	if err != nil {
		t.Fatalf("failed to create template printer: %v", err)
	}

	tests := []struct {
		name    string
		item    interface{}
		wantErr bool
	}{
		{
			name:    "single item",
			item:    map[string]interface{}{"number": 1, "title": "Test Issue"},
			wantErr: false,
		},
		{
			name:    "complex item",
			item:    map[string]interface{}{"id": "123", "name": "test"},
			wantErr: false,
		},
		{
			name:    "empty item",
			item:    map[string]interface{}{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printer.PrintOne(&buf, tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrintOne() error = %v, wantErr %v", err, tt.wantErr)
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

func TestTemplatePrinter_PrintPRs(t *testing.T) {
	template := "{{range .}}#{{.Number}}: {{.Title}}\n{{end}}"
	printer, err := NewTemplatePrinter(template, &Options{})
	if err != nil {
		t.Fatalf("failed to create template printer: %v", err)
	}

	tests := []struct {
		name    string
		prs     interface{}
		wantErr bool
	}{
		{
			name:    "single pr",
			prs:     []map[string]interface{}{{"Number": 123, "Title": "Test PR"}},
			wantErr: false,
		},
		{
			name:    "multiple prs",
			prs:     []map[string]interface{}{{"Number": 1, "Title": "PR1"}, {"Number": 2, "Title": "PR2"}},
			wantErr: false,
		},
		{
			name:    "empty list",
			prs:     []map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    "api pr pointers",
			prs:     []*api.PullRequest{{Number: 123, State: "open", Title: "Test PR"}},
			wantErr: false,
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

func TestTemplateFuncs(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     interface{}
		want     string
	}{
		{
			name:     "upper function",
			template: "{{upper .Title}}",
			data:     map[string]interface{}{"Title": "test"},
			want:     "test",
		},
		{
			name:     "upper function truncation",
			template: "{{upper .Title}}",
			data:     map[string]interface{}{"Title": "this is a very long string that exceeds fifty characters limit"},
			want:     "this is a very long string that exceeds fif",
		},
		{
			name:     "lower function",
			template: "{{lower .Title}}",
			data:     map[string]interface{}{"Title": "TEST"},
			want:     "test",
		},
		{
			name:     "lower function mixed",
			template: "{{lower .Title}}",
			data:     map[string]interface{}{"Title": "TeSt"},
			want:     "test",
		},
		{
			name:     "trunc function",
			template: "{{trunc .Title 10}}",
			data:     map[string]interface{}{"Title": "short"},
			want:     "short",
		},
		{
			name:     "trunc function truncate",
			template: "{{trunc .Title 10}}",
			data:     map[string]interface{}{"Title": "this is a long string"},
			want:     "this is...",
		},
		{
			name:     "json function",
			template: "{{json .}}",
			data:     map[string]interface{}{"Number": 123, "Title": "Test"},
			want:     "map[Number:123 Title:Test]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer, err := NewTemplatePrinter(tt.template, &Options{})
			if err != nil {
				t.Fatalf("failed to create template printer: %v", err)
			}

			var buf bytes.Buffer
			err = printer.PrintIssues(&buf, tt.data)
			if err != nil {
				t.Errorf("PrintIssues error: %v", err)
				return
			}

			result := buf.String()
			if !strings.Contains(result, tt.want) {
				t.Errorf("template output = %q, want to contain %q", result, tt.want)
			}
		})
	}
}

func TestFormatLabels(t *testing.T) {
	tests := []struct {
		name   string
		labels interface{}
		want   string
	}{
		{"nil labels", nil, ""},
		{"empty slice", []string{}, ""},
		{"string slice", []string{"bug", "enhancement"}, "bug, enhancement"},
		{"single string", []string{"bug"}, "bug"},
		{"interface slice", []interface{}{"label1", "label2"}, "[label1 label2]"}, // fmt.Sprintf format
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatLabels(tt.labels)
			if result != tt.want {
				t.Errorf("formatLabels() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestWriteJSONError(t *testing.T) {
	// Test with a type that can't be marshaled to JSON
	var buf bytes.Buffer
	err := writeJSON(&buf, func() {})
	if err == nil {
		t.Error("expected error for unmarshalable type")
	}
}

func TestExecuteTemplateError(t *testing.T) {
	template := "{{.MissingField}}"
	printer, err := NewTemplatePrinter(template, &Options{})
	if err != nil {
		t.Fatalf("failed to create template printer: %v", err)
	}

	var buf bytes.Buffer
	// Missing field should cause execution error - but Go templates don't error on missing fields
	err = printer.PrintIssues(&buf, map[string]interface{}{"Title": "Test"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWriteWriterError(t *testing.T) {
	template := "{{.Title}}"
	printer, err := NewTemplatePrinter(template, &Options{})
	if err != nil {
		t.Fatalf("failed to create template printer: %v", err)
	}

	// Use a writer that always fails
	errWriter := &errorWriter{}
	err = printer.PrintIssues(errWriter, map[string]interface{}{"Title": "Test"})
	if err == nil {
		t.Error("expected error from failing writer")
	}
}

// errorWriter is a writer that always returns an error
type errorWriter struct{}

func (w *errorWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("write error")
}

func TestFormatRelativeTimeBoundaries(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{"59 seconds", 59 * time.Second, "just now"}, // < 1 minute returns "just now"
		{"61 seconds", 61 * time.Second, "1 minute ago"},
		{"2 minutes", 2 * time.Minute, "2 minutes ago"},
		{"59 minutes", 59 * time.Minute, "59 minutes ago"},
		{"61 minutes", 61 * time.Minute, "1 hour ago"},
		{"2 hours", 2 * time.Hour, "2 hours ago"},
		{"23 hours", 23 * time.Hour, "23 hours ago"},
		{"25 hours", 25 * time.Hour, "1 day ago"},
		{"2 days", 2 * 24 * time.Hour, "2 days ago"},
		{"6 days", 6 * 24 * time.Hour, "6 days ago"},
		{"8 days", 8 * 24 * time.Hour, "1 week ago"},
		{"2 weeks", 2 * 7 * 24 * time.Hour, "2 weeks ago"},
		{"13 weeks", 13 * 7 * 24 * time.Hour, "3 months ago"},
		{"2 months", 2 * 30 * 24 * time.Hour, "2 months ago"},
		{"11 months", 11 * 30 * 24 * time.Hour, "11 months ago"},
		{"13 months", 13 * 30 * 24 * time.Hour, "1 year ago"},
		{"2 years", 2 * 365 * 24 * time.Hour, "2 years ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			past := now.Add(-tt.duration)
			result := formatRelativeTime(past)
			if result != tt.want {
				t.Errorf("formatRelativeTime() = %v, want %v", result, tt.want)
			}
		})
	}
}