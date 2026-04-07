package token

import (
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/internal/config"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
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
	if strings.TrimSpace(out.String()) != "stored-token" {
		t.Fatalf("output = %q", out.String())
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
