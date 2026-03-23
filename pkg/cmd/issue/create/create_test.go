package create

import (
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdCreate(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "create with title",
			args:    []string{"--title", "Test Issue"},
			wantErr: false,
		},
		{
			name:    "create with title and body",
			args:    []string{"--title", "Test", "--body", "Description"},
			wantErr: false,
		},
		{
			name:    "create with labels",
			args:    []string{"--title", "Test", "--label", "bug,enhancement"},
			wantErr: false,
		},
		{
			name:    "no title",
			args:    []string{},
			wantErr: false, // Command runs, error in run function
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdCreate(f, func(opts *CreateOptions) error {
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