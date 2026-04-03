// Package review implements the pr review command
package review

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

type ReviewOptions struct {
	IO              *iostreams.IOStreams
	HttpClient      func() (*http.Client, error)
	ReviewPR        func(*api.Client, string, string, int, *api.ReviewPROptions) error
	CreatePRComment func(*api.Client, string, string, int, *api.CreatePRCommentOptions) (*api.PRComment, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	Approve bool
	Request bool
	Comment string
	Force   bool // Force approval (admin only)
}

// NewCmdReview creates the review command
func NewCmdReview(f *cmdutil.Factory, runF func(*ReviewOptions) error) *cobra.Command {
	opts := &ReviewOptions{
		IO:              f.IOStreams,
		HttpClient:      f.HttpClient,
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
		`),
		Example: heredoc.Doc(`
				# Approve a PR
				$ gc pr review 123 --approve

				# Comment on a PR
				$ gc pr review 123 --comment "Looks good to me"

				# Approve a PR and leave a comment
				$ gc pr review 123 --approve --comment "LGTM"

				# Force approve a PR (admin only)
				$ gc pr review 123 --approve --force
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
			return reviewRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVarP(&opts.Approve, "approve", "a", false, "Approve the PR")
	cmd.Flags().BoolVarP(&opts.Request, "request", "r", false, "Request changes on the PR (currently unsupported by GitCode API)")
	cmd.Flags().StringVarP(&opts.Comment, "comment", "c", "", "Comment body")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Force approval (admin only)")

	return cmd
}

func reviewRun(opts *ReviewOptions) error {
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

	// Handle force approval (admin only)
	if opts.Force {
		if !opts.Approve {
			return fmt.Errorf("--force can only be used with --approve")
		}
		err := opts.ReviewPR(client, owner, repo, opts.Number, &api.ReviewPROptions{
			Force: true,
		})
		if err != nil {
			return fmt.Errorf("failed to force approve PR: %w", err)
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
		fmt.Fprintf(opts.IO.Out, "%s Commented on PR #%d\n", cs.Green("✓"), opts.Number)
		if comment.Body != "" {
			fmt.Fprintf(opts.IO.Out, "  %s\n", comment.Body)
		}
		return nil
	}

	if opts.Request {
		return fmt.Errorf("requesting changes is not supported by the current GitCode API. Use --comment to leave review feedback")
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

		fmt.Fprintf(opts.IO.Out, "%s %s PR #%d\n", cs.Green("✓"), cs.Green("approved"), opts.Number)
		if opts.Comment != "" {
			fmt.Fprintf(opts.IO.Out, "  %s\n", opts.Comment)
		}
		return nil
	}

	return fmt.Errorf("no review action specified. Use --comment, --approve, or --force with --approve")
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
