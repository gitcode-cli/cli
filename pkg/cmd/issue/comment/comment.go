// Package comment implements the issue comment command
package comment

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/issue/comment/edit"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type CommentOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	Body     string
	BodyFile string
	JSON     bool
}

// CommentResult represents the JSON output for issue comment
type CommentResult struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Author    string `json:"author,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Body      string `json:"body"`
}

// NewCmdComment creates the comment command
func NewCmdComment(f *cmdutil.Factory, runF func(*CommentOptions) error) *cobra.Command {
	opts := &CommentOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "comment <number>",
		Short: "Add a comment to an issue",
		Long: heredoc.Doc(`
			Add a comment to an issue in a GitCode repository.

			The comment body can be provided via --body flag or --body-file flag.
			Use --body-file - to read from stdin.
		`),
		Example: heredoc.Doc(`
			# Add a comment
			$ gc issue comment 123 --body "This is a comment" -R owner/repo

			# Add comment from file
			$ gc issue comment 123 --body-file comment.txt -R owner/repo

			# Add comment from stdin
			$ echo "Comment from stdin" | gc issue comment 123 --body-file - -R owner/repo

			# Output as JSON
			$ gc issue comment 123 --body "This is a comment" --json
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
			return commentRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "Comment body")
	cmd.Flags().StringVarP(&opts.BodyFile, "body-file", "F", "", "Read comment body from file (use - for stdin)")
	cmd.Flags().BoolVar(&opts.JSON, "json", false, "Output as JSON")
	cmd.AddCommand(edit.NewCmdEdit(f, nil))

	return cmd
}

func commentRun(opts *CommentOptions) error {
	cs := opts.IO.ColorScheme()

	// Validate body input
	body, err := getBody(opts)
	if err != nil {
		return err
	}
	if body == "" {
		return cmdutil.NewUsageError("comment body is required. Use --body or --body-file flag")
	}

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

	// Create comment
	comment, err := api.CreateIssueComment(client, owner, repo, opts.Number, &api.CreateCommentOptions{
		Body: body,
	})
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	result := CommentResult{
		ID:   cmdutil.FormatAPIID(comment.ID),
		URL:  fmt.Sprintf("https://gitcode.com/%s/%s/issues/%d#comment_%s", owner, repo, opts.Number, cmdutil.FormatAPIID(comment.ID)),
		Body: body,
	}
	if comment.User != nil {
		result.Author = comment.User.Login
	}
	if !comment.CreatedAt.IsZero() {
		result.CreatedAt = comment.CreatedAt.Format("2006-01-02 15:04:05")
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, result)
	}

	// Output
	fmt.Fprintf(opts.IO.Out, "%s Added comment to issue #%d\n", cs.Green("✓"), opts.Number)
	fmt.Fprintf(opts.IO.Out, "  ID: %s\n", result.ID)
	if result.Author != "" {
		fmt.Fprintf(opts.IO.Out, "  Author: %s\n", result.Author)
	}
	if result.CreatedAt != "" {
		fmt.Fprintf(opts.IO.Out, "  Created: %s\n", result.CreatedAt)
	}
	preview := body
	if len(preview) > 100 {
		preview = preview[:100] + "..."
	}
	fmt.Fprintf(opts.IO.Out, "  Body: %s\n", preview)
	return nil
}

func getBody(opts *CommentOptions) (string, error) {
	if opts.Body != "" && opts.BodyFile != "" {
		return "", fmt.Errorf("cannot use both --body and --body-file")
	}

	if opts.Body != "" {
		return opts.Body, nil
	}

	if opts.BodyFile != "" {
		if opts.BodyFile == "-" {
			body, err := cmdutil.ReadTextFromFlag(opts.IO.In, "--body-file")
			if err != nil {
				return "", fmt.Errorf("failed to read from stdin: %w", err)
			}
			return strings.TrimSpace(body), nil
		}

		// Read from file
		content, err := cmdutil.ReadTextFile(opts.BodyFile)
		if err != nil {
			return "", fmt.Errorf("failed to read file %s: %w", opts.BodyFile, err)
		}
		return strings.TrimSpace(content), nil
	}

	return "", nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
