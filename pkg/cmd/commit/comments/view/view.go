// Package view implements the commit comments view command
package view

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	ID         string
	JSON       bool
}

// NewCmdView creates the view command
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "View a commit comment",
		Long: heredoc.Doc(`
			View a specific commit comment.
		`),
		Example: heredoc.Doc(`
			# View a comment
			$ gc commit comments view 123 -R owner/repo

			# View a comment as JSON
			$ gc commit comments view 123 -R owner/repo --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ID = args[0]

			if runF != nil {
				return runF(opts)
			}
			return viewRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func viewRun(opts *ViewOptions) error {
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
	owner, repo, err := parseRepo(repository)
	if err != nil {
		return err
	}

	comment, err := api.GetCommitComment(client, owner, repo, opts.ID)
	if err != nil {
		return cmdutil.WrapNotFound(err, "commit comment %s not found in %s/%s", opts.ID, owner, repo)
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, comment)
	}

	fmt.Fprintf(opts.IO.Out, "\n%s #%s\n", cs.Bold("Comment:"), cmdutil.FormatAPIID(comment.ID))
	if comment.User != nil {
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Author:"), comment.User.Login)
	}
	fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Created:"), comment.CreatedAt)
	fmt.Fprintf(opts.IO.Out, "\n%s\n", comment.Body)
	fmt.Fprintf(opts.IO.Out, "\n")

	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
