// Package users implements the search users command.
package users

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// UsersOptions configures the search users command.
type UsersOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	Query string
	Limit int
	Page  int
	JSON  bool
}

// NewCmdUsers creates the search users command.
func NewCmdUsers(f *cmdutil.Factory, runF func(*UsersOptions) error) *cobra.Command {
	opts := &UsersOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "users <query>",
		Short: "Search for users",
		Long:  heredoc.Doc(`Search GitCode users by keyword.`),
		Example: heredoc.Doc(`
			$ gc search users "developer"
			$ gc search users "admin" --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Query = args[0]
			if runF != nil {
				return runF(opts)
			}
			return usersRun(opts)
		},
	}

	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum results per page")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func usersRun(opts *UsersOptions) error {
	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}

	results, err := api.SearchUsers(client, &api.SearchOptions{
		Q:       opts.Query,
		PerPage: opts.Limit,
		Page:    opts.Page,
	})
	if err != nil {
		return fmt.Errorf("failed to search users: %w", err)
	}

	if len(results) == 0 {
		fmt.Fprintln(opts.IO.Out, "No users found.")
		return nil
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, results)
	}

	for _, u := range results {
		name := u.Name
		if name == "" {
			name = "-"
		}
		fmt.Fprintf(opts.IO.Out, "%s\t%s\n", u.Login, name)
	}
	return nil
}
