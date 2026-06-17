package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAuthConfigPersistsAndReadsStoredToken(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_TOKEN", "")
	t.Setenv("GITCODE_TOKEN", "")

	cfg := New()
	authCfg := cfg.Authentication()

	changed, err := authCfg.Login("gitcode.com", "tester", "stored-token", "ssh", false)
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if !changed {
		t.Fatalf("Login() changed = false, want true")
	}

	token, source := authCfg.ActiveToken("gitcode.com")
	if token != "stored-token" || source != "config" {
		t.Fatalf("ActiveToken() = %q, %q", token, source)
	}

	user, err := authCfg.ActiveUser("gitcode.com")
	if err != nil {
		t.Fatalf("ActiveUser() error = %v", err)
	}
	if user != "tester" {
		t.Fatalf("ActiveUser() = %q", user)
	}

	if protocol := cfg.GitProtocol("gitcode.com"); protocol.Value != "ssh" || protocol.Source != "config" {
		t.Fatalf("GitProtocol() = %+v", protocol)
	}
}

func TestAuthConfigEnvironmentOverridesStoredToken(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_TOKEN", "env-token")
	t.Setenv("GITCODE_TOKEN", "")

	cfg := New()
	authCfg := cfg.Authentication()
	if _, err := authCfg.Login("gitcode.com", "tester", "stored-token", "https", false); err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	token, source := authCfg.ActiveToken("gitcode.com")
	if token != "env-token" || source != "GC_TOKEN" {
		t.Fatalf("ActiveToken() = %q, %q", token, source)
	}
}

func TestAuthConfigStoredTokenIgnoresEnvironment(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_TOKEN", "env-token")
	t.Setenv("GITCODE_TOKEN", "")

	cfg := New()
	authCfg := cfg.Authentication()
	if _, err := authCfg.Login("other.example.com", "tester", "stored-token", "https", false); err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	token, source := authCfg.StoredToken("other.example.com")
	if token != "stored-token" || source != "config" {
		t.Fatalf("StoredToken() = %q, %q", token, source)
	}
}

func TestAuthConfigLogoutRemovesStoredToken(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_TOKEN", "")
	t.Setenv("GITCODE_TOKEN", "")

	cfg := New()
	authCfg := cfg.Authentication()
	if _, err := authCfg.Login("gitcode.com", "tester", "stored-token", "https", false); err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	if err := authCfg.Logout("gitcode.com", "tester"); err != nil {
		t.Fatalf("Logout() error = %v", err)
	}

	token, source := authCfg.ActiveToken("gitcode.com")
	if token != "" || source != "" {
		t.Fatalf("ActiveToken() after logout = %q, %q", token, source)
	}
}

func TestAuthConfigLoginRejectsUnsupportedSecureStorage(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())

	cfg := New()
	authCfg := cfg.Authentication()
	changed, err := authCfg.Login("gitcode.com", "tester", "stored-token", "https", true)
	if err == nil {
		t.Fatal("Login() error = nil, want unsupported secure storage error")
	}
	if changed {
		t.Fatal("Login() changed = true, want false")
	}
}

func TestConfigWriteCreatesRestrictedDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix permissions not supported on Windows")
	}
	dir := t.TempDir()
	cfg := &config{configDir: filepath.Join(dir, "gc")}

	if err := os.MkdirAll(cfg.configDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	if err := cfg.Write(); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	info, err := os.Stat(cfg.configDir)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("config dir permissions = %o, want no group/other access", info.Mode().Perm())
	}
}

func TestConfigSetGetAndWritePersistValues(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_EDITOR", "")
	t.Setenv("EDITOR", "")
	t.Setenv("GC_BROWSER", "")

	cfg := New()
	if err := cfg.Set("gitcode.com", "editor", "nano"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if err := cfg.Set("gitcode.com", "browser", "firefox"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if err := cfg.Write(); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	cfg = New()
	editor, err := cfg.Get("gitcode.com", "editor")
	if err != nil {
		t.Fatalf("Get(editor) error = %v", err)
	}
	if editor != "nano" {
		t.Fatalf("Get(editor) = %q, want nano", editor)
	}
	if got := cfg.Editor("gitcode.com"); got.Value != "nano" || got.Source != "config" {
		t.Fatalf("Editor() = %+v, want config nano", got)
	}
	if got := cfg.Browser("gitcode.com"); got.Value != "firefox" || got.Source != "config" {
		t.Fatalf("Browser() = %+v, want config firefox", got)
	}
}

func TestConfigEnvironmentOverridesStoredValue(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_EDITOR", "code")

	cfg := New()
	if err := cfg.Set("gitcode.com", "editor", "nano"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	editor, err := cfg.Get("gitcode.com", "editor")
	if err != nil {
		t.Fatalf("Get(editor) error = %v", err)
	}
	if editor != "code" {
		t.Fatalf("Get(editor) = %q, want environment value code", editor)
	}
}

func TestConfigRejectsUnsupportedKeys(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())

	cfg := New()
	if err := cfg.Set("gitcode.com", "token", "secret"); err == nil {
		t.Fatal("Set(token) error = nil, want unsupported config key error")
	}
	if _, err := cfg.Get("gitcode.com", "token"); err == nil {
		t.Fatal("Get(token) error = nil, want unsupported config key error")
	}
}

func TestNormalizeTrustedHost(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		want    string
		wantErr bool
	}{
		{name: "empty defaults", want: "gitcode.com"},
		{name: "trims and lowercases", host: " Enterprise.Example.COM ", want: "enterprise.example.com"},
		{name: "rejects scheme", host: "https://gitcode.com", wantErr: true},
		{name: "rejects path", host: "gitcode.com/path", wantErr: true},
		{name: "rejects port", host: "gitcode.com:8443", wantErr: true},
		{name: "rejects userinfo", host: "user@gitcode.com", wantErr: true},
		{name: "rejects embedded whitespace", host: "git code.com", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeTrustedHost(tt.host)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NormalizeTrustedHost() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Fatalf("NormalizeTrustedHost() = %q, want %q", got, tt.want)
			}
		})
	}
}
