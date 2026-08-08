// Package list implements the discussions list command.
package list

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// ListOptions configures the discussions list command.
type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	Org string

	Page      int
	PerPage   int
	Sort      string
	Direction string
	Search    string

	JSON bool
}

// NewCmdList creates the discussions list command.
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List organization discussions",
		Long: heredoc.Doc(`
			List discussions (discuss) in a GitCode organization via the v5 API
			(GET /api/v5/orgs/{org}/discuss). Filters are applied server-side.
			Use --json for machine-readable output.
		`),
		Example: heredoc.Doc(`
			# List discussions in an organization
			$ gc discussions list --org my-org

			# Search by title/description
			$ gc discussions list --org my-org --search "release plan"

			# Sort by comment count, descending
			$ gc discussions list --org my-org --sort comment_size --direction desc

			# Output as JSON
			$ gc discussions list --org my-org --json
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Org, "org", "", "Organization path (required)")
	cmd.MarkFlagRequired("org")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number (1-based)")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "Page size (max 100, default 20)")
	cmd.Flags().StringVar(&opts.Sort, "sort", "", "Sort field: created (default) or comment_size")
	cmd.Flags().StringVar(&opts.Direction, "direction", "", "Sort direction: asc or desc (default desc)")
	cmd.Flags().StringVar(&opts.Search, "search", "", "Filter by title/description")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func listRun(opts *ListOptions) error {
	if opts.Org == "" {
		return cmdutil.NewUsageError("--org is required")
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

	discussions, err := api.ListOrgDiscussions(client, opts.Org, &api.ListOrgDiscussionsOptions{
		Page:      opts.Page,
		PerPage:   opts.PerPage,
		Sort:      opts.Sort,
		Direction: opts.Direction,
		Search:    opts.Search,
	})
	if err != nil {
		return err
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, discussions)
	}

	printDiscussions(opts.IO, discussions)
	return nil
}

func printDiscussions(io *iostreams.IOStreams, discussions []*api.Discussion) {
	cs := io.ColorScheme()
	if len(discussions) == 0 {
		fmt.Fprintf(io.Out, "No discussions found\n")
		return
	}
	fmt.Fprintf(io.Out, "%s\n", cs.Bold("Discussions"))
	for _, d := range discussions {
		state := "open"
		if d.IsClosed != 0 {
			state = "closed"
		}
		author := ""
		if d.Author != nil {
			author = d.Author.Login
		}
		pin := ""
		if d.IsPin != 0 {
			pin = cs.Yellow("📌 ") // pinned
		}
		fmt.Fprintf(io.Out, "  %s#%d  %s  (%s)  comments=%d  %s\n",
			pin,
			d.Number,
			d.Title,
			author,
			d.CommentTotal,
			cs.Cyan(state),
		)
		if d.UpdatedAt != "" {
			fmt.Fprintf(io.Out, "       updated %s\n", d.UpdatedAt)
		}
	}
	fmt.Fprintf(io.Out, "\nTotal: %d\n", len(discussions))
}
