// Package runnergroup implements the actions runner-group command.
package runnergroup

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/runner-group/list"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/runner-group/view"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdRunnerGroup creates the actions runner-group command.
func NewCmdRunnerGroup(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runner-group <command>",
		Short: "Manage organization runner groups",
		Long: heredoc.Doc(`
			Inspect organization-level runner groups for GitCode Actions.

			Runner groups organize self-hosted runners (host and K8S) at the
			organization level and can be shared with repositories.
		`),
		Example: heredoc.Doc(`
			# List runner groups in an organization
			$ gc actions runner-group list --org my-org

			# Filter by keyword
			$ gc actions runner-group list --org my-org --keyword prod

			# Output as JSON
			$ gc actions runner-group list --org my-org --json
		`),
	}

	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(view.NewCmdView(f, nil))

	return cmd
}
