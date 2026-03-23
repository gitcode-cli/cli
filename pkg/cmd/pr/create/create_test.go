package create

import (
	"testing"

	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdCreate(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "create PR with title and head",
			args:    []string{"--title", "Feature", "--head", "feature-branch", "--repo", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "create draft PR",
			args:    []string{"--title", "WIP", "--head", "draft", "--draft", "--repo", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "create with base",
			args:    []string{"--title", "Feature", "--head", "feature", "--base", "develop", "--repo", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "create cross-repo PR with fork",
			args:    []string{"--title", "Feature", "--head", "feature", "--fork", "myfork/repo", "--repo", "upstream/repo"},
			wantErr: false,
		},
		{
			name:    "missing title",
			args:    []string{"--head", "feature", "--repo", "owner/repo"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			// For validation tests, don't provide runF so actual validation runs
			var cmd *cobra.Command
			if tt.name == "missing title" {
				cmd = NewCmdCreate(f, nil)
			} else {
				cmd = NewCmdCreate(f, func(opts *CreateOptions) error {
					return nil
				})
			}
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}