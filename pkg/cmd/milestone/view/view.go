// Package view implements the milestone view command
package view

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/browser"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	Web    bool
	JSON   bool
	Issues bool // Show associated issues (default: true)
}

// MilestoneWithIssues combines milestone data with associated issues for JSON output
type MilestoneWithIssues struct {
	Number       int         `json:"number"`
	Title        string      `json:"title"`
	State        string      `json:"state"`
	DueOn        string      `json:"due_on,omitempty"`
	Description  string      `json:"description,omitempty"`
	URL          string      `json:"url"`
	Issues       []api.Issue `json:"issues,omitempty"`
	TotalIssues  int         `json:"total_issues"`
	ClosedIssues int         `json:"closed_issues"`
	OpenIssues   int         `json:"open_issues"`
}

// NewCmdView creates the view command
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "view <number>",
		Short: "View a milestone",
		Long: heredoc.Doc(`
			View a milestone in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# View a milestone
			$ gc milestone view 1 -R owner/repo

			# View a milestone as JSON
			$ gc milestone view 1 -R owner/repo --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return cmdutil.NewUsageError(fmt.Sprintf("invalid milestone number: %s", args[0]))
			}
			opts.Number = number

			if runF != nil {
				return runF(opts)
			}
			return viewRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVarP(&opts.Web, "web", "w", false, "Open in browser")
	cmd.Flags().BoolVar(&opts.Issues, "issues", true, "Show associated issues")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func viewRun(opts *ViewOptions) error {
	cs := opts.IO.ColorScheme()

	if opts.JSON && opts.Web {
		return cmdutil.NewUsageError("cannot use --json with --web")
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

	// Get milestone
	ms, err := api.GetMilestone(client, owner, repo, opts.Number)
	if err != nil {
		return cmdutil.WrapNotFound(err, "milestone #%d not found in %s/%s", opts.Number, owner, repo)
	}

	milestoneURL := fmt.Sprintf("https://gitcode.com/%s/%s/milestones/%d", owner, repo, ms.Number)

	// Open in browser if --web flag is set
	if opts.Web {
		fmt.Fprintf(opts.IO.Out, "Opening %s in your browser.\n", milestoneURL)
		if err := browser.Open(milestoneURL); err != nil {
			if opts.IO.IsStdoutTTY() {
				return err
			}
			fmt.Fprintf(opts.IO.ErrOut, "Failed to open browser: %v\n", err)
		}
		return nil
	}

	// Fetch issues if requested
	var issues []api.Issue
	if opts.Issues {
		issues, err = api.ListRepoIssuesAll(client, owner, repo, &api.IssueListOptions{
			Milestone: strconv.Itoa(opts.Number),
			State:     "all",
		})
		if err != nil {
			return fmt.Errorf("failed to list issues: %w", err)
		}
	}

	// Calculate counts
	totalIssues := len(issues)
	closedIssues := countIssuesByState(issues, "closed")
	openIssues := totalIssues - closedIssues

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, &MilestoneWithIssues{
			Number:       ms.Number,
			Title:        ms.Title,
			State:        ms.State,
			DueOn:        ms.DueOn,
			Description:  ms.Description,
			URL:          milestoneURL,
			Issues:       issues,
			TotalIssues:  totalIssues,
			ClosedIssues: closedIssues,
			OpenIssues:   openIssues,
		})
	}

	// Output milestone metadata
	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "%s #%d\n", cs.Bold(ms.Title), ms.Number)
	fmt.Fprintf(opts.IO.Out, "  State: %s\n", ms.State)
	if ms.DueOn != "" {
		fmt.Fprintf(opts.IO.Out, "  Due: %s\n", ms.DueOn)
	}
	fmt.Fprintf(opts.IO.Out, "\n")
	if ms.Description != "" {
		fmt.Fprintf(opts.IO.Out, "%s\n", ms.Description)
		fmt.Fprintf(opts.IO.Out, "\n")
	}

	// Output issues section
	if opts.Issues && totalIssues > 0 {
		fmt.Fprintf(opts.IO.Out, "Issues (%d total, %d closed, %d open)\n\n",
			totalIssues, closedIssues, openIssues)

		// Print closed issues first
		if closedIssues > 0 {
			fmt.Fprintf(opts.IO.Out, "Closed:\n")
			for _, issue := range issues {
				if strings.EqualFold(issue.State, "closed") {
					fmt.Fprintf(opts.IO.Out, "  #%s %s\n", issue.Number, issue.Title)
				}
			}
			fmt.Fprintf(opts.IO.Out, "\n")
		}

		// Print open issues
		if openIssues > 0 {
			fmt.Fprintf(opts.IO.Out, "Open:\n")
			for _, issue := range issues {
				if strings.EqualFold(issue.State, "open") {
					fmt.Fprintf(opts.IO.Out, "  #%s %s\n", issue.Number, issue.Title)
				}
			}
			fmt.Fprintf(opts.IO.Out, "\n")
		}
	} else if opts.Issues && totalIssues == 0 {
		fmt.Fprintf(opts.IO.Out, "Issues: None\n\n")
	}

	fmt.Fprintf(opts.IO.Out, "  %s\n", milestoneURL)
	fmt.Fprintf(opts.IO.Out, "\n")

	return nil
}

// countIssuesByState counts issues matching the given state
func countIssuesByState(issues []api.Issue, state string) int {
	count := 0
	for _, issue := range issues {
		if strings.EqualFold(issue.State, state) {
			count++
		}
	}
	return count
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
