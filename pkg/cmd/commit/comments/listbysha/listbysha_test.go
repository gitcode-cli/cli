package listbysha

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestListBySHARun(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v5/repos/owner/repo/commits/abc123/comments" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":1,"body":"looks good","user":{"login":"tester"}}]`))
	}))

	err := listBySHARun(&ListBySHAOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		SHA:        "abc123",
		Page:       1,
		PerPage:    20,
	})
	if err != nil {
		t.Fatalf("listBySHARun() error = %v", err)
	}
	if !strings.Contains(out.String(), "tester") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestListBySHARunJSONOutput(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v5/repos/owner/repo/commits/abc123/comments" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":1,"body":"looks good","user":{"login":"tester"}}]`))
	}))

	err := listBySHARun(&ListBySHAOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		SHA:        "abc123",
		Page:       1,
		PerPage:    20,
		JSON:       true,
	})
	if err != nil {
		t.Fatalf("listBySHARun() error = %v", err)
	}

	var comments []map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &comments); err != nil {
		t.Fatalf("output is not valid JSON: %v; output=%q", err, out.String())
	}
	if len(comments) != 1 || comments[0]["body"] != "looks good" {
		t.Fatalf("unexpected JSON output: %#v", comments)
	}
}

func TestListBySHARunJSONEmptyOutput(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))

	err := listBySHARun(&ListBySHAOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		SHA:        "abc123",
		Page:       1,
		PerPage:    20,
		JSON:       true,
	})
	if err != nil {
		t.Fatalf("listBySHARun() error = %v", err)
	}
	if got := out.String(); got != "[]\n" {
		t.Fatalf("output = %q, want JSON empty array", got)
	}
}
