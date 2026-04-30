package merge

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdMerge(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "merge PR",
			args:    []string{"123"},
			wantErr: false,
		},
		{
			name:    "merge with squash",
			args:    []string{"123", "--method", "squash"},
			wantErr: false,
		},
		{
			name:    "merge with rebase",
			args:    []string{"123", "--method", "rebase"},
			wantErr: false,
		},
		{
			name:    "merge with yes",
			args:    []string{"123", "--yes"},
			wantErr: false,
		},
		{
			name:    "merge with json output",
			args:    []string{"123", "--json", "--yes"},
			wantErr: false,
		},
		{
			name:    "no PR number",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdMerge(f, func(opts *MergeOptions) error {
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

func TestMergeRunJSONWritesMergeResult(t *testing.T) {
	t.Setenv("GC_TOKEN", "token")

	var requests []string
	httpClient := &http.Client{Transport: mergeRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		requests = append(requests, req.Method+" "+req.URL.EscapedPath())
		switch len(requests) {
		case 1:
			return mergeResponse(http.StatusOK, `{"number":123,"title":"Fix","head":{"ref":"feature/test","repo":{"full_name":"source-owner/source-repo"}}}`), nil
		case 2:
			return mergeResponse(http.StatusOK, `{"number":123,"state":"merged","title":"Fix","html_url":"https://gitcode.com/owner/repo/merge_requests/123"}`), nil
		case 3:
			return mergeResponse(http.StatusNoContent, ``), nil
		default:
			t.Fatalf("unexpected request %d: %s %s", len(requests), req.Method, req.URL.Path)
			return nil, nil
		}
	})}

	f := cmdutil.TestFactory()
	err := mergeRun(&MergeOptions{
		IO:           f.IOStreams,
		HttpClient:   func() (*http.Client, error) { return httpClient, nil },
		Repository:   "owner/repo",
		Number:       123,
		MergeMethod:  "merge",
		DeleteBranch: true,
		Yes:          true,
		JSON:         true,
	})
	if err != nil {
		t.Fatalf("mergeRun() error = %v", err)
	}

	var got map[string]interface{}
	out := f.IOStreams.Out.(*bytes.Buffer).Bytes()
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("JSON output did not parse: %v\n%s", err, string(out))
	}
	if got["deleted_branch"] != "feature/test" {
		t.Fatalf("JSON output = %#v", got)
	}
	if got["number"] != float64(123) || got["merged"] != true {
		t.Fatalf("JSON output = %#v", got)
	}
	if _, ok := got["pull_request"].(map[string]interface{}); !ok {
		t.Fatalf("JSON output missing pull_request: %#v", got)
	}
	if strings.Contains(string(out), "Merged PR") || strings.Contains(string(out), "Deleted branch") {
		t.Fatalf("JSON output contains text banner: %q", string(out))
	}
}

func TestMergeRunDeletesHeadBranch(t *testing.T) {
	t.Setenv("GC_TOKEN", "token")

	var requests []string
	httpClient := &http.Client{Transport: mergeRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		requests = append(requests, req.Method+" "+req.URL.EscapedPath())
		switch len(requests) {
		case 1:
			return mergeResponse(http.StatusOK, `{"number":123,"title":"Fix","head":{"ref":"feature/test","repo":{"full_name":"source-owner/source-repo"}}}`), nil
		case 2:
			return mergeResponse(http.StatusOK, `{"number":123,"state":"closed","merged":true}`), nil
		case 3:
			return mergeResponse(http.StatusNoContent, ``), nil
		default:
			t.Fatalf("unexpected request %d: %s %s", len(requests), req.Method, req.URL.Path)
			return nil, nil
		}
	})}

	f := cmdutil.TestFactory()
	err := mergeRun(&MergeOptions{
		IO:           f.IOStreams,
		HttpClient:   func() (*http.Client, error) { return httpClient, nil },
		Repository:   "owner/repo",
		Number:       123,
		MergeMethod:  "merge",
		DeleteBranch: true,
		Yes:          true,
	})
	if err != nil {
		t.Fatalf("mergeRun() error = %v", err)
	}

	want := []string{
		"GET /api/v5/repos/owner/repo/pulls/123",
		"PUT /api/v5/repos/owner/repo/pulls/123/merge",
		"DELETE /api/v5/repos/source-owner/source-repo/branches/feature%2Ftest",
	}
	if strings.Join(requests, "\n") != strings.Join(want, "\n") {
		t.Fatalf("requests = %#v, want %#v", requests, want)
	}
}

func TestMergeRunDeleteBranchFailsBeforeMergeWithoutHeadRepo(t *testing.T) {
	t.Setenv("GC_TOKEN", "token")

	var requests []string
	httpClient := &http.Client{Transport: mergeRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		requests = append(requests, req.Method+" "+req.URL.EscapedPath())
		if len(requests) == 1 {
			return mergeResponse(http.StatusOK, `{"number":123,"title":"Fix","head":{"ref":"feature/test"}}`), nil
		}
		t.Fatalf("unexpected request after missing head repo: %s %s", req.Method, req.URL.Path)
		return nil, nil
	})}

	f := cmdutil.TestFactory()
	err := mergeRun(&MergeOptions{
		IO:           f.IOStreams,
		HttpClient:   func() (*http.Client, error) { return httpClient, nil },
		Repository:   "owner/repo",
		Number:       123,
		MergeMethod:  "merge",
		DeleteBranch: true,
		Yes:          true,
	})
	if err == nil || !strings.Contains(err.Error(), "PR head repository is missing") {
		t.Fatalf("mergeRun() error = %v, want missing head repository", err)
	}
	if len(requests) != 1 {
		t.Fatalf("requests = %#v, want only initial GET", requests)
	}
}

type mergeRoundTripFunc func(*http.Request) (*http.Response, error)

func (fn mergeRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func mergeResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
