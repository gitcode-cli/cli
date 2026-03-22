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

	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
	"github.com/gitcode-com/gitcode-cli/pkg/iostreams"
)

type CloneOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

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

	cmd.Flags().StringVarP(&opts.GitProtocol, "git-protocol", "p", "https", "Git protocol to use (https/ssh)")
	cmd.Flags().IntVarP(&opts.Depth, "depth", "d", 0, "Create a shallow clone")
	cmd.Flags().StringVarP(&opts.Branch, "branch", "b", "", "Branch to checkout")
	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "Clone submodules")

	return cmd
}

func cloneRun(opts *CloneOptions) error {
	cs := opts.IO.ColorScheme()

	// Parse repository
	repoURL, err := parseRepoURL(opts.Repository, opts.GitProtocol)
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

func parseRepoURL(repo, protocol string) (string, error) {
	// Already a URL
	if strings.HasPrefix(repo, "http://") || strings.HasPrefix(repo, "https://") || strings.HasPrefix(repo, "git@") {
		return repo, nil
	}

	// OWNER/REPO format
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid repository format: %s", repo)
	}

	owner, name := parts[0], parts[1]
	if protocol == "ssh" {
		return fmt.Sprintf("git@gitcode.com:%s/%s.git", owner, name), nil
	}
	return fmt.Sprintf("https://gitcode.com/%s/%s.git", owner, name), nil
}