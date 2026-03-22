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
			name:    "create with tag",
			args:    []string{"v1.0.0"},
			wantErr: false,
		},
		{
			name:    "create with tag and title",
			args:    []string{"v1.0.0", "--title", "Version 1.0"},
			wantErr: false,
		},
		{
			name:    "create with draft flag",
			args:    []string{"v1.0.0", "--draft"},
			wantErr: false,
		},
		{
			name:    "create with prerelease flag",
			args:    []string{"v1.0.0-beta", "--prerelease"},
			wantErr: false,
		},
		{
			name:    "no tag specified",
			args:    []string{},
			wantErr: true,
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