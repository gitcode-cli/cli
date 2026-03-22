package diff

import (
	"testing"

	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
)

func TestNewCmdDiff(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "diff PR",
			args:    []string{"123"},
			wantErr: false,
		},
		{
			name:    "no PR number",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "invalid PR number",
			args:    []string{"abc"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdDiff(f, func(opts *DiffOptions) error {
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