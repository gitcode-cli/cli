package create

import (
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
			name:    "create with title",
			args:    []string{"--title", "Test Issue"},
			wantErr: false,
		},
		{
			name:    "create with title and body",
			args:    []string{"--title", "Test", "--body", "Description"},
			wantErr: false,
		},
		{
			name:    "create with labels",
			args:    []string{"--title", "Test", "--label", "bug,enhancement"},
			wantErr: false,
		},
		{
			name:    "no title",
			args:    []string{},
			wantErr: false, // Command runs, error in run function
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

func TestCreateRunFailsWhenAssigneesAreNotApplied(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	opts := &CreateOptions{
		IO: f.IOStreams,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					switch req.URL.Path {
					case "/api/v5/users/alice":
						return issueResponse(http.StatusOK, `{"id":"101","login":"alice"}`), nil
					case "/api/v5/repos/owner/repo/issues":
						return issueResponse(http.StatusOK, `{"number":"12","html_url":"https://gitcode.com/owner/repo/issues/12"}`), nil
					case "/api/v5/repos/owner/repo/issues/12":
						return issueResponse(http.StatusOK, `{"number":"12","assignees":[]}`), nil
					default:
						t.Fatalf("unexpected request: %s", req.URL.Path)
						return nil, nil
					}
				}),
			}, nil
		},
		Repository: "owner/repo",
		Title:      "Bug report",
		Assignees:  []string{"alice"},
	}

	err := createRun(opts)
	if err == nil {
		t.Fatal("createRun() error = nil, want assignee verification error")
	}
	if !strings.Contains(err.Error(), "did not apply the requested assignees") {
		t.Fatalf("createRun() error = %v", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func issueResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
