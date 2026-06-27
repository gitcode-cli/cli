package cmdutil

import (
	"fmt"
	"net/http"
	"os"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/config"
)

// DefaultToken returns the active token, checking environment variables first,
// then falling back to the configured token file.
func DefaultToken() string {
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

// AuthenticatedClient creates an API client using the configured default host
// and active token. Environment tokens keep their existing priority.
func AuthenticatedClient(httpClient *http.Client) (*api.Client, error) {
	cfg := config.New()
	authCfg := cfg.Authentication()
	host, _ := authCfg.DefaultHost()
	host, err := config.NormalizeTrustedHost(host)
	if err != nil {
		return nil, err
	}
	token, source := authCfg.ActiveToken(host)
	if host != "gitcode.com" {
		token, source = authCfg.StoredToken(host)
	}
	if token == "" {
		return nil, NewAuthError("not authenticated. Run: gc auth login")
	}
	client := api.NewClientFromHTTP(httpClient)
	client.SetHost(host)
	client.SetToken(token, source)
	return client, nil
}

// AuthenticatedClientFromFactory creates an authenticated API client from the
// command factory dependencies.
func AuthenticatedClientFromFactory(httpClient func() (*http.Client, error)) (*api.Client, error) {
	if httpClient == nil {
		return nil, fmt.Errorf("missing HTTP client factory")
	}
	client, err := httpClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}
	return AuthenticatedClient(client)
}
