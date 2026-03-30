package schema

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestSchemaCommand(t *testing.T) {
	root := &cobra.Command{Use: "gc", Short: "root"}
	issue := &cobra.Command{Use: "issue", Short: "issue root"}
	view := &cobra.Command{Use: "view <number>", Short: "view issue"}
	view.Flags().Bool("json", false, "Output as JSON")
	issue.AddCommand(view)
	root.AddCommand(issue)

	cmd := NewCmdSchema(root)
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"issue", "view"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "\"path\": \"gc issue view\"") {
		t.Fatalf("schema output missing command path: %s", output)
	}
	if !strings.Contains(output, "\"name\": \"json\"") {
		t.Fatalf("schema output missing json flag: %s", output)
	}
}
