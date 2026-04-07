// Package edit implements the pr edit command
package edit

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type EditOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	Title             string
	Body              string
	BodyFile          string
	Base              string
	Draft             string // "true", "false", or "" (not specified)
	Labels            []string
	Milestone         int
	CloseRelatedIssue string // "true", "false", or "" (not specified)
}

// NewCmdEdit creates the edit command
func NewCmdEdit(f *cmdutil.Factory, runF func(*EditOptions) error) *cobra.Command {
	opts := &EditOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "edit <number>",
		Short: "Edit a pull request",
		Long: heredoc.Doc(`
			Edit a pull request in a GitCode repository.

			You can update the title, body, base branch, draft status, labels, and milestone.
		`),
		Example: heredoc.Doc(`
			# Edit PR title
			$ gc pr edit 123 -R owner/repo --title "New title"

			# Edit PR body
			$ gc pr edit 123 -R owner/repo --body "New description"

			# Edit PR body from file
			$ gc pr edit 123 -R owner/repo --body-file description.md

			# Mark PR as ready for review
			$ gc pr edit 123 -R owner/repo --draft false

			# Mark PR as draft
			$ gc pr edit 123 -R owner/repo --draft true

			# Add labels
			$ gc pr edit 123 -R owner/repo --labels bug,enhancement

			# Set milestone
			$ gc pr edit 123 -R owner/repo --milestone 5
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
			return editRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "New title")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "New body")
	cmd.Flags().StringVarP(&opts.BodyFile, "body-file", "F", "", "Read body from file")
	cmd.Flags().StringVar(&opts.Base, "base", "", "New base branch")
	cmd.Flags().StringVar(&opts.Draft, "draft", "", "Mark as draft (true/false)")
	cmd.Flags().StringSliceVarP(&opts.Labels, "labels", "l", nil, "Add labels (comma-separated)")
	cmd.Flags().IntVarP(&opts.Milestone, "milestone", "m", 0, "Set milestone by number")
	cmd.Flags().StringVar(&opts.CloseRelatedIssue, "close-related-issue", "", "Close related issues when merged (true/false)")

	return cmd
}

func editRun(opts *EditOptions) error {
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

	// Build update options
	updateOpts := &api.UpdatePROptions{}

	if opts.Title != "" {
		updateOpts.Title = opts.Title
	}
	if opts.Body != "" {
		updateOpts.Body = opts.Body
	}
	if opts.BodyFile != "" {
		data, err := os.ReadFile(opts.BodyFile)
		if err != nil {
			return fmt.Errorf("failed to read body file: %w", err)
		}
		updateOpts.Body = string(data)
	}
	if opts.Base != "" {
		updateOpts.Base = opts.Base
	}
	if opts.Draft != "" {
		val := opts.Draft == "true"
		updateOpts.Draft = &val
	}
	if len(opts.Labels) > 0 {
		updateOpts.Labels = opts.Labels
	}
	if opts.Milestone > 0 {
		updateOpts.MilestoneNumber = opts.Milestone
	}
	if opts.CloseRelatedIssue != "" {
		val := opts.CloseRelatedIssue == "true"
		updateOpts.CloseRelatedIssue = &val
	}

	// Check if there's anything to update
	if updateOpts.Title == "" && updateOpts.Body == "" && opts.BodyFile == "" &&
		updateOpts.Base == "" && opts.Draft == "" &&
		len(updateOpts.Labels) == 0 && updateOpts.MilestoneNumber == 0 &&
		opts.CloseRelatedIssue == "" {
		return fmt.Errorf("no changes specified. Use flags to specify what to edit")
	}

	// Edit PR
	pr, err := api.EditPR(client, owner, repo, opts.Number, updateOpts)
	if err != nil {
		return fmt.Errorf("failed to edit PR: %w", err)
	}

	prNumber := pr.Number
	if prNumber == 0 {
		prNumber = opts.Number
	}
	fmt.Fprintf(opts.IO.Out, "%s Updated PR #%d: %s\n", cs.Green("✓"), prNumber, pr.Title)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	if token := os.Getenv("GITCODE_TOKEN"); token != "" {
		return token
	}
	return cmdutil.EnvToken()
}
