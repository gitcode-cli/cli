// Package list implements the ssh-key list command.
package list

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// ListOptions configures the ssh-key list command.
type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	Limit      int
	Page       int
	PerPage    int
	JSON       bool
}

// NewCmdList creates the ssh-key list command.
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{IO: f.IOStreams, HttpClient: f.HttpClient}
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List SSH public keys",
		Long:    heredoc.Doc("List SSH public keys for the authenticated user."),
		Example: heredoc.Doc("$ gc ssh-key list\n$ gc ssh-key list --limit 50 --page 2\n$ gc ssh-key list --json"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of SSH keys to list")
	cmd.Flags().IntVar(&opts.Page, "page", 1, "Page number")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "API page size (default: --limit)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)
	return cmd
}

func listRun(opts *ListOptions) error {
	if opts.Limit <= 0 || opts.Page <= 0 || opts.PerPage < 0 {
		return cmdutil.NewUsageError("--limit and --page must be positive; --per-page cannot be negative")
	}
	perPage := opts.PerPage
	if perPage == 0 {
		perPage = opts.Limit
	}
	if perPage > 100 {
		return cmdutil.NewUsageError("--per-page must not exceed 100")
	}
	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}
	keys, err := api.ListSSHKeys(client, opts.Page, perPage)
	if err != nil {
		return fmt.Errorf("failed to list SSH keys: %w", err)
	}
	if len(keys) > opts.Limit {
		keys = keys[:opts.Limit]
	}
	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, keys)
	}
	for _, key := range keys {
		fmt.Fprintf(opts.IO.Out, "%d\t%s\t%s\n", key.ID, key.Title, key.CreatedAt)
	}
	return nil
}
