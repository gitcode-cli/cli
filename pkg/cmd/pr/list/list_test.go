package list

import (
	"io"
	"net/http"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/output"
)

func TestNewCmdList(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "list default",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "list closed PRs",
			args:    []string{"--state", "closed"},
			wantErr: false,
		},
		{
			name:    "list with limit",
			args:    []string{"--limit", "10"},
			wantErr: false,
		},
		{
			name:    "list with base filter",
			args:    []string{"--base", "main"},
			wantErr: false,
		},
		{
			name:    "list with sort options",
			args:    []string{"--sort", "updated", "--direction", "desc", "--page", "2"},
			wantErr: false,
		},
		{
			name:    "list with table format",
			args:    []string{"--format", "table"},
			wantErr: false,
		},
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

func TestListRunBuildsFullQuery(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	var gotPath string
	opts := &ListOptions{
		IO: f.IOStreams,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					gotPath = req.URL.Path
					if req.URL.RawQuery != "" {
						gotPath += "?" + req.URL.RawQuery
					}
					return listTestResponse(http.StatusOK, `[]`), nil
				}),
			}, nil
		},
		Repository: "owner/repo",
		State:      "open",
		Head:       "feature/login",
		Base:       "main",
		Sort:       "updated",
		Direction:  "asc",
		Limit:      25,
		Page:       2,
	}

	if err := listRun(opts); err != nil {
		t.Fatalf("listRun() error = %v", err)
	}

	want := "/api/v5/repos/owner/repo/pulls?base=main&direction=asc&head=feature%2Flogin&page=2&per_page=25&sort=updated&state=open"
	if gotPath != want {
		t.Fatalf("request path = %q, want %q", gotPath, want)
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

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func listTestResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
