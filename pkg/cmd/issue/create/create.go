// Package create implements the issue create command
package create

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type CreateOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	Repository string

	// Flags
	Title     string
	Body      string
	Labels    []string
	Assignees []string
	Milestone int
	DryRun    bool
}

// NewCmdCreate creates the create command
func NewCmdCreate(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new issue",
		Long: heredoc.Doc(`
			Create a new issue in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Create an issue interactively
			$ gc issue create

			# Create an issue with title and body
			$ gc issue create --title "Bug" --body "Description"

			# Create an issue with labels
			$ gc issue create --title "Feature" --label bug,enhancement

			# Create an issue in a specific repository
			$ gc issue create -R owner/repo --title "Bug"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return createRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "Title for the issue")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "Body for the issue")
	cmd.Flags().StringSliceVarP(&opts.Labels, "label", "l", []string{}, "Labels to add")
	cmd.Flags().StringSliceVarP(&opts.Assignees, "assignee", "a", []string{}, "Assignees")
	cmd.Flags().IntVarP(&opts.Milestone, "milestone", "m", 0, "Milestone number")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Preview the create request without creating the issue")

	return cmd
}

func createRun(opts *CreateOptions) error {
	cs := opts.IO.ColorScheme()

	// Get repository
	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}

	owner, repo, err := parseRepo(repository)
	if err != nil {
		return err
	}

	// Validate title
	if opts.Title == "" {
		return cmdutil.NewUsageError("title is required. Use --title flag")
	}

	if opts.DryRun {
		fmt.Fprintf(opts.IO.Out, "Dry run: would create issue %q in %s/%s\n", opts.Title, owner, repo)
		return nil
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

	assigneeIDs, err := api.ResolveUserIDs(client, opts.Assignees)
	if err != nil {
		return fmt.Errorf("failed to resolve assignees: %w", err)
	}

	// Create issue
	issue, err := api.CreateIssue(client, owner, repo, &api.CreateIssueOptions{
		Title:       opts.Title,
		Body:        opts.Body,
		Labels:      opts.Labels,
		AssigneeIDs: assigneeIDs,
		Milestone:   opts.Milestone,
	})
	if err != nil {
		return fmt.Errorf("failed to create issue: %w", err)
	}
	fmt.Fprintf(opts.IO.Out, "%s Created issue #%s in %s/%s\n", cs.Green("✓"), issue.Number, owner, repo)
	fmt.Fprintf(opts.IO.Out, "  %s\n", issue.HTMLURL)
	if err := ensureAssigneesApplied(client, owner, repo, issue.Number, issue.HTMLURL, assigneeIDs, "created"); err != nil {
		return err
	}
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

func ensureAssigneesApplied(client *api.Client, owner, repo, issueNumber, issueURL string, expectedIDs []string, action string) error {
	if len(expectedIDs) == 0 {
		return nil
	}

	number, err := strconv.Atoi(issueNumber)
	if err != nil {
		return nil
	}

	issue, err := api.GetIssue(client, owner, repo, number)
	if err != nil {
		return nil
	}
	if hasExpectedAssignees(issue, expectedIDs) {
		return nil
	}
	if issueURL != "" {
		return fmt.Errorf("issue #%s was %s at %s, but GitCode API did not apply the requested assignees", issueNumber, action, issueURL)
	}
	return fmt.Errorf("issue #%s was %s, but GitCode API did not apply the requested assignees", issueNumber, action)
}

func hasExpectedAssignees(issue *api.Issue, expectedIDs []string) bool {
	if issue == nil || len(expectedIDs) == 0 {
		return true
	}

	actual := make(map[string]struct{}, len(issue.Assignees))
	for _, assignee := range issue.Assignees {
		if assignee == nil || assignee.ID == nil {
			continue
		}
		actual[fmt.Sprint(assignee.ID)] = struct{}{}
	}

	for _, expectedID := range expectedIDs {
		if _, ok := actual[expectedID]; !ok {
			return false
		}
	}
	return true
}
