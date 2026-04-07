// Package view implements the pr view command
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

	// Arguments
	Repository string
	Number     int

	// Flags
	Web        bool
	Comments   bool
	JSON       bool
	TimeFormat string
}

// NewCmdView creates the view command
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "view <number>",
		Short: "View a pull request",
		Long: heredoc.Doc(`
			View a pull request in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# View a PR
			$ gc pr view 123 -R owner/repo

			# View PR in browser
			$ gc pr view 123 -R owner/repo --web

			# View PR with relative time display
			$ gc pr view 123 -R owner/repo --time-format relative
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid PR number: %s", args[0])
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
	cmd.Flags().BoolVarP(&opts.Web, "web", "w", false, "Open in browser")
	cmd.Flags().BoolVarP(&opts.Comments, "comments", "c", false, "View comments")
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
	owner, repo, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	// Get PR
	pr, err := api.GetPullRequest(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to get PR: %w", err)
	}

	// Open in browser if --web flag is set
	if opts.Web {
		fmt.Fprintf(opts.IO.Out, "Opening %s in your browser.\n", pr.HTMLURL)
		return browser.Open(pr.HTMLURL)
	}

	if opts.JSON {
		if opts.Comments {
			comments, err := api.ListPRComments(client, owner, repo, opts.Number)
			if err != nil {
				return fmt.Errorf("failed to get comments: %w", err)
			}
			return cmdutil.WriteJSON(opts.IO.Out, map[string]interface{}{
				"pull_request": pr,
				"comments":     comments,
			})
		}
		return cmdutil.WriteJSON(opts.IO.Out, pr)
	}

	if err := renderPRDetails(opts.IO.Out, cs, pr, timeFormat, now); err != nil {
		return err
	}

	if opts.Comments {
		comments, err := api.ListPRComments(client, owner, repo, opts.Number)
		if err != nil {
			fmt.Fprintf(opts.IO.ErrOut, "%s Failed to get comments: %v\n", cs.Yellow("!"), err)
		} else if err := renderPRComments(opts.IO.Out, comments, timeFormat, now); err != nil {
			return err
		}
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

func renderPRView(out io.Writer, cs *iostreams.ColorScheme, pr *api.PullRequest, comments []api.PRComment, timeFormat output.TimeFormat, now time.Time) error {
	if err := renderPRDetails(out, cs, pr, timeFormat, now); err != nil {
		return err
	}
	return renderPRComments(out, comments, timeFormat, now)
}

func renderPRDetails(out io.Writer, cs *iostreams.ColorScheme, pr *api.PullRequest, timeFormat output.TimeFormat, now time.Time) error {
	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, "%s #%d\n", cs.Bold(pr.Title), pr.Number)
	if pr.Draft {
		fmt.Fprintf(out, "  State: %s (draft)\n", pr.State)
	} else if pr.Merged {
		fmt.Fprintf(out, "  State: merged\n")
	} else {
		fmt.Fprintf(out, "  State: %s\n", pr.State)
	}
	if pr.User != nil {
		fmt.Fprintf(out, "  Author: %s\n", pr.User.Login)
	}
	if pr.Head != nil && pr.Base != nil {
		fmt.Fprintf(out, "  Branch: %s -> %s\n", pr.Head.Ref, pr.Base.Ref)
	}
	fmt.Fprintf(out, "  Created: %s\n", formatViewTime(pr.CreatedAt.Time, timeFormat, now))
	fmt.Fprintf(out, "  Updated: %s\n", formatViewTime(pr.UpdatedAt.Time, timeFormat, now))
	fmt.Fprintf(out, "  Additions: +%d  Deletions: -%d  Files: %d\n", pr.Additions, pr.Deletions, pr.ChangedFiles)
	fmt.Fprintf(out, "  Commits: %d  Comments: %d\n", pr.Commits, pr.Comments)
	if len(pr.Assignees) > 0 {
		fmt.Fprintf(out, "  Assignees: %s\n", joinUsers(pr.Assignees))
	}
	if len(pr.Reviewers) > 0 {
		fmt.Fprintf(out, "  Reviewers: %s\n", joinUsers(pr.Reviewers))
	}
	if len(pr.Labels) > 0 {
		fmt.Fprintf(out, "  Labels: %s\n", joinLabels(pr.Labels))
	}
	fmt.Fprintf(out, "\n")
	if pr.Body != "" {
		fmt.Fprintf(out, "%s\n", pr.Body)
		fmt.Fprintf(out, "\n")
	}
	fmt.Fprintf(out, "  %s\n", pr.HTMLURL)

	return nil
}

func renderPRComments(out io.Writer, comments []api.PRComment, timeFormat output.TimeFormat, now time.Time) error {
	if len(comments) == 0 {
		fmt.Fprintf(out, "\n--- No comments ---\n\n")
		return nil
	}

	fmt.Fprintf(out, "\n--- Comments (%d) ---\n", len(comments))
	for _, c := range comments {
		author := "unknown"
		if c.User != nil && c.User.Login != "" {
			author = c.User.Login
		}
		fmt.Fprintf(out, "\n%s at %s", author, formatViewTime(c.CreatedAt.Time, timeFormat, now))
		if c.CommentType != "" {
			fmt.Fprintf(out, " [%s]", c.CommentType)
		}
		if c.DiffFile != "" {
			fmt.Fprintf(out, " (%s)", c.DiffFile)
		}
		fmt.Fprintf(out, ":\n%s\n", c.Body)
	}

	fmt.Fprintf(out, "\n")
	return nil
}
