// Package comment implements the issue comment command
package comment

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/gitcode-com/gitcode-cli/api"
	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
	"github.com/gitcode-com/gitcode-cli/pkg/iostreams"
)

type CommentOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	Body string
}

// NewCmdComment creates the comment command
func NewCmdComment(f *cmdutil.Factory, runF func(*CommentOptions) error) *cobra.Command {
	opts := &CommentOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "comment <number>",
		Short: "Add a comment to an issue",
		Long: heredoc.Doc(`
			Add a comment to an issue in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Add a comment
			$ gc issue comment 123 --body "This is a comment"

			# Add comment from stdin
			$ gc issue comment 123 --body-file -
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
			return commentRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "Comment body")

	return cmd
}

func commentRun(opts *CommentOptions) error {
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

	// Validate body
	if opts.Body == "" {
		return fmt.Errorf("comment body is required. Use --body flag")
	}

	// Create comment
	comment, err := api.CreateIssueComment(client, owner, repo, opts.Number, &api.CreateCommentOptions{
		Body: opts.Body,
	})
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Added comment to issue #%d\n", cs.Green("✓"), opts.Number)
	fmt.Fprintf(opts.IO.Out, "  Comment ID: %v\n", comment.ID)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	if repo == "" {
		// TODO: get from current git repo
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