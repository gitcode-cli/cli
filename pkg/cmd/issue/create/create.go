// Package create implements the issue create command
package create

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type CreateOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string

	// Flags
	Title       string
	Body        string
	Labels      []string
	Assignees   []string
	Milestone   int
}

// NewCmdCreate creates the create command
func NewCmdCreate(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new issue",
		Long: heredoc.Doc(`
			Create a new issue in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Create an issue interactively
			$ gc issue create

			# Create an issue with title and body
			$ gc issue create --title "Bug" --body "Description"

			# Create an issue with labels
			$ gc issue create --title "Feature" --label bug,enhancement

			# Create an issue in a specific repository
			$ gc issue create -R owner/repo --title "Bug"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return createRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "Title for the issue")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "Body for the issue")
	cmd.Flags().StringSliceVarP(&opts.Labels, "label", "l", []string{}, "Labels to add")
	cmd.Flags().StringSliceVarP(&opts.Assignees, "assignee", "a", []string{}, "Assignees")
	cmd.Flags().IntVarP(&opts.Milestone, "milestone", "m", 0, "Milestone number")

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

	// Validate title
	if opts.Title == "" {
		return fmt.Errorf("title is required. Use --title flag")
	}

	// Create issue
	issue, err := api.CreateIssue(client, owner, repo, &api.CreateIssueOptions{
		Title:     opts.Title,
		Body:      opts.Body,
		Labels:    opts.Labels,
		Assignees: opts.Assignees,
		Milestone: opts.Milestone,
	})
	if err != nil {
		return fmt.Errorf("failed to create issue: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Created issue #%s in %s/%s\n", cs.Green("✓"), issue.Number, owner, repo)
	fmt.Fprintf(opts.IO.Out, "  %s\n", issue.HTMLURL)
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