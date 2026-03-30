package config

import "testing"

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
