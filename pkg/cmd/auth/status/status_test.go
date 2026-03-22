package status

import (
	"testing"

	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
)

func TestNewCmdStatus(t *testing.T) {
	f := cmdutil.TestFactory()
	cmd := NewCmdStatus(f, nil)
	if cmd == nil {
		t.Fatal("NewCmdStatus returned nil")
	}
	if cmd.Use != "status" {
		t.Errorf("Expected Use 'status', got %q", cmd.Use)
	}
}