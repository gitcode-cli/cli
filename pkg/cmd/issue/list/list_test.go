package list

import (
	"io"
	"net/http"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

func TestNewCmdList(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "list default",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "list with state",
			args:    []string{"--state", "closed"},
			wantErr: false,
		},
		{
			name:    "list with limit",
			args:    []string{"--limit", "10"},
			wantErr: false,
		},
		{
			name:    "list with labels",
			args:    []string{"--label", "bug"},
			wantErr: false,
		},
		{
			name:    "list with milestone",
			args:    []string{"--milestone", "v1.0"},
			wantErr: false,
		},
		{
			name:    "list with assignee",
			args:    []string{"--assignee", "username"},
			wantErr: false,
		},
		{
			name:    "list with creator",
			args:    []string{"--creator", "username"},
			wantErr: false,
		},
		{
			name:    "list with sort",
			args:    []string{"--sort", "updated"},
			wantErr: false,
		},
		{
			name:    "list with direction",
			args:    []string{"--direction", "asc"},
			wantErr: false,
		},
		{
			name:    "list with search",
			args:    []string{"--search", "bug"},
			wantErr: false,
		},
		{
			name:    "list with created-after",
			args:    []string{"--created-after", "2024-01-01"},
			wantErr: false,
		},
		{
			name:    "list with updated-after",
			args:    []string{"--updated-after", "2024-01-01"},
			wantErr: false,
		},
		{
			name:    "list with combined filters",
			args:    []string{"--state", "open", "--sort", "updated", "--direction", "desc"},
			wantErr: false,
		},
		{
			name:    "list with json compatibility",
			args:    []string{"--json"},
			wantErr: false,
		},
		{
			name:    "list with format flag",
			args:    []string{"--format", "table"},
			wantErr: false,
		},
		{
			name:    "list with time format flag",
			args:    []string{"--time-format", "relative"},
			wantErr: false,
		},
		{
			name:    "list with template flag",
			args:    []string{"--template", "{{.Title}}"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdList(f, func(opts *ListOptions) error {
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

func TestResolveOutputOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ListOptions
		wantErr bool
	}{
		{
			name: "json compatibility",
			opts: &ListOptions{JSON: true},
		},
		{
			name: "format json",
			opts: &ListOptions{Format: "json"},
		},
		{
			name: "time format relative",
			opts: &ListOptions{TimeFormat: "relative"},
		},
		{
			name:    "invalid format",
			opts:    &ListOptions{Format: "yaml"},
			wantErr: true,
		},
		{
			name:    "invalid time format",
			opts:    &ListOptions{TimeFormat: "iso"},
			wantErr: true,
		},
		{
			name:    "json with incompatible format",
			opts:    &ListOptions{JSON: true, Format: "table"},
			wantErr: true,
		},
		{
			name:    "json with template",
			opts:    &ListOptions{JSON: true, Template: "{{.Title}}"},
			wantErr: true,
		},
		{
			name:    "template with format",
			opts:    &ListOptions{Format: "simple", Template: "{{.Title}}"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := resolveOutputOptions(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Fatalf("resolveOutputOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNormalizeIssueListTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		endOfDay bool
		want     string
		wantErr  bool
	}{
		{
			name:     "date only start of day",
			input:    "2026-03-31",
			endOfDay: false,
			want:     "2026-03-31T00:00:00Z",
		},
		{
			name:     "date only end of day",
			input:    "2026-03-31",
			endOfDay: true,
			want:     "2026-03-31T23:59:59Z",
		},
		{
			name:     "rfc3339",
			input:    "2026-03-31T12:30:00+08:00",
			endOfDay: true,
			want:     "2026-03-31T12:30:00+08:00",
		},
		{
			name:     "invalid",
			input:    "2026/03/31",
			endOfDay: false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeIssueListTime(tt.input, tt.endOfDay)
			if tt.wantErr {
				if err == nil {
					t.Fatal("normalizeIssueListTime() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("normalizeIssueListTime() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("normalizeIssueListTime() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestListRunRejectsOutputUsageErrorsBeforeHTTP(t *testing.T) {
	httpCalled := false
	opts := &ListOptions{
		IO: cmdutil.TestFactory().IOStreams,
		HttpClient: func() (*http.Client, error) {
			httpCalled = true
			return &http.Client{}, nil
		},
		Format: "yaml",
	}

	err := listRun(opts)
	if err == nil {
		t.Fatal("listRun() error = nil, want usage error")
	}
	if httpCalled {
		t.Fatal("listRun() called HttpClient before validating output flags")
	}
}

func TestListRunAllowsTemplateOutputForEmptyResults(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	ioStreams, _, stdout, _ := iostreams.Test()
	opts := &ListOptions{
		IO:         ioStreams,
		HttpClient: func() (*http.Client, error) { return &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     http.StatusText(http.StatusOK),
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`[]`)),
			}, nil
		})}, nil },
		Repository: "owner/repo",
		Template:   "{{len .}} issues",
	}

	if err := listRun(opts); err != nil {
		t.Fatalf("listRun() error = %v", err)
	}

	if got := stdout.String(); got != "0 issues" {
		t.Fatalf("stdout = %q, want %q", got, "0 issues")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
