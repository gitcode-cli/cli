package edit

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestEditRunWithNotesFile(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	notesPath := filepath.Join(t.TempDir(), "notes.md")
	if err := os.WriteFile(notesPath, []byte("updated notes"), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v5/repos/owner/repo/releases/tags/v1.0.0":
			_, _ = w.Write([]byte(`{"id":1,"tag_name":"v1.0.0","name":"v1.0.0"}`))
		case r.Method == http.MethodPatch && r.URL.Path == "/api/v5/repos/owner/repo/releases/1":
			_, _ = w.Write([]byte(`{"id":1,"tag_name":"v1.0.0","name":"Updated Release","html_url":"https://gitcode.com/owner/repo/releases/v1.0.0"}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))

	err := editRun(&EditOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		TagName:    "v1.0.0",
		NotesFile:  notesPath,
	})
	if err != nil {
		t.Fatalf("editRun() error = %v", err)
	}

	if !strings.Contains(out.String(), "Updated release Updated Release") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}
