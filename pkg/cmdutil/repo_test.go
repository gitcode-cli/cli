package cmdutil

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestResolveRepo(t *testing.T) {
	tests := []struct {
		name     string
		repo     string
		baseRepo func() (string, error)
		wantRepo string
		wantErr  string
	}{
		{
			name:     "uses explicit repo",
			repo:     "owner/repo",
			baseRepo: func() (string, error) { return "ignored/repo", nil },
			wantRepo: "owner/repo",
		},
		{
			name:     "falls back to current repo",
			baseRepo: func() (string, error) { return "detected/repo", nil },
			wantRepo: "detected/repo",
		},
		{
			name:     "missing repo and git context",
			baseRepo: func() (string, error) { return "", errors.New("not in a git repository") },
			wantErr:  "not in a git repository",
		},
		{
			name:    "missing repo and resolver",
			wantErr: "no repository specified. Use -R owner/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveRepo(tt.repo, tt.baseRepo)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("ResolveRepo() error = nil, want %q", tt.wantErr)
				}
				if got := ExitCode(err); got != ExitUsage {
					t.Fatalf("ExitCode() = %d, want %d", got, ExitUsage)
				}
				if err.Error() == tt.wantErr || strings.Contains(err.Error(), tt.wantErr) {
					return
				}
				t.Fatalf("ResolveRepo() error = %q, want containing %q", err.Error(), tt.wantErr)
			}
			if err != nil {
				t.Fatalf("ResolveRepo() unexpected error = %v", err)
			}
			if got != tt.wantRepo {
				t.Fatalf("ResolveRepo() = %q, want %q", got, tt.wantRepo)
			}
		})
	}
}

func TestParseRepo(t *testing.T) {
	tests := []struct {
		name      string
		repo      string
		wantOwner string
		wantRepo  string
		wantErr   string
	}{
		{
			name:      "owner repo",
			repo:      "owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
		},
		{
			name:      "https url",
			repo:      "https://gitcode.com/owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
		},
		{
			name:      "ssh url",
			repo:      "git@gitcode.com:owner/repo.git",
			wantOwner: "owner",
			wantRepo:  "repo",
		},
		{
			name:    "invalid repo",
			repo:    "owner/repo/extra",
			wantErr: "invalid repository format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := ParseRepo(tt.repo)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("ParseRepo() error = nil, want %q", tt.wantErr)
				}
				if got := ExitCode(err); got != ExitUsage {
					t.Fatalf("ExitCode() = %d, want %d", got, ExitUsage)
				}
				if err.Error() == tt.wantErr || strings.Contains(err.Error(), tt.wantErr) {
					return
				}
				t.Fatalf("ParseRepo() error = %q, want containing %q", err.Error(), tt.wantErr)
			}
			if err != nil {
				t.Fatalf("ParseRepo() unexpected error = %v", err)
			}
			if owner != tt.wantOwner || repo != tt.wantRepo {
				t.Fatalf("ParseRepo() = (%q, %q), want (%q, %q)", owner, repo, tt.wantOwner, tt.wantRepo)
			}
		})
	}
}

func TestParseRepoOutsideGitRepoRequiresExplicitRepo(t *testing.T) {
	dir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}

	_, _, err = ParseRepo("")
	if err == nil {
		t.Fatal("ParseRepo() error = nil, want missing repo error")
	}
	if !strings.Contains(err.Error(), "no repository specified. Use -R owner/repo") {
		t.Fatalf("ParseRepo() error = %q", err.Error())
	}
	if got := ExitCode(err); got != ExitUsage {
		t.Fatalf("ExitCode() = %d, want %d", got, ExitUsage)
	}
}

func TestResolvePRURL(t *testing.T) {
	tests := []struct {
		name    string
		htmlURL string
		owner   string
		repo    string
		number  int
		want    string
	}{
		{
			name:    "API URL provided",
			htmlURL: "https://gitcode.com/owner/repo/merge_requests/123",
			owner:   "owner",
			repo:    "repo",
			number:  123,
			want:    "https://gitcode.com/owner/repo/merge_requests/123",
		},
		{
			name:    "Empty URL fallback",
			htmlURL: "",
			owner:   "owner",
			repo:    "repo",
			number:  123,
			want:    "https://gitcode.com/owner/repo/merge_requests/123",
		},
		{
			name:    "Whitespace URL fallback",
			htmlURL: "   ",
			owner:   "owner",
			repo:    "repo",
			number:  123,
			want:    "https://gitcode.com/owner/repo/merge_requests/123",
		},
		{
			name:    "Zero number returns empty",
			htmlURL: "",
			owner:   "owner",
			repo:    "repo",
			number:  0,
			want:    "",
		},
		{
			name:    "Empty owner returns empty",
			htmlURL: "",
			owner:   "",
			repo:    "repo",
			number:  123,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolvePRURL(tt.htmlURL, tt.owner, tt.repo, tt.number)
			if got != tt.want {
				t.Errorf("ResolvePRURL(%q, %q, %q, %d) = %q, want %q",
					tt.htmlURL, tt.owner, tt.repo, tt.number, got, tt.want)
			}
		})
	}
}
