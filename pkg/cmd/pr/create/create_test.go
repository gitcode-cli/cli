package create

import (
	"testing"

	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
)

func TestNewCmdCreate(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "create PR with title",
			args:    []string{"--title", "Feature", "--head", "feature-branch"},
			wantErr: false,
		},
		{
			name:    "create draft PR",
			args:    []string{"--title", "WIP", "--head", "draft", "--draft"},
			wantErr: false,
		},
		{
			name:    "create with base",
			args:    []string{"--title", "Feature", "--head", "feature", "--base", "develop"},
			wantErr: false,
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