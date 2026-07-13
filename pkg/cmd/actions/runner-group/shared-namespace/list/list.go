// Package list implements the actions runner-group shared-namespace list command.
package list

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// ListOptions configures the actions runner-group shared-namespace list command.
type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	Org           string
	RunnerGroupID string

	Keyword    string
	Limit      int
	Page       int
	Paginate   bool
	PerPage    int
	LimitSet   bool
	PerPageSet bool

	JSON bool
}

// NewCmdList creates the actions runner-group shared-namespace list command.
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "list <runner-group-id>",
		Short: "List shared namespaces for a runner group",
		Long: heredoc.Doc(`
			List namespaces that have access to an organization runner group.

			Use --json for machine-readable output.
		`),
		Example: heredoc.Doc(`
			# List shared namespaces for a runner group
			$ gc actions runner-group shared-namespace list <runner-group-id> --org my-org

			# Filter by keyword
			$ gc actions runner-group shared-namespace list <runner-group-id> --org my-org --keyword prod

			# Fetch all pages
			$ gc actions runner-group shared-namespace list <runner-group-id> --org my-org --paginate --per-page 100

			# Output as JSON
			$ gc actions runner-group shared-namespace list <runner-group-id> --org my-org --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.RunnerGroupID = strings.TrimSpace(args[0])
			if opts.RunnerGroupID == "" {
				return cmdutil.NewUsageError("runner group id is required")
			}
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
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of shared namespaces to list")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number to fetch")
	cmd.Flags().BoolVar(&opts.Paginate, "paginate", false, "Fetch all pages")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "API page size (default: --limit, or 100 with --paginate)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func listRun(opts *ListOptions) error {
	if opts.Org == "" {
		return cmdutil.NewUsageError("--org is required")
	}
	if opts.RunnerGroupID == "" {
		return cmdutil.NewUsageError("runner group id is required")
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

	namespaces, err := listSharedNamespaces(client, opts.Org, opts.RunnerGroupID, opts)
	if err != nil {
		return fmt.Errorf("failed to list shared namespaces: %w", err)
	}

	if len(namespaces) == 0 {
		if opts.JSON {
			return cmdutil.WriteJSON(opts.IO.Out, namespaces)
		}
		fmt.Fprintf(opts.IO.Out, "No shared namespaces found\n")
		return nil
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, namespaces)
	}

	printSharedNamespaces(opts.IO, namespaces)
	return nil
}

func listSharedNamespaces(client *api.Client, org, runnerGroupID string, opts *ListOptions) ([]api.SharedNamespace, error) {
	perPage := resolvePerPage(opts)
	if !opts.Paginate {
		resp, err := api.ListRunnerGroupSharedNamespaces(client, org, runnerGroupID, &api.ListRunnerGroupRunnersOptions{
			Keyword: opts.Keyword,
			PerPage: perPage,
			Page:    opts.Page,
		})
		if err != nil {
			return nil, err
		}
		return trimNamespaces(resp.SharedNamespaces, opts), nil
	}

	var all []api.SharedNamespace
	for page := 1; ; page++ {
		resp, err := api.ListRunnerGroupSharedNamespaces(client, org, runnerGroupID, &api.ListRunnerGroupRunnersOptions{
			Keyword: opts.Keyword,
			PerPage: perPage,
			Page:    page,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, resp.SharedNamespaces...)
		if opts.LimitSet && len(all) >= opts.Limit {
			return all[:opts.Limit], nil
		}
		if len(resp.SharedNamespaces) < perPage {
			break
		}
	}
	if all == nil {
		return []api.SharedNamespace{}, nil
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

func trimNamespaces(namespaces []api.SharedNamespace, opts *ListOptions) []api.SharedNamespace {
	if namespaces == nil {
		return []api.SharedNamespace{}
	}
	if opts.PerPageSet && opts.LimitSet && len(namespaces) > opts.Limit {
		return namespaces[:opts.Limit]
	}
	return namespaces
}

func printSharedNamespaces(io *iostreams.IOStreams, namespaces []api.SharedNamespace) {
	cs := io.ColorScheme()
	fmt.Fprintf(io.Out, "%s\n", cs.Bold("Shared Namespaces"))
	for _, ns := range namespaces {
		fmt.Fprintf(io.Out, "  %s  from=%s  to=%s  type=%s\n",
			cs.Blue(ns.ID),
			ns.FromNamespaceID,
			ns.ToNamespaceID,
			ns.Type,
		)
	}
	fmt.Fprintf(io.Out, "\nTotal: %d\n", len(namespaces))
}
