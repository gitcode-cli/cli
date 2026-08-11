// Package repos implements the search repos command.
package repos

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// ReposOptions configures the search repos command.
type ReposOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	Query string
	Limit int
	Page  int
	JSON  bool
}

// NewCmdRepos creates the search repos command.
func NewCmdRepos(f *cmdutil.Factory, runF func(*ReposOptions) error) *cobra.Command {
	opts := &ReposOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "repos <query>",
		Short: "Search for repositories",
		Long:  heredoc.Doc(`Search GitCode repositories by keyword.`),
		Example: heredoc.Doc(`
			$ gc search repos "gitcode"
			$ gc search repos "cli tool" --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Query = args[0]
			if runF != nil {
				return runF(opts)
			}
			return reposRun(opts)
		},
	}

	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum results per page")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func reposRun(opts *ReposOptions) error {
	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}

	results, err := api.SearchRepos(client, &api.SearchOptions{
		Q:       opts.Query,
		PerPage: opts.Limit,
		Page:    opts.Page,
	})
	if err != nil {
		return fmt.Errorf("failed to search repositories: %w", err)
	}

	if len(results) == 0 {
		fmt.Fprintln(opts.IO.Out, "No repositories found.")
		return nil
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, results)
	}

	for _, r := range results {
		desc := r.Description
		if desc == "" {
			desc = "-"
		}
		fmt.Fprintf(opts.IO.Out, "%s\t%s\n", r.FullName, desc)
	}
	return nil
}
