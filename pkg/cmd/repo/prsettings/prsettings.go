// Package prsettings implements the repo pr-settings command.
package prsettings

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/repo/prsettings/update"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/repo/prsettings/view"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdPRSettings creates the repo pr-settings command.
func NewCmdPRSettings(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pr-settings <command>",
		Short: "Manage pull request merge gate settings",
		Long: heredoc.Doc(`
			View and edit pull request settings for a GitCode repository,
			including merge gate rules, approval requirements, and pipeline
			gating.
		`),
		Example: heredoc.Doc(`
			# View PR settings
			$ gc repo pr-settings view -R owner/repo

			# Require pipeline success before merge
			$ gc repo pr-settings update --pipeline-required -R owner/repo

			# Set minimum reviewers to 2
			$ gc repo pr-settings update --reviewers 2 -R owner/repo
		`),
	}

	cmd.AddCommand(view.NewCmdView(f, nil))
	cmd.AddCommand(update.NewCmdUpdate(f, nil))

	return cmd
}
