// Package list implements the commit comments list command
package list

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

type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	Repository string
	Page       int
	PerPage    int
	Order      string
}

// NewCmdList creates the list command
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List commit comments",
		Long: heredoc.Doc(`
			List all commit comments in a repository.
		`),
		Example: heredoc.Doc(`
			# List all commit comments
			$ gc commit comments list -R owner/repo

			# List with pagination
			$ gc commit comments list -R owner/repo --page 1 --per-page 50
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().IntVarP(&opts.Page, "page", "p", 1, "Page number")
	cmd.Flags().IntVarP(&opts.PerPage, "per-page", "P", 20, "Results per page (max 100)")
	cmd.Flags().StringVar(&opts.Order, "order", "", "Sort order: asc or desc")

	return cmd
}

func listRun(opts *ListOptions) error {
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
		Order:   opts.Order,
	}

	comments, err := api.ListCommitComments(client, owner, repo, listOpts)
	if err != nil {
		return fmt.Errorf("failed to list comments: %w", err)
	}

	if len(comments) == 0 {
		fmt.Fprintf(opts.IO.Out, "No comments found.\n")
		return nil
	}

	// Output
	fmt.Fprintf(opts.IO.Out, "\n")
	for _, c := range comments {
		author := "unknown"
		if c.User != nil {
			author = c.User.Login
		}
		body := c.Body
		if len(body) > 50 {
			body = body[:47] + "..."
		}
		fmt.Fprintf(opts.IO.Out, "#%-6s %s  %s\n", cmdutil.FormatAPIID(c.ID), cs.Bold(author), body)
	}
	fmt.Fprintf(opts.IO.Out, "\n")

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
