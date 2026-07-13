// Package list implements the actions runner list command.
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

// ListOptions configures the actions runner list command.
type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string

	Keyword    string
	Limit      int
	Page       int
	Paginate   bool
	PerPage    int
	LimitSet   bool
	PerPageSet bool

	JSON bool
}

// NewCmdList creates the actions runner list command.
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List repository host runners",
		Long: heredoc.Doc(`
			List host runners for a GitCode repository.

			Use --json for machine-readable output.
		`),
		Example: heredoc.Doc(`
			# List host runners
			$ gc actions runner list -R owner/repo

			# Filter by keyword
			$ gc actions runner list -R owner/repo --keyword prod

			# Fetch all pages
			$ gc actions runner list -R owner/repo --paginate --per-page 100

			# Output as JSON
			$ gc actions runner list -R owner/repo --json
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

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVar(&opts.Keyword, "keyword", "", "Filter by keyword")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of runners to list")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number to fetch")
	cmd.Flags().BoolVar(&opts.Paginate, "paginate", false, "Fetch all pages")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "API page size (default: --limit, or 100 with --paginate)")
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

	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}
	owner, repo, err := cmdutil.ParseRepo(repository)
	if err != nil {
		return err
	}

	runners, err := listRunners(client, owner, repo, opts)
	if err != nil {
		return fmt.Errorf("failed to list runners: %w", err)
	}

	if len(runners) == 0 {
		if opts.JSON {
			return cmdutil.WriteJSON(opts.IO.Out, runners)
		}
		fmt.Fprintf(opts.IO.Out, "No runners found\n")
		return nil
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, runners)
	}

	printRunners(opts.IO, runners)
	return nil
}

func listRunners(client *api.Client, owner, repo string, opts *ListOptions) ([]api.Runner, error) {
	perPage := resolvePerPage(opts)
	if !opts.Paginate {
		resp, err := api.ListRepoRunners(client, owner, repo, &api.ListRunnerGroupRunnersOptions{
			Keyword: opts.Keyword,
			PerPage: perPage,
			Page:    opts.Page,
		})
		if err != nil {
			return nil, err
		}
		return trimRunners(resp.Runners, opts), nil
	}

	var all []api.Runner
	for page := 1; ; page++ {
		resp, err := api.ListRepoRunners(client, owner, repo, &api.ListRunnerGroupRunnersOptions{
			Keyword: opts.Keyword,
			PerPage: perPage,
			Page:    page,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, resp.Runners...)
		if opts.LimitSet && len(all) >= opts.Limit {
			return all[:opts.Limit], nil
		}
		if len(resp.Runners) < perPage {
			break
		}
	}
	if all == nil {
		return []api.Runner{}, nil
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

func trimRunners(runners []api.Runner, opts *ListOptions) []api.Runner {
	if runners == nil {
		return []api.Runner{}
	}
	if opts.PerPageSet && opts.LimitSet && len(runners) > opts.Limit {
		return runners[:opts.Limit]
	}
	return runners
}

func printRunners(io *iostreams.IOStreams, runners []api.Runner) {
	cs := io.ColorScheme()
	fmt.Fprintf(io.Out, "%s\n", cs.Bold("Runners"))
	for _, r := range runners {
		name := r.RunnerName
		if name == "" {
			name = r.Name
		}
		labels := make([]string, 0, len(r.Labels))
		for _, l := range r.Labels {
			labels = append(labels, l.LabelName)
		}
		fmt.Fprintf(io.Out, "  %s  %s  labels=%s\n",
			cs.Blue(r.ID),
			name,
			strings.Join(labels, ","),
		)
	}
	fmt.Fprintf(io.Out, "\nTotal: %d\n", len(runners))
}
