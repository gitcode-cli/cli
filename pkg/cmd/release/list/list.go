// Package list implements the release list command
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
)

type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Flags
	Repository string
	Limit      int
	JSON       bool
}

// NewCmdList creates the list command
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List releases",
		Long: heredoc.Doc(`
			List releases in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# List releases
			$ gc release list

			# List releases in a specific repository
			$ gc release list -R owner/repo

			# Limit the number of results
			$ gc release list --limit 10
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of releases to list")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func listRun(opts *ListOptions) error {
	cs := opts.IO.ColorScheme()

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

	// List releases
	releases, err := api.ListReleases(client, owner, repo, &api.ReleaseListOptions{
		PerPage: opts.Limit,
	})
	if err != nil {
		return fmt.Errorf("failed to list releases: %w", err)
	}

	if len(releases) == 0 {
		if opts.JSON {
			return cmdutil.WriteJSON(opts.IO.Out, releases)
		}
		fmt.Fprintf(opts.IO.Out, "No releases found in %s/%s\n", owner, repo)
		return nil
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, releases)
	}

	// Output
	fmt.Fprintf(opts.IO.Out, "\n")
	latestPublishedMarked := false
	for _, r := range releases {
		tag := r.TagName
		if r.Name != "" {
			tag = r.Name
		}

		// Status indicators
		var status string
		if r.Draft {
			status = cs.Gray("(draft)")
		} else if r.Prerelease {
			status = cs.Yellow("(pre-release)")
		} else if !latestPublishedMarked {
			status = cs.Green("(latest)")
			latestPublishedMarked = true
		} else {
			status = cs.Green("(published)")
		}

		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold(tag), status)
		if r.Body != "" {
			// Show first line of body
			lines := strings.Split(r.Body, "\n")
			if len(lines) > 0 && lines[0] != "" {
				fmt.Fprintf(opts.IO.Out, "  %s\n", truncate(lines[0], 60))
			}
		}
		fmt.Fprintf(opts.IO.Out, "  %s\n", r.HTMLURL)
		fmt.Fprintf(opts.IO.Out, "\n")
	}

	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
