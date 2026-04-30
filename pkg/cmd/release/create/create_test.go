package create

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdCreate(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "create with tag",
			args:    []string{"v1.0.0"},
			wantErr: false,
		},
		{
			name:    "create with tag and title",
			args:    []string{"v1.0.0", "--title", "Version 1.0"},
			wantErr: false,
		},
		{
			name:    "create with draft flag",
			args:    []string{"v1.0.0", "--draft"},
			wantErr: false,
		},
		{
			name:    "create with prerelease flag",
			args:    []string{"v1.0.0-beta", "--prerelease"},
			wantErr: false,
		},
		{
			name:    "create with json output",
			args:    []string{"v1.0.0", "--json"},
			wantErr: false,
		},
		{
			name:    "no tag specified",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdCreate(f, func(opts *CreateOptions) error {
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

func TestCreateRunJSONWritesCreatedRelease(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	opts := &CreateOptions{
		IO:         f.IOStreams,
		Repository: "owner/repo",
		TagName:    "v1.0.0",
		Title:      "Version 1.0.0",
		Notes:      "Release notes",
		JSON:       true,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: releaseRoundTripFunc(func(req *http.Request) (*http.Response, error) {
					if req.Method != http.MethodPost || req.URL.Path != "/api/v5/repos/owner/repo/releases" {
						t.Fatalf("unexpected request: %s %s", req.Method, req.URL.Path)
					}
					return releaseResponse(http.StatusOK, `{"id":1,"tag_name":"v1.0.0","name":"Version 1.0.0","html_url":"https://gitcode.com/owner/repo/releases/v1.0.0"}`), nil
				}),
			}, nil
		},
	}

	if err := createRun(opts); err != nil {
		t.Fatalf("createRun() error = %v", err)
	}

	var got map[string]interface{}
	out := f.IOStreams.Out.(*bytes.Buffer).Bytes()
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("JSON output did not parse: %v\n%s", err, string(out))
	}
	if got["tag_name"] != "v1.0.0" || got["html_url"] != "https://gitcode.com/owner/repo/releases/v1.0.0" {
		t.Fatalf("JSON output = %#v", got)
	}
	if strings.Contains(string(out), "Created release") {
		t.Fatalf("JSON output contains text banner: %q", string(out))
	}
}

type releaseRoundTripFunc func(*http.Request) (*http.Response, error)

func (fn releaseRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func releaseResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
