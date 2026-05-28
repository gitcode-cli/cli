package log

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestRunBuildsFileBranchQuery(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	var gotPath string
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		if r.URL.RawQuery != "" {
			gotPath += "?" + r.URL.RawQuery
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"sha":"abcdef123456","commit":{"message":"fix target file","author":{"date":"2026-05-27T10:00:00+08:00"}}}]`))
	}))

	err := run(&Options{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		BaseRepo:   func() (string, error) { return "owner/repo", nil },
		File:       "src/main.go",
		Branch:     "main",
		Limit:      20,
		Page:       2,
	})
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}
	for _, want := range []string{
		"/api/v5/repos/owner/repo/commits?",
		"path=src%2Fmain.go",
		"sha=main",
		"page=2",
		"per_page=20",
	} {
		if !strings.Contains(gotPath, want) {
			t.Fatalf("path = %q, missing %q", gotPath, want)
		}
	}
	if !strings.Contains(out.String(), "abcdef12") || !strings.Contains(out.String(), "fix target file") {
		t.Fatalf("output = %q", out.String())
	}
}

func TestRunJSON(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"sha":"abcdef123456","commit":{"message":"fix target file"}}]`))
	}))

	err := run(&Options{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		Limit:      30,
		Page:       1,
		JSON:       true,
	})
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}
	var commits []map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &commits); err != nil {
		t.Fatalf("output is not JSON: %v; output=%q", err, out.String())
	}
	if len(commits) != 1 || commits[0]["sha"] != "abcdef123456" {
		t.Fatalf("unexpected JSON: %#v", commits)
	}
}
