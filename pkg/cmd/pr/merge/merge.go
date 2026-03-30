// Package merge implements the pr merge command
package merge

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

type MergeOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	MergeMethod  string
	DeleteBranch bool
}

// NewCmdMerge creates the merge command
func NewCmdMerge(f *cmdutil.Factory, runF func(*MergeOptions) error) *cobra.Command {
	opts := &MergeOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "merge <number>",
		Short: "Merge a pull request",
		Long: heredoc.Doc(`
			Merge a pull request in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Merge a PR
			$ gc pr merge 123

			# Merge with squash
			$ gc pr merge 123 --squash

			# Merge and delete branch
			$ gc pr merge 123 --delete-branch
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
			return mergeRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.MergeMethod, "method", "m", "merge", "Merge method (merge/squash/rebase)")
	cmd.Flags().BoolVarP(&opts.DeleteBranch, "delete-branch", "d", false, "Delete branch after merge")

	return cmd
}

func mergeRun(opts *MergeOptions) error {
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

	// Get PR first
	pr, err := api.GetPullRequest(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to get PR: %w", err)
	}

	// Merge PR
	_, err = api.MergePullRequest(client, owner, repo, opts.Number, &api.MergePROptions{
		MergeMethod: opts.MergeMethod,
	})
	if err != nil {
		return fmt.Errorf("failed to merge PR: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Merged PR #%d\n", cs.Green("✓"), opts.Number)

	// Delete branch if requested
	if opts.DeleteBranch && pr.Head != nil {
		// TODO: delete branch via API
		fmt.Fprintf(opts.IO.Out, "  Branch %s can be deleted\n", pr.Head.Ref)
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
