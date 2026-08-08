// Package list implements the discussions project list command.
package list

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/discussions/render"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// ListOptions configures the discussions project list command.
type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string

	Limit     int
	LimitSet  bool
	Page      int
	PerPage   int
	Sort      string
	Direction string
	Search    string

	JSON bool
}

// NewCmdList creates the discussions project list command.
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List repository discussions",
		Long: heredoc.Doc(`
			List discussions (discuss) in a GitCode repository via the v5 API
			(GET /api/v5/repos/{owner}/{repo}/discuss). Filters are applied
			server-side. Use --json for machine-readable output.
		`),
		Example: heredoc.Doc(`
			# List discussions in a repository
			$ gc discussions project list -R owner/repo

			# Search by title/description
			$ gc discussions project list -R owner/repo --search "release plan"

			# Sort by comment count, descending
			$ gc discussions project list -R owner/repo --sort comment_size --direction desc

			# Output as JSON
			$ gc discussions project list -R owner/repo --json
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			opts.LimitSet = cmd.Flags().Changed("limit")
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of discussions to list")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number (1-based)")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "Page size (max 100, default 20)")
	cmd.Flags().StringVar(&opts.Sort, "sort", "", "Sort field: created (default) or comment_size")
	cmd.Flags().StringVar(&opts.Direction, "direction", "", "Sort direction: asc or desc (default desc)")
	cmd.Flags().StringVar(&opts.Search, "search", "", "Filter by title/description")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func listRun(opts *ListOptions) error {
	if opts.Limit <= 0 {
		return cmdutil.NewUsageError("--limit must be greater than 0")
	}
	if opts.Page < 0 {
		return cmdutil.NewUsageError("--page must be greater than or equal to 0")
	}
	if opts.PerPage < 0 || opts.PerPage > 100 {
		return cmdutil.NewUsageError("--per-page must be between 0 and 100")
	}
	if opts.Sort != "" && opts.Sort != "created" && opts.Sort != "comment_size" {
		return cmdutil.NewUsageError("--sort must be created or comment_size")
	}
	if opts.Direction != "" && opts.Direction != "asc" && opts.Direction != "desc" {
		return cmdutil.NewUsageError("--direction must be asc or desc")
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

	discussions, err := api.ListRepoDiscussions(client, owner, repo, &api.ListOrgDiscussionsOptions{
		Page:      opts.Page,
		PerPage:   opts.PerPage,
		Sort:      opts.Sort,
		Direction: opts.Direction,
		Search:    opts.Search,
	})
	if err != nil {
		return err
	}
	discussions = trimDiscussions(discussions, opts)

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, discussions)
	}

	render.PrintDiscussions(opts.IO, discussions)
	return nil
}

// trimDiscussions caps the result to --limit. Discussions --per-page defaults
// to the API default (20), so the client-side cap is applied after fetch.
func trimDiscussions(discussions []*api.Discussion, opts *ListOptions) []*api.Discussion {
	if len(discussions) > opts.Limit {
		return discussions[:opts.Limit]
	}
	return discussions
}
