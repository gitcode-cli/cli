// Package root implements the root command for gc
package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/gitcode-com/gitcode-cli/pkg/cmd/version"
)

// Execute runs the root command
func Execute(ver, commit, date string) error {
	rootCmd := NewRootCmd(ver, commit, date)
	return rootCmd.Execute()
}

// NewRootCmd creates the root command
func NewRootCmd(ver, commit, date string) *cobra.Command {
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

	return cmd
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}