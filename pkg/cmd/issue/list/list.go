// Package list implements the issue list command
package list

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/gitcode-com/gitcode-cli/api"
	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
	"github.com/gitcode-com/gitcode-cli/pkg/iostreams"
)

type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string

	// Flags
	State    string
	Limit    int
	Labels   string
	Assignee string
}

// NewCmdList creates the list command
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List issues in a repository",
		Long: heredoc.Doc(`
			List issues in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# List open issues
			$ gc issue list

			# List closed issues
			$ gc issue list --state closed

			# List issues with specific labels
			$ gc issue list --label bug,enhancement

			# List issues in a specific repository
			$ gc issue list -R owner/repo
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
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of issues to list")
	cmd.Flags().StringVarP(&opts.Labels, "label", "l", "", "Filter by labels (comma separated)")
	cmd.Flags().StringVarP(&opts.Assignee, "assignee", "a", "", "Filter by assignee")

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

	// Get repository
	owner, repo, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	// List issues
	issues, err := api.ListRepoIssues(client, owner, repo, &api.IssueListOptions{
		State:   opts.State,
		Labels:  opts.Labels,
		PerPage: opts.Limit,
	})
	if err != nil {
		return fmt.Errorf("failed to list issues: %w", err)
	}

	// Output
	if len(issues) == 0 {
		fmt.Fprintf(opts.IO.Out, "No issues found\n")
		return nil
	}

	fmt.Fprintf(opts.IO.Out, "\n")
	for _, issue := range issues {
		state := "open"
		if issue.State == "closed" {
			state = cs.Red("closed")
		} else {
			state = cs.Green("open")
		}
		fmt.Fprintf(opts.IO.Out, "#%-6s %s  %s\n", issue.Number, state, issue.Title)
	}
	fmt.Fprintf(opts.IO.Out, "\n")

	return nil
}

func parseRepo(repo string) (string, string, error) {
	if repo == "" {
		// TODO: get from current git repo
		return "", "", fmt.Errorf("no repository specified. Use -R owner/repo")
	}

	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repository format: %s", repo)
	}
	return parts[0], parts[1], nil
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITCODE_TOKEN")
}