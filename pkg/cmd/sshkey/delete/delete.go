// Package delete implements the ssh-key delete command.
package delete

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

// DeleteOptions configures the ssh-key delete command.
type DeleteOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	ID         int64
	Yes        bool
	DryRun     bool
	JSON       bool
}

type deleteResult struct {
	ID     int64  `json:"id"`
	Action string `json:"action"`
}

// NewCmdDelete creates the ssh-key delete command.
func NewCmdDelete(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{IO: f.IOStreams, HttpClient: f.HttpClient}
	cmd := &cobra.Command{
		Use:     "delete <id>",
		Short:   "Delete an SSH public key",
		Long:    heredoc.Doc("Delete an SSH public key by ID. Non-interactive mode requires --yes."),
		Example: heredoc.Doc("$ gc ssh-key delete 123\n$ gc ssh-key delete 123 --yes --json"),
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
			return deleteRun(opts)
		},
	}
	cmd.Flags().BoolVarP(&opts.Yes, "yes", "y", false, "Skip confirmation")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Preview the deletion without making changes")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)
	return cmd
}

func deleteRun(opts *DeleteOptions) error {
	id := strconv.FormatInt(opts.ID, 10)
	if opts.DryRun {
		result := deleteResult{ID: opts.ID, Action: "would-delete"}
		if opts.JSON {
			return cmdutil.WriteJSON(opts.IO.Out, result)
		}
		fmt.Fprintf(opts.IO.Out, "Would delete SSH key %d.\n", opts.ID)
		return nil
	}
	if err := cmdutil.ConfirmOrAbort(cmdutil.ConfirmOptions{
		IO: opts.IO, Yes: opts.Yes, Expected: id,
		Prompt: fmt.Sprintf("Type the SSH key id to confirm deletion: %s\n", id),
	}); err != nil {
		return err
	}
	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}
	if err := api.DeleteSSHKey(client, opts.ID); err != nil {
		return cmdutil.WrapNotFound(err, "SSH key %d not found", opts.ID)
	}
	result := deleteResult{ID: opts.ID, Action: "deleted"}
	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, result)
	}
	fmt.Fprintf(opts.IO.Out, "Deleted SSH key %d.\n", opts.ID)
	return nil
}

func parseID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, cmdutil.NewUsageError("SSH key id must be a positive integer")
	}
	return id, nil
}
