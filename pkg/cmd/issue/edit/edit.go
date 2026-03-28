// Package edit implements the issue edit command
package edit

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type EditOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	Repository string
	Number     int

	// Flags
	Title        string
	Body         string
	State        string
	Assignees    []string
	Labels       []string
	Milestone    int
	SecurityHole bool
}

// NewCmdEdit creates the edit command
func NewCmdEdit(f *cmdutil.Factory, runF func(*EditOptions) error) *cobra.Command {
	opts := &EditOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "edit <number>",
		Short: "Edit an issue",
		Long: heredoc.Doc(`
			Edit an existing issue in a GitCode repository.

			You can update the title, body, state, assignees, labels, milestone, and visibility.
		`),
		Example: heredoc.Doc(`
			# Edit issue title
			$ gc issue edit 123 --title "New title" -R owner/repo

			# Edit issue body
			$ gc issue edit 123 --body "New description" -R owner/repo

			# Close an issue
			$ gc issue edit 123 --state close -R owner/repo

			# Reopen an issue
			$ gc issue edit 123 --state reopen -R owner/repo

			# Assign users
			$ gc issue edit 123 --assignee user1 --assignee user2 -R owner/repo

			# Add labels
			$ gc issue edit 123 --label bug,enhancement -R owner/repo

			# Set milestone
			$ gc issue edit 123 --milestone 5 -R owner/repo

			# Set as private issue
			$ gc issue edit 123 --security-hole -R owner/repo

			# Combine multiple edits
			$ gc issue edit 123 --title "Bug fix" --assignee user1 --label bug --milestone 1 -R owner/repo
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid issue number: %s", args[0])
			}
			opts.Number = number

			if runF != nil {
				return runF(opts)
			}
			return editRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "New title")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "New body/description")
	cmd.Flags().StringVarP(&opts.State, "state", "s", "", "State: open, closed, reopen, close")
	cmd.Flags().StringSliceVarP(&opts.Assignees, "assignee", "a", []string{}, "Assignees (comma-separated)")
	cmd.Flags().StringSliceVarP(&opts.Labels, "label", "l", []string{}, "Labels (comma-separated)")
	cmd.Flags().IntVarP(&opts.Milestone, "milestone", "m", 0, "Milestone number")
	cmd.Flags().BoolVar(&opts.SecurityHole, "security-hole", false, "Mark as private issue")

	return cmd
}

func editRun(opts *EditOptions) error {
	cs := opts.IO.ColorScheme()

	// Validate at least one edit option is provided
	if opts.Title == "" && opts.Body == "" && opts.State == "" &&
		len(opts.Assignees) == 0 && len(opts.Labels) == 0 &&
		opts.Milestone == 0 && !opts.SecurityHole {
		return fmt.Errorf("at least one edit option is required (e.g., --title, --body, --state, --assignee, --label, --milestone, --security-hole)")
	}

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := getEnvToken()
	if token == "" {
		return fmt.Errorf("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "environment")

	owner, repo, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	// Normalize state value
	state := opts.State
	if state == "closed" {
		state = "close"
	}

	// Build update options
	updateOpts := &api.UpdateIssueOptions{
		Title:     opts.Title,
		Body:      opts.Body,
		State:     state,
		Assignees: opts.Assignees,
		Labels:    opts.Labels,
		Milestone: opts.Milestone,
	}

	if opts.SecurityHole {
		updateOpts.SecurityHole = "true"
	}

	issue, err := api.UpdateIssue(client, owner, repo, opts.Number, updateOpts)
	if err != nil {
		return fmt.Errorf("failed to update issue: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Updated issue #%s in %s/%s\n", cs.Green("✓"), issue.Number, owner, repo)
	fmt.Fprintf(opts.IO.Out, "  %s\n", issue.HTMLURL)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	if repo == "" {
		return "", "", fmt.Errorf("no repository specified. Use -R owner/repo")
	}

	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repository format: %s", repo)
	}
	return parts[0], parts[1], nil
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITCODE_TOKEN")
}
