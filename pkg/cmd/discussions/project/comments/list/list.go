// Package list implements the discussions project comments list command (repo scope).
package list

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

type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	Number     int

	Page    int
	PerPage int
	Order   string

	JSON bool
}

// NewCmdList creates the discussions project comments list command (repo scope).
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{IO: f.IOStreams, HttpClient: f.HttpClient, BaseRepo: f.BaseRepo}
	cmd := &cobra.Command{
		Use:   "list <number>",
		Short: "List repository discussion comments",
		Long: heredoc.Doc(`
			List comments on a GitCode repository (project-level) discussion via
			the v5 API (GET /api/v5/repos/{owner}/{repo}/discuss/{number}/comment).
		`),
		Example: heredoc.Doc(`
			# List comments on a repo discussion
			$ gc discussions project comments list 42 -R owner/repo
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return cmdutil.NewUsageError(fmt.Sprintf("invalid discussion number: %s", args[0]))
			}
			opts.Number = n
			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number (1-based)")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "Page size (max 100, default 20)")
	cmd.Flags().StringVar(&opts.Order, "order", "", "Sort: time_asc, time_desc, hot_desc")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)
	return cmd
}

func listRun(opts *ListOptions) error {
	if opts.Number < 1 {
		return cmdutil.NewUsageError("discussion number must be a positive integer")
	}
	if opts.Page < 0 {
		return cmdutil.NewUsageError("--page must be >= 0")
	}
	if opts.PerPage < 0 || opts.PerPage > 100 {
		return cmdutil.NewUsageError("--per-page must be between 0 and 100")
	}
	if opts.Order != "" && opts.Order != "time_asc" && opts.Order != "time_desc" && opts.Order != "hot_desc" {
		return cmdutil.NewUsageError("--order must be time_asc, time_desc, or hot_desc")
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

	cs, err := api.ListRepoDiscussionComments(client, owner, repo, opts.Number, &api.ListDiscussionCommentsOptions{
		Page: opts.Page, PerPage: opts.PerPage, Order: opts.Order,
	})
	if err != nil {
		return cmdutil.WrapNotFound(err, "discussion #%d not found in %s", opts.Number, repository)
	}
	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, cs)
	}
	render.PrintComments(opts.IO, cs)
	return nil
}
