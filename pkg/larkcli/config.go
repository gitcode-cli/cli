package larkcli

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gitcode.com/gitcode-cli/cli/pkg/config"
)

// larkConfig is the persisted lark integration state stored in lark.json.
// It is intentionally separate from pkg/config's host-scoped config.json
// because a Feishu chat id is not tied to a GitCode host.
type larkConfig struct {
	DefaultChatID string `json:"default_chat_id,omitempty"`
}

// ConfigPath returns the lark config file path (~/.config/gc/lark.json or
// $GC_CONFIG_DIR/lark.json).
func ConfigPath() string {
	return filepath.Join(configDir(), "lark.json")
}

// DefaultChatID resolves the default Feishu chat id from, in priority order:
// the GC_LARK_DEFAULT_CHAT_ID environment variable, then the persisted
// lark.json. Returns an empty string when nothing is configured.
func DefaultChatID() string {
	if v := strings.TrimSpace(os.Getenv(EnvDefaultChat)); v != "" {
		return v
	}
	cfg, err := readConfig()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(cfg.DefaultChatID)
}

// SaveDefaultChat persists the default chat id to lark.json.
func SaveDefaultChat(chatID string) error {
	chatID = strings.TrimSpace(chatID)
	cfg, _ := readConfig()
	cfg.DefaultChatID = chatID
	return writeConfig(cfg)
}

// ClearDefaultChat removes the persisted default chat id.
func ClearDefaultChat() error {
	cfg, _ := readConfig()
	cfg.DefaultChatID = ""
	return writeConfig(cfg)
}

func configDir() string {
	if dir := os.Getenv("GC_CONFIG_DIR"); dir != "" {
		return dir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}
	return filepath.Join(home, ".config", "gc")
}

func readConfig() (larkConfig, error) {
	data, err := os.ReadFile(ConfigPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return larkConfig{}, nil
		}
		return larkConfig{}, err
	}
	var cfg larkConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return larkConfig{}, err
	}
	return cfg, nil
}

func writeConfig(cfg larkConfig) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	if err := os.Chmod(dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return config.SecureWriteFile(ConfigPath(), data, 0o600)
}
