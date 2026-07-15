// Package listbysha implements the commit comments list-by-sha command
package listbysha

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type ListBySHAOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	SHA        string
	Page       int
	PerPage    int
	JSON       bool
}

// NewCmdListBySHA creates the list-by-sha command
func NewCmdListBySHA(f *cmdutil.Factory, runF func(*ListBySHAOptions) error) *cobra.Command {
	opts := &ListBySHAOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "list-by-sha <sha>",
		Short: "List comments for a specific commit",
		Long: heredoc.Doc(`
			List comments for a specific commit in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# List comments for a commit
			$ gc commit comments list-by-sha abc123 -R owner/repo

			# List with pagination
			$ gc commit comments list-by-sha abc123 -R owner/repo --page 1 --per-page 50

			# List comments as JSON
			$ gc commit comments list-by-sha abc123 -R owner/repo --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.SHA = args[0]

			if runF != nil {
				return runF(opts)
			}
			return listBySHARun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().IntVarP(&opts.Page, "page", "p", 1, "Page number")
	cmd.Flags().IntVarP(&opts.PerPage, "per-page", "P", 20, "Results per page (max 100)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func listBySHARun(opts *ListBySHAOptions) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}
	client, err := cmdutil.AuthenticatedClient(httpClient)
	if err != nil {
		return err
	}

	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}
	owner, repo, err := parseRepo(repository)
	if err != nil {
		return err
	}

	listOpts := &api.ListOptions{
		Page:    opts.Page,
		PerPage: opts.PerPage,
	}

	comments, err := api.ListCommentsForCommit(client, owner, repo, opts.SHA, listOpts)
	if err != nil {
		return fmt.Errorf("failed to list comments: %w", err)
	}

	if opts.JSON {
		if comments == nil {
			comments = []api.CommitComment{}
		}
		return cmdutil.WriteJSON(opts.IO.Out, comments)
	}

	if len(comments) == 0 {
		fmt.Fprintf(opts.IO.Out, "No comments found for commit %s.\n", opts.SHA)
		return nil
	}

	// Calculate max ID width for alignment
	maxIDWidth := 0
	for _, c := range comments {
		w := len(fmt.Sprintf("#%s", cmdutil.FormatAPIID(c.ID)))
		if w > maxIDWidth {
			maxIDWidth = w
		}
	}

	// Output
	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "Comments for commit %s:\n\n", opts.SHA)
	for _, c := range comments {
		author := "unknown"
		if c.User != nil {
			author = c.User.Login
		}
		fmt.Fprintf(opts.IO.Out, "%-*s  %s at %s\n", maxIDWidth, fmt.Sprintf("#%s", cmdutil.FormatAPIID(c.ID)), cs.Bold(author), c.CreatedAt)
		fmt.Fprintf(opts.IO.Out, "  %s\n\n", c.Body)
	}

	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
