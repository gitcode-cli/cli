// Package view implements the issue view command
package view

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	Number     int

	// Flags
	Comments bool
	Web      bool
}

// NewCmdView creates the view command
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "view <number>",
		Short: "View an issue",
		Long: heredoc.Doc(`
			View an issue in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# View an issue
			$ gc issue view 123

			# View issue with comments
			$ gc issue view 123 --comments

			# View issue in browser
			$ gc issue view 123 --web
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid issue number: %s", args[0])
			}
			opts.Number = number

			if runF != nil {
				return runF(opts)
			}
			return viewRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVarP(&opts.Comments, "comments", "c", false, "View issue comments")
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
	token := getEnvToken()
	if token == "" {
		return fmt.Errorf("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "environment")

	// Get repository
	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}

	owner, repo, err := parseRepo(repository)
	if err != nil {
		return err
	}

	// Get issue
	issue, err := api.GetIssue(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to get issue: %w", err)
	}

	// Open in browser if --web flag is set
	if opts.Web {
		fmt.Fprintf(opts.IO.Out, "Opening %s in your browser.\n", issue.HTMLURL)
		return browser.Open(issue.HTMLURL)
	}

	// Output
	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "%s #%s\n", cs.Bold(issue.Title), issue.Number)
	fmt.Fprintf(opts.IO.Out, "  State: %s\n", issue.State)
	if issue.User != nil {
		fmt.Fprintf(opts.IO.Out, "  Author: %s\n", issue.User.Login)
	}
	fmt.Fprintf(opts.IO.Out, "  Created: %s\n", issue.CreatedAt.Format("2006-01-02 15:04"))
	if len(issue.Labels) > 0 {
		labels := make([]string, len(issue.Labels))
		for i, l := range issue.Labels {
			labels[i] = l.Name
		}
		fmt.Fprintf(opts.IO.Out, "  Labels: %s\n", strings.Join(labels, ", "))
	}
	fmt.Fprintf(opts.IO.Out, "\n")
	if issue.Body != "" {
		fmt.Fprintf(opts.IO.Out, "%s\n", issue.Body)
		fmt.Fprintf(opts.IO.Out, "\n")
	}
	fmt.Fprintf(opts.IO.Out, "  %s\n", issue.HTMLURL)
	fmt.Fprintf(opts.IO.Out, "\n")

	// Show comments if requested
	if opts.Comments && issue.Comments > 0 {
		comments, err := api.ListIssueComments(client, owner, repo, opts.Number)
		if err != nil {
			return fmt.Errorf("failed to get comments: %w", err)
		}

		fmt.Fprintf(opts.IO.Out, "--- Comments (%d) ---\n\n", len(comments))
		for _, c := range comments {
			fmt.Fprintf(opts.IO.Out, "%s at %s:\n", cs.Bold(c.User.Login), c.CreatedAt.Format("2006-01-02 15:04"))
			fmt.Fprintf(opts.IO.Out, "%s\n\n", c.Body)
		}
	}

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
