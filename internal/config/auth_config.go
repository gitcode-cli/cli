package config

import "os"

// authConfig implements AuthConfig interface
type authConfig struct {
	config *config
}

// ActiveToken returns the active token for a host
func (a *authConfig) ActiveToken(hostname string) (string, string) {
	// Check environment variables first
	if token := getEnvToken(); token != "" {
		return token, "environment"
	}
	// TODO: check keyring and config file
	return "", ""
}

// HasActiveToken checks if a token exists for a host
func (a *authConfig) HasActiveToken(hostname string) bool {
	token, _ := a.ActiveToken(hostname)
	return token != ""
}

// ActiveUser returns the active user for a host
func (a *authConfig) ActiveUser(hostname string) (string, error) {
	// TODO: implement
	return "", nil
}

// Hosts returns all configured hosts
func (a *authConfig) Hosts() []string {
	// TODO: read from config
	return []string{"gitcode.com"}
}

// DefaultHost returns the default host
func (a *authConfig) DefaultHost() (string, string) {
	if host := getEnvHost(); host != "" {
		return host, "environment"
	}
	return "gitcode.com", "default"
}

// Login creates or updates authentication for a host
func (a *authConfig) Login(hostname, username, token, gitProtocol string, secureStorage bool) (bool, error) {
	// TODO: implement
	return false, nil
}

// Logout removes authentication for a host
func (a *authConfig) Logout(hostname, username string) error {
	// TODO: implement
	return nil
}

// SwitchUser switches the active user for a host
func (a *authConfig) SwitchUser(hostname, user string) error {
	// TODO: implement
	return nil
}

// getEnvToken retrieves token from environment variables
func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITCODE_TOKEN")
}

// getEnvHost retrieves host from environment variable
func getEnvHost() string {
	return os.Getenv("GC_HOST")
}