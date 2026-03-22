package checkout

import (
	"testing"

	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
)

func TestNewCmdCheckout(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "checkout PR",
			args:    []string{"123"},
			wantErr: false,
		},
		{
			name:    "checkout with branch name",
			args:    []string{"123", "--branch", "my-feature"},
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
			cmd := NewCmdCheckout(f, func(opts *CheckoutOptions) error {
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