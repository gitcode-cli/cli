// Package runnerset implements the actions runner-set command.
package runnerset

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/runner-set/list"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/runner-set/shared-runner-sets"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdRunnerSet creates the actions runner-set command.
func NewCmdRunnerSet(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runner-set <command>",
		Short: "Manage repository K8S runner sets",
		Long: heredoc.Doc(`
			List K8S runner sets for a GitCode repository.

			These are repo-level runner sets, separate from organization runner
			groups (see gc actions runner-group runner-set).
		`),
		Example: heredoc.Doc(`
			# List K8S runner sets for a repository
			$ gc actions runner-set list -R owner/repo

			# Output as JSON
			$ gc actions runner-set list -R owner/repo --json
		`),
	}

	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(sharedrunnersets.NewCmdList(f, nil))

	return cmd
}
