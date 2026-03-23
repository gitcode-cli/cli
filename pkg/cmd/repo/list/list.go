// Package list implements the repo list command
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

	// Flags
	Limit      int
	Visibility string
	Owner      string
}

// NewCmdList creates the list command
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List repositories",
		Long: heredoc.Doc(`
			List repositories for the authenticated user or an organization.
		`),
		Example: heredoc.Doc(`
			# List your repositories
			$ gc repo list

			# List with limit
			$ gc repo list --limit 50

			# List only public repos
			$ gc repo list --visibility public
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}

	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of repos to list")
	cmd.Flags().StringVarP(&opts.Visibility, "visibility", "v", "", "Filter by visibility (public/private)")
	cmd.Flags().StringVarP(&opts.Owner, "owner", "o", "", "List repos for an organization")

	return cmd
}

func listRun(opts *ListOptions) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	// Create API client (token from env)
	client := api.NewClientFromHTTP(httpClient)

	// Get token from environment
	token := getEnvToken()
	if token == "" {
		return fmt.Errorf("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "environment")

	// List repos
	repos, err := api.ListUserRepos(client, &api.RepoListOptions{
		PerPage:    opts.Limit,
		Visibility: opts.Visibility,
	})
	if err != nil {
		return fmt.Errorf("failed to list repositories: %w", err)
	}

	// Output
	if len(repos) == 0 {
		fmt.Fprintf(opts.IO.Out, "No repositories found\n")
		return nil
	}

	fmt.Fprintf(opts.IO.Out, "\n")
	for _, repo := range repos {
		visibility := "public"
		if repo.Private {
			visibility = "private"
		}
		fmt.Fprintf(opts.IO.Out, "%s  %s  %s\n", cs.Bold(repo.FullName), visibility, repo.Description)
	}
	fmt.Fprintf(opts.IO.Out, "\n")

	return nil
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITCODE_TOKEN")
}