package list

import (
	"net/http"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestListRun(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v5/repos/owner/repo/comments" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":1,"body":"looks good","user":{"login":"tester"}}]`))
	}))

	err := listRun(&ListOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		Page:       1,
		PerPage:    20,
	})
	if err != nil {
		t.Fatalf("listRun() error = %v", err)
	}
	if !strings.Contains(out.String(), "tester") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}
