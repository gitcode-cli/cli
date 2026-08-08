// Package replies implements the discussions project comments replies command (repo scope).
package replies

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/discussions/render"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type RepliesOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	Number     int
	CommentID  string

	Page    int
	PerPage int

	JSON bool
}

// NewCmdReplies creates the discussions project comments replies command (repo scope).
func NewCmdReplies(f *cmdutil.Factory, runF func(*RepliesOptions) error) *cobra.Command {
	opts := &RepliesOptions{IO: f.IOStreams, HttpClient: f.HttpClient, BaseRepo: f.BaseRepo}
	cmd := &cobra.Command{
		Use:   "replies <number> <comment-id>",
		Short: "List replies to a repository discussion comment",
		Long: heredoc.Doc(`
			List replies to a comment on a GitCode repository (project-level)
			discussion via the v5 API.
		`),
		Example: heredoc.Doc(`
			$ gc discussions project comments replies 42 <comment-id> -R owner/repo
		`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return cmdutil.NewUsageError(fmt.Sprintf("invalid discussion number: %s", args[0]))
			}
			opts.Number = n
			opts.CommentID = args[1]
			if runF != nil {
				return runF(opts)
			}
			return repliesRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number (1-based)")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "Page size (max 100, default 20)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)
	return cmd
}

func repliesRun(opts *RepliesOptions) error {
	if opts.Number < 1 {
		return cmdutil.NewUsageError("discussion number must be a positive integer")
	}
	if opts.CommentID == "" {
		return cmdutil.NewUsageError("comment-id is required")
	}
	if opts.Page < 0 {
		return cmdutil.NewUsageError("--page must be >= 0")
	}
	if opts.PerPage < 0 || opts.PerPage > 100 {
		return cmdutil.NewUsageError("--per-page must be between 0 and 100")
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

	rs, err := api.ListRepoDiscussionCommentReplies(client, owner, repo, opts.Number, opts.CommentID, &api.ListDiscussionCommentsOptions{
		Page: opts.Page, PerPage: opts.PerPage,
	})
	if err != nil {
		return cmdutil.WrapNotFound(err, "discussion #%d or comment %s not found in %s", opts.Number, opts.CommentID, repository)
	}
	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, rs)
	}
	render.PrintComments(opts.IO, rs)
	return nil
}
