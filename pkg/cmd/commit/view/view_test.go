package view

import (
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdView(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "with sha and repo",
			args:    []string{"abc123", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "with show-diff flag",
			args:    []string{"abc123", "-R", "owner/repo", "--show-diff"},
			wantErr: false,
		},
		{
			name:    "with json flag",
			args:    []string{"abc123", "-R", "owner/repo", "--json"},
			wantErr: false,
		},
		{
			name:    "missing sha",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdView(f, func(opts *ViewOptions) error {
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