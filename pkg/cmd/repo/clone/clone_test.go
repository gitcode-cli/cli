package clone

import (
	"testing"

	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
)

func TestNewCmdClone(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "clone with owner/repo format",
			args:    []string{"owner/repo"},
			wantErr: false,
		},
		{
			name:    "clone with full URL",
			args:    []string{"https://gitcode.com/owner/repo.git"},
			wantErr: false,
		},
		{
			name:    "clone with SSH URL",
			args:    []string{"git@gitcode.com:owner/repo.git"},
			wantErr: false,
		},
		{
			name:    "no arguments",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdClone(f, func(opts *CloneOptions) error {
				return nil
			})
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseRepoURL(t *testing.T) {
	tests := []struct {
		name      string
		repo      string
		protocol  string
		wantURL   string
		wantErr   bool
	}{
		{
			name:     "https URL",
			repo:     "https://gitcode.com/owner/repo.git",
			protocol: "https",
			wantURL:  "https://gitcode.com/owner/repo.git",
			wantErr:  false,
		},
		{
			name:     "ssh URL",
			repo:     "git@gitcode.com:owner/repo.git",
			protocol: "ssh",
			wantURL:  "git@gitcode.com:owner/repo.git",
			wantErr:  false,
		},
		{
			name:     "owner/repo format with https",
			repo:     "owner/repo",
			protocol: "https",
			wantURL:  "https://gitcode.com/owner/repo.git",
			wantErr:  false,
		},
		{
			name:     "owner/repo format with ssh",
			repo:     "owner/repo",
			protocol: "ssh",
			wantURL:  "git@gitcode.com:owner/repo.git",
			wantErr:  false,
		},
		{
			name:     "invalid format",
			repo:     "invalid",
			protocol: "https",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, err := parseRepoURL(tt.repo, tt.protocol)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRepoURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotURL != tt.wantURL {
				t.Errorf("parseRepoURL() = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}