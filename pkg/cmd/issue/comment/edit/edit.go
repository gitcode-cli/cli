// Package edit implements the issue comment edit command.
package edit

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type EditOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	ID         string
	Body       string
	BodyFile   string
}

// NewCmdEdit creates the issue comment edit command.
func NewCmdEdit(f *cmdutil.Factory, runF func(*EditOptions) error) *cobra.Command {
	opts := &EditOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "edit [id]",
		Short: "Edit an issue comment",
		Long: heredoc.Doc(`
			Edit an existing issue comment in a GitCode repository.

			The comment body can be provided via --body flag or --body-file flag.
			Use --body-file - to read from stdin.
		`),
		Example: heredoc.Doc(`
			# Edit a comment by argument
			$ gc issue comment edit 12345 -R owner/repo --body "Updated comment"

			# Edit a comment by flag
			$ gc issue comment edit --id 12345 -R owner/repo --body "Updated comment"

			# Edit from file
			$ gc issue comment edit 12345 -R owner/repo --body-file comment.md
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.ID = args[0]
			}
			if runF != nil {
				return runF(opts)
			}
			return editRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVar(&opts.ID, "id", "", "Issue comment ID")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "New comment body")
	cmd.Flags().StringVarP(&opts.BodyFile, "body-file", "F", "", "Read comment body from file (use - for stdin)")

	return cmd
}

func editRun(opts *EditOptions) error {
	cs := opts.IO.ColorScheme()

	if opts.ID == "" {
		return fmt.Errorf("comment ID is required. Use an argument or --id flag")
	}

	body, err := getBody(opts)
	if err != nil {
		return err
	}
	if body == "" {
		return fmt.Errorf("comment body is required. Use --body or --body-file flag")
	}

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

	owner, repo, err := cmdutil.ParseRepo(repository)
	if err != nil {
		return err
	}

	comment, err := api.UpdateIssueComment(client, owner, repo, opts.ID, &api.UpdateCommentOptions{
		Body: body,
	})
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Updated issue comment %s\n", cs.Green("✓"), opts.ID)
	if comment.UpdatedAt.IsZero() {
		return nil
	}
	fmt.Fprintf(opts.IO.Out, "  Updated: %s\n", comment.UpdatedAt.Format("2006-01-02 15:04"))
	return nil
}

func getBody(opts *EditOptions) (string, error) {
	if opts.Body != "" && opts.BodyFile != "" {
		return "", fmt.Errorf("cannot use both --body and --body-file")
	}

	if opts.Body != "" {
		return opts.Body, nil
	}

	if opts.BodyFile != "" {
		if opts.BodyFile == "-" {
			reader := bufio.NewReader(opts.IO.In)
			var sb strings.Builder
			for {
				line, err := reader.ReadString('\n')
				if err != nil && err != io.EOF {
					return "", fmt.Errorf("failed to read from stdin: %w", err)
				}
				sb.WriteString(line)
				if err == io.EOF {
					break
				}
			}
			return strings.TrimSpace(sb.String()), nil
		}

		content, err := os.ReadFile(opts.BodyFile)
		if err != nil {
			return "", fmt.Errorf("failed to read file %s: %w", opts.BodyFile, err)
		}
		return strings.TrimSpace(string(content)), nil
	}

	return "", nil
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITCODE_TOKEN")
}
