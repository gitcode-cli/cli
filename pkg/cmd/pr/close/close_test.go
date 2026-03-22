package close

import (
	"testing"

	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
)

func TestNewCmdClose(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "close PR",
			args:    []string{"123"},
			wantErr: false,
		},
		{
			name:    "close with comment",
			args:    []string{"123", "--comment", "Not needed"},
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
			cmd := NewCmdClose(f, func(opts *CloseOptions) error {
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