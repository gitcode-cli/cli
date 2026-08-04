package lark

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/larkcli"
)

// newCmdConfig creates the "gc lark config" command group for managing the
// default Feishu chat id stored in ~/.config/gc/lark.json.
func newCmdConfig(f *cmdutil.Factory, runSet func(*configSetOptions) error, runGet func(*configGetOptions) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Manage gc-lark integration settings (default chat)",
		Long: heredoc.Doc(`
			Manage gc-lark integration settings stored in ~/.config/gc/lark.json.
			The default chat id is resolved in priority order: the
			GC_LARK_DEFAULT_CHAT_ID environment variable, then the persisted
			default_chat_id. Use --chat-id on the send command to override.
		`),
		Args: cobra.NoArgs,
	}

	cmd.AddCommand(newCmdConfigSet(f, runSet))
	cmd.AddCommand(newCmdConfigGet(f, runGet))
	return cmd
}

type configSetOptions struct {
	IO          *iostreams.IOStreams
	DefaultChat string
}

func newCmdConfigSet(f *cmdutil.Factory, runF func(*configSetOptions) error) *cobra.Command {
	opts := &configSetOptions{IO: f.IOStreams}

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set the default Feishu/Lark chat id",
		Example: heredoc.Doc(`
			# Set the default chat
			$ gc lark config set --default-chat oc_xxx
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return configSetRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.DefaultChat, "default-chat", "", "Feishu/Lark chat id (oc_xxx) to use as default")
	_ = cmd.MarkFlagRequired("default-chat")

	return cmd
}

func configSetRun(opts *configSetOptions) error {
	cs := opts.IO.ColorScheme()
	if err := larkcli.SaveDefaultChat(opts.DefaultChat); err != nil {
		return fmt.Errorf("failed to save default chat: %w", err)
	}
	fmt.Fprintf(opts.IO.Out, "%s Default chat set to %s\n", cs.Green("✓"), opts.DefaultChat)
	return nil
}

type configGetOptions struct {
	IO   *iostreams.IOStreams
	JSON bool
}

func newCmdConfigGet(f *cmdutil.Factory, runF func(*configGetOptions) error) *cobra.Command {
	opts := &configGetOptions{IO: f.IOStreams}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Show the effective default Feishu/Lark chat id",
		Example: heredoc.Doc(`
			# Show the effective default chat
			$ gc lark config get

			# JSON output
			$ gc lark config get --json
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return configGetRun(opts)
		},
	}

	cmdutil.AddJSONFlag(cmd, &opts.JSON)
	return cmd
}

func configGetRun(opts *configGetOptions) error {
	chatID := larkcli.DefaultChatID()
	cs := opts.IO.ColorScheme()
	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, map[string]interface{}{
			"default_chat": chatID,
		})
	}
	if chatID == "" {
		fmt.Fprintf(opts.IO.Out, "%s no default chat configured (set with: gc lark config set --default-chat <oc_xxx>)\n", cs.Blue("i"))
	} else {
		fmt.Fprintf(opts.IO.Out, "%s\n", chatID)
	}
	return nil
}
