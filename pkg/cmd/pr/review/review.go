// Package review implements the pr review command
package review

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type ReviewOptions struct {
	IO              *iostreams.IOStreams
	HttpClient      func() (*http.Client, error)
	BaseRepo        func() (string, error)
	ReviewPR        func(*api.Client, string, string, int, *api.ReviewPROptions) error
	CreatePRComment func(*api.Client, string, string, int, *api.CreatePRCommentOptions) (*api.PRComment, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	Approve     bool
	Request     bool
	Comment     string
	CommentFile string
	Force       bool // Force approval (admin only)
	JSON        bool
}

// ReviewResult represents the JSON output for pr review
type ReviewResult struct {
	Number  int    `json:"number"`
	Action  string `json:"action"`
	Owner   string `json:"owner"`
	Repo    string `json:"repo"`
	URL     string `json:"url"`
	Comment string `json:"comment,omitempty"`
}

// NewCmdReview creates the review command
func NewCmdReview(f *cmdutil.Factory, runF func(*ReviewOptions) error) *cobra.Command {
	opts := &ReviewOptions{
		IO:              f.IOStreams,
		HttpClient:      f.HttpClient,
		BaseRepo:        f.BaseRepo,
		ReviewPR:        api.ReviewPR,
		CreatePRComment: api.CreatePRComment,
	}

	cmd := &cobra.Command{
		Use:   "review <number>",
		Short: "Review a pull request",
		Long: heredoc.Doc(`
			Review a pull request in a GitCode repository.

			You can approve or comment on a PR. GitCode's current API does not
			support "request changes", so --request returns a clear error.

			Note: --approve requires GitCode's "approval permission", which is
			separate from the "merge permission" used by 'gc pr merge'. Users
			with merge permission cannot approve PRs without explicit approval
			permission granted by the repository administrators. If you receive
			a 403 Forbidden error, use --comment to leave review feedback instead.
		`),
		Example: heredoc.Doc(`
			# Approve a PR
			$ gc pr review 123 -R owner/repo --approve

			# Comment on a PR
			$ gc pr review 123 -R owner/repo --comment "Looks good to me"

			# Comment from a file
			$ gc pr review 123 -R owner/repo --comment-file review-notes.md

			# Comment from stdin
			$ gc pr review 123 -R owner/repo --comment-file -

			# Approve a PR and leave a comment
			$ gc pr review 123 -R owner/repo --approve --comment "LGTM"

			# Approve a PR with comment from file
			$ gc pr review 123 -R owner/repo --approve --comment-file self-check.md

			# Force approve a PR (admin only)
			$ gc pr review 123 -R owner/repo --approve --force

			# Output as JSON
			$ gc pr review 123 -R owner/repo --approve --json
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
			return reviewRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVarP(&opts.Approve, "approve", "a", false, "Approve the PR")
	cmd.Flags().BoolVarP(&opts.Request, "request", "r", false, "Request changes on the PR (currently unsupported by GitCode API)")
	cmd.Flags().StringVarP(&opts.Comment, "comment", "c", "", "Comment body")
	cmd.Flags().StringVarP(&opts.CommentFile, "comment-file", "F", "", "Read comment from file (use - for stdin)")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Force approval (admin only)")
	cmd.Flags().BoolVar(&opts.JSON, "json", false, "Output as JSON")

	return cmd
}

func reviewRun(opts *ReviewOptions) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}
	client, err := cmdutil.AuthenticatedClient(httpClient)
	if err != nil {
		return err
	}

	// Read comment from file if specified
	if opts.CommentFile != "" {
		if opts.Comment != "" {
			return cmdutil.NewUsageError("--comment and --comment-file are mutually exclusive")
		}
		if opts.CommentFile == "-" {
			comment, err := cmdutil.ReadTextFromFlag(opts.IO.In, "--comment-file")
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			opts.Comment = strings.TrimSpace(comment)
		} else {
			comment, err := cmdutil.ReadTextFile(opts.CommentFile)
			if err != nil {
				return fmt.Errorf("failed to read comment file: %w", err)
			}
			opts.Comment = strings.TrimSpace(comment)
		}
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

	prURL := fmt.Sprintf("https://gitcode.com/%s/%s/pulls/%d", owner, repo, opts.Number)

	// Handle force approval (admin only)
	if opts.Force {
		if !opts.Approve {
			return cmdutil.NewUsageError("--force can only be used with --approve")
		}
		err := opts.ReviewPR(client, owner, repo, opts.Number, &api.ReviewPROptions{
			Force: true,
		})
		if err != nil {
			return fmt.Errorf("failed to force approve PR: %w", err)
		}

		result := ReviewResult{
			Number: opts.Number,
			Action: "force_approved",
			Owner:  owner,
			Repo:   repo,
			URL:    prURL,
		}

		if opts.JSON {
			return cmdutil.WriteJSON(opts.IO.Out, result)
		}
		fmt.Fprintf(opts.IO.Out, "%s %s PR #%d\n", cs.Green("✓"), cs.Green("force approved"), opts.Number)
		return nil
	}

	// Handle comment only (GitCode uses /comments endpoint, not /reviews)
	if opts.Comment != "" && !opts.Approve && !opts.Request {
		comment, err := opts.CreatePRComment(client, owner, repo, opts.Number, &api.CreatePRCommentOptions{
			Body: opts.Comment,
		})
		if err != nil {
			return fmt.Errorf("failed to comment on PR: %w", err)
		}

		result := ReviewResult{
			Number:  opts.Number,
			Action:  "commented",
			Owner:   owner,
			Repo:    repo,
			URL:     fmt.Sprintf("%s#comment_%s", prURL, cmdutil.FormatAPIID(comment.ID)),
			Comment: opts.Comment,
		}

		if opts.JSON {
			return cmdutil.WriteJSON(opts.IO.Out, result)
		}
		fmt.Fprintf(opts.IO.Out, "%s Commented on PR #%d\n", cs.Green("✓"), opts.Number)
		if comment.Body != "" {
			fmt.Fprintf(opts.IO.Out, "  %s\n", comment.Body)
		}
		return nil
	}

	if opts.Request {
		return cmdutil.NewUsageError("requesting changes is not supported by the current GitCode API. Use --comment to leave review feedback")
	}

	if opts.Approve {
		if opts.Comment != "" {
			if _, err := opts.CreatePRComment(client, owner, repo, opts.Number, &api.CreatePRCommentOptions{
				Body: opts.Comment,
			}); err != nil {
				return fmt.Errorf("failed to comment on PR before approval: %w", err)
			}
		}

		if err := opts.ReviewPR(client, owner, repo, opts.Number, &api.ReviewPROptions{}); err != nil {
			return fmt.Errorf("failed to approve PR: %w", err)
		}

		result := ReviewResult{
			Number:  opts.Number,
			Action:  "approved",
			Owner:   owner,
			Repo:    repo,
			URL:     prURL,
			Comment: opts.Comment,
		}

		if opts.JSON {
			return cmdutil.WriteJSON(opts.IO.Out, result)
		}

		fmt.Fprintf(opts.IO.Out, "%s %s PR #%d\n", cs.Green("✓"), cs.Green("approved"), opts.Number)
		if opts.Comment != "" {
			fmt.Fprintf(opts.IO.Out, "  %s\n", opts.Comment)
		}
		return nil
	}

	return cmdutil.NewUsageError("no review action specified. Use --comment, --approve, or --force with --approve")
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
