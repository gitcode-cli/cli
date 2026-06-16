// Package checkout implements the pr checkout command
package checkout

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	gitpkg "gitcode.com/gitcode-cli/cli/git"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type CheckoutOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	BranchName string
}

// NewCmdCheckout creates the checkout command
func NewCmdCheckout(f *cmdutil.Factory, runF func(*CheckoutOptions) error) *cobra.Command {
	opts := &CheckoutOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "checkout <number>",
		Short: "Check out a pull request locally",
		Long: heredoc.Doc(`
			Check out a pull request branch locally.
		`),
		Example: heredoc.Doc(`
			# Checkout PR branch
			$ gc pr checkout 123 -R owner/repo

			# Checkout with custom branch name
			$ gc pr checkout 123 -R owner/repo --branch my-feature
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return cmdutil.NewUsageError(fmt.Sprintf("invalid PR number: %s", args[0]))
			}
			opts.Number = number

			if runF != nil {
				return runF(opts)
			}
			return checkoutRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.BranchName, "branch", "b", "", "Custom branch name")

	return cmd
}

func checkoutRun(opts *CheckoutOptions) error {
	cs := opts.IO.ColorScheme()

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

	// Get PR
	pr, err := api.GetPullRequest(client, owner, repo, opts.Number)
	if err != nil {
		return cmdutil.WrapNotFound(err, "PR #%d not found in %s/%s", opts.Number, owner, repo)
	}

	branchName := opts.BranchName
	if branchName == "" {
		branchName = pr.Head.Ref
	}

	// Validate branch name before using it in git commands
	if err := gitpkg.ValidateRef(branchName); err != nil {
		return fmt.Errorf("invalid branch name %q: %w", branchName, err)
	}

	// Use SafeFetch and SafeFetchFromURL for validated git fetch operations.
	// These prevent option-injection attacks when ref or URL comes from the API.

	// Fetch the branch from origin
	err = gitpkg.SafeFetchWithOutput(opts.IO.Out, opts.IO.ErrOut, "", "origin", pr.Head.Ref, branchName)
	if err != nil {
		// Try fetching from head repo if different (fork)
		if pr.Head.Repo != nil && pr.Head.Repo.FullName != owner+"/"+repo {
			fetchURL := pr.Head.Repo.CloneURL
			err = gitpkg.SafeFetchFromURLWithOutput(opts.IO.Out, opts.IO.ErrOut, "", fetchURL, pr.Head.Ref, branchName)
			if err != nil {
				return fmt.Errorf("failed to fetch branch: %w", err)
			}
		} else {
			return fmt.Errorf("failed to fetch branch: %w", err)
		}
	}

	// Checkout the branch
	if err := gitpkg.SafeCheckoutWithOutput(opts.IO.Out, opts.IO.ErrOut, "", branchName); err != nil {
		return fmt.Errorf("failed to checkout branch: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Checked out PR #%d to branch %s\n", cs.Green("✓"), pr.Number, branchName)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
