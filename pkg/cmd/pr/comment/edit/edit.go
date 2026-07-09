// Package edit implements the PR comment edit command
package edit

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

type EditOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	ID         int
	Body       string
	BodyFile   string
}

// NewCmdEdit creates the edit command
func NewCmdEdit(f *cmdutil.Factory, runF func(*EditOptions) error) *cobra.Command {
	opts := &EditOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "edit <comment-id>",
		Short: "Edit a PR comment",
		Long: heredoc.Doc(`
			Edit an existing pull request comment.

			The comment body can be provided via --body flag or --body-file flag.
			Use --body-file - to read from stdin.
		`),
		Example: heredoc.Doc(`
			# Edit a comment
			$ gc pr comment edit 123 --body "Updated text" -R owner/repo

			# Edit from file
			$ gc pr comment edit 123 --body-file comment.md -R owner/repo

			# Edit from stdin
			$ echo "Updated text" | gc pr comment edit 123 --body-file - -R owner/repo
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return cmdutil.NewUsageError(fmt.Sprintf("invalid comment ID: %q, expected a numeric ID", args[0]))
			}
			if id <= 0 {
				return cmdutil.NewUsageError(fmt.Sprintf("comment ID must be a positive integer, got: %s", args[0]))
			}
			opts.ID = id

			if runF != nil {
				return runF(opts)
			}
			return editRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "New comment body")
	cmd.Flags().StringVarP(&opts.BodyFile, "body-file", "F", "", "Read comment body from file (use - for stdin)")

	return cmd
}

func editRun(opts *EditOptions) error {
	cs := opts.IO.ColorScheme()

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

	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}

	owner, repo, err := cmdutil.ParseRepo(repository)
	if err != nil {
		return err
	}

	editOpts := &api.EditPRCommentOptions{
		Body: body,
	}
	comment, err := api.EditPRComment(client, owner, repo, opts.ID, editOpts)
	if err != nil {
		return fmt.Errorf("failed to edit comment: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Edited comment #%s\n", cs.Green("✓"), cmdutil.FormatAPIID(comment.ID))
	if !comment.UpdatedAt.IsZero() {
		fmt.Fprintf(opts.IO.Out, "  Updated: %s\n", comment.UpdatedAt.Format("2006-01-02 15:04"))
	}
	return nil
}

func getBody(opts *EditOptions) (string, error) {
	body, err := cmdutil.ReadBody(opts.Body, opts.BodyFile, opts.IO.In)
	if err != nil {
		return "", err
	}
	if opts.BodyFile != "" {
		return strings.TrimSpace(body), nil
	}
	return body, nil
}
