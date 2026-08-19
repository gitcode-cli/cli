// Package view implements the ssh-key view command.
package view

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// ViewOptions configures the ssh-key view command.
type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	ID         int64
	JSON       bool
}

// NewCmdView creates the ssh-key view command.
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{IO: f.IOStreams, HttpClient: f.HttpClient}
	cmd := &cobra.Command{
		Use:     "view <id>",
		Short:   "View an SSH public key",
		Long:    heredoc.Doc("View an SSH public key by its numeric ID."),
		Example: heredoc.Doc("$ gc ssh-key view 123 --json"),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			opts.ID = id
			if runF != nil {
				return runF(opts)
			}
			return viewRun(opts)
		},
	}
	cmdutil.AddJSONFlag(cmd, &opts.JSON)
	return cmd
}

func viewRun(opts *ViewOptions) error {
	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}
	key, err := api.GetSSHKey(client, opts.ID)
	if err != nil {
		return cmdutil.WrapNotFound(err, "SSH key %d not found", opts.ID)
	}
	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, key)
	}
	fmt.Fprintf(opts.IO.Out, "%s\n  ID:      %d\n  Key:     %s\n", cmdutil.SanitizeTerminalText(key.Title), key.ID,
		cmdutil.SanitizeTerminalText(key.Key))
	if key.CreatedAt != "" {
		fmt.Fprintf(opts.IO.Out, "  Created: %s\n", cmdutil.SanitizeTerminalText(key.CreatedAt))
	}
	return nil
}

func parseID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, cmdutil.NewUsageError("SSH key id must be a positive integer")
	}
	return id, nil
}
