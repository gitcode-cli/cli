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
			name:    "create milestone",
			args:    []string{"v1.0"},
			wantErr: false,
		},
		{
			name:    "create with description",
			args:    []string{"v1.0", "--description", "First release"},
			wantErr: false,
		},
		{
			name:    "create with due date",
			args:    []string{"v2.0", "--due-date", "2024-12-31"},
			wantErr: false,
		},
		{
			name:    "no title",
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