// Package list implements the milestone list command
package list

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string

	// Flags
	State string
	Limit int
	JSON  bool
}

// NewCmdList creates the list command
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List milestones",
		Long: heredoc.Doc(`
			List milestones in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# List open milestones
			$ gc milestone list -R owner/repo

			# List closed milestones
			$ gc milestone list --state closed

			# List milestones as JSON
			$ gc milestone list -R owner/repo --json
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
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of milestones to list")
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

	// List milestones
	milestones, err := api.ListRepoMilestones(client, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to list milestones: %w", err)
	}

	// Output
	if len(milestones) == 0 {
		if opts.JSON {
			return cmdutil.WriteJSON(opts.IO.Out, milestones)
		}
		fmt.Fprintf(opts.IO.Out, "No milestones found\n")
		return nil
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, milestones)
	}

	// Calculate max number width for alignment
	maxNumWidth := 0
	for _, ms := range milestones {
		w := len(fmt.Sprintf("#%d", ms.Number))
		if w > maxNumWidth {
			maxNumWidth = w
		}
	}

	fmt.Fprintf(opts.IO.Out, "\n")
	for _, ms := range milestones {
		state := cs.Green("open")
		if ms.State == "closed" {
			state = cs.Red("closed")
		}
		fmt.Fprintf(opts.IO.Out, "%-*s  %s  %s\n", maxNumWidth, fmt.Sprintf("#%d", ms.Number), state, ms.Title)
		if ms.Description != "" {
			fmt.Fprintf(opts.IO.Out, "%-*s  %s\n", maxNumWidth, "", ms.Description)
		}
	}
	fmt.Fprintf(opts.IO.Out, "\n")

	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
