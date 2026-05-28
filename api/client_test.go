package api

import (
	"testing"
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

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
