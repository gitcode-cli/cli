// Package runnerset implements the actions runner-group runner-set command.
package runnerset

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/runner-group/runner-set/list"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdRunnerSet creates the actions runner-group runner-set command.
func NewCmdRunnerSet(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runner-set <command>",
		Short: "Manage K8S runner sets in a runner group",
		Long: heredoc.Doc(`
			List K8S runner sets within an organization runner group.
		`),
		Example: heredoc.Doc(`
			# List K8S runner sets in a runner group
			$ gc actions runner-group runner-set list <runner-group-id> --org my-org

			# Output as JSON
			$ gc actions runner-group runner-set list <runner-group-id> --org my-org --json
		`),
	}

	cmd.AddCommand(list.NewCmdList(f, nil))

	return cmd
}
