// Package sshkey implements the ssh-key command.
package sshkey

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/sshkey/add"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/sshkey/delete"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/sshkey/list"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/sshkey/view"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdSSHKey creates the ssh-key command.
func NewCmdSSHKey(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssh-key <command>",
		Short: "Manage SSH public keys",
		Long: heredoc.Doc(`
			List, add, view, and delete SSH public keys for the authenticated user.
		`),
		Example: heredoc.Doc(`
			# List SSH keys
			$ gc ssh-key list

			# Add an SSH public key from a file
			$ gc ssh-key add --title laptop --key ~/.ssh/id_ed25519.pub

			# View or delete a key
			$ gc ssh-key view 123
			$ gc ssh-key delete 123 --yes
		`),
	}

	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(add.NewCmdAdd(f, nil))
	cmd.AddCommand(view.NewCmdView(f, nil))
	cmd.AddCommand(delete.NewCmdDelete(f, nil))
	return cmd
}
