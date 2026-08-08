// Package log implements the repo log command.
package log

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type Options struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	File       string
	Branch     string
	Limit      int
	Page       int
	PerPage    int
	LimitSet   bool
	PerPageSet bool
	JSON       bool
}

// NewCmdLog creates the repo log command.
func NewCmdLog(f *cmdutil.Factory, runF func(*Options) error) *cobra.Command {
	opts := &Options{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
		Limit:      30,
		Page:       1,
	}

	cmd := &cobra.Command{
		Use:   "log",
		Short: "List repository commits",
		Long: heredoc.Doc(`
			List repository commits, optionally scoped to a file path and branch.
		`),
		Example: heredoc.Doc(`
			# List recent commits
			$ gc repo log -R owner/repo

			# List commits that touched a file on a branch
			$ gc repo log -R owner/repo --file src/main.go --branch main

			# Output as JSON
			$ gc repo log -R owner/repo --file src/main.go --branch main --json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			opts.LimitSet = cmd.Flags().Changed("limit")
			opts.PerPageSet = cmd.Flags().Changed("per-page")
			return run(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVar(&opts.File, "file", "", "Filter commits that touched this file path")
	cmd.Flags().StringVarP(&opts.Branch, "branch", "b", "", "Branch, tag, or commit SHA")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of commits to list")
	cmd.Flags().IntVar(&opts.Page, "page", 1, "Page number to fetch")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "API page size (default: --limit)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func run(opts *Options) error {
	if opts.Limit <= 0 {
		return cmdutil.NewUsageError("--limit must be greater than 0")
	}
	if opts.Page <= 0 {
		return cmdutil.NewUsageError("--page must be greater than 0")
	}
	if opts.PerPage < 0 {
		return cmdutil.NewUsageError("--per-page must be greater than or equal to 0")
	}

	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}

	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}
	owner, repo, err := parseRepo(repository)
	if err != nil {
		return err
	}

	commits, err := api.ListCommits(client, owner, repo, &api.CommitListOptions{
		Path:    opts.File,
		SHA:     opts.Branch,
		Page:    opts.Page,
		PerPage: resolvePerPage(opts),
	})
	if err != nil {
		return fmt.Errorf("failed to list commits: %w", err)
	}
	commits = trimCommits(commits, opts)

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, commits)
	}
	if len(commits) == 0 {
		fmt.Fprintln(opts.IO.Out, "No commits found")
		return nil
	}
	for _, commit := range commits {
		fmt.Fprintf(opts.IO.Out, "%s  %s  %s\n", shortSHA(commit.SHA), commitDate(commit), firstLine(commitMessage(commit)))
	}
	return nil
}

func commitMessage(commit api.RepositoryCommit) string {
	if commit.Commit == nil {
		return ""
	}
	return commit.Commit.Message
}

func commitDate(commit api.RepositoryCommit) string {
	if commit.Commit == nil || commit.Commit.Author == nil || commit.Commit.Author.Date.IsZero() {
		return "unknown"
	}
	return commit.Commit.Author.Date.Format("2006-01-02")
}

func firstLine(text string) string {
	line, _, _ := strings.Cut(text, "\n")
	return line
}

func shortSHA(sha string) string {
	if len(sha) <= 8 {
		return sha
	}
	return sha[:8]
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

// resolvePerPage returns the API page size: the explicit --per-page when set,
// otherwise --limit. Mirrors the pr list / issue list pattern.
func resolvePerPage(opts *Options) int {
	if opts.PerPageSet && opts.PerPage > 0 {
		return opts.PerPage
	}
	return opts.Limit
}

// trimCommits truncates to --limit when both --limit and --per-page were set.
func trimCommits(commits []api.RepositoryCommit, opts *Options) []api.RepositoryCommit {
	if opts.PerPageSet && opts.LimitSet && len(commits) > opts.Limit {
		return commits[:opts.Limit]
	}
	return commits
}
