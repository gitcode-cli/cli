package resolve

import (
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdResolve(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "resolve with PR number and discussion ID",
			args:    []string{"123", "d1", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "resolve without repo flag",
			args:    []string{"123", "d1"},
			wantErr: false,
		},
		{
			name:    "missing PR number",
			args:    []string{"123"},
			wantErr: true,
		},
		{
			name:    "missing all args",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "invalid PR number",
			args:    []string{"abc", "d1"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdResolve(f, func(opts *resolveOptions) error {
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

func TestNewCmdUnresolve(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "unresolve with PR number and discussion ID",
			args:    []string{"123", "d1", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "missing discussion ID",
			args:    []string{"123"},
			wantErr: true,
		},
		{
			name:    "invalid PR number",
			args:    []string{"abc", "d1"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdUnresolve(f, func(opts *resolveOptions) error {
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
