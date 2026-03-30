// Package stats implements the repo stats command
package stats

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type StatsOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	Repository string
	Branch     string
	Author     string
	OnlySelf   bool
	Since      string
	Until      string
}

// NewCmdStats creates the stats command
func NewCmdStats(f *cmdutil.Factory, runF func(*StatsOptions) error) *cobra.Command {
	opts := &StatsOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Get code contribution statistics",
		Long: heredoc.Doc(`
			Get code contribution statistics for a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Get stats for main branch
			$ gc repo stats --branch main -R owner/repo

			# Get stats for a specific author
			$ gc repo stats --branch main --author username -R owner/repo

			# Get only your own stats
			$ gc repo stats --branch main --only-self -R owner/repo

			# Get stats for a date range
			$ gc repo stats --branch main --since 2024-01-01 --until 2024-12-31 -R owner/repo
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return statsRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Branch, "branch", "b", "main", "Branch name (required)")
	cmd.Flags().StringVarP(&opts.Author, "author", "a", "", "Filter by author username")
	cmd.Flags().BoolVar(&opts.OnlySelf, "only-self", false, "Only show your own stats")
	cmd.Flags().StringVar(&opts.Since, "since", "", "Only commits after this date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&opts.Until, "until", "", "Only commits before this date (YYYY-MM-DD)")

	cmd.MarkFlagRequired("branch")

	return cmd
}

func statsRun(opts *StatsOptions) error {
	cs := opts.IO.ColorScheme()

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

	statsOpts := &api.CommitStatsOptions{
		BranchName: opts.Branch,
		Author:     opts.Author,
		OnlySelf:   opts.OnlySelf,
		Since:      opts.Since,
		Until:      opts.Until,
	}

	stats, err := api.GetCommitStatistics(client, owner, repo, statsOpts)
	if err != nil {
		return fmt.Errorf("failed to get commit statistics: %w", err)
	}

	// Output
	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "%s: %d\n", cs.Bold("Total commits"), stats.Total)
	fmt.Fprintf(opts.IO.Out, "\n")

	// Output statistics by author
	if len(stats.Statistics) > 0 {
		fmt.Fprintf(opts.IO.Out, "%s\n", cs.Bold("By Author:"))
		fmt.Fprintf(opts.IO.Out, "\n")
		fmt.Fprintf(opts.IO.Out, "%-20s %10s %10s %10s\n", "Author", "Additions", "Deletions", "Total")
		fmt.Fprintf(opts.IO.Out, "%s\n", strings.Repeat("-", 52))
		for _, s := range stats.Statistics {
			fmt.Fprintf(opts.IO.Out, "%-20s %10d %10d %10d\n", s.Author, s.Additions, s.Deletions, s.Total)
		}
		fmt.Fprintf(opts.IO.Out, "\n")
	}

	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITCODE_TOKEN")
}
