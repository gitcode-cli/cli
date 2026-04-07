package diff

import (
	"net/http"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestDiffRun(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v5/repos/owner/repo/commit/abc123/diff" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte("diff --git a/file b/file"))
	}))

	err := diffRun(&DiffOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		SHA:        "abc123",
	})
	if err != nil {
		t.Fatalf("diffRun() error = %v", err)
	}
	if !strings.Contains(out.String(), "diff --git") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}
