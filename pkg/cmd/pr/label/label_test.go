package label

import (
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdLabel(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "add labels",
			args:    []string{"123", "--add", "bug,enhancement"},
			wantErr: false,
		},
		{
			name:    "remove label",
			args:    []string{"123", "--remove", "bug"},
			wantErr: false,
		},
		{
			name:    "list labels",
			args:    []string{"123", "--list"},
			wantErr: false,
		},
		{
			name:    "no issue number",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "invalid issue number",
			args:    []string{"abc"},
			wantErr: true,
		},
		{
			name:    "no action specified",
			args:    []string{"123"},
			wantErr: false, // Command runs, error in run function
		},
		{
			name:    "add with repo",
			args:    []string{"123", "--add", "bug", "-R", "owner/repo"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdLabel(f, func(opts *LabelOptions) error {
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
