package merge

import (
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdMerge(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "merge PR",
			args:    []string{"123"},
			wantErr: false,
		},
		{
			name:    "merge with squash",
			args:    []string{"123", "--method", "squash"},
			wantErr: false,
		},
		{
			name:    "merge with rebase",
			args:    []string{"123", "--method", "rebase"},
			wantErr: false,
		},
		{
			name:    "merge with yes",
			args:    []string{"123", "--yes"},
			wantErr: false,
		},
		{
			name:    "no PR number",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdMerge(f, func(opts *MergeOptions) error {
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
