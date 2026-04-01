package relations

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

func TestRelationsRunText(t *testing.T) {
	io, _, out, _ := iostreams.Test()
	opts := &RelationsOptions{
		IO: io,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					switch req.URL.Path {
					case "/api/v5/repos/owner/repo/issues":
						return jsonResponse(`[
							{"number":"7","title":"Issue Seven","state":"open","html_url":"https://example.test/issues/7"},
							{"number":"8","title":"Issue Eight","state":"closed","html_url":"https://example.test/issues/8"}
						]`), nil
					case "/api/v5/repos/owner/repo/issues/7/pull_requests":
						return jsonResponse(`[
							{"number":2,"title":"PR Two","state":"open","html_url":"https://example.test/pr/2"}
						]`), nil
					case "/api/v5/repos/owner/repo/issues/8/pull_requests":
						return jsonResponse(`[
							{"number":2,"title":"PR Two","state":"open","html_url":"https://example.test/pr/2"}
						]`), nil
					default:
						t.Fatalf("unexpected path: %s", req.URL.Path)
						return nil, nil
					}
				}),
			}, nil
		},
		BaseRepo: func() (string, error) { return "owner/repo", nil },
		State:    "all",
		Limit:    100,
		Mode:     1,
	}

	t.Setenv("GC_TOKEN", "test-token")

	if err := relationsRun(opts); err != nil {
		t.Fatalf("relationsRun() error = %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "PR #2 [open] PR Two") {
		t.Fatalf("expected PR header in output, got %q", got)
	}
	if !strings.Contains(got, "Issue #7 [open] Issue Seven") || !strings.Contains(got, "Issue #8 [closed] Issue Eight") {
		t.Fatalf("expected issue rows in output, got %q", got)
	}
}

func TestRelationsRunJSON(t *testing.T) {
	io, _, out, _ := iostreams.Test()
	opts := &RelationsOptions{
		IO: io,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					switch req.URL.Path {
					case "/api/v5/repos/owner/repo/issues":
						return jsonResponse(`[{"number":"7","title":"Issue Seven","state":"open","html_url":"https://example.test/issues/7"}]`), nil
					case "/api/v5/repos/owner/repo/issues/7/pull_requests":
						return jsonResponse(`[{"number":2,"title":"PR Two","state":"open","html_url":"https://example.test/pr/2"}]`), nil
					default:
						t.Fatalf("unexpected path: %s", req.URL.Path)
						return nil, nil
					}
				}),
			}, nil
		},
		BaseRepo: func() (string, error) { return "owner/repo", nil },
		State:    "all",
		Limit:    100,
		Mode:     1,
		JSON:     true,
	}

	t.Setenv("GC_TOKEN", "test-token")

	if err := relationsRun(opts); err != nil {
		t.Fatalf("relationsRun() error = %v", err)
	}

	var rows []RelationRow
	if err := json.Unmarshal(out.Bytes(), &rows); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if len(rows) != 1 || rows[0].Issue.Number != "7" || rows[0].PR.Number != 2 {
		t.Fatalf("unexpected rows: %#v", rows)
	}
}

func TestRelationsRunValidation(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := &RelationsOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return &http.Client{}, nil },
		BaseRepo:   func() (string, error) { return "owner/repo", nil },
		State:      "bad",
		Limit:      100,
		Mode:       1,
	}

	t.Setenv("GC_TOKEN", "test-token")

	err := relationsRun(opts)
	if err == nil || !strings.Contains(err.Error(), "invalid state") {
		t.Fatalf("relationsRun() error = %v, want invalid state error", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       ioNopCloser(body),
	}
}

func ioNopCloser(body string) *readCloser {
	return &readCloser{Reader: strings.NewReader(body)}
}

type readCloser struct {
	*strings.Reader
}

func (r *readCloser) Close() error { return nil }
