// Package config provides configuration management for gitcode-cli
package config

import (
	"os"
	"path/filepath"
	"strings"
)

// Config interface defines configuration operations
type Config interface {
	// Get retrieves a configuration value
	Get(host, key string) (string, error)
	// Set stores a configuration value
	Set(host, key, value string) error

	// GitProtocol returns the preferred git protocol
	GitProtocol(host string) ConfigEntry
	// Editor returns the preferred editor
	Editor(host string) ConfigEntry
	// Browser returns the preferred browser
	Browser(host string) ConfigEntry
	// Pager returns the preferred pager
	Pager(host string) ConfigEntry

	// Authentication returns authentication configuration
	Authentication() AuthConfig

	// Write persists the configuration
	Write() error
}

// ConfigEntry represents a configuration value with its source
type ConfigEntry struct {
	Value  string
	Source string // "environment", "config", "default"
}

// AuthConfig interface defines authentication operations
type AuthConfig interface {
	// ActiveToken returns the active token for a host
	ActiveToken(hostname string) (string, string)
	// StoredToken returns the stored token for a host without environment overrides
	StoredToken(hostname string) (string, string)
	// HasActiveToken checks if a token exists for a host
	HasActiveToken(hostname string) bool
	// ActiveUser returns the active user for a host
	ActiveUser(hostname string) (string, error)
	// Hosts returns all configured hosts
	Hosts() []string
	// DefaultHost returns the default host
	DefaultHost() (string, string)
	// Login creates or updates authentication for a host
	Login(hostname, username, token, gitProtocol string, secureStorage bool) (bool, error)
	// Logout removes authentication for a host
	Logout(hostname, username string) error
	// SwitchUser switches the active user for a host
	SwitchUser(hostname, user string) error
}

// config implements Config interface
type config struct {
	configDir string
}

// New creates a new Config
func New() Config {
	return &config{
		configDir: configDir(),
	}
}

// configDir returns the configuration directory
func configDir() string {
	if dir := os.Getenv("GC_CONFIG_DIR"); dir != "" {
		return dir
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "~"
	}

	return filepath.Join(homeDir, ".config", "gc")
}

// Get retrieves a configuration value
func (c *config) Get(host, key string) (string, error) {
	// Check environment variable first
	envKey := "GC_" + strings.ToUpper(key)
	if val := os.Getenv(envKey); val != "" {
		return val, nil
	}

	// TODO: read from config file
	return "", nil
}

// Set stores a configuration value
func (c *config) Set(host, key, value string) error {
	// TODO: implement config file writing
	return nil
}

// GitProtocol returns the preferred git protocol
func (c *config) GitProtocol(host string) ConfigEntry {
	if val := os.Getenv("GC_GIT_PROTOCOL"); val != "" {
		return ConfigEntry{Value: val, Source: "environment"}
	}

	state, err := c.readAuthState()
	if err == nil {
		if user := state.activeUser(host); user != nil && user.GitProtocol != "" {
			return ConfigEntry{Value: user.GitProtocol, Source: "config"}
		}
	}
	return ConfigEntry{Value: "https", Source: "default"}
}

// Editor returns the preferred editor
func (c *config) Editor(host string) ConfigEntry {
	if val := os.Getenv("GC_EDITOR"); val != "" {
		return ConfigEntry{Value: val, Source: "environment"}
	}
	if val := os.Getenv("EDITOR"); val != "" {
		return ConfigEntry{Value: val, Source: "environment"}
	}
	return ConfigEntry{Value: "vim", Source: "default"}
}

// Browser returns the preferred browser
func (c *config) Browser(host string) ConfigEntry {
	if val := os.Getenv("GC_BROWSER"); val != "" {
		return ConfigEntry{Value: val, Source: "environment"}
	}
	return ConfigEntry{Value: "", Source: "default"}
}

// Pager returns the preferred pager
func (c *config) Pager(host string) ConfigEntry {
	if val := os.Getenv("GC_PAGER"); val != "" {
		return ConfigEntry{Value: val, Source: "environment"}
	}
	if val := os.Getenv("PAGER"); val != "" {
		return ConfigEntry{Value: val, Source: "environment"}
	}
	return ConfigEntry{Value: "less", Source: "default"}
}

// Authentication returns authentication configuration
func (c *config) Authentication() AuthConfig {
	return &authConfig{config: c}
}

// Write persists the configuration
func (c *config) Write() error {
	// Ensure config directory exists
	if err := os.MkdirAll(c.configDir, 0o700); err != nil {
		return err
	}
	if err := os.Chmod(c.configDir, 0o700); err != nil {
		return err
	}
	// TODO: write config files
	return nil
}
