// Package auth implements authentication commands
package auth

import (
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/auth/login"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/auth/logout"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/auth/status"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/auth/token"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdAuth creates the auth command
func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Authenticate with GitCode",
		Long: `Authenticate with a GitCode account.

Available commands:
  login    Log in to a GitCode account
  logout   Log out of a GitCode account
  status   View authentication status
  token    Print an authentication token`,
	}

	cmd.AddCommand(login.NewCmdLogin(f, nil))
	cmd.AddCommand(logout.NewCmdLogout(f, nil))
	cmd.AddCommand(status.NewCmdStatus(f, nil))
	cmd.AddCommand(token.NewCmdToken(f, nil))

	return cmd
}