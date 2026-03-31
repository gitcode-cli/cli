// Package comments implements the issue comments command.
package comments

import (
	"encoding/json"
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

type CommentsOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	Number     int

	Limit int
	Order string
	Since string
	JSON  bool
}

// NewCmdComments creates the issue comments command.
func NewCmdComments(f *cmdutil.Factory, runF func(*CommentsOptions) error) *cobra.Command {
	opts := &CommentsOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "comments <number>",
		Short: "List comments on an issue",
		Long: heredoc.Doc(`
			List comments on an issue in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# List issue comments
			$ gc issue comments 123 -R owner/repo

			# Limit the number of comments returned
			$ gc issue comments 123 -R owner/repo --limit 10

			# Show newest comments first
			$ gc issue comments 123 -R owner/repo --order desc

			# Output as JSON
			$ gc issue comments 123 -R owner/repo --json
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
			return commentsRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 0, "Maximum number of comments to list (0 = all)")
	cmd.Flags().StringVar(&opts.Order, "order", "asc", "Sort order (asc/desc)")
	cmd.Flags().StringVar(&opts.Since, "since", "", "Filter comments updated after this time (ISO 8601 format)")
	cmd.Flags().BoolVar(&opts.JSON, "json", false, "Output as JSON")

	return cmd
}

func commentsRun(opts *CommentsOptions) error {
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

	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}

	owner, repo, err := parseRepo(repository)
	if err != nil {
		return err
	}

	if opts.Order != "" && opts.Order != "asc" && opts.Order != "desc" {
		return fmt.Errorf("invalid order %q: must be asc or desc", opts.Order)
	}

	comments, err := api.ListIssueComments(client, owner, repo, opts.Number, &api.IssueCommentListOptions{
		PerPage: opts.Limit,
		Order:   opts.Order,
		Since:   opts.Since,
	})
	if err != nil {
		return fmt.Errorf("failed to list comments: %w", err)
	}

	if opts.JSON {
		data, err := json.MarshalIndent(comments, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal comments: %w", err)
		}
		fmt.Fprintf(opts.IO.Out, "%s\n", data)
		return nil
	}

	if len(comments) == 0 {
		fmt.Fprintf(opts.IO.Out, "No comments on issue #%d\n", opts.Number)
		return nil
	}

	fmt.Fprintf(opts.IO.Out, "\nComments on issue #%d (%d total):\n\n", opts.Number, len(comments))
	for i, comment := range comments {
		author := "unknown"
		if comment.User != nil && comment.User.Login != "" {
			author = comment.User.Login
		}

		fmt.Fprintf(opts.IO.Out, "%s) ID: %s\n", cs.Gray(fmt.Sprintf("#%d", i+1)), cmdutil.FormatAPIID(comment.ID))
		fmt.Fprintf(opts.IO.Out, "   Author: %s", cs.Bold(author))
		if !comment.CreatedAt.IsZero() {
			fmt.Fprintf(opts.IO.Out, " at %s", comment.CreatedAt.Format("2006-01-02 15:04"))
		}
		fmt.Fprintln(opts.IO.Out)

		if !comment.UpdatedAt.IsZero() && comment.UpdatedAt.Time != comment.CreatedAt.Time {
			fmt.Fprintf(opts.IO.Out, "   Updated: %s\n", comment.UpdatedAt.Format("2006-01-02 15:04"))
		}

		if comment.Body != "" {
			fmt.Fprintln(opts.IO.Out)
			for _, line := range strings.Split(comment.Body, "\n") {
				fmt.Fprintf(opts.IO.Out, "   %s\n", line)
			}
		}
		fmt.Fprintln(opts.IO.Out)
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
