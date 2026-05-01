package root

import (
	"bytes"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestRootHelpMentionsWindowsPowerShellAlias(t *testing.T) {
	cmd := NewRootCmd("dev", "none", "unknown", cmdutil.TestFactory())
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	help := out.String()
	for _, want := range []string{"Windows PowerShell", "Get-Content", "gitcode", "gc.exe", "python -m gc_cli"} {
		if !strings.Contains(help, want) {
			t.Fatalf("help missing %q: %s", want, help)
		}
	}
}
