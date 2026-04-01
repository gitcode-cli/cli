// Package relations implements the issue relations command.
package relations

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type RelationsOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	State      string
	Limit      int
	Mode       int
	JSON       bool
}

type RelationRow struct {
	Issue RelationIssue `json:"issue"`
	PR    RelationPR    `json:"pr"`
}

type RelationIssue struct {
	Number  string `json:"number"`
	Title   string `json:"title"`
	State   string `json:"state"`
	HTMLURL string `json:"html_url"`
}

type RelationPR struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	State   string `json:"state"`
	HTMLURL string `json:"html_url"`
}

// NewCmdRelations creates the issue relations command.
func NewCmdRelations(f *cmdutil.Factory, runF func(*RelationsOptions) error) *cobra.Command {
	opts := &RelationsOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "relations",
		Short: "List issue and pull request relations in a repository",
		Long: heredoc.Doc(`
			List the issue and pull request relation table for a repository.

			The command walks repository issues, fetches the pull requests linked to
			each issue, and outputs relation rows that include both PR and issue state.
		`),
		Example: heredoc.Doc(`
			# List all issue/PR relations in a repository
			$ gc issue relations -R owner/repo

			# Output machine-readable relation rows
			$ gc issue relations -R owner/repo --json

			# Only inspect open issues
			$ gc issue relations -R owner/repo --state open --limit 50
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return relationsRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.State, "state", "s", "all", "Filter source issues by state (open/closed/all)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 100, "Maximum number of issues to inspect")
	cmd.Flags().IntVar(&opts.Mode, "mode", 1, "Issue PR lookup mode: 0 (default), 1 (enhanced with mergeable status)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func relationsRun(opts *RelationsOptions) error {
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

	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}

	owner, repo, err := cmdutil.ParseRepo(repository)
	if err != nil {
		return err
	}

	if opts.State != "open" && opts.State != "closed" && opts.State != "all" {
		return cmdutil.NewUsageError("invalid state: must be open, closed, or all")
	}
	if opts.Limit <= 0 {
		return cmdutil.NewUsageError("limit must be greater than 0")
	}
	if opts.Mode != 0 && opts.Mode != 1 {
		return cmdutil.NewUsageError("mode must be 0 or 1")
	}

	issues, err := api.ListRepoIssues(client, owner, repo, &api.IssueListOptions{
		State:   opts.State,
		PerPage: opts.Limit,
	})
	if err != nil {
		return fmt.Errorf("failed to list issues: %w", err)
	}

	rows := make([]RelationRow, 0)
	seen := make(map[string]struct{})
	for _, issue := range issues {
		number, err := issueNumberToInt(issue.Number)
		if err != nil {
			continue
		}

		prs, err := api.GetIssuePullRequests(client, owner, repo, number, opts.Mode)
		if err != nil {
			return fmt.Errorf("failed to get pull requests for issue #%s: %w", issue.Number, err)
		}

		for _, pr := range prs {
			key := fmt.Sprintf("%d:%s", pr.Number, issue.Number)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}

			rows = append(rows, RelationRow{
				Issue: RelationIssue{
					Number:  issue.Number,
					Title:   issue.Title,
					State:   issue.State,
					HTMLURL: issue.HTMLURL,
				},
				PR: RelationPR{
					Number:  pr.Number,
					Title:   pr.Title,
					State:   pr.State,
					HTMLURL: pr.HTMLURL,
				},
			})
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].PR.Number == rows[j].PR.Number {
			return rows[i].Issue.Number < rows[j].Issue.Number
		}
		return rows[i].PR.Number < rows[j].PR.Number
	})

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, rows)
	}

	if len(rows) == 0 {
		fmt.Fprintf(opts.IO.Out, "No issue and pull request relations found\n")
		return nil
	}

	currentPR := -1
	for _, row := range rows {
		if row.PR.Number != currentPR {
			if currentPR != -1 {
				fmt.Fprintln(opts.IO.Out)
			}
			currentPR = row.PR.Number
			fmt.Fprintf(opts.IO.Out, "PR #%d [%s] %s\n", row.PR.Number, row.PR.State, row.PR.Title)
			fmt.Fprintf(opts.IO.Out, "  %s\n", row.PR.HTMLURL)
		}
		fmt.Fprintf(opts.IO.Out, "  Issue #%s [%s] %s\n", row.Issue.Number, row.Issue.State, row.Issue.Title)
		fmt.Fprintf(opts.IO.Out, "    %s\n", row.Issue.HTMLURL)
	}

	return nil
}

func issueNumberToInt(number string) (int, error) {
	var value int
	for _, c := range number {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid issue number %q", number)
		}
		value = value*10 + int(c-'0')
	}
	if value <= 0 {
		return 0, fmt.Errorf("invalid issue number %q", number)
	}
	return value, nil
}
