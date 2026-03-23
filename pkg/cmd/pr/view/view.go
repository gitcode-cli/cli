// Package view implements the pr view command
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
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	Web      bool
	Comments bool
}

// NewCmdView creates the view command
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "view <number>",
		Short: "View a pull request",
		Long: heredoc.Doc(`
			View a pull request in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# View a PR
			$ gc pr view 123

			# View PR in browser
			$ gc pr view 123 --web
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid PR number: %s", args[0])
			}
			opts.Number = number

			if runF != nil {
				return runF(opts)
			}
			return viewRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVarP(&opts.Web, "web", "w", false, "Open in browser")
	cmd.Flags().BoolVarP(&opts.Comments, "comments", "c", false, "View comments")

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
	owner, repo, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	// Get PR
	pr, err := api.GetPullRequest(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to get PR: %w", err)
	}

	// Output
	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "%s #%d\n", cs.Bold(pr.Title), pr.Number)
	if pr.Draft {
		fmt.Fprintf(opts.IO.Out, "  State: %s (draft)\n", pr.State)
	} else if pr.Merged {
		fmt.Fprintf(opts.IO.Out, "  State: merged\n")
	} else {
		fmt.Fprintf(opts.IO.Out, "  State: %s\n", pr.State)
	}
	if pr.User != nil {
		fmt.Fprintf(opts.IO.Out, "  Author: %s\n", pr.User.Login)
	}
	fmt.Fprintf(opts.IO.Out, "  Branch: %s -> %s\n", pr.Head.Ref, pr.Base.Ref)
	fmt.Fprintf(opts.IO.Out, "  Created: %s\n", pr.CreatedAt.Format("2006-01-02 15:04"))
	fmt.Fprintf(opts.IO.Out, "  Additions: +%d  Deletions: -%d\n", pr.Additions, pr.Deletions)
	fmt.Fprintf(opts.IO.Out, "\n")
	if pr.Body != "" {
		fmt.Fprintf(opts.IO.Out, "%s\n", pr.Body)
		fmt.Fprintf(opts.IO.Out, "\n")
	}
	fmt.Fprintf(opts.IO.Out, "  %s\n", pr.HTMLURL)
	fmt.Fprintf(opts.IO.Out, "\n")

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