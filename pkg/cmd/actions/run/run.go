// Package run implements the actions run command.
package run

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/run/list"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdRun creates the actions run command.
func NewCmdRun(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <command>",
		Short: "Manage pipeline runs",
		Long: heredoc.Doc(`
			Inspect pipeline run records for a repository.
		`),
		Example: heredoc.Doc(`
			# List recent pipeline runs
			$ gc actions run list -R owner/repo
		`),
	}

	cmd.AddCommand(list.NewCmdList(f, nil))

	return cmd
}
