// Package artifact implements the actions artifact command.
package artifact

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/artifact/delete"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/artifact/download"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/artifact/list"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/artifact/view"
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

			# View an artifact detail
			$ gc actions artifact view <artifact-id> -R owner/repo

			# Download an artifact
			$ gc actions artifact download <artifact-id> -R owner/repo --output artifact.zip

			# Delete an artifact
			$ gc actions artifact delete <artifact-id> -R owner/repo --yes
		`),
	}

	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(view.NewCmdView(f, nil))
	cmd.AddCommand(download.NewCmdDownload(f, nil))
	cmd.AddCommand(delete.NewCmdDelete(f, nil))

	return cmd
}
