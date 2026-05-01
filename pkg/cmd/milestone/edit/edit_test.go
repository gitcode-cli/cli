package edit

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdEdit(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "edit title",
			args:    []string{"5", "--title", "New title", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "edit description",
			args:    []string{"5", "--description", "New description", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "edit state open",
			args:    []string{"5", "--state", "open", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "edit state closed",
			args:    []string{"5", "--state", "closed", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "edit due date",
			args:    []string{"5", "--due-date", "2024-12-31", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "edit multiple fields",
			args:    []string{"5", "--title", "Title", "--description", "Desc", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "edit with json output",
			args:    []string{"5", "--title", "Title", "--json", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "missing milestone number",
			args:    []string{"-R", "owner/repo"},
			wantErr: true,
		},
		{
			name:    "invalid milestone number",
			args:    []string{"abc", "-R", "owner/repo"},
			wantErr: true,
		},
		{
			name:    "no edit options",
			args:    []string{"5", "-R", "owner/repo"},
			wantErr: false, // Command parses successfully, error happens in runF
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdEdit(f, func(opts *EditOptions) error {
				// Mock run function - just validate
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

func TestEditRunJSONWritesUpdatedMilestone(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	opts := &EditOptions{
		IO: f.IOStreams,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					if req.URL.Path != "/api/v5/repos/owner/repo/milestones/5" {
						t.Fatalf("unexpected request: %s", req.URL.Path)
					}
					return milestoneResponse(http.StatusOK, `{"number":5,"title":"updated","state":"open"}`), nil
				}),
			}, nil
		},
		Repository: "owner/repo",
		Number:     5,
		Title:      "updated",
		JSON:       true,
	}

	if err := editRun(opts); err != nil {
		t.Fatalf("editRun() error = %v", err)
	}

	var got map[string]interface{}
	out := f.IOStreams.Out.(*bytes.Buffer).Bytes()
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("JSON output did not parse: %v\n%s", err, string(out))
	}
	if got["number"] != float64(5) {
		t.Fatalf("JSON output = %#v", got)
	}
	if strings.Contains(string(out), "Updated milestone") {
		t.Fatalf("JSON output contains text banner: %q", string(out))
	}
}

func TestEditRun_InvalidStateValue(t *testing.T) {
	f := cmdutil.TestFactory()
	opts := &EditOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Number:     5,
		Repository: "owner/repo",
		State:      "pending",
	}

	err := editRun(opts)
	if err == nil {
		t.Error("Expected error for invalid state value")
	}
	if !strings.Contains(err.Error(), "invalid state value") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestEditRun_InvalidDueDateFormat(t *testing.T) {
	f := cmdutil.TestFactory()
	opts := &EditOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Number:     5,
		Repository: "owner/repo",
		DueDate:    "not-a-date",
	}

	err := editRun(opts)
	if err == nil {
		t.Error("Expected error for invalid due date format")
	}
	if !strings.Contains(err.Error(), "invalid due date format") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestEditRun_NoEditOptions(t *testing.T) {
	f := cmdutil.TestFactory()
	opts := &EditOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Number:     5,
		Repository: "owner/repo",
	}

	err := editRun(opts)
	if err == nil {
		t.Error("Expected error when no edit options provided")
	}
	if !strings.Contains(err.Error(), "at least one edit option is required") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestParseRepo(t *testing.T) {
	tests := []struct {
		repo      string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		{"owner/repo", "owner", "repo", false},
		{"gitcode-cli/cli", "gitcode-cli", "cli", false},
		{"", "", "", true},
		{"invalid", "", "", true},
		{"too/many/parts", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.repo, func(t *testing.T) {
			owner, repo, err := parseRepo(tt.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if owner != tt.wantOwner {
				t.Errorf("parseRepo() owner = %v, want %v", owner, tt.wantOwner)
			}
			if repo != tt.wantRepo {
				t.Errorf("parseRepo() repo = %v, want %v", repo, tt.wantRepo)
			}
		})
	}
}

func TestEditRun_AuthError(t *testing.T) {
	// Clear any existing token
	os.Unsetenv("GC_TOKEN")
	os.Unsetenv("GITCODE_TOKEN")

	f := cmdutil.TestFactory()
	// Override EnvToken to return empty
	opts := &EditOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Number:     5,
		Repository: "owner/repo",
		Title:      "test",
	}

	// No token set - should fail before API call
	err := editRun(opts)
	if err == nil {
		t.Error("Expected auth error when no token provided")
	}
	if !strings.Contains(err.Error(), "not authenticated") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func milestoneResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
