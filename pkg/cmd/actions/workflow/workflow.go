// Package workflow implements the actions workflow command.
package workflow

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/workflow/list"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/workflow/run"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdWorkflow creates the actions workflow command.
func NewCmdWorkflow(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow <command>",
		Short: "Manage workflow files",
		Long: heredoc.Doc(`
			List and inspect GitCode Actions workflow files in a repository.
		`),
		Example: heredoc.Doc(`
			# List workflows
			$ gc actions workflow list -R owner/repo
		`),
	}

	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(run.NewCmdRun(f, nil))

	return cmd
}
