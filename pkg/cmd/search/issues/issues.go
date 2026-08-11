// Package issues implements the search issues command.
package issues

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// IssuesOptions configures the search issues command.
type IssuesOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	Query string
	Repo  string
	State string
	Limit int
	Page  int
	JSON  bool
}

// NewCmdIssues creates the search issues command.
func NewCmdIssues(f *cmdutil.Factory, runF func(*IssuesOptions) error) *cobra.Command {
	opts := &IssuesOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "issues <query>",
		Short: "Search for issues",
		Long:  heredoc.Doc(`Search GitCode issues by keyword.`),
		Example: heredoc.Doc(`
			$ gc search issues "bug report"
			$ gc search issues "feature" --repo owner/repo --state open
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Query = args[0]
			if runF != nil {
				return runF(opts)
			}
			return issuesRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repo, "repo", "R", "", "Limit search to a repository (owner/repo)")
	cmd.Flags().StringVar(&opts.State, "state", "", "Filter by state (open/closed)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum results per page")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func issuesRun(opts *IssuesOptions) error {
	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}

	results, err := api.SearchIssues(client, &api.SearchOptions{
		Q:       opts.Query,
		PerPage: opts.Limit,
		Page:    opts.Page,
	}, opts.Repo, opts.State)
	if err != nil {
		return fmt.Errorf("failed to search issues: %w", err)
	}

	if len(results) == 0 {
		fmt.Fprintln(opts.IO.Out, "No issues found.")
		return nil
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, results)
	}

	for _, r := range results {
		fmt.Fprintf(opts.IO.Out, "#%s\t%s\t%s\n", r.Number, r.State, r.Title)
	}
	return nil
}
