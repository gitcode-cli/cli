// Package view implements the issue view command
package view

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/browser"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/output"
)

type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	Comments   bool
	Web        bool
	JSON       bool
	TimeFormat string
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
		Short: "View an issue",
		Long: heredoc.Doc(`
			View an issue in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# View an issue
			$ gc issue view 123

			# View issue with comments
			$ gc issue view 123 --comments

			# View issue with relative time display
			$ gc issue view 123 --time-format relative

			# View issue in browser
			$ gc issue view 123 --web
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid issue number: %s", args[0])
			}
			opts.Number = number
			if _, err := parseTimeFormat(opts.TimeFormat); err != nil {
				return err
			}

			if runF != nil {
				return runF(opts)
			}
			return viewRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVarP(&opts.Comments, "comments", "c", false, "View issue comments")
	cmd.Flags().BoolVarP(&opts.Web, "web", "w", false, "Open in browser")
	cmd.Flags().StringVar(&opts.TimeFormat, "time-format", "absolute", "Time format for dates (absolute/relative)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func viewRun(opts *ViewOptions) error {
	cs := opts.IO.ColorScheme()
	timeFormat, err := parseTimeFormat(opts.TimeFormat)
	if err != nil {
		return err
	}
	now := time.Now()

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

	// Get issue
	issue, err := api.GetIssue(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to get issue: %w", err)
	}

	// Open in browser if --web flag is set
	if opts.Web {
		fmt.Fprintf(opts.IO.Out, "Opening %s in your browser.\n", issue.HTMLURL)
		return browser.Open(issue.HTMLURL)
	}

	if opts.JSON {
		if opts.Comments && issue.Comments > 0 {
			comments, err := api.ListIssueComments(client, owner, repo, opts.Number, nil)
			if err != nil {
				return fmt.Errorf("failed to get comments: %w", err)
			}
			return cmdutil.WriteJSON(opts.IO.Out, map[string]interface{}{
				"issue":    issue,
				"comments": comments,
			})
		}
		return cmdutil.WriteJSON(opts.IO.Out, issue)
	}

	if err := renderIssueDetails(opts.IO.Out, cs, issue, timeFormat, now); err != nil {
		return err
	}

	if opts.Comments && issue.Comments > 0 {
		comments, err := api.ListIssueComments(client, owner, repo, opts.Number, nil)
		if err != nil {
			return fmt.Errorf("failed to get comments: %w", err)
		}
		return renderIssueComments(opts.IO.Out, cs, comments, timeFormat, now)
	}

	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

func parseTimeFormat(value string) (output.TimeFormat, error) {
	format, err := output.ParseTimeFormat(strings.ToLower(strings.TrimSpace(value)))
	if err != nil {
		return "", cmdutil.NewUsageError(err.Error())
	}
	return format, nil
}

func formatViewTime(t time.Time, format output.TimeFormat, now time.Time) string {
	if t.IsZero() {
		return "unknown"
	}
	if format != output.TimeFormatRelative {
		return output.FormatTime(t, format)
	}
	return formatRelativeTime(t, now)
}

func formatRelativeTime(t time.Time, now time.Time) string {
	delta := now.Sub(t)
	if delta < 0 {
		delta = -delta
		return "in " + formatDuration(delta)
	}
	return formatDuration(delta)
}

func formatDuration(delta time.Duration) string {
	switch {
	case delta < time.Minute:
		return "just now"
	case delta < time.Hour:
		mins := int(delta.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case delta < 24*time.Hour:
		hours := int(delta.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case delta < 7*24*time.Hour:
		days := int(delta.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	case delta < 30*24*time.Hour:
		weeks := int(delta.Hours() / (24 * 7))
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	case delta < 365*24*time.Hour:
		months := int(delta.Hours() / (24 * 30))
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	default:
		years := int(delta.Hours() / (24 * 365))
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

func joinLabels(labels []*api.Label) string {
	names := make([]string, 0, len(labels))
	for _, label := range labels {
		if label != nil && label.Name != "" {
			names = append(names, label.Name)
		}
	}
	return strings.Join(names, ", ")
}

func joinUsers(users []*api.User) string {
	names := make([]string, 0, len(users))
	for _, user := range users {
		if user != nil && user.Login != "" {
			names = append(names, user.Login)
		}
	}
	return strings.Join(names, ", ")
}

func renderIssueView(out io.Writer, cs *iostreams.ColorScheme, issue *api.Issue, comments []api.IssueComment, timeFormat output.TimeFormat, now time.Time) error {
	if err := renderIssueDetails(out, cs, issue, timeFormat, now); err != nil {
		return err
	}
	return renderIssueComments(out, cs, comments, timeFormat, now)
}

func renderIssueDetails(out io.Writer, cs *iostreams.ColorScheme, issue *api.Issue, timeFormat output.TimeFormat, now time.Time) error {
	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, "%s #%s\n", cs.Bold(issue.Title), issue.Number)
	fmt.Fprintf(out, "  State: %s\n", issue.State)
	if issue.User != nil {
		fmt.Fprintf(out, "  Author: %s\n", issue.User.Login)
	}
	fmt.Fprintf(out, "  Created: %s\n", formatViewTime(issue.CreatedAt.Time, timeFormat, now))
	fmt.Fprintf(out, "  Updated: %s\n", formatViewTime(issue.UpdatedAt.Time, timeFormat, now))
	if issue.ClosedAt != nil && !issue.ClosedAt.IsZero() {
		fmt.Fprintf(out, "  Closed: %s\n", formatViewTime(issue.ClosedAt.Time, timeFormat, now))
	}
	if issue.Milestone != nil && issue.Milestone.Title != "" {
		fmt.Fprintf(out, "  Milestone: %s\n", issue.Milestone.Title)
	}
	if len(issue.Assignees) > 0 {
		fmt.Fprintf(out, "  Assignees: %s\n", joinUsers(issue.Assignees))
	}
	if len(issue.Labels) > 0 {
		fmt.Fprintf(out, "  Labels: %s\n", joinLabels(issue.Labels))
	}
	fmt.Fprintf(out, "  Comments: %d\n", issue.Comments)
	fmt.Fprintf(out, "\n")
	if issue.Body != "" {
		fmt.Fprintf(out, "%s\n", issue.Body)
		fmt.Fprintf(out, "\n")
	}
	fmt.Fprintf(out, "  %s\n", issue.HTMLURL)
	fmt.Fprintf(out, "\n")

	return nil
}

func renderIssueComments(out io.Writer, cs *iostreams.ColorScheme, comments []api.IssueComment, timeFormat output.TimeFormat, now time.Time) error {
	if len(comments) == 0 {
		return nil
	}

	fmt.Fprintf(out, "--- Comments (%d) ---\n\n", len(comments))
	for _, c := range comments {
		author := "unknown"
		if c.User != nil && c.User.Login != "" {
			author = c.User.Login
		}
		fmt.Fprintf(out, "%s at %s:\n", cs.Bold(author), formatViewTime(c.CreatedAt.Time, timeFormat, now))
		if c.Body != "" {
			fmt.Fprintf(out, "%s\n\n", c.Body)
		} else {
			fmt.Fprintln(out)
		}
	}

	return nil
}
