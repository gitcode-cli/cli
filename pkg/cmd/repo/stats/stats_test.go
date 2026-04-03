package stats

import (
	"net/http"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestStatsRun(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v5/owner/repo/repository/commit_statistics" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"total":2,"statistics":[{"author":"alice","additions":10,"deletions":2,"total":12}]}`))
	}))

	err := statsRun(&StatsOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		Branch:     "main",
	})
	if err != nil {
		t.Fatalf("statsRun() error = %v", err)
	}

	body := out.String()
	if !strings.Contains(body, "Total commits: 2") {
		t.Fatalf("output missing total commits: %q", body)
	}
	if !strings.Contains(body, "alice") {
		t.Fatalf("output missing author stats: %q", body)
	}
}
