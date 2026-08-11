// Package search implements the search command.
package search

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/search/issues"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/search/repos"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/search/users"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdSearch creates the search command.
func NewCmdSearch(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search <command>",
		Short: "Search for repositories, issues, and users",
		Long: heredoc.Doc(`
			Search GitCode for repositories, issues, and users.

			Use --json for machine-readable output.
		`),
		Example: heredoc.Doc(`
			# Search repositories
			$ gc search repos "gitcode"

			# Search issues
			$ gc search issues "bug report"

			# Search users
			$ gc search users "developer"
		`),
	}

	cmd.AddCommand(repos.NewCmdRepos(f, nil))
	cmd.AddCommand(issues.NewCmdIssues(f, nil))
	cmd.AddCommand(users.NewCmdUsers(f, nil))

	return cmd
}
