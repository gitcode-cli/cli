package list

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/output"
	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestNewCmdList(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{name: "list default", args: []string{}, wantErr: false},
		{name: "list by status", args: []string{"--status", "FAILED"}, wantErr: false},
		{name: "list by event", args: []string{"--event", "Push"}, wantErr: false},
		{name: "list with limit", args: []string{"--limit", "10"}, wantErr: false},
		{name: "list with branch", args: []string{"--branch", "main"}, wantErr: false},
		{name: "list with executor", args: []string{"--executor", "dev"}, wantErr: false},
		{name: "list with workflow", args: []string{"--workflow", "CI"}, wantErr: false},
		{name: "list with workflow-id", args: []string{"--workflow-id", "wf-1"}, wantErr: false},
		{name: "list with pr", args: []string{"--pr", "42"}, wantErr: false},
		{name: "list with paginate", args: []string{"--paginate", "--per-page", "100"}, wantErr: false},
		{name: "list with page", args: []string{"--page", "2"}, wantErr: false},
		{name: "list with table format", args: []string{"--format", "table"}, wantErr: false},
		{name: "list with json", args: []string{"--json"}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdList(f, func(opts *ListOptions) error {
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

func TestNewCmdListStatusEnum(t *testing.T) {
	cmd := NewCmdList(cmdutil.TestFactory(), func(opts *ListOptions) error {
		return nil
	})

	flag := cmd.Flags().Lookup("status")
	if flag == nil {
		t.Fatal("status flag missing")
	}
	got := flag.Annotations[cmdutil.FlagEnumAnnotation]
	for _, want := range []string{"COMPLETED", "RUNNING", "FAILED", "CANCELED", "IGNORED", "PAUSED", "SUSPEND"} {
		if !slices.Contains(got, want) {
			t.Fatalf("status enum = %v, missing %q", got, want)
		}
	}
}

func TestNewCmdListEventEnum(t *testing.T) {
	cmd := NewCmdList(cmdutil.TestFactory(), func(opts *ListOptions) error {
		return nil
	})

	flag := cmd.Flags().Lookup("event")
	if flag == nil {
		t.Fatal("event flag missing")
	}
	got := flag.Annotations[cmdutil.FlagEnumAnnotation]
	for _, want := range []string{"MR", "Push", "Manual"} {
		if !slices.Contains(got, want) {
			t.Fatalf("event enum = %v, missing %q", got, want)
		}
	}
}

func TestListRunBuildsV8Query(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	var gotPath string
	opts := &ListOptions{
		IO: f.IOStreams,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: testutil.NewRoundTripFunc(func(req *http.Request) (*http.Response, error) {
					gotPath = req.URL.Path
					if req.URL.RawQuery != "" {
						gotPath += "?" + req.URL.RawQuery
					}
					return listTestResponse(http.StatusOK, `{"total_count":0,"workflow_runs":[]}`), nil
				}),
			}, nil
		},
		Repository:    "owner/repo",
		Status:        "FAILED",
		Event:         "Push",
		Branch:        "main",
		Executor:      "dev",
		WorkflowID:    "wf-1",
		WorkflowName:  "ci",
		PullRequestID: "42",
		Limit:         25,
		Page:          2,
	}

	if err := listRun(opts); err != nil {
		t.Fatalf("listRun() error = %v", err)
	}

	if !strings.HasPrefix(gotPath, "/api/v8/repos/owner/repo/actions/runs?") {
		t.Fatalf("request path = %q, want v8 prefix", gotPath)
	}
	assertNoAccessTokenQuery(t, gotPath)

	parsed, err := url.Parse(gotPath)
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}
	q := parsed.Query()
	for _, key := range []string{"status", "event", "branch", "executor", "workflow_id", "workflow_name", "pull_request_id", "per_page", "page"} {
		if _, ok := q[key]; !ok {
			t.Fatalf("query missing %s in %s", key, q.Encode())
		}
	}
	if q.Get("status") != "FAILED" {
		t.Fatalf("status = %q, want FAILED", q.Get("status"))
	}
	if q.Get("per_page") != "25" {
		t.Fatalf("per_page = %q, want 25", q.Get("per_page"))
	}
	if q.Get("page") != "2" {
		t.Fatalf("page = %q, want 2", q.Get("page"))
	}
}

func TestListRunDefaultPerPage(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	var gotPath string
	opts := &ListOptions{
		IO: f.IOStreams,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: testutil.NewRoundTripFunc(func(req *http.Request) (*http.Response, error) {
					gotPath = req.URL.Path
					if req.URL.RawQuery != "" {
						gotPath += "?" + req.URL.RawQuery
					}
					return listTestResponse(http.StatusOK, `{"total_count":0,"workflow_runs":[]}`), nil
				}),
			}, nil
		},
		Repository: "owner/repo",
		Limit:      30,
	}

	if err := listRun(opts); err != nil {
		t.Fatalf("listRun() error = %v", err)
	}

	if gotPath != "/api/v8/repos/owner/repo/actions/runs?per_page=30" {
		t.Fatalf("request path = %q, want default per_page=30", gotPath)
	}
}

