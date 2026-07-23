package login

import (
	"bytes"
	"gitcode.com/gitcode-cli/cli/pkg/testutil"
	"net/http"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/config"
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

func TestNewCmdLoginBindsWebAndHostnameFlags(t *testing.T) {
	f := cmdutil.TestFactory()
	var gotWeb bool
	var gotHostname string
	cmd := NewCmdLogin(f, func(opts *LoginOptions) error {
		gotWeb = opts.Web
		gotHostname = opts.Hostname
		return nil
	})
	cmd.SetArgs([]string{"--web", "--hostname", "gitcode.com"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !gotWeb {
		t.Fatal("Web = false, want true")
	}
	if gotHostname != "gitcode.com" {
		t.Fatalf("Hostname = %q, want gitcode.com", gotHostname)
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
				Transport: testutil.NewRoundTripFunc(func(req *http.Request) (*http.Response, error) {
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

	if openedURL != "https://gitcode.com/setting/token-classic/create" {
		t.Fatalf("opened URL = %q", openedURL)
	}
	if !strings.Contains(out.String(), "Opening https://gitcode.com/setting/token-classic/create in your browser.") {
		t.Fatalf("output = %q", out.String())
	}
}

func TestLoginWithWebRejectsCustomHost(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := &LoginOptions{
		IO:       io,
		Hostname: "enterprise.example.com",
		OpenBrowser: func(url string) error {
			t.Fatalf("unexpected browser open: %s", url)
			return nil
		},
	}

	err := loginWithWeb(opts)
	if err == nil {
		t.Fatal("loginWithWeb() error = nil, want unsupported host")
	}
	if !strings.Contains(err.Error(), "--web only supports gitcode.com") {
		t.Fatalf("loginWithWeb() error = %q, want unsupported host", err.Error())
	}
}

func TestLoginWithTokenRejectsMalformedHost(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	io, _, _, _ := iostreams.Test()

	opts := &LoginOptions{
		IO:       io,
		Hostname: "https://gitcode.com",
		Token:    "test-token",
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: testutil.NewRoundTripFunc(func(req *http.Request) (*http.Response, error) {
					t.Fatalf("unexpected request to %s", req.URL.String())
					return nil, nil
				}),
			}, nil
		},
		Config: func() (config.Config, error) {
			return config.New(), nil
		},
	}

	err := loginWithTokenFlag(opts)
	if err == nil {
		t.Fatal("loginWithTokenFlag() error = nil, want invalid host")
	}
	if !strings.Contains(err.Error(), "invalid host") {
		t.Fatalf("loginWithTokenFlag() error = %q, want invalid host", err.Error())
	}
}

func TestLoginWithWebRejectsMalformedHost(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := &LoginOptions{
		IO:       io,
		Hostname: "bad/host",
		OpenBrowser: func(url string) error {
			t.Fatalf("unexpected browser open: %s", url)
			return nil
		},
	}

	err := loginWithWeb(opts)
	if err == nil {
		t.Fatal("loginWithWeb() error = nil, want invalid host")
	}
	if !strings.Contains(err.Error(), "invalid host") {
		t.Fatalf("loginWithWeb() error = %q, want invalid host", err.Error())
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

func ioNopCloser(body string) *readCloser {
	return &readCloser{Reader: strings.NewReader(body)}
}

type readCloser struct {
	*strings.Reader
}

func (r *readCloser) Close() error { return nil }
