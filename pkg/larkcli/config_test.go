package larkcli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigPath_UsesEnvDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("GC_CONFIG_DIR", dir)
	if got := ConfigPath(); got != filepath.Join(dir, "lark.json") {
		t.Errorf("ConfigPath() = %q, want %q", got, filepath.Join(dir, "lark.json"))
	}
}

func TestDefaultChatID_EnvOverridesConfig(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("GC_CONFIG_DIR", dir)
	t.Setenv(EnvDefaultChat, "")
	if err := SaveDefaultChat("oc_from_config"); err != nil {
		t.Fatalf("SaveDefaultChat err = %v", err)
	}
	t.Setenv(EnvDefaultChat, "oc_from_env")
	if got := DefaultChatID(); got != "oc_from_env" {
		t.Errorf("DefaultChatID() = %q, want oc_from_env (env overrides config)", got)
	}
}

func TestDefaultChatID_FallsBackToConfig(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("GC_CONFIG_DIR", dir)
	t.Setenv(EnvDefaultChat, "")
	if err := SaveDefaultChat("oc_persisted"); err != nil {
		t.Fatalf("SaveDefaultChat err = %v", err)
	}
	if got := DefaultChatID(); got != "oc_persisted" {
		t.Errorf("DefaultChatID() = %q, want oc_persisted", got)
	}
}

func TestDefaultChatID_NothingConfigured(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("GC_CONFIG_DIR", dir)
	t.Setenv(EnvDefaultChat, "")
	if got := DefaultChatID(); got != "" {
		t.Errorf("DefaultChatID() = %q, want empty", got)
	}
}

func TestClearDefaultChat(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("GC_CONFIG_DIR", dir)
	t.Setenv(EnvDefaultChat, "")
	if err := SaveDefaultChat("oc_x"); err != nil {
		t.Fatalf("SaveDefaultChat err = %v", err)
	}
	if err := ClearDefaultChat(); err != nil {
		t.Fatalf("ClearDefaultChat err = %v", err)
	}
	if got := DefaultChatID(); got != "" {
		t.Errorf("DefaultChatID() = %q, want empty after clear", got)
	}
}

func TestSaveDefaultChat_TrimsAndPersists(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("GC_CONFIG_DIR", dir)
	t.Setenv(EnvDefaultChat, "")
	if err := SaveDefaultChat("  oc_trimmed  "); err != nil {
		t.Fatalf("SaveDefaultChat err = %v", err)
	}
	if got := DefaultChatID(); got != "oc_trimmed" {
		t.Errorf("DefaultChatID() = %q, want oc_trimmed", got)
	}
	// File must be created with 0600 on unix and exist.
	info, err := os.Stat(ConfigPath())
	if err != nil {
		t.Fatalf("Stat config: %v", err)
	}
	if mode := info.Mode().Perm(); mode != 0o600 {
		t.Errorf("config perm = %o, want 0600", mode)
	}
}
