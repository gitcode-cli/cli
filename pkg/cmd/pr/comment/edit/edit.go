// Package edit implements the PR comment edit command
package edit

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

type EditOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	Repository string
	ID         int
	Body       string
}

// NewCmdEdit creates the edit command
func NewCmdEdit(f *cmdutil.Factory, runF func(*EditOptions) error) *cobra.Command {
	opts := &EditOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "edit <comment-id>",
		Short: "Edit a PR comment",
		Long: heredoc.Doc(`
			Edit an existing pull request comment.
		`),
		Example: heredoc.Doc(`
			# Edit a comment
			$ gc pr comment edit 123 --body "Updated text" -R owner/repo
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return cmdutil.NewUsageError(fmt.Sprintf("invalid comment ID: %s", args[0]))
			}
			opts.ID = id

			if runF != nil {
				return runF(opts)
			}
			return editRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "New comment body (required)")
	_ = cmd.MarkFlagRequired("body")

	return cmd
}

func editRun(opts *EditOptions) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}
	client, err := cmdutil.AuthenticatedClient(httpClient)
	if err != nil {
		return err
	}

	owner, repo, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	editOpts := &api.EditPRCommentOptions{
		Body: opts.Body,
	}
	comment, err := api.EditPRComment(client, owner, repo, opts.ID, editOpts)
	if err != nil {
		return fmt.Errorf("failed to edit comment: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Edited comment #%s\n", cs.Green("✓"), cmdutil.FormatAPIID(comment.ID))
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
