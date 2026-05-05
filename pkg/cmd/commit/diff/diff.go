// Package diff implements the commit diff command
package diff

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type DiffOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	SHA        string
}

// NewCmdDiff creates the diff command
func NewCmdDiff(f *cmdutil.Factory, runF func(*DiffOptions) error) *cobra.Command {
	opts := &DiffOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "diff <sha>",
		Short: "Get commit diff",
		Long: heredoc.Doc(`
			Get the diff of a commit in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Get commit diff
			$ gc commit diff abc123 -R owner/repo
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.SHA = args[0]

			if runF != nil {
				return runF(opts)
			}
			return diffRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")

	return cmd
}

func diffRun(opts *DiffOptions) error {
	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}
	client, err := cmdutil.AuthenticatedClient(httpClient)
	if err != nil {
		return err
	}

	// Get repository
	owner, repo, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	// Get commit diff
	diff, err := api.GetCommitDiff(client, owner, repo, opts.SHA)
	if err != nil {
		return fmt.Errorf("failed to get commit diff: %w", err)
	}

	// Output
	fmt.Fprintf(opts.IO.Out, "%s\n", diff)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
