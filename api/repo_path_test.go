package api

import (
	"net/url"
	"testing"
)

func TestEscapedRepoPath(t *testing.T) {
	tests := []struct {
		name  string
		owner string
		repo  string
	}{
		{"normal", "owner", "repo"},
		{"slash in owner", "a/b", "repo"},
		{"question in repo", "owner", "c?d"},
		{"hash in repo", "owner", "e#f"},
		{"percent in owner", "o%ner", "repo"},
		{"empty owner", "", "repo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapedRepoPath(tt.owner, tt.repo)
			want := "/repos/" + url.PathEscape(tt.owner) + "/" + url.PathEscape(tt.repo)
			if got != want {
				t.Errorf("escapedRepoPath(%q, %q) = %q, want %q", tt.owner, tt.repo, got, want)
			}
		})
	}
}

// TestEscapedRepoPath_PreventsPathInjection verifies that crafted owner/repo
// with path separators cannot manipulate the URL path structure (#316).
func TestEscapedRepoPath_PreventsPathInjection(t *testing.T) {
	// owner containing "/" should NOT create an extra path segment
	got := escapedRepoPath("evil/owner", "repo")
	if got == "/repos/evil/owner/repo" {
		t.Fatalf("path injection: owner '/' not escaped: %q", got)
	}
	// repo containing "?" should NOT start a query string
	got = escapedRepoPath("owner", "evil?repo")
	if got == "/repos/owner/evil?repo" {
		t.Fatalf("path injection: repo '?' not escaped: %q", got)
	}
	// repo containing "#" should NOT start a fragment
	got = escapedRepoPath("owner", "evil#repo")
	if got == "/repos/owner/evil#repo" {
		t.Fatalf("path injection: repo '#' not escaped: %q", got)
	}
}
