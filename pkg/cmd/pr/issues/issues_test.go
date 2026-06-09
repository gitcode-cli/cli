package issues

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestNewCmdIssues(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "list PR issues",
			args:    []string{"123"},
			wantErr: false,
		},
		{
			name:    "list PR issues with repo",
			args:    []string{"123", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "no PR number",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "invalid PR number",
			args:    []string{"abc"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdIssues(f, func(opts *IssuesOptions) error {
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

func TestIssuesRunJSON(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/issues") {
			_, _ = w.Write([]byte(`[{"id":"1","number":"42","title":"Fix login bug","state":"open","html_url":"https://gitcode.com/owner/repo/issues/42"}]`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))

	err := issuesRun(&IssuesOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		BaseRepo:   func() (string, error) { return "owner/repo", nil },
		Number:     1,
		JSON:       true,
	})
	if err != nil {
		t.Fatalf("issuesRun() error = %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("output is not JSON: %v; output=%q", err, out.String())
	}
	if len(result) != 1 || result[0]["number"] != "42" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestIssuesRunEmpty(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/issues") {
			_, _ = w.Write([]byte(`[]`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))

	err := issuesRun(&IssuesOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		BaseRepo:   func() (string, error) { return "owner/repo", nil },
		Number:     1,
		JSON:       false,
	})
	if err != nil {
		t.Fatalf("issuesRun() error = %v", err)
	}
	if !strings.Contains(out.String(), "No linked issues") {
		t.Fatalf("expected empty message, got: %q", out.String())
	}
}

func TestIssuesRunText(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/issues") {
			_, _ = w.Write([]byte(`[{"id":"1","number":"42","title":"Fix login bug","state":"open","html_url":"https://gitcode.com/owner/repo/issues/42"}]`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))

	err := issuesRun(&IssuesOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		BaseRepo:   func() (string, error) { return "owner/repo", nil },
		Number:     1,
		JSON:       false,
	})
	if err != nil {
		t.Fatalf("issuesRun() error = %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "#42") || !strings.Contains(output, "Fix login bug") {
		t.Fatalf("expected issue info in output, got: %q", output)
	}
}
