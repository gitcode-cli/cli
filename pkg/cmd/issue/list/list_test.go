package list

import (
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdList(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "list default",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "list with state",
			args:    []string{"--state", "closed"},
			wantErr: false,
		},
		{
			name:    "list with limit",
			args:    []string{"--limit", "10"},
			wantErr: false,
		},
		{
			name:    "list with labels",
			args:    []string{"--label", "bug"},
			wantErr: false,
		},
		{
			name:    "list with milestone",
			args:    []string{"--milestone", "v1.0"},
			wantErr: false,
		},
		{
			name:    "list with assignee",
			args:    []string{"--assignee", "username"},
			wantErr: false,
		},
		{
			name:    "list with creator",
			args:    []string{"--creator", "username"},
			wantErr: false,
		},
		{
			name:    "list with sort",
			args:    []string{"--sort", "updated"},
			wantErr: false,
		},
		{
			name:    "list with direction",
			args:    []string{"--direction", "asc"},
			wantErr: false,
		},
		{
			name:    "list with search",
			args:    []string{"--search", "bug"},
			wantErr: false,
		},
		{
			name:    "list with created-after",
			args:    []string{"--created-after", "2024-01-01"},
			wantErr: false,
		},
		{
			name:    "list with updated-after",
			args:    []string{"--updated-after", "2024-01-01"},
			wantErr: false,
		},
		{
			name:    "list with combined filters",
			args:    []string{"--state", "open", "--sort", "updated", "--direction", "desc"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdList(f, func(opts *ListOptions) error {
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