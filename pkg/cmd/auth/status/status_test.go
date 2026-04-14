package status

import (
	"net/http"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/internal/config"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdStatus(t *testing.T) {
	f := cmdutil.TestFactory()
	cmd := NewCmdStatus(f, nil)
	if cmd == nil {
		t.Fatal("NewCmdStatus returned nil")
	}
	if cmd.Use != "status" {
		t.Errorf("Expected Use 'status', got %q", cmd.Use)
	}
	if strings.Contains(cmd.Long, "keyring") {
		t.Fatalf("status help unexpectedly references keyring: %q", cmd.Long)
	}
	if !strings.Contains(cmd.Long, "auth.json") {
		t.Fatalf("status help should reference local config storage: %q", cmd.Long)
	}
}

func TestStatusRunUsesStoredToken(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_TOKEN", "")
	t.Setenv("GITCODE_TOKEN", "")

	cfg := config.New()
	if _, err := cfg.Authentication().Login("gitcode.com", "tester", "stored-token", "ssh", false); err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	f := cmdutil.TestFactory()
	out := &strings.Builder{}
	f.IOStreams.Out = out

	opts := &StatusOptions{
		IO: f.IOStreams,
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
	}

	if err := statusRun(opts); err != nil {
		t.Fatalf("statusRun() error = %v", err)
	}
	if !strings.Contains(out.String(), "Logged in as tester (config)") {
		t.Fatalf("output = %q", out.String())
	}
	if !strings.Contains(out.String(), "Git operations protocol: ssh") {
		t.Fatalf("output = %q", out.String())
	}
}

func TestStatusRunWithHostnameUsesStoredTokenInsteadOfEnvOverride(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_TOKEN", "env-token")
	t.Setenv("GITCODE_TOKEN", "")

	cfg := config.New()
	if _, err := cfg.Authentication().Login("other.example.com", "stored-user", "stored-token", "ssh", false); err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	f := cmdutil.TestFactory()
	out := &strings.Builder{}
	f.IOStreams.Out = out

	opts := &StatusOptions{
		IO:          f.IOStreams,
		Hostname:    "other.example.com",
		HostnameSet: true,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					auth := req.Header.Get("Authorization")
					if auth != "Bearer stored-token" {
						t.Fatalf("Authorization = %q, want %q", auth, "Bearer stored-token")
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     make(http.Header),
						Body:       ioNopCloser(`{"login":"stored-user"}`),
					}, nil
				}),
			}, nil
		},
		Config: func() (config.Config, error) {
			return config.New(), nil
		},
	}

	if err := statusRun(opts); err != nil {
		t.Fatalf("statusRun() error = %v", err)
	}
	if !strings.Contains(out.String(), "Logged in as stored-user (config)") {
		t.Fatalf("output = %q", out.String())
	}
}

func TestStatusRunShowTokenDisplaysFullToken(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_TOKEN", "secret-token")
	t.Setenv("GITCODE_TOKEN", "")

	f := cmdutil.TestFactory()
	out := &strings.Builder{}
	f.IOStreams.Out = out

	opts := &StatusOptions{
		IO:        f.IOStreams,
		ShowToken: true,
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
	}

	if err := statusRun(opts); err != nil {
		t.Fatalf("statusRun() error = %v", err)
	}
	if !strings.Contains(out.String(), "Token: secret-token") {
		t.Fatalf("output = %q", out.String())
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
