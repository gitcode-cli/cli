// Package reopen implements the issue reopen command
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

// ReopenResult represents the JSON output for issue reopen
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
		Short: "Reopen a closed issue",
		Long: heredoc.Doc(`
			Reopen a closed issue in a GitCode repository.

				Non-interactive mode: Requires --yes to skip confirmation.
		`),
		Example: heredoc.Doc(`
			# Reopen an issue
			$ gc issue reopen 123

			# Reopen with a comment
			$ gc issue reopen 123 --comment "Still reproduces on latest version"

			# Reopen in a specific repository
			$ gc issue reopen 123 -R owner/repo --yes

			# Output as JSON
			$ gc issue reopen 123 --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return cmdutil.NewUsageError(fmt.Sprintf("invalid issue number: %s", args[0]))
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
	if opts.Comment != "" {
		if err := cmdutil.ScanContentForSecrets(opts.Comment); err != nil {
			return err
		}
	}
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
		Prompt:   fmt.Sprintf("! This will reopen issue #%d in %s/%s\nType the issue number to confirm: ", opts.Number, owner, repo),
	}); err != nil {
		return err
	}

	// Add comment if provided
	if opts.Comment != "" {
		_, err := api.CreateIssueComment(client, owner, repo, opts.Number, &api.CreateCommentOptions{
			Body: opts.Comment,
		})
		if err != nil {
			return fmt.Errorf("failed to add comment: %w", err)
		}
	}

	// Reopen issue
	issue, err := api.ReopenIssue(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to reopen issue: %w", err)
	}

	result := ReopenResult{
		Number: opts.Number,
		State:  issue.State,
		Owner:  owner,
		Repo:   repo,
		URL:    issue.HTMLURL,
	}
	if opts.Comment != "" {
		result.Comment = opts.Comment
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, result)
	}

	fmt.Fprintf(opts.IO.Out, "%s Reopened issue #%s in %s/%s\n", cs.Green("✓"), issue.Number, owner, repo)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
