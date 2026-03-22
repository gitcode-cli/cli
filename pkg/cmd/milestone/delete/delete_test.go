package delete

import (
	"testing"

	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
)

func TestNewCmdDelete(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "delete milestone",
			args:    []string{"1"},
			wantErr: false,
		},
		{
			name:    "delete with yes flag",
			args:    []string{"1", "--yes"},
			wantErr: false,
		},
		{
			name:    "no milestone number",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "invalid milestone number",
			args:    []string{"abc"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdDelete(f, func(opts *DeleteOptions) error {
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