func TestListRunPaginatesUntilLimit(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := iostreams.Test()
	var gotPaths []string
	opts := &ListOptions{
		IO: io,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: testutil.NewRoundTripFunc(func(req *http.Request) (*http.Response, error) {
					gotPath := req.URL.Path
					if req.URL.RawQuery != "" {
						gotPath += "?" + req.URL.RawQuery
					}
					gotPaths = append(gotPaths, gotPath)
					switch req.URL.Query().Get("page") {
					case "1":
						return listTestResponse(http.StatusOK, `{"total_count":4,"workflow_runs":[{"workflow_run_id":"r1","run_number":1,"status":"FAILED"},{"workflow_run_id":"r2","run_number":2,"status":"COMPLETED"}]}`), nil
					case "2":
						return listTestResponse(http.StatusOK, `{"total_count":4,"workflow_runs":[{"workflow_run_id":"r3","run_number":3,"status":"RUNNING"},{"workflow_run_id":"r4","run_number":4,"status":"CANCELED"}]}`), nil
					default:
						t.Fatalf("unexpected page %q", req.URL.Query().Get("page"))
						return nil, nil
					}
				}),
			}, nil
		},
		Repository: "owner/repo",
		Limit:      3,
		LimitSet:   true,
		Paginate:   true,
		PerPage:    2,
		PerPageSet: true,
		JSON:       true,
	}

	if err := listRun(opts); err != nil {
		t.Fatalf("listRun() error = %v", err)
	}
	for _, want := range []string{
		"/api/v8/repos/owner/repo/actions/runs?page=1&per_page=2",
		"/api/v8/repos/owner/repo/actions/runs?page=2&per_page=2",
	} {
		if !slices.Contains(gotPaths, want) {
			t.Fatalf("paths = %#v, missing %q", gotPaths, want)
		}
	}
	var runs []map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &runs); err != nil {
		t.Fatalf("output is not JSON: %v; output=%q", err, out.String())
	}
	if len(runs) != 3 {
		t.Fatalf("len(runs) = %d, want 3; output=%s", len(runs), out.String())
	}
}

func TestListRunJSON(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := iostreams.Test()
	opts := &ListOptions{
		IO: io,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: testutil.NewRoundTripFunc(func(req *http.Request) (*http.Response, error) {
					return listTestResponse(http.StatusOK, `{"total_count":1,"workflow_runs":[{"workflow_run_id":"run-1","workflow_id":"wf-1","workflow_name":"CI","title":"run CI","status":"COMPLETED","event":"Push","run_number":7,"head_branch":"main","head_sha":"abc","actor":{"id":"1","object_id":"u1","login":"dev","name":"Dev"},"start_time":1700000000,"end_time":1700000100,"pause_time":0}]}`), nil
				}),
			}, nil
		},
		Repository: "owner/repo",
		Limit:      30,
		JSON:       true,
	}

	if err := listRun(opts); err != nil {
		t.Fatalf("listRun() error = %v", err)
	}

	var runs []map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &runs); err != nil {
		t.Fatalf("output is not JSON: %v; output=%q", err, out.String())
	}
	if len(runs) != 1 {
		t.Fatalf("len(runs) = %d, want 1", len(runs))
	}
	if runs[0]["status"] != "COMPLETED" {
		t.Fatalf("status = %v, want COMPLETED", runs[0]["status"])
	}
	if runs[0]["run_number"] != float64(7) {
		t.Fatalf("run_number = %v, want 7", runs[0]["run_number"])
	}
}

func TestListRunEmpty(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := iostreams.Test()
	opts := &ListOptions{
		IO: io,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: testutil.NewRoundTripFunc(func(req *http.Request) (*http.Response, error) {
					return listTestResponse(http.StatusOK, `{"total_count":0,"workflow_runs":[]}`), nil
				}),
			}, nil
		},
		Repository: "owner/repo",
		Limit:      30,
	}

	if err := listRun(opts); err != nil {
		t.Fatalf("listRun() error = %v", err)
	}
	if out.String() != "No pipeline runs found\n" {
		t.Fatalf("output = %q, want empty message", out.String())
	}
}

func TestListRunEmptyJSON(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := iostreams.Test()
	opts := &ListOptions{
		IO: io,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: testutil.NewRoundTripFunc(func(req *http.Request) (*http.Response, error) {
					return listTestResponse(http.StatusOK, `{"total_count":0,"workflow_runs":[]}`), nil
				}),
			}, nil
		},
		Repository: "owner/repo",
		Limit:      30,
		JSON:       true,
	}

	if err := listRun(opts); err != nil {
		t.Fatalf("listRun() error = %v", err)
	}
	if out.String() != "[]\n" {
		t.Fatalf("output = %q, want []\\n", out.String())
	}
}

func TestListRunInvalidLimit(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Repository: "owner/repo",
		Limit:      0,
	}

	err := listRun(opts)
	if err == nil {
		t.Fatal("listRun() error = nil, want error")
	}
}

func TestListRunPaginateWithPage(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Repository: "owner/repo",
		Limit:      30,
		Paginate:   true,
		Page:       2,
	}

	err := listRun(opts)
	if err == nil {
		t.Fatal("listRun() error = nil, want error for --paginate with --page")
	}
}

func TestResolveOutputFormat(t *testing.T) {
	tests := []struct {
		name     string
		jsonFlag bool
		raw      string
		want     output.Format
		wantErr  bool
	}{
		{name: "default", want: output.FormatSimple},
		{name: "table", raw: "table", want: output.FormatTable},
		{name: "json flag", jsonFlag: true, want: output.FormatJSON},
		{name: "json with format json", jsonFlag: true, raw: "json", want: output.FormatJSON},
		{name: "invalid format", raw: "yaml", wantErr: true},
		{name: "json incompatible format", jsonFlag: true, raw: "table", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveOutputFormat(tt.jsonFlag, tt.raw)
			if (err != nil) != tt.wantErr {
				t.Fatalf("resolveOutputFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Fatalf("resolveOutputFormat() = %q, want %q", got, tt.want)
			}
		})
	}
}

func listTestResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func assertNoAccessTokenQuery(t *testing.T, path string) {
	t.Helper()
	if strings.Contains(path, "access_token=") {
		t.Fatalf("request path unexpectedly contains access_token query: %q", path)
	}
}
