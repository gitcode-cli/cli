package token

import (
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/config"
)

func TestTokenRunUsesStoredToken(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_TOKEN", "")
	t.Setenv("GITCODE_TOKEN", "")

	cfg := config.New()
	if _, err := cfg.Authentication().Login("gitcode.com", "tester", "stored-token", "https", false); err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	f := cmdutil.TestFactory()
	out := &strings.Builder{}
	errOut := &strings.Builder{}
	f.IOStreams.Out = out
	f.IOStreams.ErrOut = errOut

	opts := &TokenOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Config: func() (config.Config, error) {
			return config.New(), nil
		},
	}

	if err := tokenRun(opts); err != nil {
		t.Fatalf("tokenRun() error = %v", err)
	}
	if strings.TrimSpace(out.String()) != "stored-token" {
		t.Fatalf("output = %q", out.String())
	}
	if !strings.Contains(errOut.String(), "Warning: displaying authentication token") {
		t.Fatalf("stderr should contain token warning, got: %q", errOut.String())
	}
}

func TestTokenRunWithHostnameUsesStoredTokenInsteadOfEnvOverride(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_TOKEN", "env-token")
	t.Setenv("GITCODE_TOKEN", "")

	cfg := config.New()
	if _, err := cfg.Authentication().Login("other.example.com", "tester", "stored-token", "https", false); err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	f := cmdutil.TestFactory()
	out := &strings.Builder{}
	f.IOStreams.Out = out

	opts := &TokenOptions{
		IO:          f.IOStreams,
		HttpClient:  f.HttpClient,
		Hostname:    "other.example.com",
		HostnameSet: true,
		Config: func() (config.Config, error) {
			return config.New(), nil
		},
	}

	if err := tokenRun(opts); err != nil {
		t.Fatalf("tokenRun() error = %v", err)
	}
	if got := strings.TrimSpace(out.String()); got != "stored-token" {
		t.Fatalf("output = %q", got)
	}
}

func TestTokenRunDefaultCustomHostUsesStoredTokenInsteadOfEnvOverride(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_HOST", "other.example.com")
	t.Setenv("GC_TOKEN", "env-token")
	t.Setenv("GITCODE_TOKEN", "")

	cfg := config.New()
	if _, err := cfg.Authentication().Login("other.example.com", "tester", "stored-token", "https", false); err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	f := cmdutil.TestFactory()
	out := &strings.Builder{}
	f.IOStreams.Out = out

	opts := &TokenOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Config: func() (config.Config, error) {
			return config.New(), nil
		},
	}

	if err := tokenRun(opts); err != nil {
		t.Fatalf("tokenRun() error = %v", err)
	}
	if got := strings.TrimSpace(out.String()); got != "stored-token" {
		t.Fatalf("output = %q", got)
	}
}

func TestTokenRunRejectsMalformedDefaultHost(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_HOST", "https://gitcode.com")
	t.Setenv("GC_TOKEN", "env-token")
	t.Setenv("GITCODE_TOKEN", "")

	f := cmdutil.TestFactory()
	opts := &TokenOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Config: func() (config.Config, error) {
			return config.New(), nil
		},
	}

	err := tokenRun(opts)
	if err == nil {
		t.Fatal("tokenRun() error = nil, want invalid host")
	}
	if !strings.Contains(err.Error(), "invalid host") {
		t.Fatalf("tokenRun() error = %q, want invalid host", err.Error())
	}
}

func TestTokenRunJSONOutputsWarningToStderr(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_TOKEN", "")
	t.Setenv("GITCODE_TOKEN", "")

	cfg := config.New()
	if _, err := cfg.Authentication().Login("gitcode.com", "tester", "stored-token", "https", false); err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	f := cmdutil.TestFactory()
	out := &strings.Builder{}
	errOut := &strings.Builder{}
	f.IOStreams.Out = out
	f.IOStreams.ErrOut = errOut

	opts := &TokenOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		JSON:       true,
		Config: func() (config.Config, error) {
			return config.New(), nil
		},
	}

	if err := tokenRun(opts); err != nil {
		t.Fatalf("tokenRun() error = %v", err)
	}
	if !strings.Contains(errOut.String(), "Warning: displaying authentication token") {
		t.Fatalf("stderr should contain token warning in JSON mode, got: %q", errOut.String())
	}
	if !strings.Contains(out.String(), "stored-token") {
		t.Fatalf("stdout should contain token in JSON output, got: %q", out.String())
	}
}
