// Package ready implements the pr ready command
package ready

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/gitcode-com/gitcode-cli/api"
	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
	"github.com/gitcode-com/gitcode-cli/pkg/iostreams"
)

type ReadyOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	WIP    bool
	Ready  bool
}

// NewCmdReady creates the ready command
func NewCmdReady(f *cmdutil.Factory, runF func(*ReadyOptions) error) *cobra.Command {
	opts := &ReadyOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "ready <number>",
		Short: "Toggle PR ready/wip status",
		Long: heredoc.Doc(`
			Toggle a pull request between ready and work-in-progress (draft) status.
		`),
		Example: heredoc.Doc(`
			# Mark PR as ready for review
			$ gc pr ready 123

			# Mark PR as work-in-progress
			$ gc pr ready 123 --wip
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid PR number: %s", args[0])
			}
			opts.Number = number

			if runF != nil {
				return runF(opts)
			}
			return readyRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVarP(&opts.WIP, "wip", "w", false, "Mark as work-in-progress")
	cmd.Flags().BoolVarP(&opts.Ready, "ready", "r", false, "Mark as ready for review")

	return cmd
}

func readyRun(opts *ReadyOptions) error {
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

	// Get repository
	owner, repo, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	// Get current PR to preserve title
	pr, err := api.GetPullRequest(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to get PR: %w", err)
	}

	// Determine draft status
	draft := opts.WIP
	if opts.Ready {
		draft = false
	}

	// Update PR
	pr, err = api.UpdatePullRequest(client, owner, repo, opts.Number, &api.UpdatePROptions{
		Draft: &draft,
		Title: pr.Title,
	})
	if err != nil {
		return fmt.Errorf("failed to update PR: %w", err)
	}

	// Output
	if draft {
		fmt.Fprintf(opts.IO.Out, "%s Marked PR #%d as work-in-progress\n", cs.Yellow("✓"), opts.Number)
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Marked PR #%d as ready for review\n", cs.Green("✓"), opts.Number)
	}
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