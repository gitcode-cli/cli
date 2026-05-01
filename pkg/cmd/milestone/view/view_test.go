package view

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestNewCmdView(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "view milestone",
			args:    []string{"1"},
			wantErr: false,
		},
		{
			name:    "view with web flag",
			args:    []string{"1", "--web"},
			wantErr: false,
		},
		{
			name:    "view with json flag",
			args:    []string{"1", "--json"},
			wantErr: false,
		},
		{
			name:    "view with issues flag",
			args:    []string{"1", "--issues"},
			wantErr: false,
		},
		{
			name:    "view without issues",
			args:    []string{"1", "--issues=false"},
			wantErr: false,
		},
		{
			name:    "no milestone number",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "invalid milestone number",
			args:    []string{"abc"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdView(f, func(opts *ViewOptions) error {
				return nil
			})
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewCmdViewIssuesFlagParsing(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantIssues bool
	}{
		{
			name:       "default shows issues",
			args:       []string{"1"},
			wantIssues: true,
		},
		{
			name:       "explicit --issues=true",
			args:       []string{"1", "--issues=true"},
			wantIssues: true,
		},
		{
			name:       "explicit --issues=false",
			args:       []string{"1", "--issues=false"},
			wantIssues: false,
		},
		{
			name:       "--issues flag",
			args:       []string{"1", "--issues"},
			wantIssues: true,
		},
		{
			name:       "--no-issues style",
			args:       []string{"1", "--issues=false"},
			wantIssues: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			var capturedOpts *ViewOptions
			cmd := NewCmdView(f, func(opts *ViewOptions) error {
				capturedOpts = opts
				return nil
			})
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			if capturedOpts.Issues != tt.wantIssues {
				t.Errorf("Issues flag = %v, want %v", capturedOpts.Issues, tt.wantIssues)
			}
		})
	}
}

func TestViewRunJSONOutput(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v5/repos/owner/repo/milestones/1" {
			_, _ = w.Write([]byte(`{"id":1,"number":1,"title":"v1","state":"open","description":"release"}`))
			return
		}
		if r.URL.Path == "/api/v5/repos/owner/repo/issues" {
			_, _ = w.Write([]byte(`[]`))
			return
		}
		t.Fatalf("unexpected path: %s", r.URL.Path)
	}))

	err := viewRun(&ViewOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		BaseRepo:   func() (string, error) { return "owner/repo", nil },
		Repository: "owner/repo",
		Number:     1,
		JSON:       true,
		Issues:     true,
	})
	if err != nil {
		t.Fatalf("viewRun() error = %v", err)
	}

	var milestone map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &milestone); err != nil {
		t.Fatalf("output is not valid JSON: %v; output=%q", err, out.String())
	}
	if milestone["title"] != "v1" {
		t.Fatalf("unexpected JSON output: %#v", milestone)
	}
	if milestone["total_issues"] != float64(0) {
		t.Fatalf("unexpected total_issues: %v", milestone["total_issues"])
	}
}

func TestViewRunJSONOutputWithIssues(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v5/repos/owner/repo/milestones/1" {
			_, _ = w.Write([]byte(`{"id":1,"number":1,"title":"v1","state":"open","description":"release"}`))
			return
		}
		if r.URL.Path == "/api/v5/repos/owner/repo/issues" {
			_, _ = w.Write([]byte(`[
				{"id":"1","number":"10","title":"Bug fix","state":"closed"},
				{"id":"2","number":"11","title":"New feature","state":"open"},
				{"id":"3","number":"12","title":"Another bug","state":"open"}
			]`))
			return
		}
		t.Fatalf("unexpected path: %s", r.URL.Path)
	}))

	err := viewRun(&ViewOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		BaseRepo:   func() (string, error) { return "owner/repo", nil },
		Repository: "owner/repo",
		Number:     1,
		JSON:       true,
		Issues:     true,
	})
	if err != nil {
		t.Fatalf("viewRun() error = %v", err)
	}

	var result MilestoneWithIssues
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v; output=%q", err, out.String())
	}
	if result.Title != "v1" {
		t.Fatalf("unexpected title: %v", result.Title)
	}
	if result.TotalIssues != 3 {
		t.Fatalf("unexpected total_issues: %v", result.TotalIssues)
	}
	if result.ClosedIssues != 1 {
		t.Fatalf("unexpected closed_issues: %v, want 1", result.ClosedIssues)
	}
	if result.OpenIssues != 2 {
		t.Fatalf("unexpected open_issues: %v, want 2", result.OpenIssues)
	}
	if len(result.Issues) != 3 {
		t.Fatalf("unexpected issues length: %v", len(result.Issues))
	}
}

