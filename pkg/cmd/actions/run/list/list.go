// Package list implements the actions run list command.
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

// ListOptions configures the actions run list command.
type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string

	Status        string
	Event         string
	Branch        string
	Executor      string
	WorkflowID    string
	WorkflowName  string
	PullRequestID string

	Limit      int
	Page       int
	Paginate   bool
	PerPage    int
	LimitSet   bool
	PerPageSet bool

	JSON   bool
	Format string
}

// NewCmdList creates the actions run list command.
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List pipeline runs",
		Long: heredoc.Doc(`
			List pipeline (workflow) run records for a GitCode repository.

			Filters are applied server-side via the Actions v8 API. Use --json for
			machine-readable output.
		`),
		Example: heredoc.Doc(`
			# List recent pipeline runs
			$ gc actions run list -R owner/repo

			# Filter by status and event
			$ gc actions run list -R owner/repo --status FAILED --event Push

			# Filter by branch and executor
			$ gc actions run list -R owner/repo --branch main --executor dev

			# Filter by workflow name or id
			$ gc actions run list -R owner/repo --workflow "CI"
			$ gc actions run list -R owner/repo --workflow-id wf-1

			# Fetch all pages
			$ gc actions run list -R owner/repo --paginate --per-page 100

			# Render as a table
			$ gc actions run list -R owner/repo --format table

			# Output as JSON
			$ gc actions run list -R owner/repo --json
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
	cmd.Flags().StringVar(&opts.Status, "status", "", "Filter by run status (COMPLETED/RUNNING/FAILED/CANCELED/IGNORED/PAUSED/SUSPEND)")
	cmdutil.SetFlagEnum(cmd, "status", "COMPLETED", "RUNNING", "FAILED", "CANCELED", "IGNORED", "PAUSED", "SUSPEND")
	cmd.Flags().StringVar(&opts.Event, "event", "", "Filter by trigger event (MR/Push/Manual)")
	cmdutil.SetFlagEnum(cmd, "event", "MR", "Push", "Manual")
	cmd.Flags().StringVar(&opts.Branch, "branch", "", "Filter by branch")
	cmd.Flags().StringVar(&opts.Executor, "executor", "", "Filter by executor username")
	cmd.Flags().StringVar(&opts.WorkflowID, "workflow-id", "", "Filter by workflow id")
	cmd.Flags().StringVar(&opts.WorkflowName, "workflow", "", "Filter by workflow name")
	cmd.Flags().StringVar(&opts.PullRequestID, "pr", "", "Filter by PR number")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of runs to list")
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

	runs, err := listRuns(client, owner, repo, opts)
	if err != nil {
		return fmt.Errorf("failed to list pipeline runs: %w", err)
	}

	if len(runs) == 0 {
		if format == output.FormatJSON {
			return cmdutil.WriteJSON(opts.IO.Out, runs)
		}
		fmt.Fprintf(opts.IO.Out, "No pipeline runs found\n")
		return nil
	}

	if format == output.FormatJSON {
		return cmdutil.WriteJSON(opts.IO.Out, runs)
	}

	printer, err := output.NewWorkflowRunListPrinter(output.WorkflowRunListOptions{
		Format: format,
		Color:  opts.IO.ColorScheme(),
	})
	if err != nil {
		return err
	}
	return printer.Print(opts.IO.Out, runs)
}

func listRuns(client *api.Client, owner, repo string, opts *ListOptions) ([]api.WorkflowRun, error) {
	perPage := resolvePerPage(opts)
	if !opts.Paginate {
		resp, err := api.ListActionsRuns(client, owner, repo, &api.ActionsListRunsOptions{
			Status:        opts.Status,
			Event:         opts.Event,
			Branch:        opts.Branch,
			Executor:      opts.Executor,
			PullRequestID: opts.PullRequestID,
			WorkflowID:    opts.WorkflowID,
			WorkflowName:  opts.WorkflowName,
			PerPage:       perPage,
			Page:          opts.Page,
		})
		if err != nil {
			return nil, err
		}
		return trimRuns(resp.WorkflowRuns, opts), nil
	}

	var all []api.WorkflowRun
	for page := 1; ; page++ {
		resp, err := api.ListActionsRuns(client, owner, repo, &api.ActionsListRunsOptions{
			Status:        opts.Status,
			Event:         opts.Event,
			Branch:        opts.Branch,
			Executor:      opts.Executor,
			PullRequestID: opts.PullRequestID,
			WorkflowID:    opts.WorkflowID,
			WorkflowName:  opts.WorkflowName,
			PerPage:       perPage,
			Page:          page,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, resp.WorkflowRuns...)
		if opts.LimitSet && len(all) >= opts.Limit {
			return all[:opts.Limit], nil
		}
		if len(resp.WorkflowRuns) < perPage {
			break
		}
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

func trimRuns(runs []api.WorkflowRun, opts *ListOptions) []api.WorkflowRun {
	if runs == nil {
		return []api.WorkflowRun{}
	}
	if opts.PerPageSet && opts.LimitSet && len(runs) > opts.Limit {
		return runs[:opts.Limit]
	}
	return runs
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
