// Package list implements the pr list command
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

type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string

	// Flags
	State string
	Limit int
	Head  string
	Base  string
	JSON  bool
}

// NewCmdList creates the list command
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List pull requests",
		Long: heredoc.Doc(`
			List pull requests in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# List open PRs
			$ gc pr list

			# List closed PRs
			$ gc pr list --state closed

			# List PRs in a specific repository
			$ gc pr list -R owner/repo
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.State, "state", "s", "open", "Filter by state (open/closed/all)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of PRs to list")
	cmd.Flags().StringVarP(&opts.Head, "head", "H", "", "Filter by head branch")
	cmd.Flags().StringVarP(&opts.Base, "base", "B", "", "Filter by base branch")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func listRun(opts *ListOptions) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := cmdutil.EnvToken()
	if token == "" {
		return cmdutil.NewAuthError("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "environment")

	// Get repository
	owner, repo, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	// List PRs
	prs, err := api.ListPullRequests(client, owner, repo, &api.PRListOptions{
		State:   opts.State,
		Head:    opts.Head,
		Base:    opts.Base,
		PerPage: opts.Limit,
	})
	if err != nil {
		return fmt.Errorf("failed to list PRs: %w", err)
	}

	// Output
	if len(prs) == 0 {
		if opts.JSON {
			return cmdutil.WriteJSON(opts.IO.Out, prs)
		}
		fmt.Fprintf(opts.IO.Out, "No pull requests found\n")
		return nil
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, prs)
	}

	fmt.Fprintf(opts.IO.Out, "\n")
	for _, pr := range prs {
		var state string
		switch pr.State {
		case "merged":
			state = cs.Magenta("merged")
		case "closed":
			state = cs.Red("closed")
		default:
			if pr.Draft {
				state = cs.Gray("draft")
			} else {
				state = cs.Green("open")
			}
		}
		fmt.Fprintf(opts.IO.Out, "#%-6d %s  %s\n", pr.Number, state, pr.Title)
	}
	fmt.Fprintf(opts.IO.Out, "\n")

	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
