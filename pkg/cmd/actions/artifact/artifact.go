// Package artifact implements the actions artifact command.
package artifact

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/artifact/list"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdArtifact creates the actions artifact command.
func NewCmdArtifact(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "artifact <command>",
		Short: "Manage workflow artifacts",
		Long: heredoc.Doc(`
			Inspect workflow run artifacts for a repository.
		`),
		Example: heredoc.Doc(`
			# List repository artifacts
			$ gc actions artifact list -R owner/repo
		`),
	}

	cmd.AddCommand(list.NewCmdList(f, nil))

	return cmd
}
