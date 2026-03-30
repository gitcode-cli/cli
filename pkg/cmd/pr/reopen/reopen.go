// Package reopen implements the pr reopen command
package reopen

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type ReopenOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	Comment string
}

// NewCmdReopen creates the reopen command
func NewCmdReopen(f *cmdutil.Factory, runF func(*ReopenOptions) error) *cobra.Command {
	opts := &ReopenOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "reopen <number>",
		Short: "Reopen a closed pull request",
		Long: heredoc.Doc(`
			Reopen a closed pull request in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Reopen a PR
			$ gc pr reopen 123

			# Reopen with a comment
			$ gc pr reopen 123 --comment "Reopening for further changes"
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
			return reopenRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Comment, "comment", "c", "", "Add a comment before reopening")

	return cmd
}

func reopenRun(opts *ReopenOptions) error {
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

	// Add comment if provided
	if opts.Comment != "" {
		_, err := api.CreatePRComment(client, owner, repo, opts.Number, &api.CreatePRCommentOptions{
			Body: opts.Comment,
		})
		if err != nil {
			return fmt.Errorf("failed to add comment: %w", err)
		}
	}

	// Reopen PR
	_, err = api.ReopenPullRequest(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to reopen PR: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Reopened PR #%d in %s/%s\n", cs.Green("✓"), opts.Number, owner, repo)
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
