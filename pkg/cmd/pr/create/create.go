// Package create implements the pr create command
package create

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

type CreateOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string

	// Flags
	Title     string
	Body      string
	Head      string
	Base      string
	Draft     bool
	Fill      bool
	Web       bool
}

// NewCmdCreate creates the create command
func NewCmdCreate(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a pull request",
		Long: heredoc.Doc(`
			Create a new pull request in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Create a PR interactively
			$ gc pr create

			# Create a PR with title and body
			$ gc pr create --title "Feature" --body "Description"

			# Create a draft PR
			$ gc pr create --draft

			# Create a PR from current branch
			$ gc pr create --fill
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return createRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "Title for the PR")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "Body for the PR")
	cmd.Flags().StringVarP(&opts.Head, "head", "H", "", "Head branch")
	cmd.Flags().StringVarP(&opts.Base, "base", "B", "main", "Base branch")
	cmd.Flags().BoolVarP(&opts.Draft, "draft", "d", false, "Create as draft")
	cmd.Flags().BoolVarP(&opts.Fill, "fill", "f", false, "Fill from last commit")
	cmd.Flags().BoolVarP(&opts.Web, "web", "w", false, "Open in browser")

	return cmd
}

func createRun(opts *CreateOptions) error {
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

	// Validate required fields
	if opts.Title == "" {
		return fmt.Errorf("title is required. Use --title flag")
	}
	if opts.Head == "" {
		return fmt.Errorf("head branch is required. Use --head flag")
	}

	// Create PR
	pr, err := api.CreatePullRequest(client, owner, repo, &api.CreatePROptions{
		Title: opts.Title,
		Body:  opts.Body,
		Head:  opts.Head,
		Base:  opts.Base,
		Draft: opts.Draft,
	})
	if err != nil {
		return fmt.Errorf("failed to create PR: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Created PR #%d in %s/%s\n", cs.Green("✓"), pr.Number, owner, repo)
	fmt.Fprintf(opts.IO.Out, "  %s\n", pr.HTMLURL)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	if repo == "" {
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