// Package list implements the pr list command
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
	"gitcode.com/gitcode-cli/cli/pkg/output"
)

type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	Repository string

	// Flags
	State         string
	Limit         int
	Head          string
	Base          string
	Sort          string
	Direction     string
	Page          int
	Paginate      bool
	PerPage       int
	LimitSet      bool
	PerPageSet    bool
	JSON          bool
	Format        string
	Milestone     string
	CommitMessage string
}

// NewCmdList creates the list command
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List pull requests",
		Long: heredoc.Doc(`
			List pull requests in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# List open PRs
			$ gc pr list -R owner/repo

			# List closed PRs
			$ gc pr list -R owner/repo --state closed

			# Filter by head and base branches
			$ gc pr list -R owner/repo --head feature/login --base main

			# Sort results
			$ gc pr list -R owner/repo --sort updated --direction desc

			# Fetch all pages
			$ gc pr list -R owner/repo --paginate --per-page 100

			# Filter by commit message text
			$ gc pr list -R owner/repo --commit-message "fix login"

			# Render as a table
			$ gc pr list -R owner/repo --format table

			# Output as JSON
			$ gc pr list -R owner/repo --json
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
	cmd.Flags().StringVarP(&opts.State, "state", "s", "open", "Filter by state (open/closed/merged/all)")
	cmdutil.SetFlagEnum(cmd, "state", "open", "closed", "merged", "all")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of PRs to list")
	cmd.Flags().StringVarP(&opts.Head, "head", "H", "", "Filter by head branch")
	cmd.Flags().StringVarP(&opts.Base, "base", "B", "", "Filter by base branch")
	cmd.Flags().StringVarP(&opts.Milestone, "milestone", "m", "", "Filter by milestone title")
	cmd.Flags().StringVar(&opts.Sort, "sort", "", "Sort by created/updated/popularity/long-running")
	cmdutil.SetFlagEnum(cmd, "sort", "created", "updated", "popularity", "long-running")
	cmd.Flags().StringVar(&opts.Direction, "direction", "", "Sort direction (asc/desc)")
	cmdutil.SetFlagEnum(cmd, "direction", "asc", "desc")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number to fetch")
	cmd.Flags().BoolVar(&opts.Paginate, "paginate", false, "Fetch all pages")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "API page size (default: --limit, or 100 with --paginate)")
	cmd.Flags().StringVar(&opts.CommitMessage, "commit-message", "", "Filter PRs whose commit messages contain text")
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

	// Get repository
	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}
	owner, repo, err := parseRepo(repository)
	if err != nil {
		return err
	}

	// Validate pagination parameters
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

	// List PRs
	prs, err := listPullRequests(client, owner, repo, opts)
	if err != nil {
		return fmt.Errorf("failed to list PRs: %w", err)
	}
	if strings.TrimSpace(opts.CommitMessage) != "" {
		prs, err = filterByCommitMessage(client, owner, repo, prs, opts.CommitMessage)
		if err != nil {
			return err
		}
	}

	// Output
	if len(prs) == 0 {
		if format == output.FormatJSON {
			return cmdutil.WriteJSON(opts.IO.Out, prs)
		}
		fmt.Fprintf(opts.IO.Out, "No pull requests found\n")
		return nil
	}

	if format == output.FormatJSON {
		return cmdutil.WriteJSON(opts.IO.Out, prs)
	}
	printer, err := output.NewPRListPrinter(output.PRListOptions{
		Format: format,
		Color:  opts.IO.ColorScheme(),
	})
	if err != nil {
		return err
	}
	return printer.Print(opts.IO.Out, prs)
}

func listPullRequests(client *api.Client, owner, repo string, opts *ListOptions) ([]api.PullRequest, error) {
	perPage := resolvePerPage(opts)
	if !opts.Paginate {
		prs, err := api.ListPullRequests(client, owner, repo, &api.PRListOptions{
			State:     opts.State,
			Head:      opts.Head,
			Base:      opts.Base,
			Sort:      opts.Sort,
			Direction: opts.Direction,
			PerPage:   perPage,
			Page:      opts.Page,
			Milestone: opts.Milestone,
		})
		if err != nil {
			return nil, err
		}
		return trimPRs(prs, opts), nil
	}

	var all []api.PullRequest
	for page := 1; ; page++ {
		prs, err := api.ListPullRequests(client, owner, repo, &api.PRListOptions{
			State:     opts.State,
			Head:      opts.Head,
			Base:      opts.Base,
			Sort:      opts.Sort,
			Direction: opts.Direction,
			PerPage:   perPage,
			Page:      page,
			Milestone: opts.Milestone,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, prs...)
		if opts.LimitSet && len(all) >= opts.Limit {
			return all[:opts.Limit], nil
		}
		if len(prs) < perPage {
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

func trimPRs(prs []api.PullRequest, opts *ListOptions) []api.PullRequest {
	if opts.PerPageSet && opts.LimitSet && len(prs) > opts.Limit {
		return prs[:opts.Limit]
	}
	return prs
}

func filterByCommitMessage(client *api.Client, owner, repo string, prs []api.PullRequest, query string) ([]api.PullRequest, error) {
	needle := strings.ToLower(strings.TrimSpace(query))
	var filtered []api.PullRequest
	for _, pr := range prs {
		commits, err := api.ListPRCommits(client, owner, repo, pr.Number)
		if err != nil {
			return nil, fmt.Errorf("failed to list commits for PR #%d: %w", pr.Number, err)
		}
		for _, commit := range commits {
			if strings.Contains(strings.ToLower(commit.MessageText()), needle) {
				filtered = append(filtered, pr)
				break
			}
		}
	}
	return filtered, nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
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
