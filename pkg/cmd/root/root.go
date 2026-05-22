// Package root implements the root command for gc.
package root

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/auth"
	commitcmd "gitcode.com/gitcode-cli/cli/pkg/cmd/commit"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/help"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/issue"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/label"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/milestone"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/pr"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/release"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/repo"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/schema"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/version"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

const commandNameEnv = "GITCODE_CLI_COMMAND_NAME"

// Execute runs the root command.
func Execute(ver, commit, date string) error {
	f := cmdutil.NewFactory()
	rootCmd := NewRootCmd(ver, commit, date, f)
	return rootCmd.Execute()
}

// NewRootCmd creates the root command.
func NewRootCmd(ver, commit, date string, f *cmdutil.Factory) *cobra.Command {
	commandName := resolveCommandName()

	cmd := &cobra.Command{
		Use:           commandName,
		Short:         "GitCode CLI - Command line tool for GitCode",
		Long:          rootLong(commandName),
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands.
	cmd.AddCommand(version.NewCmdVersion(ver, commit, date, commandName))
	cmd.AddCommand(auth.NewCmdAuth(f))
	cmd.AddCommand(repo.NewCmdRepo(f))
	cmd.AddCommand(issue.NewCmdIssue(f))
	cmd.AddCommand(pr.NewCmdPR(f))
	cmd.AddCommand(commitcmd.NewCmdCommit(f))
	cmd.AddCommand(label.NewCmdLabel(f))
	cmd.AddCommand(milestone.NewCmdMilestone(f))
	cmd.AddCommand(release.NewCmdRelease(f))
	cmd.AddCommand(schema.NewCmdSchema(cmd))

	rewriteExamples(cmd, commandName)

	// Set custom help command with search and discovery features.
	cmd.SetHelpCommand(help.NewCmdHelp(cmd))

	return cmd
}

func resolveCommandName() string {
	if name := normalizedCommandName(os.Getenv(commandNameEnv)); name != "" {
		return name
	}

	executable := strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))
	if name := normalizedCommandName(executable); name != "" {
		return name
	}

	if runtime.GOOS == "windows" {
		return "gitcode"
	}
	return "gc"
}

func normalizedCommandName(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "gc":
		return "gc"
	case "gitcode":
		return "gitcode"
	default:
		return ""
	}
}

func rootLong(commandName string) string {
	long := fmt.Sprintf(`%s is a command line tool for GitCode.

It provides convenient access to GitCode features including:
  - Authentication management (auth login, auth status)
  - Repository operations (repo clone, repo create)
  - Issue management (issue create, issue list)
  - Pull Request management (pr create, pr review)`, commandName)

	if runtime.GOOS == "windows" && commandName == "gc" {
		long += `

Windows PowerShell defines "gc" as an alias for Get-Content.
If "gc" is intercepted by PowerShell, use "gitcode", "gc.exe", or
"python -m gc_cli" instead.`
	}

	return long
}

func rewriteExamples(cmd *cobra.Command, rootName string) {
	cmd.Example = rewriteCommandReferences(cmd.Example, rootName)
	for _, child := range cmd.Commands() {
		rewriteExamples(child, rootName)
	}
}

func rewriteCommandReferences(text, rootName string) string {
	if text == "" || rootName == "" || rootName == "gc" {
		return text
	}

	replacer := strings.NewReplacer(
		"$ gc ", "$ "+rootName+" ",
		"| gc ", "| "+rootName+" ",
		"Use \"gc ", "Use \""+rootName+" ",
	)
	return replacer.Replace(text)
}
