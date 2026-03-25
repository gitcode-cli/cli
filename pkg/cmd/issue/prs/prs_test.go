package prs

import (
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdPrs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid issue number",
			args:    []string{"123", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "no issue number",
			args:    []string{"-R", "owner/repo"},
			wantErr: true,
		},
		{
			name:    "invalid issue number",
			args:    []string{"abc", "-R", "owner/repo"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdPrs(f, func(opts *PrsOptions) error {
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