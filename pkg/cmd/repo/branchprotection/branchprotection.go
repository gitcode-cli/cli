// Package branchprotection implements the repo branch-protection command.
package branchprotection

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/repo/branchprotection/create"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/repo/branchprotection/delete"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/repo/branchprotection/list"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/repo/branchprotection/update"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdBranchProtection creates the repo branch-protection command.
func NewCmdBranchProtection(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "branch-protection <command>",
		Short: "Manage branch protection rules",
		Long: heredoc.Doc(`
			Manage branch protection rules for a GitCode repository.

			Protection rules control who can push and merge to specific
			branches or branch patterns.
		`),
		Example: heredoc.Doc(`
			# List protection rules
			$ gc repo branch-protection list -R owner/repo

			# Create a protection rule
			$ gc repo branch-protection create --wildcard "main" --pusher admin --merger admin -R owner/repo

			# Update a protection rule
			$ gc repo branch-protection update "main" --pusher develop -R owner/repo

			# Delete a protection rule
			$ gc repo branch-protection delete "main" -R owner/repo
		`),
	}

	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(create.NewCmdCreate(f, nil))
	cmd.AddCommand(update.NewCmdUpdate(f, nil))
	cmd.AddCommand(delete.NewCmdDelete(f, nil))

	return cmd
}
