package list

import (
	"io"
	"net/http"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
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
			name:    "list with limit",
			args:    []string{"--limit", "10"},
			wantErr: false,
		},
		{
			name:    "list with repo",
			args:    []string{"-R", "owner/repo"},
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

func TestListRunMarksOnlyFirstPublishedReleaseAsLatest(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	out := &strings.Builder{}
	f.IOStreams.Out = out

	err := listRun(&ListOptions{
		IO: f.IOStreams,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     http.StatusText(http.StatusOK),
						Header:     make(http.Header),
						Body: io.NopCloser(strings.NewReader(`[
							{"tag_name":"v2.0.0","html_url":"https://gitcode.com/owner/repo/-/releases/v2.0.0","draft":false,"prerelease":false},
							{"tag_name":"v1.9.0","html_url":"https://gitcode.com/owner/repo/-/releases/v1.9.0","draft":false,"prerelease":false},
							{"tag_name":"v2.1.0-rc1","html_url":"https://gitcode.com/owner/repo/-/releases/v2.1.0-rc1","draft":false,"prerelease":true}
						]`)),
					}, nil
				}),
			}, nil
		},
		Repository: "owner/repo",
	})
	if err != nil {
		t.Fatalf("listRun() error = %v", err)
	}

	output := out.String()
	if strings.Count(output, "(latest)") != 1 {
		t.Fatalf("output = %q, want exactly one latest marker", output)
	}
	if !strings.Contains(output, "(published)") {
		t.Fatalf("output = %q, want published marker", output)
	}
	if !strings.Contains(output, "(pre-release)") {
		t.Fatalf("output = %q, want pre-release marker", output)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
