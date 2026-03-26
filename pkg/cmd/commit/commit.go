// Package commit implements the commit command
package commit

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/commit/diff"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/commit/patch"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/commit/view"
)

// NewCmdCommit creates the commit command
func NewCmdCommit(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit <command>",
		Short: "Manage commits",
		Long: heredoc.Doc(`
			Work with GitCode commits.

			View commit details, diffs, and manage commit comments.
		`),
		Example: heredoc.Doc(`
			# View a commit
			$ gc commit view abc123 -R owner/repo

			# View commit with diff files
			$ gc commit view abc123 -R owner/repo --show-diff

			# Get commit diff
			$ gc commit diff abc123 -R owner/repo
		`),
		Annotations: map[string]string{
			"IsCore": "true",
		},
	}

	cmd.AddCommand(view.NewCmdView(f, nil))
	cmd.AddCommand(diff.NewCmdDiff(f, nil))
	cmd.AddCommand(patch.NewCmdPatch(f, nil))

	return cmd
}