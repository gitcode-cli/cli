// Package list implements the actions runner-group list command.
package list

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/output"
)

// ListOptions configures the actions runner-group list command.
type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	Org string

	Keyword    string
	Limit      int
	Page       int
	Paginate   bool
	PerPage    int
	LimitSet   bool
	PerPageSet bool

	JSON   bool
	Format string
}

// NewCmdList creates the actions runner-group list command.
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List organization runner groups",
		Long: heredoc.Doc(`
			List runner groups in a GitCode organization.

			Filters are applied server-side via the Actions v8 API. Use --json
			for machine-readable output.
		`),
		Example: heredoc.Doc(`
			# List runner groups in an organization
			$ gc actions runner-group list --org my-org

			# Filter by keyword
			$ gc actions runner-group list --org my-org --keyword prod

			# Fetch all pages
			$ gc actions runner-group list --org my-org --paginate --per-page 100

			# Output as JSON
			$ gc actions runner-group list --org my-org --json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			opts.LimitSet = cmd.Flags().Changed("limit")
			opts.PerPageSet = cmd.Flags().Changed("per-page")
			return listRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Org, "org", "", "Organization path (required)")
	cmd.MarkFlagRequired("org")
	cmd.Flags().StringVar(&opts.Keyword, "keyword", "", "Filter by keyword")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of runner groups to list")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number to fetch")
	cmd.Flags().BoolVar(&opts.Paginate, "paginate", false, "Fetch all pages")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "API page size (default: --limit, or 100 with --paginate)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)
	cmdutil.AddFormatFlag(cmd, &opts.Format)

	return cmd
}

func listRun(opts *ListOptions) error {
	format, err := resolveOutputFormat(opts.JSON, opts.Format)
	if err != nil {
		return err
	}

	if opts.Org == "" {
		return cmdutil.NewUsageError("--org is required")
	}
	if opts.Limit <= 0 {
		return cmdutil.NewUsageError("--limit must be greater than 0")
	}
	if opts.Page < 0 {
		return cmdutil.NewUsageError("--page must be greater than or equal to 0")
	}
	if opts.PerPage < 0 {
		return cmdutil.NewUsageError("--per-page must be greater than or equal to 0")
	}
	if opts.Paginate && opts.Page > 0 {
		return cmdutil.NewUsageError("--paginate cannot be combined with --page")
	}

	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}

	groups, err := listRunnerGroups(client, opts.Org, opts)
	if err != nil {
		return fmt.Errorf("failed to list runner groups: %w", err)
	}

	if len(groups) == 0 {
		if format == output.FormatJSON {
			return cmdutil.WriteJSON(opts.IO.Out, groups)
		}
		fmt.Fprintf(opts.IO.Out, "No runner groups found\n")
		return nil
	}

	if format == output.FormatJSON {
		return cmdutil.WriteJSON(opts.IO.Out, groups)
	}

	printRunnerGroups(opts.IO, groups)
	return nil
}

func listRunnerGroups(client *api.Client, org string, opts *ListOptions) ([]api.RunnerGroup, error) {
	perPage := resolvePerPage(opts)
	if !opts.Paginate {
		resp, err := api.ListOrgRunnerGroups(client, org, &api.ListOrgRunnerGroupsOptions{
			Keyword: opts.Keyword,
			PerPage: perPage,
			Page:    opts.Page,
		})
		if err != nil {
			return nil, err
		}
		return trimGroups(resp.RunnerGroups, opts), nil
	}

	var all []api.RunnerGroup
	for page := 1; ; page++ {
		resp, err := api.ListOrgRunnerGroups(client, org, &api.ListOrgRunnerGroupsOptions{
			Keyword: opts.Keyword,
			PerPage: perPage,
			Page:    page,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, resp.RunnerGroups...)
		if opts.LimitSet && len(all) >= opts.Limit {
			return all[:opts.Limit], nil
		}
		if len(resp.RunnerGroups) < perPage {
			break
		}
	}
	if all == nil {
		return []api.RunnerGroup{}, nil
	}
	return all, nil
}

func resolvePerPage(opts *ListOptions) int {
	if opts.PerPageSet && opts.PerPage > 0 {
		return opts.PerPage
	}
	if opts.Paginate {
		return 100
	}
	return opts.Limit
}

func trimGroups(groups []api.RunnerGroup, opts *ListOptions) []api.RunnerGroup {
	if groups == nil {
		return []api.RunnerGroup{}
	}
	if opts.PerPageSet && opts.LimitSet && len(groups) > opts.Limit {
		return groups[:opts.Limit]
	}
	return groups
}

func printRunnerGroups(io *iostreams.IOStreams, groups []api.RunnerGroup) {
	cs := io.ColorScheme()
	fmt.Fprintf(io.Out, "%s\n", cs.Bold("Runner Groups"))
	for _, g := range groups {
		name := g.RunnerGroupName
		if name == "" {
			name = g.Name
		}
		share := "private"
		if g.ShareAll {
			share = "shared-all"
		}
		fmt.Fprintf(io.Out, "  %s  %s  runners=%d  %s\n",
			cs.Blue(g.ID),
			name,
			g.RunnerCount,
			share,
		)
	}
	fmt.Fprintf(io.Out, "\nTotal: %d\n", len(groups))
}

func resolveOutputFormat(jsonFlag bool, raw string) (output.Format, error) {
	format, err := output.ParseFormat(raw)
	if err != nil {
		return "", cmdutil.NewUsageError(err.Error())
	}
	if jsonFlag {
		if raw != "" && format != output.FormatJSON {
			return "", cmdutil.NewUsageError("--json cannot be combined with --format unless --format json")
		}
		return output.FormatJSON, nil
	}
	return format, nil
}
