// Package fork implements the repo fork command
package fork

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/gitcode-com/gitcode-cli/api"
	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
	"github.com/gitcode-com/gitcode-cli/pkg/iostreams"
)

type ForkOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string

	// Flags
	Clone bool
}

// NewCmdFork creates the fork command
func NewCmdFork(f *cmdutil.Factory, runF func(*ForkOptions) error) *cobra.Command {
	opts := &ForkOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "fork [<repository>]",
		Short: "Fork a repository",
		Long: heredoc.Doc(`
			Fork a repository to your account.
		`),
		Example: heredoc.Doc(`
			# Fork a repository
			$ gc repo fork owner/repo

			# Fork and clone
			$ gc repo fork owner/repo --clone
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Repository = args[0]
			}

			if runF != nil {
				return runF(opts)
			}
			return forkRun(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Clone, "clone", "c", false, "Clone the fork after creating")

	return cmd
}

func forkRun(opts *ForkOptions) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := os.Getenv("GC_TOKEN")
	if token == "" {
		token = os.Getenv("GITCODE_TOKEN")
	}
	if token == "" {
		return fmt.Errorf("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "environment")

	// Fork repo
	repo, err := api.ForkRepo(client, "owner", "repo") // TODO: parse repository
	if err != nil {
		return fmt.Errorf("failed to fork repository: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Forked repository %s\n", cs.Green("✓"), repo.FullName)
	fmt.Fprintf(opts.IO.Out, "  %s\n", repo.HTMLURL)

	return nil
}