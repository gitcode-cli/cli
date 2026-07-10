// Package list implements the actions artifact list command.
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

// ListOptions configures the actions artifact list command.
type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	RunID      string

	Name       string
	Sort       string
	Direction  string
	Limit      int
	Page       int
	Paginate   bool
	PerPage    int
	LimitSet   bool
	PerPageSet bool

	JSON   bool
	Format string
}

// NewCmdList creates the actions artifact list command.
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List repository artifacts",
		Long: heredoc.Doc(`
			List workflow run artifacts for a GitCode repository.

			Filters are applied server-side via the Actions v8 API. Use --json for
			machine-readable output.
		`),
		Example: heredoc.Doc(`
			# List repository artifacts
			$ gc actions artifact list -R owner/repo

			# List artifacts of a specific run
			$ gc actions artifact list -R owner/repo --run <run-id>

			# Filter by name (fuzzy)
			$ gc actions artifact list -R owner/repo --name build

			# Sort by creation time
			$ gc actions artifact list -R owner/repo --sort created --direction desc

			# Fetch all pages
			$ gc actions artifact list -R owner/repo --paginate --per-page 100

			# Render as a table
			$ gc actions artifact list -R owner/repo --format table

			# Output as JSON
			$ gc actions artifact list -R owner/repo --json
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
	cmd.Flags().StringVar(&opts.RunID, "run", "", "Filter by workflow run id (list run-scoped artifacts)")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Filter by artifact name (fuzzy)")
	cmd.Flags().StringVar(&opts.Sort, "sort", "", "Sort by created")
	cmdutil.SetFlagEnum(cmd, "sort", "created")
	cmd.Flags().StringVar(&opts.Direction, "direction", "", "Sort direction (asc/desc)")
	cmdutil.SetFlagEnum(cmd, "direction", "asc", "desc")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of artifacts to list")
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

	artifacts, err := listArtifacts(client, owner, repo, opts)
	if err != nil {
		return fmt.Errorf("failed to list artifacts: %w", err)
	}

	if len(artifacts) == 0 {
		if format == output.FormatJSON {
			return cmdutil.WriteJSON(opts.IO.Out, artifacts)
		}
		fmt.Fprintf(opts.IO.Out, "No artifacts found\n")
		return nil
	}

	if format == output.FormatJSON {
		return cmdutil.WriteJSON(opts.IO.Out, artifacts)
	}

	printer, err := output.NewArtifactListPrinter(output.ArtifactListOptions{
		Format: format,
		Color:  opts.IO.ColorScheme(),
	})
	if err != nil {
		return err
	}
	return printer.Print(opts.IO.Out, artifacts)
}

func listArtifacts(client *api.Client, owner, repo string, opts *ListOptions) ([]api.Artifact, error) {
	perPage := resolvePerPage(opts)
	if !opts.Paginate {
		resp, err := fetchArtifactsPage(client, owner, repo, opts, opts.Page)
		if err != nil {
			return nil, err
		}
		return trimArtifacts(resp.Artifacts, opts), nil
	}

	var all []api.Artifact
	for page := 1; ; page++ {
		resp, err := fetchArtifactsPage(client, owner, repo, opts, page)
		if err != nil {
			return nil, err
		}
		all = append(all, resp.Artifacts...)
		if opts.LimitSet && len(all) >= opts.Limit {
			return all[:opts.Limit], nil
		}
		if len(resp.Artifacts) < perPage {
			break
		}
	}
	if all == nil {
		all = []api.Artifact{}
	}
	return all, nil
}

// fetchArtifactsPage fetches one page of artifacts from either the repository
// endpoint (default) or the run-scoped endpoint (when --run is set).
func fetchArtifactsPage(client *api.Client, owner, repo string, opts *ListOptions, page int) (*api.ArtifactsResponse, error) {
	o := &api.ActionsListArtifactsOptions{
		Name:      opts.Name,
		Sort:      opts.Sort,
		Direction: opts.Direction,
		PerPage:   resolvePerPage(opts),
		Page:      page,
	}
	if opts.RunID != "" {
		return api.ListActionsRunArtifacts(client, owner, repo, opts.RunID, o)
	}
	return api.ListActionsArtifacts(client, owner, repo, o)
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

func trimArtifacts(artifacts []api.Artifact, opts *ListOptions) []api.Artifact {
	if artifacts == nil {
		return []api.Artifact{}
	}
	if opts.PerPageSet && opts.LimitSet && len(artifacts) > opts.Limit {
		return artifacts[:opts.Limit]
	}
	return artifacts
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
