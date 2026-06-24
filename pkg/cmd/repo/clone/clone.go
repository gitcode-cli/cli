// Package clone implements the repo clone command
package clone

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	gitpkg "gitcode.com/gitcode-cli/cli/git"
	"gitcode.com/gitcode-cli/cli/internal/config"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type CloneOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	Config     func() (config.Config, error)

	// Arguments
	Repository string
	Directory  string

	// Flags
	GitProtocol string
	Depth       int
	Branch      string
	Recursive   bool
}

// NewCmdClone creates the clone command
func NewCmdClone(f *cmdutil.Factory, runF func(*CloneOptions) error) *cobra.Command {
	opts := &CloneOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "clone <repository> [<directory>]",
		Short: "Clone a repository locally",
		Long: heredoc.Doc(`
			Clone a GitCode repository to your local machine.

			The repository can be specified as:
			- OWNER/REPO format
			- Full URL (https://gitcode.com/OWNER/REPO)
		`),
		Example: heredoc.Doc(`
			# Clone a repository
			$ gc repo clone owner/repo

			# Clone to a specific directory
			$ gc repo clone owner/repo my-project

			# Clone with SSH
			$ gc repo clone owner/repo --git-protocol ssh

			# Shallow clone
			$ gc repo clone owner/repo --depth 1
		`),
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Repository = args[0]
			if len(args) > 1 {
				opts.Directory = args[1]
			}

			if runF != nil {
				return runF(opts)
			}
			return cloneRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.GitProtocol, "git-protocol", "p", "", "Git protocol to use (https/ssh, default: auth config)")
	cmdutil.SetFlagEnum(cmd, "git-protocol", "https", "ssh")
	cmd.Flags().IntVarP(&opts.Depth, "depth", "d", 0, "Create a shallow clone")
	cmd.Flags().StringVarP(&opts.Branch, "branch", "b", "", "Branch to checkout")
	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "Clone submodules")

	return cmd
}

func cloneRun(opts *CloneOptions) error {
	cs := opts.IO.ColorScheme()

	// Design note: repo clone does not use git.SafeFetch / SafeCheckout wrappers
	// because its inputs (repository URL, --branch, directory) are user-supplied,
	// not server-controlled. A user who can inject arguments into their own git
	// clone command could simply run git directly instead.
	// We do validate --branch via git.ValidateRef for consistency with pr checkout.

	// Validate depth
	if opts.Depth < 0 {
		return cmdutil.NewUsageError("--depth must be greater than 0")
	}

	// Validate branch name
	if opts.Branch != "" {
		if err := gitpkg.ValidateRef(opts.Branch); err != nil {
			return cmdutil.NewUsageError(fmt.Sprintf("invalid branch name: %v", err))
		}
	}

	// Parse repository
	protocol, err := resolveGitProtocol(opts)
	if err != nil {
		return err
	}

	repoURL, err := parseRepoURL(opts.Repository, protocol)
	if err != nil {
		return err
	}

	// Build git command
	gitArgs := []string{"clone"}
	if opts.Depth > 0 {
		gitArgs = append(gitArgs, "--depth", fmt.Sprintf("%d", opts.Depth))
	}
	if opts.Branch != "" {
		gitArgs = append(gitArgs, "--branch", opts.Branch)
	}
	if opts.Recursive {
		gitArgs = append(gitArgs, "--recursive")
	}
	gitArgs = append(gitArgs, repoURL)
	if opts.Directory != "" {
		gitArgs = append(gitArgs, opts.Directory)
	}

	// Execute git clone
	gitCmd := exec.Command("git", gitArgs...)
	gitCmd.Stdin = os.Stdin
	gitCmd.Stdout = opts.IO.Out
	gitCmd.Stderr = opts.IO.ErrOut

	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Cloned repository %s\n", cs.Green("✓"), opts.Repository)
	return nil
}

func resolveGitProtocol(opts *CloneOptions) (string, error) {
	if opts.GitProtocol != "" {
		return opts.GitProtocol, nil
	}
	if opts.Config == nil {
		return "ssh", nil
	}
	cfg, err := opts.Config()
	if err != nil {
		return "", fmt.Errorf("failed to read config: %w", err)
	}
	protocol := cfg.GitProtocol("gitcode.com").Value
	if protocol == "" {
		return "ssh", nil
	}
	return protocol, nil
}

func parseRepoURL(repo, protocol string) (string, error) {
	// Already a URL
	if strings.HasPrefix(repo, "http://") || strings.HasPrefix(repo, "https://") || strings.HasPrefix(repo, "git@") {
		return repo, nil
	}

	// OWNER/REPO format
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", cmdutil.NewUsageError(fmt.Sprintf("invalid repository format: %s", repo))
	}

	owner, name := parts[0], parts[1]
	if protocol == "ssh" {
		return fmt.Sprintf("git@gitcode.com:%s/%s.git", owner, name), nil
	}
	return fmt.Sprintf("https://gitcode.com/%s/%s.git", owner, name), nil
}
