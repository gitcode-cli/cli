package cmdutil

import (
	"os"

	"gitcode.com/gitcode-cli/cli/internal/config"
)

// EnvToken returns the active token from supported environment variables.
func EnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	if token := os.Getenv("GITCODE_TOKEN"); token != "" {
		return token
	}

	cfg := config.New()
	host, _ := cfg.Authentication().DefaultHost()
	token, _ := cfg.Authentication().ActiveToken(host)
	return token
}
