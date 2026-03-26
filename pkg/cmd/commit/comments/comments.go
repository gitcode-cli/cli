// Package comments implements the commit comments command
package comments

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/commit/comments/create"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/commit/comments/edit"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/commit/comments/list"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/commit/comments/listbysha"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/commit/comments/view"
)

// NewCmdComments creates the comments command
func NewCmdComments(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comments <command>",
		Short: "Manage commit comments",
		Long: heredoc.Doc(`
			Manage comments on commits in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Create a comment on a commit
			$ gc commit comments create abc123 --body "Comment text" -R owner/repo

			# List comments for a commit
			$ gc commit comments list-by-sha abc123 -R owner/repo

			# View a specific comment
			$ gc commit comments view 123 -R owner/repo
		`),
	}

	cmd.AddCommand(create.NewCmdCreate(f, nil))
	cmd.AddCommand(view.NewCmdView(f, nil))
	cmd.AddCommand(edit.NewCmdEdit(f, nil))
	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(listbysha.NewCmdListBySHA(f, nil))

	return cmd
}