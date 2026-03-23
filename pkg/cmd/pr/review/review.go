// Package review implements the pr review command
package review

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

type ReviewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	Approve  bool
	Request  bool
	Comment  string
	Force    bool // Force approval (admin only)
}

// NewCmdReview creates the review command
func NewCmdReview(f *cmdutil.Factory, runF func(*ReviewOptions) error) *cobra.Command {
	opts := &ReviewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "review <number>",
		Short: "Review a pull request",
		Long: heredoc.Doc(`
			Review a pull request in a GitCode repository.

			You can approve, request changes, or comment on a PR.
		`),
		Example: heredoc.Doc(`
			# Approve a PR
			$ gc pr review 123 --approve

			# Request changes
			$ gc pr review 123 --request

			# Comment on a PR
			$ gc pr review 123 --comment "Looks good to me"

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
	cmd.Flags().BoolVarP(&opts.Request, "request", "r", false, "Request changes on the PR")
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
		err := api.ReviewPR(client, owner, repo, opts.Number, &api.ReviewPROptions{
			Force: true,
		})
		if err != nil {
			return fmt.Errorf("failed to force approve PR: %w", err)
		}
		fmt.Fprintf(opts.IO.Out, "%s %s PR #%d\n", cs.Green("✓"), cs.Green("force approved"), opts.Number)
		return nil
	}

	// Determine review event
	event := "COMMENT"
	if opts.Approve {
		event = "APPROVE"
	} else if opts.Request {
		event = "REQUEST_CHANGES"
	}

	// Create review
	review, err := api.CreatePRReview(client, owner, repo, opts.Number, &api.CreatePRReviewOptions{
		Body:  opts.Comment,
		Event: event,
	})
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	// Output
	action := "commented on"
	if opts.Approve {
		action = cs.Green("approved")
	} else if opts.Request {
		action = cs.Yellow("requested changes on")
	}

	fmt.Fprintf(opts.IO.Out, "%s %s PR #%d\n", cs.Green("✓"), action, opts.Number)
	if review.Body != "" {
		fmt.Fprintf(opts.IO.Out, "  %s\n", review.Body)
	}
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