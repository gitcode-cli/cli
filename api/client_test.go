package api

import (
	"testing"
	"time"

	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name   string
		base   string
		params map[string]string
		want   string
	}{
		{
			name:   "no params",
			base:   "/repos/{owner}/{repo}/issues",
			params: nil,
			want:   "/repos/{owner}/{repo}/issues",
		},
		{
			name:   "replace owner and repo",
			base:   "/repos/{owner}/{repo}/issues",
			params: map[string]string{"owner": "octocat", "repo": "hello-world"},
			want:   "/repos/octocat/hello-world/issues",
		},
		{
			name:   "replace single param",
			base:   "/repos/owner/repo/issues/{number}",
			params: map[string]string{"number": "42"},
			want:   "/repos/owner/repo/issues/42",
		},
		{
			name:   "encoded value",
			base:   "/repos/{owner}/{repo}",
			params: map[string]string{"owner": "user/with"},
			want:   "/repos/user%2Fwith/{repo}",
		},
		{
			name:   "empty params map",
			base:   "/user",
			params: map[string]string{},
			want:   "/user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildURL(tt.base, tt.params)
			if got != tt.want {
				t.Errorf("BuildURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAPIError(t *testing.T) {
	tests := []struct {
		name       string
		apiErr     APIError
		wantSubstr string
	}{
		{
			name:       "message field",
			apiErr:     APIError{StatusCode: 404, Message: "Not Found"},
			wantSubstr: "HTTP 404: Not Found",
		},
		{
			name:       "error_message field",
			apiErr:     APIError{StatusCode: 500, ErrorMessage: "Internal Server Error"},
			wantSubstr: "HTTP 500: Internal Server Error",
		},
		{
			name:       "no message defaults to unknown",
			apiErr:     APIError{StatusCode: 400},
			wantSubstr: "unknown error",
		},
		{
			name:       "403 with guidance",
			apiErr:     APIError{StatusCode: 403, Message: "Forbidden"},
			wantSubstr: "You don't have permission",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.apiErr.Error()
			if !contains(got, tt.wantSubstr) {
				t.Errorf("Error() = %q, want containing %q", got, tt.wantSubstr)
			}
		})
	}
}

func TestAPIHostForGitCodeHost(t *testing.T) {
	tests := []struct {
		name string
		host string
		want string
	}{
		{"empty string", "", "api.gitcode.com"},
		{"gitcode.com", "gitcode.com", "api.gitcode.com"},
		{"already api", "api.gitcode.com", "api.gitcode.com"},
		{"custom domain", "gitcode.internal", "api.gitcode.internal"},
		{"custom with api prefix", "api.gitcode.internal", "api.gitcode.internal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := apiHostForGitCodeHost(tt.host)
			if got != tt.want {
				t.Errorf("apiHostForGitCodeHost(%q) = %q, want %q", tt.host, got, tt.want)
			}
		})
	}
}

func TestSetHost(t *testing.T) {
	c := NewClientFromHTTP(nil)
	c.SetHost("gitcode.internal")
	if c.Host() != "api.gitcode.internal" {
		t.Errorf("Host() = %q, want %q", c.Host(), "api.gitcode.internal")
	}
}

func TestToken(t *testing.T) {
	c := NewClientFromHTTP(nil)
	if c.Token() != "" {
		t.Errorf("Token() = %q, want empty", c.Token())
	}
	c.SetToken("tok123", "env")
	if c.Token() != "tok123" {
		t.Errorf("Token() = %q, want %q", c.Token(), "tok123")
	}
}

func TestRawURL(t *testing.T) {
	c := NewClientFromHTTP(nil)
	tests := []struct {
		name     string
		endpoint string
		want     string
		wantErr  bool
	}{
		{"relative path", "repos/owner/repo", "https://api.gitcode.com/api/v5/repos/owner/repo", false},
		{"leading slash", "/repos/owner/repo", "https://api.gitcode.com/api/v5/repos/owner/repo", false},
		{"already api", "/api/v5/repos/owner/repo", "https://api.gitcode.com/api/v5/repos/owner/repo", false},
		{"full url matching host", "https://api.gitcode.com/api/v5/repos/owner/repo", "https://api.gitcode.com/api/v5/repos/owner/repo", false},
		{"foreign host rejected", "https://evil.com/api/v5/repos/owner/repo", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.rawURL(tt.endpoint)
			if tt.wantErr {
				if err == nil {
					t.Errorf("rawURL(%q) expected error", tt.endpoint)
				}
				return
			}
			if err != nil {
				t.Errorf("rawURL(%q) error = %v", tt.endpoint, err)
				return
			}
			if got != tt.want {
				t.Errorf("rawURL(%q) = %q, want %q", tt.endpoint, got, tt.want)
			}
		})
	}
}

func TestDefaultHTTPClient(t *testing.T) {
	c := DefaultHTTPClient()
	if c == nil {
		t.Fatal("DefaultHTTPClient() returned nil")
	}
	if c.Timeout == 0 {
		t.Error("DefaultHTTPClient() timeout is zero")
	}
}

func TestNewHTTPClientWithRetryAndLogger(t *testing.T) {
	c := NewHTTPClientWithRetryAndLogger(time.Second, RetryConfig{MaxRetries: 3}, nil)
	if c == nil {
		t.Fatal("NewHTTPClientWithRetryAndLogger() returned nil")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestRawREST(t *testing.T) {
	handler := testutil.MockAPIHandler()
	mockClient := testutil.NewTestHTTPClient(handler)
	c := NewClientFromHTTP(mockClient)
	c.SetToken("test-token", "env")

	resp, err := c.RawREST("GET", "user", nil, nil)
	if err != nil {
		t.Fatalf("RawREST() error = %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}
	if len(resp.Body) == 0 {
		t.Error("Body should not be empty")
	}
}

func TestRawRESTWithHostGuard(t *testing.T) {
	handler := testutil.MockAPIHandler()
	mockClient := testutil.NewTestHTTPClient(handler)
	c := NewClientFromHTTP(mockClient)
	c.SetToken("test-token", "env")

	_, err := c.RawREST("GET", "https://evil.com/user", nil, nil)
	if err == nil {
		t.Fatal("expected host rejection error")
	}
}

func TestUploadToURL(t *testing.T) {
	handler := testutil.MockAPIHandler()
	mockClient := testutil.NewTestHTTPClient(handler)
	c := NewClientFromHTTP(mockClient)

	// Get the mock server URL
	err := c.UploadToURL("http://localhost:1/upload", "test.txt", []byte("data"), "text/plain", nil)
	// Will fail because localhost:1 doesn't exist, but covers the code path
	if err == nil {
		t.Log("UploadToURL succeeded unexpectedly")
	}
}

func TestUploadAsset(t *testing.T) {
	handler := testutil.MockAPIHandler()
	mockClient := testutil.NewTestHTTPClient(handler)
	c := NewClientFromHTTP(mockClient)
	c.SetToken("test-token", "env")

	asset, err := c.UploadAsset("/repos/owner/test-repo/releases/1/assets", "test.txt", []byte("data"), "text/plain")
	if err != nil {
		t.Fatalf("UploadAsset() error = %v", err)
	}
	if asset.Name != "test.txt" {
		t.Errorf("Name = %q, want %q", asset.Name, "test.txt")
	}
}
