package review

import (
	"testing"

	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
)

func TestNewCmdReview(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "approve PR",
			args:    []string{"123", "--approve"},
			wantErr: false,
		},
		{
			name:    "request changes",
			args:    []string{"123", "--request"},
			wantErr: false,
		},
		{
			name:    "review with comment",
			args:    []string{"123", "--comment", "Looks good"},
			wantErr: false,
		},
		{
			name:    "approve with body",
			args:    []string{"123", "--approve", "--comment", "LGTM"},
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
			cmd := NewCmdReview(f, func(opts *ReviewOptions) error {
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