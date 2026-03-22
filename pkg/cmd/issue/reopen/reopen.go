// Package reopen implements the issue reopen command
package reopen

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
		Short: "Reopen a closed issue",
		Long: heredoc.Doc(`
			Reopen a closed issue in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Reopen an issue
			$ gc issue reopen 123

			# Reopen with a comment
			$ gc issue reopen 123 --comment "Still reproduces on latest version"
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
		_, err := api.CreateIssueComment(client, owner, repo, opts.Number, &api.CreateCommentOptions{
			Body: opts.Comment,
		})
		if err != nil {
			return fmt.Errorf("failed to add comment: %w", err)
		}
	}

	// Reopen issue
	issue, err := api.ReopenIssue(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to reopen issue: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Reopened issue #%s in %s/%s\n", cs.Green("✓"), issue.Number, owner, repo)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	if repo == "" {
		// TODO: get from current git repo
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