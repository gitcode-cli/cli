package logout

import (
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/internal/config"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestLogoutRunRemovesStoredToken(t *testing.T) {
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

	opts := &LogoutOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Config: func() (config.Config, error) {
			return config.New(), nil
		},
	}

	if err := logoutRun(opts); err != nil {
		t.Fatalf("logoutRun() error = %v", err)
	}

	token, source := config.New().Authentication().ActiveToken("gitcode.com")
	if token != "" || source != "" {
		t.Fatalf("ActiveToken() after logout = %q, %q", token, source)
	}
	if !strings.Contains(out.String(), "Cleared stored authentication for gitcode.com") {
		t.Fatalf("output = %q", out.String())
	}
}
