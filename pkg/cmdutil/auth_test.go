package cmdutil

import (
	"gitcode.com/gitcode-cli/cli/pkg/testutil"
	"io"
	"net/http"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/pkg/config"
)

func TestAuthenticatedClientRejectsEnvTokenForCustomHost(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_HOST", "enterprise.example.com")
	t.Setenv("GC_TOKEN", "env-token")
	t.Setenv("GITCODE_TOKEN", "")

	_, err := AuthenticatedClient(&http.Client{
		Transport: testutil.NewRoundTripFunc(func(req *http.Request) (*http.Response, error) {
			t.Fatalf("unexpected request to %s", req.URL.String())
			return nil, nil
		}),
	})
	if err == nil {
		t.Fatal("AuthenticatedClient() error = nil, want auth error")
	}
	if !strings.Contains(err.Error(), "not authenticated") {
		t.Fatalf("AuthenticatedClient() error = %q, want not authenticated", err.Error())
	}
}

func TestAuthenticatedClientUsesStoredTokenForCustomHost(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_HOST", "enterprise.example.com")
	t.Setenv("GC_TOKEN", "env-token")
	t.Setenv("GITCODE_TOKEN", "")

	cfg := config.New()
	if _, err := cfg.Authentication().Login("enterprise.example.com", "tester", "stored-token", "https", false); err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	client, err := AuthenticatedClient(&http.Client{
		Transport: testutil.NewRoundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Host != "api.enterprise.example.com" {
				t.Fatalf("request host = %q, want api.enterprise.example.com", req.URL.Host)
			}
			if got := req.Header.Get("Authorization"); got != "Bearer stored-token" {
				t.Fatalf("Authorization = %q, want stored token", got)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     http.StatusText(http.StatusOK),
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{}`)),
			}, nil
		}),
	})
	if err != nil {
		t.Fatalf("AuthenticatedClient() error = %v", err)
	}
	if err := client.Get("/user", &struct{}{}); err != nil {
		t.Fatalf("Get() error = %v", err)
	}
}

func TestAuthenticatedClientRejectsMalformedHost(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	t.Setenv("GC_HOST", "https://gitcode.com")
	t.Setenv("GC_TOKEN", "env-token")

	_, err := AuthenticatedClient(&http.Client{})
	if err == nil {
		t.Fatal("AuthenticatedClient() error = nil, want invalid host error")
	}
	if !strings.Contains(err.Error(), "invalid host") {
		t.Fatalf("AuthenticatedClient() error = %q, want invalid host", err.Error())
	}
}
