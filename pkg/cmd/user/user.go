// Package user implements the user command.
package user

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/user/edit"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/user/view"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdUser creates the user command.
func NewCmdUser(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user <command>",
		Short: "View and edit user profile",
		Long: heredoc.Doc(`
			View and edit GitCode user profile information.
		`),
		Example: heredoc.Doc(`
			# View current user
			$ gc user view

			# View a specific user
			$ gc user view <username>

			# Edit your profile
			$ gc user edit --name "New Name" --bio "New bio"
		`),
	}

	cmd.AddCommand(view.NewCmdView(f, nil))
	cmd.AddCommand(edit.NewCmdEdit(f, nil))

	return cmd
}
