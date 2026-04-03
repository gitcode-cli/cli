package cmdutil

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
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
			wantErr:  "could not determine current repository",
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
}

func TestParseRepoFallsBackToCurrentRepo(t *testing.T) {
	dir := t.TempDir()
	run := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %v (%s)", args, err, string(out))
		}
	}

	run("init")
	run("remote", "add", "origin", "https://gitcode.com/fallback/repo.git")

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()
	if err := os.Chdir(filepath.Clean(dir)); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}

	owner, repo, err := ParseRepo("")
	if err != nil {
		t.Fatalf("ParseRepo() unexpected error = %v", err)
	}
	if owner != "fallback" || repo != "repo" {
		t.Fatalf("ParseRepo() = (%q, %q), want (%q, %q)", owner, repo, "fallback", "repo")
	}
}
