// Package job implements the actions job command.
package job

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/job/list"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdJob creates the actions job command.
func NewCmdJob(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "job <command>",
		Short: "Manage workflow jobs",
		Long: heredoc.Doc(`
			Inspect workflow jobs for a pipeline run.
		`),
		Example: heredoc.Doc(`
			# List jobs of a pipeline run
			$ gc actions job list <run-id> -R owner/repo
		`),
	}

	cmd.AddCommand(list.NewCmdList(f, nil))

	return cmd
}
