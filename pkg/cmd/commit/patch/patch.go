// Package patch implements the commit patch command
package patch

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type PatchOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	SHA        string
}

// NewCmdPatch creates the patch command
func NewCmdPatch(f *cmdutil.Factory, runF func(*PatchOptions) error) *cobra.Command {
	opts := &PatchOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "patch <sha>",
		Short: "Get commit patch",
		Long: heredoc.Doc(`
			Get the patch of a commit in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Get commit patch
			$ gc commit patch abc123 -R owner/repo
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.SHA = args[0]

			if runF != nil {
				return runF(opts)
			}
			return patchRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")

	return cmd
}

func patchRun(opts *PatchOptions) error {
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

	// Get commit patch
	patch, err := api.GetCommitPatch(client, owner, repo, opts.SHA)
	if err != nil {
		return fmt.Errorf("failed to get commit patch: %w", err)
	}

	// Output
	fmt.Fprintf(opts.IO.Out, "%s\n", patch)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	if token := os.Getenv("GITCODE_TOKEN"); token != "" {
		return token
	}
	return cmdutil.EnvToken()
}
