package root

import (
	"bytes"
	"runtime"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestRootHelpMentionsWindowsPowerShellAliasForGC(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("PowerShell alias guidance is Windows-specific")
	}

	t.Setenv(commandNameEnv, "gc")

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

func TestRootCommandNameFromEnvironment(t *testing.T) {
	t.Setenv(commandNameEnv, "gitcode")

	cmd := NewRootCmd("dev", "none", "unknown", cmdutil.TestFactory())
	if cmd.Use != "gitcode" {
		t.Fatalf("NewRootCmd().Use = %q, want gitcode", cmd.Use)
	}

	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Usage:\n  gitcode [command]") {
		t.Fatalf("help should use gitcode command name: %s", output)
	}
	if strings.Contains(output, "Usage:\n  gc [command]") {
		t.Fatalf("help should not use gc usage when gitcode is selected: %s", output)
	}
}

func TestGitcodeCommandNamePropagatesToDiscoveryCommands(t *testing.T) {
	t.Setenv(commandNameEnv, "gitcode")

	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "version",
			args: []string{"version"},
			want: []string{"gitcode version dev"},
		},
		{
			name: "help json",
			args: []string{"help", "--json"},
			want: []string{`"path": "gitcode"`, `"path": "gitcode pr create"`},
		},
		{
			name: "schema",
			args: []string{"schema"},
			want: []string{`"name": "gitcode"`, `"path": "gitcode"`, `"path": "gitcode pr create"`},
		},
		{
			name: "powershell completion",
			args: []string{"completion", "powershell"},
			want: []string{"powershell completion for gitcode"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewRootCmd("dev", "none", "unknown", cmdutil.TestFactory())
			out := &bytes.Buffer{}
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("Execute() error = %v\n%s", err, out.String())
			}

			output := out.String()
			for _, want := range tt.want {
				if !strings.Contains(output, want) {
					t.Fatalf("output missing %q: %s", want, output)
				}
			}
		})
	}
}

func TestGitcodeHelpRewritesExamples(t *testing.T) {
	t.Setenv(commandNameEnv, "gitcode")

	cmd := NewRootCmd("dev", "none", "unknown", cmdutil.TestFactory())
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"help", "issue"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v\n%s", err, out.String())
	}

	output := out.String()
	if !strings.Contains(output, `$ gitcode issue create`) {
		t.Fatalf("help output missing rewritten gitcode example: %s", output)
	}
	if strings.Contains(output, `$ gc issue create`) {
		t.Fatalf("help output should not expose gc example when root is gitcode: %s", output)
	}
	if strings.Contains(output, `Use "gc issue`) {
		t.Fatalf("help output should not expose gc use hint when root is gitcode: %s", output)
	}
}
