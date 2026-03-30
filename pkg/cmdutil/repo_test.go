package cmdutil

import (
	"errors"
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
