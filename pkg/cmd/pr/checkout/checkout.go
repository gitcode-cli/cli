// Package checkout implements the pr checkout command
package checkout

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
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
			$ gc pr checkout 123

			# Checkout with custom branch name
			$ gc pr checkout 123 --branch my-feature
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

	// Get PR
	pr, err := api.GetPullRequest(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to get PR: %w", err)
	}

	branchName := opts.BranchName
	if branchName == "" {
		branchName = pr.Head.Ref
	}

	// Fetch the branch
	fetchCmd := exec.Command("git", "fetch", "origin", pr.Head.Ref+":"+branchName)
	fetchCmd.Stdout = opts.IO.Out
	fetchCmd.Stderr = opts.IO.ErrOut
	if err := fetchCmd.Run(); err != nil {
		// Try fetching from head repo if different
		if pr.Head.Repo != nil && pr.Head.Repo.FullName != owner+"/"+repo {
			fetchURL := pr.Head.Repo.CloneURL
			fetchCmd = exec.Command("git", "fetch", fetchURL, pr.Head.Ref+":"+branchName)
			fetchCmd.Stdout = opts.IO.Out
			fetchCmd.Stderr = opts.IO.ErrOut
			if err := fetchCmd.Run(); err != nil {
				return fmt.Errorf("failed to fetch branch: %w", err)
			}
		} else {
			return fmt.Errorf("failed to fetch branch: %w", err)
		}
	}

	// Checkout the branch
	checkoutCmd := exec.Command("git", "checkout", branchName)
	checkoutCmd.Stdout = opts.IO.Out
	checkoutCmd.Stderr = opts.IO.ErrOut
	if err := checkoutCmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout branch: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Checked out PR #%d to branch %s\n", cs.Green("✓"), pr.Number, branchName)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	if repo == "" {
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