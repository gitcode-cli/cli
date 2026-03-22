// Package root implements the root command for gc
package root

import (
	"github.com/spf13/cobra"

	"github.com/gitcode-com/gitcode-cli/pkg/cmd/auth"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/issue"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/label"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/milestone"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/pr"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/release"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/repo"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/version"
	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
)

// Execute runs the root command
func Execute(ver, commit, date string) error {
	f := cmdutil.NewFactory()
	rootCmd := NewRootCmd(ver, commit, date, f)
	return rootCmd.Execute()
}

// NewRootCmd creates the root command
func NewRootCmd(ver, commit, date string, f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gc",
		Short: "GitCode CLI - Command line tool for GitCode",
		Long: `gc is a command line tool for GitCode.

It provides convenient access to GitCode features including:
  • Authentication management (auth login, auth status)
  • Repository operations (repo clone, repo create)
  • Issue management (issue create, issue list)
  • Pull Request management (pr create, pr review)`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(version.NewCmdVersion(ver, commit, date))
	cmd.AddCommand(auth.NewCmdAuth(f))
	cmd.AddCommand(repo.NewCmdRepo(f))
	cmd.AddCommand(issue.NewCmdIssue(f))
	cmd.AddCommand(pr.NewCmdPR(f))
	cmd.AddCommand(label.NewCmdLabel(f))
	cmd.AddCommand(milestone.NewCmdMilestone(f))
	cmd.AddCommand(release.NewCmdRelease(f))

	return cmd
}