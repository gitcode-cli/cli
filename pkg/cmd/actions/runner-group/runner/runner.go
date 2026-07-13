// Package runner implements the actions runner-group runner command.
package runner

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/runner-group/runner/list"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdRunner creates the actions runner-group runner command.
func NewCmdRunner(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runner <command>",
		Short: "Manage runners in a runner group",
		Long: heredoc.Doc(`
			List runners (host or K8S) within an organization runner group.
		`),
		Example: heredoc.Doc(`
			# List host runners in a runner group
			$ gc actions runner-group runner list <runner-group-id> --org my-org

			# Output as JSON
			$ gc actions runner-group runner list <runner-group-id> --org my-org --json
		`),
	}

	cmd.AddCommand(list.NewCmdList(f, nil))

	return cmd
}
