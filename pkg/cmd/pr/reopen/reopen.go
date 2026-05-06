// Package reopen implements the pr reopen command
package reopen

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type ReopenOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	Comment string
	Yes     bool
	JSON    bool
}

// ReopenResult represents the JSON output for pr reopen
type ReopenResult struct {
	Number  int    `json:"number"`
	State   string `json:"state"`
	Owner   string `json:"owner"`
	Repo    string `json:"repo"`
	URL     string `json:"url"`
	Comment string `json:"comment,omitempty"`
}

// NewCmdReopen creates the reopen command
func NewCmdReopen(f *cmdutil.Factory, runF func(*ReopenOptions) error) *cobra.Command {
	opts := &ReopenOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "reopen <number>",
		Short: "Reopen a closed pull request",
		Long: heredoc.Doc(`
			Reopen a closed pull request in a GitCode repository.

				Non-interactive mode: Requires --yes to skip confirmation.
		`),
		Example: heredoc.Doc(`
			# Reopen a PR
			$ gc pr reopen 123 -R owner/repo

			# Reopen with a comment
			$ gc pr reopen 123 -R owner/repo --comment "Reopening for further changes"

			# Non-interactive mode
			$ gc pr reopen 123 -R owner/repo --yes

			# Output as JSON
			$ gc pr reopen 123 -R owner/repo --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return cmdutil.NewUsageError(fmt.Sprintf("invalid PR number: %s", args[0]))
			}
			opts.Number = number

			if runF != nil {
				return runF(opts)
			}
			return reopenRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Comment, "comment", "c", "", "Add a comment before reopening")
	cmd.Flags().BoolVar(&opts.Yes, "yes", false, "Skip confirmation prompt")
	cmd.Flags().BoolVar(&opts.JSON, "json", false, "Output as JSON")

	return cmd
}

func reopenRun(opts *ReopenOptions) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}
	client, err := cmdutil.AuthenticatedClient(httpClient)
	if err != nil {
		return err
	}

	// Get repository
	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}
	owner, repo, err := parseRepo(repository)
	if err != nil {
		return err
	}

	if err := cmdutil.ConfirmOrAbort(cmdutil.ConfirmOptions{
		IO:       opts.IO,
		Yes:      opts.Yes,
		Expected: strconv.Itoa(opts.Number),
		Prompt:   fmt.Sprintf("! This will reopen PR #%d in %s/%s\nType the PR number to confirm: ", opts.Number, owner, repo),
	}); err != nil {
		return err
	}

	// Add comment if provided
	if opts.Comment != "" {
		_, err := api.CreatePRComment(client, owner, repo, opts.Number, &api.CreatePRCommentOptions{
			Body: opts.Comment,
		})
		if err != nil {
			return fmt.Errorf("failed to add comment: %w", err)
		}
	}

	// Reopen PR
	pr, err := api.ReopenPullRequest(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to reopen PR: %w", err)
	}

	result := ReopenResult{
		Number: opts.Number,
		State:  pr.State,
		Owner:  owner,
		Repo:   repo,
		URL:    cmdutil.ResolvePRURL(pr.HTMLURL, owner, repo, opts.Number),
	}
	if opts.Comment != "" {
		result.Comment = opts.Comment
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, result)
	}

	fmt.Fprintf(opts.IO.Out, "%s Reopened PR #%d in %s/%s\n", cs.Green("✓"), opts.Number, owner, repo)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
