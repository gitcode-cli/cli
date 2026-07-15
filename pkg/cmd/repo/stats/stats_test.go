package stats

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
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
		_, _ = w.Write([]byte(`{"total":2,"statistics":[{"user_name":"alice","add_lines":10,"delete_lines":2,"total":12}]}`))
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

func TestStatsRunJSONOutput(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v5/owner/repo/repository/commit_statistics" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"total":2,"statistics":[{"user_name":"alice","add_lines":10,"delete_lines":2,"total":12}]}`))
	}))

	err := statsRun(&StatsOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		Branch:     "main",
		JSON:       true,
	})
	if err != nil {
		t.Fatalf("statsRun() error = %v", err)
	}

	var stats map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &stats); err != nil {
		t.Fatalf("output is not valid JSON: %v; output=%q", err, out.String())
	}
	if stats["total"] != float64(2) {
		t.Fatalf("unexpected JSON output: %#v", stats)
	}
}

func TestStatsRunUsesBaseRepoWhenRepoOmitted(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, _, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "auto/detected") {
			t.Fatalf("expected auto/detected in path, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"total":0,"statistics":[]}`))
	}))

	err := statsRun(&StatsOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "",
		BaseRepo:   func() (string, error) { return "auto/detected", nil },
		Branch:     "main",
	})
	if err != nil {
		t.Fatalf("statsRun() error = %v", err)
	}
}

func TestStatsRunBaseRepoError(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, _, _ := testutil.NewTestIOStreams()
	err := statsRun(&StatsOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return &http.Client{}, nil },
		Repository: "",
		BaseRepo:   func() (string, error) { return "", cmdutil.NewUsageError("not in a git repository") },
		Branch:     "main",
	})
	if err == nil {
		t.Fatal("statsRun() error = nil, want usage error")
	}
	if got := cmdutil.ExitCode(err); got != cmdutil.ExitUsage {
		t.Fatalf("ExitCode() = %d, want %d (ExitUsage)", got, cmdutil.ExitUsage)
	}
}
