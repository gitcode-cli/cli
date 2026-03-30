// Package listbysha implements the commit comments list-by-sha command
package listbysha

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type ListBySHAOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	Repository string
	SHA        string
	Page       int
	PerPage    int
}

// NewCmdListBySHA creates the list-by-sha command
func NewCmdListBySHA(f *cmdutil.Factory, runF func(*ListBySHAOptions) error) *cobra.Command {
	opts := &ListBySHAOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
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

	return cmd
}

func listBySHARun(opts *ListBySHAOptions) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := getEnvToken()
	if token == "" {
		return fmt.Errorf("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "environment")

	owner, repo, err := parseRepo(opts.Repository)
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

	if len(comments) == 0 {
		fmt.Fprintf(opts.IO.Out, "No comments found for commit %s.\n", opts.SHA)
		return nil
	}

	// Output
	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "Comments for commit %s:\n\n", opts.SHA)
	for _, c := range comments {
		author := "unknown"
		if c.User != nil {
			author = c.User.Login
		}
		fmt.Fprintf(opts.IO.Out, "#%-6v %s at %s\n", c.ID, cs.Bold(author), c.CreatedAt)
		fmt.Fprintf(opts.IO.Out, "  %s\n\n", c.Body)
	}

	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITCODE_TOKEN")
}
