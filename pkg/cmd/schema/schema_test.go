package schema

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
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

func TestSchemaCommandIncludesEnumMetadata(t *testing.T) {
	root := &cobra.Command{Use: "gc", Short: "root"}
	issue := &cobra.Command{Use: "issue", Short: "issue root"}
	list := &cobra.Command{Use: "list", Short: "list issues"}
	cmdutil.AddFormatFlag(list, new(string))
	cmdutil.AddTimeFormatFlag(list, new(string))
	list.Flags().String("state", "open", "Filter by state (open/closed/all)")
	cmdutil.SetFlagEnum(list, "state", "open", "closed", "all")
	issue.AddCommand(list)
	root.AddCommand(issue)

	cmd := NewCmdSchema(root)
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"issue", "list"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"name": "format"`) || !strings.Contains(output, `"enum": [`) {
		t.Fatalf("schema output missing enum payload: %s", output)
	}
	for _, want := range []string{`"json"`, `"simple"`, `"table"`, `"absolute"`, `"relative"`, `"open"`, `"closed"`, `"all"`} {
		if !strings.Contains(output, want) {
			t.Fatalf("schema output missing enum value %s: %s", want, output)
		}
	}
}

func TestSchemaCommandRewritesExamplesForGitcodeRoot(t *testing.T) {
	root := &cobra.Command{Use: "gitcode", Short: "root"}
	issue := &cobra.Command{
		Use:     "issue",
		Short:   "issue root",
		Example: "# View an issue\n$ gc issue view 1 -R owner/repo\n$ echo body | gc issue comment 1 --body-file -",
	}
	root.AddCommand(issue)

	cmd := NewCmdSchema(root)
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"issue"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `$ gitcode issue view 1 -R owner/repo`) {
		t.Fatalf("schema output missing rewritten gitcode example: %s", output)
	}
	if !strings.Contains(output, `| gitcode issue comment 1 --body-file -`) {
		t.Fatalf("schema output missing rewritten piped gitcode example: %s", output)
	}
	if strings.Contains(output, `$ gc issue view`) {
		t.Fatalf("schema output should not expose gc example when root is gitcode: %s", output)
	}
	if strings.Contains(output, `| gc issue comment`) {
		t.Fatalf("schema output should not expose piped gc example when root is gitcode: %s", output)
	}
}
