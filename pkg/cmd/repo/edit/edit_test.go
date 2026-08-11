package edit

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestNewCmdEdit(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{name: "with description", args: []string{"--description", "test"}, wantErr: false},
		{name: "with private", args: []string{"--private"}, wantErr: false},
		{name: "with name", args: []string{"--name", "newname"}, wantErr: false},
		{name: "with json", args: []string{"--description", "x", "--json"}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCmdEdit(cmdutil.TestFactory(), func(opts *EditOptions) error {
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

func TestNewCmdEditFlagsExist(t *testing.T) {
	cmd := NewCmdEdit(cmdutil.TestFactory(), nil)
	for _, flag := range []string{"repo", "description", "homepage", "default-branch", "name", "private", "public", "json"} {
		if cmd.Flags().Lookup(flag) == nil {
			t.Fatalf("flag %q missing", flag)
		}
	}
}

func TestEditRunSuccess(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")
	io, _, _, errOut := iostreams.Test()
	var gotMethod string
	opts := &EditOptions{
		IO: io,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: testutil.NewRoundTripFunc(func(req *http.Request) (*http.Response, error) {
					gotMethod = req.Method
					return editTestResponse(http.StatusOK, `{"name":"repo","full_name":"owner/repo"}`), nil
				}),
			}, nil
		},
		Repository:  "owner/repo",
		Description: "New description",
	}

	if err := editRun(opts); err != nil {
		t.Fatalf("editRun() error = %v", err)
	}
	if gotMethod != "PATCH" {
		t.Fatalf("method = %q, want PATCH", gotMethod)
	}
	if !strings.Contains(errOut.String(), "Updated repository owner/repo") {
		t.Fatalf("stderr should contain confirmation; got: %s", errOut.String())
	}
}

func TestEditRunJSON(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")
	io, _, out, _ := iostreams.Test()
	opts := &EditOptions{
		IO: io,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: testutil.NewRoundTripFunc(func(req *http.Request) (*http.Response, error) {
					return editTestResponse(http.StatusOK, `{"name":"repo","description":"updated"}`), nil
				}),
			}, nil
		},
		Repository:  "owner/repo",
		Description: "updated",
		JSON:        true,
	}

	if err := editRun(opts); err != nil {
		t.Fatalf("editRun() error = %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(out.String()), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v; raw: %s", err, out.String())
	}
	if result["name"] != "repo" {
		t.Fatalf("name = %v, want repo", result["name"])
	}
}

func TestEditRunPrivatePublicConflict(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")
	io, _, _, _ := iostreams.Test()
	opts := &EditOptions{
		IO:         io,
		Repository: "owner/repo",
		Private:    true,
		Public:     true,
	}

	err := editRun(opts)
	if err == nil {
		t.Fatal("editRun() error = nil, want error for --private + --public")
	}
}

func TestEditRunNoFields(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")
	io, _, _, _ := iostreams.Test()
	opts := &EditOptions{
		IO:         io,
		Repository: "owner/repo",
	}

	err := editRun(opts)
	if err == nil {
		t.Fatal("editRun() error = nil, want error for no fields")
	}
}

func TestEditRunPrivate(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")
	io, _, _, _ := iostreams.Test()
	var gotBody string
	opts := &EditOptions{
		IO: io,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: testutil.NewRoundTripFunc(func(req *http.Request) (*http.Response, error) {
					buf := make([]byte, 1024)
					n, _ := req.Body.Read(buf)
					gotBody = string(buf[:n])
					return editTestResponse(http.StatusOK, `{"name":"repo","private":true}`), nil
				}),
			}, nil
		},
		Repository: "owner/repo",
		Private:    true,
	}

	if err := editRun(opts); err != nil {
		t.Fatalf("editRun() error = %v", err)
	}
	if !strings.Contains(gotBody, `"private":true`) {
		t.Fatalf("request body should contain private:true; got: %s", gotBody)
	}
}

func TestEditRunError(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")
	io, _, _, _ := iostreams.Test()
	opts := &EditOptions{
		IO: io,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: testutil.NewRoundTripFunc(func(req *http.Request) (*http.Response, error) {
					return editTestResponse(http.StatusUnauthorized, `{"message":"unauthorized"}`), nil
				}),
			}, nil
		},
		Repository:  "owner/repo",
		Description: "test",
	}

	err := editRun(opts)
	if err == nil {
		t.Fatal("editRun() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "failed to update repository") {
		t.Fatalf("error = %q, want to wrap update failure", err.Error())
	}
}

func TestEditRunSecretScan(t *testing.T) {
	t.Setenv("GC_TOKEN", "leaked-secret-xyz")
	io, _, _, _ := iostreams.Test()
	opts := &EditOptions{
		IO: io,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{Transport: testutil.NewRoundTripFunc(func(req *http.Request) (*http.Response, error) {
				t.Fatal("HTTP request should not be made when secret is detected")
				return nil, nil
			})}, nil
		},
		Repository:  "owner/repo",
		Description: "token: leaked-secret-xyz",
	}

	err := editRun(opts)
	if err == nil {
		t.Fatal("editRun() error = nil, want error for secret in description")
	}
}

func editTestResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
