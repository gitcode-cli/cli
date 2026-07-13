// Package runner implements the actions runner command.
package runner

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/runner/list"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdRunner creates the actions runner command.
func NewCmdRunner(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runner <command>",
		Short: "Manage repository runners",
		Long: heredoc.Doc(`
			List runners (host or K8S) for a GitCode repository.

			These are repo-level runners, separate from organization runner
			groups (see gc actions runner-group).
		`),
		Example: heredoc.Doc(`
			# List host runners for a repository
			$ gc actions runner list -R owner/repo

			# Output as JSON
			$ gc actions runner list -R owner/repo --json
		`),
	}

	cmd.AddCommand(list.NewCmdList(f, nil))

	return cmd
}
