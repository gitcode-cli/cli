// Package sharednamespace implements the actions runner-group shared-namespace command.
package sharednamespace

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/runner-group/shared-namespace/list"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdSharedNamespace creates the actions runner-group shared-namespace command.
func NewCmdSharedNamespace(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shared-namespace <command>",
		Short: "Manage shared namespaces for a runner group",
		Long: heredoc.Doc(`
			List namespaces that have access to an organization runner group.
		`),
		Example: heredoc.Doc(`
			# List shared namespaces for a runner group
			$ gc actions runner-group shared-namespace list <runner-group-id> --org my-org

			# Output as JSON
			$ gc actions runner-group shared-namespace list <runner-group-id> --org my-org --json
		`),
	}

	cmd.AddCommand(list.NewCmdList(f, nil))

	return cmd
}
