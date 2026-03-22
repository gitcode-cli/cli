// Package view implements the repo view command
package view

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

type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string

	// Flags
	Web bool
}

// NewCmdView creates the view command
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "view [<repository>]",
		Short: "View a repository",
		Long: heredoc.Doc(`
			View information about a repository.
		`),
		Example: heredoc.Doc(`
			# View a repository
			$ gc repo view owner/repo

			# View current repository
			$ gc repo view

			# Open in browser
			$ gc repo view owner/repo --web
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Repository = args[0]
			}

			if runF != nil {
				return runF(opts)
			}
			return viewRun(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Web, "web", "w", false, "Open in browser")

	return cmd
}

func viewRun(opts *ViewOptions) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := os.Getenv("GC_TOKEN")
	if token == "" {
		token = os.Getenv("GITCODE_TOKEN")
	}
	if token != "" {
		client.SetToken(token, "environment")
	}

	// Parse repository
	owner, name, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	repo, err := api.GetRepo(client, owner, name)
	if err != nil {
		return fmt.Errorf("failed to get repository: %w", err)
	}

	// Output
	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "%s\n", cs.Bold(repo.FullName))
	fmt.Fprintf(opts.IO.Out, "  %s\n", repo.Description)
	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "  Language: %s\n", repo.Language)
	fmt.Fprintf(opts.IO.Out, "  Stars: %d  Forks: %d  Issues: %d\n", repo.StargazersCount, repo.ForksCount, repo.OpenIssuesCount)
	fmt.Fprintf(opts.IO.Out, "  Default branch: %s\n", repo.DefaultBranch)
	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "  %s\n", repo.HTMLURL)
	fmt.Fprintf(opts.IO.Out, "\n")

	return nil
}

func parseRepo(repo string) (string, string, error) {
	if repo == "" {
		// TODO: get from current git repo
		return "", "", fmt.Errorf("no repository specified")
	}

	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repository format: %s", repo)
	}
	return parts[0], parts[1], nil
}