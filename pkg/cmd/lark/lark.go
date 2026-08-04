// Package lark implements the lark command group, a subprocess bridge to the
// official Feishu/Lark CLI (lark-cli). gc delegates messaging and auth status
// to lark-cli so that Feishu OAuth credentials stay in the OS keychain and are
// never seen by gc.
package lark

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdLark creates the lark command group.
func NewCmdLark(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lark <command>",
		Short: "Send messages and manage Feishu/Lark integration via lark-cli",
		Long: heredoc.Doc(`
			Work with Feishu/Lark through the official lark-cli tool.

			gc delegates messaging and auth status to lark-cli so that Feishu
			OAuth credentials stay in the OS keychain and are never seen by gc.
			When lark-cli is missing, run "gc lark install" to bootstrap it.
		`),
		Example: heredoc.Doc(`
			# Send a text message to a Feishu group
			$ gc lark send --chat-id oc_xxx --text "deploy done"

			# Use the configured default group
			$ gc lark config set --default-chat oc_xxx
			$ gc lark send --text "build green"

			# Check lark-cli setup and login
			$ gc lark doctor
		`),
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "lark",
		},
	}

	cmd.AddCommand(newCmdSend(f, nil))
	cmd.AddCommand(newCmdAuth(f, nil))
	cmd.AddCommand(newCmdInstall(f, nil))
	cmd.AddCommand(newCmdDoctor(f, nil))
	cmd.AddCommand(newCmdConfig(f, nil, nil))

	return cmd
}
