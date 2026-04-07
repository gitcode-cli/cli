package patch

import (
	"net/http"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestPatchRun(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v5/repos/owner/repo/commit/abc123/patch" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte("@@ -1 +1 @@"))
	}))

	err := patchRun(&PatchOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		SHA:        "abc123",
	})
	if err != nil {
		t.Fatalf("patchRun() error = %v", err)
	}
	if !strings.Contains(out.String(), "@@ -1 +1 @@") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}
