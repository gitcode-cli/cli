package login

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/internal/config"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

func TestNewCmdLogin(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "login with with-token flag",
			args:    []string{"--with-token"},
			wantErr: false,
		},
		{
			name:    "login with hostname flag",
			args:    []string{"--hostname", "gitcode.com", "--with-token"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdLogin(f, func(opts *LoginOptions) error {
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

func TestLoginWithWebOpensBrowser(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	io, _, out, _ := iostreams.Test()
	io.In = bytes.NewBufferString("test-token\n")

	var openedURL string
	opts := &LoginOptions{
		IO: io,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     make(http.Header),
						Body:       ioNopCloser(`{"login":"tester"}`),
					}, nil
				}),
			}, nil
		},
		Config: func() (config.Config, error) {
			return config.New(), nil
		},
		OpenBrowser: func(url string) error {
			openedURL = url
			return nil
		},
		Web: true,
	}

	if err := loginWithWeb(opts); err != nil {
		t.Fatalf("loginWithWeb() error = %v", err)
	}

	if openedURL != "https://gitcode.com/-/profile/personal_access_tokens" {
		t.Fatalf("opened URL = %q", openedURL)
	}
	if !strings.Contains(out.String(), "Opening https://gitcode.com/-/profile/personal_access_tokens in your browser.") {
		t.Fatalf("output = %q", out.String())
	}
}

func TestNewCmdLoginRejectsInteractiveLoginWithoutTTY(t *testing.T) {
	f := cmdutil.TestFactory()
	cmd := NewCmdLogin(f, nil)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want usage error")
	}
	if !strings.Contains(err.Error(), "use --with-token") {
		t.Fatalf("Execute() error = %q", err.Error())
	}
}

func TestNewCmdLoginRejectsRemovedTokenFlag(t *testing.T) {
	f := cmdutil.TestFactory()
	cmd := NewCmdLogin(f, nil)
	cmd.SetArgs([]string{"--token", "test-token"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want unknown flag error")
	}
	if !strings.Contains(err.Error(), "unknown flag: --token") {
		t.Fatalf("Execute() error = %q", err.Error())
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func ioNopCloser(body string) *readCloser {
	return &readCloser{Reader: strings.NewReader(body)}
}

type readCloser struct {
	*strings.Reader
}

func (r *readCloser) Close() error { return nil }
