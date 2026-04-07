package config

import (
	"os"
	"path/filepath"
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
