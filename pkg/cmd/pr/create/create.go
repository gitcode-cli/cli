// Package create implements the pr create command
package create

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/browser"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type CreateOptions struct {
	IO             *iostreams.IOStreams
	HttpClient     func() (*http.Client, error)
	Branch         func() (string, error)
	ExecGitCommand func(string, ...string) (string, error)
	CreatePR       func(*api.Client, string, string, *api.CreatePROptions) (*api.PullRequest, error)
	OpenBrowser    func(string) error

	// Arguments
	Repository string

	// Flags
	Title string
	Body  string
	Head  string
	Base  string
	Draft bool
	Fill  bool
	Web   bool
	Fork  string // 跨仓库 PR：fork 项目路径【owner/repo】
}

// NewCmdCreate creates the create command
func NewCmdCreate(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:             f.IOStreams,
		HttpClient:     f.HttpClient,
		Branch:         f.Branch,
		ExecGitCommand: execGitCommand,
		CreatePR:       api.CreatePullRequest,
		OpenBrowser:    browser.Open,
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a pull request",
		Long: heredoc.Doc(`
			Create a new pull request in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Create a PR with title and body (uses current branch as head)
			$ gc pr create -R owner/repo --title "Feature" --body "Description"

			# Create a PR with specific head and base branches
			$ gc pr create -R owner/repo --head feature-branch --base main --title "Feature" --body "Description"

			# Create a draft PR
			$ gc pr create -R owner/repo --title "Feature" --draft

			# Create a cross-repo PR (from fork to upstream)
			$ gc pr create -R upstream/repo --fork myfork/repo --head feature-branch --title "Feature"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return createRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "Title for the PR")
	cmd.Flags().StringVarP(&opts.Body, "body", "b", "", "Body for the PR")
	cmd.Flags().StringVarP(&opts.Head, "head", "H", "", "Head branch (default: current branch)")
	cmd.Flags().StringVarP(&opts.Base, "base", "B", "main", "Base branch")
	cmd.Flags().BoolVarP(&opts.Draft, "draft", "d", false, "Create as draft")
	cmd.Flags().BoolVarP(&opts.Fill, "fill", "f", false, "Fill from last commit")
	cmd.Flags().BoolVarP(&opts.Web, "web", "w", false, "Open in browser")
	cmd.Flags().StringVarP(&opts.Fork, "fork", "F", "", "Fork repository path for cross-repo PR (owner/repo)")

	return cmd
}

func createRun(opts *CreateOptions) error {
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

	// Validate required fields
	if opts.Fill {
		if err := fillFromLastCommit(opts); err != nil {
			return err
		}
	}

	if opts.Title == "" {
		return fmt.Errorf("title is required. Use --title flag")
	}

	// Auto-detect head branch if not specified
	head := opts.Head
	if head == "" {
		if opts.Branch == nil {
			return fmt.Errorf("head branch is required. Use --head flag")
		}
		output, err := opts.Branch()
		if err != nil {
			return fmt.Errorf("could not determine current branch. Use --head flag: %w", err)
		}
		head = strings.TrimSpace(output)
		if head == "" || head == "HEAD" {
			return fmt.Errorf("could not determine current branch. Use --head flag")
		}
	}

	// Create PR
	pr, err := opts.CreatePR(client, owner, repo, &api.CreatePROptions{
		Title:    opts.Title,
		Body:     opts.Body,
		Head:     head,
		Base:     opts.Base,
		Draft:    opts.Draft,
		ForkPath: opts.Fork,
	})
	if err != nil {
		return fmt.Errorf("failed to create PR: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Created PR #%d in %s/%s\n", cs.Green("✓"), pr.Number, owner, repo)
	fmt.Fprintf(opts.IO.Out, "  %s\n", pr.HTMLURL)
	if opts.Web {
		fmt.Fprintf(opts.IO.Out, "Opening %s in your browser.\n", pr.HTMLURL)
		if err := opts.OpenBrowser(pr.HTMLURL); err != nil {
			return fmt.Errorf("failed to open browser: %w", err)
		}
	}
	return nil
}

func fillFromLastCommit(opts *CreateOptions) error {
	output, err := opts.ExecGitCommand("git", "log", "-1", "--pretty=%B")
	if err != nil {
		return fmt.Errorf("failed to fill from last commit: %w", err)
	}

	message := strings.TrimSpace(output)
	if message == "" {
		return fmt.Errorf("failed to fill from last commit: commit message is empty")
	}

	subject, body := splitCommitMessage(message)
	if opts.Title == "" {
		opts.Title = subject
	}
	if opts.Body == "" {
		opts.Body = body
	}

	return nil
}

func splitCommitMessage(message string) (string, string) {
	lines := strings.Split(strings.ReplaceAll(message, "\r\n", "\n"), "\n")
	subject := strings.TrimSpace(lines[0])
	if len(lines) == 1 {
		return subject, ""
	}

	body := strings.TrimSpace(strings.Join(lines[1:], "\n"))
	return subject, body
}

func execGitCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
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
