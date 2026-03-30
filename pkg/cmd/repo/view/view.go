// Package view implements the repo view command
package view

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/browser"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	Repository string

	// Flags
	Web  bool
	JSON bool
}

// NewCmdView creates the view command
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
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
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func viewRun(opts *ViewOptions) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := cmdutil.EnvToken()
	if token != "" {
		client.SetToken(token, "environment")
	}

	// Parse repository
	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}

	owner, name, err := parseRepo(repository)
	if err != nil {
		return err
	}

	repo, err := api.GetRepo(client, owner, name)
	if err != nil {
		return fmt.Errorf("failed to get repository: %w", err)
	}

	// Open in browser if --web flag is set
	if opts.Web {
		fmt.Fprintf(opts.IO.Out, "Opening %s in your browser.\n", repo.HTMLURL)
		return browser.Open(repo.HTMLURL)
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, repo)
	}

	// Output
	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "%s\n", cs.Bold(repo.FullName))
	if repo.Description != "" {
		fmt.Fprintf(opts.IO.Out, "  %s\n", repo.Description)
	}
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
	return cmdutil.ParseRepo(repo)
}
