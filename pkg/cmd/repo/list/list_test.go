package list

import (
	"testing"

	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
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

func TestListOptions(t *testing.T) {
	opts := &ListOptions{
		Limit:      30,
		Visibility: "public",
	}

	if opts.Limit != 30 {
		t.Errorf("Expected Limit 30, got %d", opts.Limit)
	}
	if opts.Visibility != "public" {
		t.Errorf("Expected Visibility 'public', got %q", opts.Visibility)
	}
}