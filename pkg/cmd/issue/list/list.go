// Package list implements the issue list command
package list

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
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
	Labels        string
	Assignee      string
	Milestone     string
	Creator       string
	Sort          string
	Direction     string
	Since         string
	CreatedAfter  string
	CreatedBefore string
	UpdatedAfter  string
	UpdatedBefore string
	Search        string
	Page          int
	JSON          bool
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
		Short: "List issues in a repository",
		Long: heredoc.Doc(`
			List issues in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# List open issues
			$ gc issue list

			# List closed issues
			$ gc issue list --state closed

			# List issues with specific labels
			$ gc issue list --label bug,enhancement

			# List issues in a specific repository
			$ gc issue list -R owner/repo

			# Filter by milestone
			$ gc issue list --milestone "v1.0"

			# Filter by assignee
			$ gc issue list --assignee username

			# Filter by creator
			$ gc issue list --creator username

			# Sort by updated time
			$ gc issue list --sort updated --direction asc

			# Filter by creation time
			$ gc issue list --created-after "2024-01-01"

			# Search by keyword
			$ gc issue list --search "bug"

			# Combine multiple filters
			$ gc issue list --state open --milestone "v1.0" --assignee username --sort updated
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.State, "state", "s", "open", "Filter by state (open/closed/all)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of issues to list")
	cmd.Flags().StringVarP(&opts.Labels, "label", "l", "", "Filter by labels (comma separated)")
	cmd.Flags().StringVarP(&opts.Assignee, "assignee", "a", "", "Filter by assignee username")
	cmd.Flags().StringVarP(&opts.Milestone, "milestone", "m", "", "Filter by milestone title")
	cmd.Flags().StringVar(&opts.Creator, "creator", "", "Filter by creator username")
	cmd.Flags().StringVar(&opts.Sort, "sort", "created", "Sort by (created/updated)")
	cmd.Flags().StringVar(&opts.Direction, "direction", "desc", "Sort direction (asc/desc)")
	cmd.Flags().StringVar(&opts.Since, "since", "", "Filter by update time (ISO 8601 format)")
	cmd.Flags().StringVar(&opts.CreatedAfter, "created-after", "", "Filter issues created after this time")
	cmd.Flags().StringVar(&opts.CreatedBefore, "created-before", "", "Filter issues created before this time")
	cmd.Flags().StringVar(&opts.UpdatedAfter, "updated-after", "", "Filter issues updated after this time")
	cmd.Flags().StringVar(&opts.UpdatedBefore, "updated-before", "", "Filter issues updated before this time")
	cmd.Flags().StringVar(&opts.Search, "search", "", "Search by keyword in title or body")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number for pagination")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func listRun(opts *ListOptions) error {
	cs := opts.IO.ColorScheme()

	if err := normalizeDateFilters(opts); err != nil {
		return err
	}

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := cmdutil.EnvToken()
	if token == "" {
		return cmdutil.NewAuthError("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "environment")

	// Get repository
	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}

	owner, repo, err := parseRepo(repository)
	if err != nil {
		return err
	}

	// List issues
	issues, err := api.ListRepoIssues(client, owner, repo, &api.IssueListOptions{
		State:         opts.State,
		Labels:        opts.Labels,
		PerPage:       opts.Limit,
		Page:          opts.Page,
		Milestone:     opts.Milestone,
		Assignee:      opts.Assignee,
		Creator:       opts.Creator,
		Sort:          opts.Sort,
		Direction:     opts.Direction,
		Since:         opts.Since,
		CreatedAfter:  opts.CreatedAfter,
		CreatedBefore: opts.CreatedBefore,
		UpdatedAfter:  opts.UpdatedAfter,
		UpdatedBefore: opts.UpdatedBefore,
		Search:        opts.Search,
	})
	if err != nil {
		return fmt.Errorf("failed to list issues: %w", err)
	}

	// Output
	if len(issues) == 0 {
		if opts.JSON {
			return cmdutil.WriteJSON(opts.IO.Out, issues)
		}
		fmt.Fprintf(opts.IO.Out, "No issues found\n")
		return nil
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, issues)
	}

	// Calculate max number width for alignment
	maxNumWidth := 0
	for _, issue := range issues {
		w := len(fmt.Sprintf("#%d", issue.Number))
		if w > maxNumWidth {
			maxNumWidth = w
		}
	}

	fmt.Fprintf(opts.IO.Out, "\n")
	for _, issue := range issues {
		state := "open"
		if issue.State == "closed" {
			state = cs.Red("closed")
		} else {
			state = cs.Green("open")
		}
		fmt.Fprintf(opts.IO.Out, "%-*s  %s  %s\n", maxNumWidth, fmt.Sprintf("#%d", issue.Number), state, issue.Title)
	}
	fmt.Fprintf(opts.IO.Out, "\n")

	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

func normalizeDateFilters(opts *ListOptions) error {
	if opts == nil {
		return nil
	}

	var err error
	if opts.Since, err = normalizeIssueListTime(opts.Since, false); err != nil {
		return fmt.Errorf("invalid --since: %w", err)
	}
	if opts.CreatedAfter, err = normalizeIssueListTime(opts.CreatedAfter, false); err != nil {
		return fmt.Errorf("invalid --created-after: %w", err)
	}
	if opts.CreatedBefore, err = normalizeIssueListTime(opts.CreatedBefore, true); err != nil {
		return fmt.Errorf("invalid --created-before: %w", err)
	}
	if opts.UpdatedAfter, err = normalizeIssueListTime(opts.UpdatedAfter, false); err != nil {
		return fmt.Errorf("invalid --updated-after: %w", err)
	}
	if opts.UpdatedBefore, err = normalizeIssueListTime(opts.UpdatedBefore, true); err != nil {
		return fmt.Errorf("invalid --updated-before: %w", err)
	}
	return nil
}

func normalizeIssueListTime(value string, endOfDay bool) (string, error) {
	if value == "" {
		return "", nil
	}

	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}
	for _, format := range formats {
		if ts, err := time.Parse(format, value); err == nil {
			return ts.Format(time.RFC3339), nil
		}
	}

	if ts, err := time.Parse("2006-01-02", value); err == nil {
		if endOfDay {
			ts = ts.Add(24*time.Hour - time.Second)
		}
		return ts.Format(time.RFC3339), nil
	}

	return "", fmt.Errorf("expected YYYY-MM-DD or ISO 8601 datetime")
}