func TestViewRunTextOutputWithIssues(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v5/repos/owner/repo/milestones/1" {
			_, _ = w.Write([]byte(`{"id":1,"number":1,"title":"v1.0","state":"open","description":"First release"}`))
			return
		}
		if r.URL.Path == "/api/v5/repos/owner/repo/issues" {
			_, _ = w.Write([]byte(`[
				{"id":"1","number":"10","title":"Bug fix","state":"closed"},
				{"id":"2","number":"11","title":"New feature","state":"open"},
				{"id":"3","number":"12","title":"Another bug","state":"open"}
			]`))
			return
		}
		t.Fatalf("unexpected path: %s", r.URL.Path)
	}))

	err := viewRun(&ViewOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		BaseRepo:   func() (string, error) { return "owner/repo", nil },
		Repository: "owner/repo",
		Number:     1,
		Issues:     true,
	})
	if err != nil {
		t.Fatalf("viewRun() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Issues (3 total, 1 closed, 2 open)") {
		t.Errorf("output missing issue counts: %s", output)
	}
	if !strings.Contains(output, "Closed:") {
		t.Errorf("output missing closed section: %s", output)
	}
	if !strings.Contains(output, "#10 Bug fix") {
		t.Errorf("output missing closed issue: %s", output)
	}
	if !strings.Contains(output, "Open:") {
		t.Errorf("output missing open section: %s", output)
	}
	if !strings.Contains(output, "#11 New feature") {
		t.Errorf("output missing open issue: %s", output)
	}
}

func TestViewRunWithoutIssues(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Only milestone request should be made
		if r.URL.Path == "/api/v5/repos/owner/repo/milestones/1" {
			_, _ = w.Write([]byte(`{"id":1,"number":1,"title":"v1.0","state":"open"}`))
			return
		}
		// Issues request should NOT be made when Issues=false
		t.Fatalf("unexpected issues request when Issues=false: %s", r.URL.Path)
	}))

	err := viewRun(&ViewOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		BaseRepo:   func() (string, error) { return "owner/repo", nil },
		Repository: "owner/repo",
		Number:     1,
		Issues:     false,
	})
	if err != nil {
		t.Fatalf("viewRun() error = %v", err)
	}

	output := out.String()
	if strings.Contains(output, "Issues") {
		t.Errorf("output should not contain issues section: %s", output)
	}
}

func TestViewRunEmptyMilestone(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v5/repos/owner/repo/milestones/1" {
			_, _ = w.Write([]byte(`{"id":1,"number":1,"title":"v1.0","state":"open"}`))
			return
		}
		if r.URL.Path == "/api/v5/repos/owner/repo/issues" {
			_, _ = w.Write([]byte(`[]`)) // Empty array - no issues
			return
		}
		t.Fatalf("unexpected path: %s", r.URL.Path)
	}))

	err := viewRun(&ViewOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		BaseRepo:   func() (string, error) { return "owner/repo", nil },
		Repository: "owner/repo",
		Number:     1,
		Issues:     true,
	})
	if err != nil {
		t.Fatalf("viewRun() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Issues: None") {
		t.Errorf("output should show 'Issues: None': %s", output)
	}
}

func TestViewRunRejectsJSONWithWebBeforeAuth(t *testing.T) {
	io, _, _, _ := testutil.NewTestIOStreams()
	err := viewRun(&ViewOptions{
		IO:     io,
		Web:    true,
		JSON:   true,
		Number: 1,
		HttpClient: func() (*http.Client, error) {
			t.Fatal("HttpClient should not be called when --json and --web conflict")
			return nil, nil
		},
	})
	if err == nil || !strings.Contains(err.Error(), "cannot use --json with --web") {
		t.Fatalf("viewRun() error = %v, want json/web usage error", err)
	}
}

func TestCountIssuesByState(t *testing.T) {
	issues := []api.Issue{
		{Number: "1", Title: "A", State: "open"},
		{Number: "2", Title: "B", State: "closed"},
		{Number: "3", Title: "C", State: "open"},
		{Number: "4", Title: "D", State: "CLOSED"}, // test case insensitivity
	}

	openCount := countIssuesByState(issues, "open")
	if openCount != 2 {
		t.Errorf("open count = %d, want 2", openCount)
	}

	closedCount := countIssuesByState(issues, "closed")
	if closedCount != 2 {
		t.Errorf("closed count = %d, want 2", closedCount)
	}
}
