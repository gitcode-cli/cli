// Package comment implements the PR comment command
package comment

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmtDelete "gitcode.com/gitcode-cli/cli/pkg/cmd/pr/comment/delete"
	cmtEdit "gitcode.com/gitcode-cli/cli/pkg/cmd/pr/comment/edit"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	cmtResolve "gitcode.com/gitcode-cli/cli/pkg/cmd/pr/comment/resolve"
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
	Path     string // File path for inline comment
	Position int    // Line number in the new version of the file (new side of the diff) for inline comment
	JSON     bool
}

// CommentResult represents the JSON output for pr comment
type CommentResult struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Author    string `json:"author,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Path      string `json:"path,omitempty"`
	Position  int    `json:"position,omitempty"`
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
		Short: "Add a comment to a pull request",
		Long: heredoc.Doc(`
			Add a comment to a pull request in a GitCode repository.

			The comment body can be provided via --body flag or --body-file flag.
			Use --body-file - to read from stdin.

			For inline comments on specific lines, use --path and --position flags.
			--path specifies the file path (e.g., "api/auth.go").
			--position specifies the line number in the new version of the file
			(the new/right side of the diff), not an offset within the diff hunk.
			Use the new-side line number shown by "gc pr diff".
		`),
		Example: heredoc.Doc(`
			# Add a general comment
			$ gc pr comment 123 --body "This looks good" -R owner/repo

			# Add comment from file
			$ gc pr comment 123 --body-file comment.txt -R owner/repo

			# Add comment from stdin
			$ echo "Comment from stdin" | gc pr comment 123 --body-file - -R owner/repo

			# Add inline comment on a specific file and line.
			# Run the diff first and read the new-side line number of the target line:
			$ gc pr diff 123 -R owner/repo
			# Here 42 is the line number in the new version of api/auth.go:
			$ gc pr comment 123 --body "Consider renaming this" --path api/auth.go --position 42 -R owner/repo

			# Output as JSON
			$ gc pr comment 123 --body "This looks good" --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return cmdutil.NewUsageError(fmt.Sprintf("invalid pull request number: %s", args[0]))
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
	cmd.Flags().StringVar(&opts.Path, "path", "", "File path for inline comment")
	cmd.Flags().IntVar(&opts.Position, "position", 0, "Line number in the new version of the file (new side of the diff) for inline comment")
	cmd.Flags().BoolVar(&opts.JSON, "json", false, "Output as JSON")

cmd.AddCommand(cmtEdit.NewCmdEdit(f, nil))
	cmd.AddCommand(cmtDelete.NewCmdDelete(f, nil))
	cmd.AddCommand(cmtResolve.NewCmdResolve(f, nil))
	cmd.AddCommand(cmtResolve.NewCmdUnresolve(f, nil))

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

	// Validate inline comment flags
	if opts.Path != "" && opts.Position == 0 {
		return cmdutil.NewUsageError("--position is required when --path is specified for inline comments")
	}
	if opts.Path != "" && opts.Position < 1 {
		return cmdutil.NewUsageError("--position must be greater than 0")
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

	owner, repo, err := cmdutil.ParseRepo(repository)
	if err != nil {
		return err
	}

	// Create comment
	commentOpts := &api.CreatePRCommentOptions{
		Body: body,
	}
	if opts.Path != "" {
		commentOpts.Path = opts.Path
		commentOpts.Position = opts.Position
	}

	comment, err := api.CreatePRComment(client, owner, repo, opts.Number, commentOpts)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	result := CommentResult{
		ID:   cmdutil.FormatAPIID(comment.ID),
		URL:  fmt.Sprintf("https://gitcode.com/%s/%s/pulls/%d#comment_%s", owner, repo, opts.Number, cmdutil.FormatAPIID(comment.ID)),
		Body: body,
	}
	if comment.User != nil {
		result.Author = comment.User.Login
	}
	if !comment.CreatedAt.IsZero() {
		result.CreatedAt = comment.CreatedAt.Format("2006-01-02 15:04:05")
	}
	if comment.DiffFile != "" {
		result.Path = comment.DiffFile
		// Handle DiffPosition as interface{} - convert to int if possible
		if pos, ok := comment.DiffPosition.(int); ok {
			result.Position = pos
		} else if pos, ok := comment.DiffPosition.(float64); ok {
			result.Position = int(pos)
		}
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, result)
	}

	// Output
	fmt.Fprintf(opts.IO.Out, "%s Added comment to pull request #%d\n", cs.Green("✓"), opts.Number)
	fmt.Fprintf(opts.IO.Out, "  ID: %s\n", result.ID)
	if result.Author != "" {
		fmt.Fprintf(opts.IO.Out, "  Author: %s\n", result.Author)
	}
	if result.CreatedAt != "" {
		fmt.Fprintf(opts.IO.Out, "  Created: %s\n", result.CreatedAt)
	}
	if result.Path != "" {
		fmt.Fprintf(opts.IO.Out, "  File: %s (position %d)\n", result.Path, result.Position)
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
