// Package delete implements the PR comment delete command
package delete

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

type DeleteOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	ID         int
	Yes        bool
}

// NewCmdDelete creates the delete command
func NewCmdDelete(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "delete <comment-id>",
		Short: "Delete a PR comment",
		Long: heredoc.Doc(`
			Delete a pull request comment.

				Non-interactive mode: Requires --yes to skip confirmation.
		`),
		Example: heredoc.Doc(`
			# Delete a comment
			$ gc pr comment delete 123 -R owner/repo

			# Delete without confirmation
			$ gc pr comment delete 123 -R owner/repo --yes
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
			return deleteRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVarP(&opts.Yes, "yes", "y", false, "Skip confirmation")

	return cmd
}

func deleteRun(opts *DeleteOptions) error {
	cs := opts.IO.ColorScheme()

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

	if err := cmdutil.ConfirmOrAbort(cmdutil.ConfirmOptions{
		IO:       opts.IO,
		Yes:      opts.Yes,
		Expected: strconv.Itoa(opts.ID),
		Prompt:   fmt.Sprintf("! This will delete PR comment #%d\nType the comment ID to confirm: ", opts.ID),
	}); err != nil {
		return err
	}

	if err := api.DeletePRComment(client, owner, repo, opts.ID); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Deleted comment #%d\n", cs.Red("✗"), opts.ID)
	return nil
}
