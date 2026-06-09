package branch

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestNewCmdBranchView(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "view branch",
			args:    []string{"view", "main"},
			wantErr: false,
		},
		{
			name:    "view branch with repo",
			args:    []string{"view", "main", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "no branch name",
			args:    []string{"view"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdBranch(f, func(opts *ViewOptions) error {
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

func TestViewRunJSON(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(r.URL.Path, "/branches/main") {
			_, _ = w.Write([]byte(`{"name":"main","protected":false,"commit":{"id":"abc123def456","short_id":"abc123de","title":"Initial commit","author":{"id":"1","login":"dev"}}}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))

	err := viewRun(&ViewOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		BaseRepo:   func() (string, error) { return "owner/repo", nil },
		Branch:     "main",
		JSON:       true,
	})
	if err != nil {
		t.Fatalf("viewRun() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("output is not JSON: %v; output=%q", err, out.String())
	}
	if result["name"] != "main" {
		t.Fatalf("unexpected name: %#v", result["name"])
	}
	commit, ok := result["commit"].(map[string]interface{})
	if !ok || commit["id"] != "abc123def456" {
		t.Fatalf("unexpected commit: %#v", result["commit"])
	}
}

func TestViewRunText(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(r.URL.Path, "/branches/main") {
			_, _ = w.Write([]byte(`{"name":"main","protected":true,"commit":{"id":"abc123def456","short_id":"abc123de","title":"Initial commit","author":{"id":"1","login":"dev"}}}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))

	err := viewRun(&ViewOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		BaseRepo:   func() (string, error) { return "owner/repo", nil },
		Branch:     "main",
		JSON:       false,
	})
	if err != nil {
		t.Fatalf("viewRun() error = %v", err)
	}
	output := out.String()
	for _, want := range []string{"main", "Protected: yes", "abc123def456", "Initial commit", "dev"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q: %s", want, output)
		}
	}
}

func TestViewRunNotFound(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, _, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"404 Branch Not Found"}`))
	}))

	err := viewRun(&ViewOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		BaseRepo:   func() (string, error) { return "owner/repo", nil },
		Branch:     "nonexistent",
		JSON:       false,
	})
	if err == nil {
		t.Fatal("expected error for nonexistent branch")
	}
}

func TestViewRunNoCommit(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(r.URL.Path, "/branches/orphan") {
			_, _ = w.Write([]byte(`{"name":"orphan","protected":false}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))

	err := viewRun(&ViewOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		BaseRepo:   func() (string, error) { return "owner/repo", nil },
		Branch:     "orphan",
		JSON:       false,
	})
	if err != nil {
		t.Fatalf("viewRun() error = %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "orphan") {
		t.Fatalf("output missing branch name: %s", output)
	}
	if strings.Contains(output, "Commit:") {
		t.Fatalf("output should not contain Commit when nil: %s", output)
	}
}
