// Package root implements the root command for gc
package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/gitcode-com/gitcode-cli/pkg/cmd/auth"
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
  • Merge Request management (mr create, mr review)`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(version.NewCmdVersion(ver, commit, date))
	cmd.AddCommand(auth.NewCmdAuth(f))

	return cmd
}