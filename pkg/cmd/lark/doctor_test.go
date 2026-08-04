package lark

import (
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/larkcli"
)

// notInstalledRun is a RunFunc that simulates lark-cli missing (FindLarkCLI
// resolves nothing in the doctor path because we stub the run directly).
func notInstalledRun(args []string) (*larkcli.Result, error) {
	return nil, larkcli.ErrNotInstalled
}

func TestDoctorRun_NotInstalled(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_LARK_DEFAULT_CHAT_ID", "")
	t.Setenv("GC_LARK_CLI_BIN", "/nonexistent/missing-lark")
	t.Setenv("PATH", t.TempDir())

	io, _, out, _ := iostreams.Test()
	opts := &doctorOptions{IO: io, LarkRun: notInstalledRun}
	if err := doctorRun(opts); err != nil {
		t.Fatalf("err = %v", err)
	}
	if !strings.Contains(out.String(), "not installed") {
		t.Errorf("out = %q, want not installed", out.String())
	}
}

func TestDoctorRun_InstalledAndReadyJSON(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_LARK_DEFAULT_CHAT_ID", "oc_default")
	t.Setenv("GC_LARK_CLI_BIN", "/bin/true") // any existing executable so FindLarkCLI returns a path

	io, _, out, _ := iostreams.Test()
	opts := &doctorOptions{
		IO:      io,
		JSON:    true,
		LarkRun: stubRun("", `{"ok":true,"identity":"user"}`, 0),
	}
	if err := doctorRun(opts); err != nil {
		t.Fatalf("err = %v", err)
	}
	s := out.String()
	if !strings.Contains(s, `"installed": true`) {
		t.Errorf("out = %q, want installed true", s)
	}
	if !strings.Contains(s, `"login_ready": true`) {
		t.Errorf("out = %q, want login_ready true", s)
	}
	if !strings.Contains(s, `"default_chat": "oc_default"`) {
		t.Errorf("out = %q, want default_chat oc_default", s)
	}
}

func TestDoctorRun_InstalledButNotLoggedIn(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_LARK_DEFAULT_CHAT_ID", "")
	t.Setenv("GC_LARK_CLI_BIN", "/bin/true")
	t.Setenv("PATH", t.TempDir())

	io, _, out, _ := iostreams.Test()
	opts := &doctorOptions{
		IO:      io,
		LarkRun: stubRun("", "", 1), // non-zero exit -> not ready
	}
	if err := doctorRun(opts); err != nil {
		t.Fatalf("err = %v", err)
	}
	if !strings.Contains(out.String(), "login not ready") {
		t.Errorf("out = %q, want login not ready", out.String())
	}
}
