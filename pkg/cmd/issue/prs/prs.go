// Package prs implements the issue prs command
package prs

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

type PrsOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	Mode int
}

// NewCmdPrs creates the prs command
func NewCmdPrs(f *cmdutil.Factory, runF func(*PrsOptions) error) *cobra.Command {
	opts := &PrsOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "prs <number>",
		Short: "List Pull Requests associated with an issue",
		Long: heredoc.Doc(`
			List Pull Requests associated with an issue in a GitCode repository.

			Use --mode 1 to get enhanced information including mergeable status.
		`),
		Example: heredoc.Doc(`
			# List PRs for an issue
			$ gc issue prs 123 -R owner/repo

			# Get enhanced info including mergeable status
			$ gc issue prs 123 --mode 1 -R owner/repo
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
			return prsRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().IntVar(&opts.Mode, "mode", 0, "Mode: 0 (default), 1 (enhanced with mergeable status)")

	return cmd
}

func prsRun(opts *PrsOptions) error {
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

	// Get issue PRs
	prs, err := api.GetIssuePullRequests(client, owner, repo, opts.Number, opts.Mode)
	if err != nil {
		return fmt.Errorf("failed to get issue pull requests: %w", err)
	}

	if len(prs) == 0 {
		fmt.Fprintf(opts.IO.Out, "No pull requests associated with issue #%d\n", opts.Number)
		return nil
	}

	// Output
	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "Pull Requests for Issue #%d:\n\n", opts.Number)
	for _, pr := range prs {
		var stateIndicator string
		if pr.State == "open" {
			stateIndicator = cs.Green("open")
		} else if pr.State == "closed" {
			stateIndicator = cs.Red("closed")
		} else if pr.State == "merged" {
			stateIndicator = cs.Magenta("merged")
		} else {
			stateIndicator = pr.State
		}

		fmt.Fprintf(opts.IO.Out, "#%s %s %s\n", cs.Bold(strconv.Itoa(pr.Number)), stateIndicator, pr.Title)
		if pr.User != nil {
			fmt.Fprintf(opts.IO.Out, "  Author: %s\n", pr.User.Login)
		}
		if pr.Head != nil && pr.Base != nil {
			fmt.Fprintf(opts.IO.Out, "  Branch: %s -> %s\n", pr.Head.Ref, pr.Base.Ref)
		}
		// Show assignees if present
		if len(pr.Assignees) > 0 {
			assignees := make([]string, len(pr.Assignees))
			for i, a := range pr.Assignees {
				assignees[i] = a.Login
			}
			fmt.Fprintf(opts.IO.Out, "  Assignees: %s\n", strings.Join(assignees, ", "))
		}
		// Show mergeable status in enhanced mode
		if opts.Mode == 1 {
			mergeable := "unknown"
			if pr.CanMergeCheck {
				mergeable = cs.Green("can merge")
			} else {
				mergeable = cs.Yellow("cannot merge")
			}
			fmt.Fprintf(opts.IO.Out, "  Mergeable: %s\n", mergeable)
		}
		fmt.Fprintf(opts.IO.Out, "  %s\n", pr.HTMLURL)
		fmt.Fprintf(opts.IO.Out, "\n")
	}

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
