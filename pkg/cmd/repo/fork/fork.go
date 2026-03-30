// Package fork implements the repo fork command
package fork

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/git"
	"gitcode.com/gitcode-cli/cli/internal/config"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type ForkOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	Config     func() (config.Config, error)
	ParseRepo  func(string) (*git.Repo, error)
	ForkRepo   func(*api.Client, string, string) (*api.Repository, error)
	CloneRepo  func(*git.Repo, string, string, int) error

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
		Config:     f.Config,
		ParseRepo:  git.ParseRepo,
		ForkRepo:   api.ForkRepo,
		CloneRepo:  git.CloneWithProtocol,
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

	if opts.Repository == "" {
		return fmt.Errorf("repository is required. Usage: gc repo fork owner/repo")
	}

	sourceRepo, err := opts.ParseRepo(opts.Repository)
	if err != nil {
		return fmt.Errorf("invalid repository: %w", err)
	}

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
	repo, err := opts.ForkRepo(client, sourceRepo.Owner, sourceRepo.Name)
	if err != nil {
		return fmt.Errorf("failed to fork repository: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Forked repository %s\n", cs.Green("✓"), repo.FullName)
	fmt.Fprintf(opts.IO.Out, "  %s\n", repo.HTMLURL)

	if opts.Clone {
		cfg, err := opts.Config()
		if err != nil {
			return fmt.Errorf("failed to read config: %w", err)
		}

		forkedRepo, err := opts.ParseRepo(repo.FullName)
		if err != nil {
			return fmt.Errorf("failed to parse forked repository: %w", err)
		}

		protocol := cfg.GitProtocol(sourceRepo.Host).Value
		if protocol == "" {
			protocol = "https"
		}

		if err := opts.CloneRepo(forkedRepo, "", protocol, 0); err != nil {
			return fmt.Errorf("failed to clone forked repository: %w", err)
		}

		fmt.Fprintf(opts.IO.Out, "%s Cloned fork to ./%s\n", cs.Green("✓"), forkedRepo.Name)
	}

	return nil
}
