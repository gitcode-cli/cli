package list

import (
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdList(t *testing.T) {
	f := cmdutil.TestFactory()
	cmd := NewCmdList(f, func(opts *ListOptions) error {
		return nil
	})

	if cmd == nil {
		t.Fatal("NewCmdList returned nil")
	}
	if cmd.Use != "list" {
		t.Errorf("Expected Use 'list', got %q", cmd.Use)
	}
}