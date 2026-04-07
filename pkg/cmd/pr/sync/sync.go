// Package sync implements the pr sync command.
package sync

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	gitpkg "gitcode.com/gitcode-cli/cli/git"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type SyncResult struct {
	SourcePR      string `json:"source_pr"`
	SourcePRURL   string `json:"source_pr_url"`
	TargetRepo    string `json:"target_repo"`
	TargetBranch  string `json:"target_branch"`
	SyncBranch    string `json:"sync_branch"`
	PRNumber      int    `json:"pr_number,omitempty"`
	PRURL         string `json:"pr_url,omitempty"`
	CommitsSynced int    `json:"commits_synced"`
	ConflictError string `json:"conflict_error,omitempty"`
}

type SyncOptions struct {
	IO            *iostreams.IOStreams
	HttpClient    func() (*http.Client, error)
	GetPR         func(*api.Client, string, string, int) (*api.PullRequest, error)
	ListPRCommits func(*api.Client, string, string, int) ([]api.Commit, error)
	GetRepo       func(*api.Client, string, string) (*api.Repository, error)
	CreatePR      func(*api.Client, string, string, *api.CreatePROptions) (*api.PullRequest, error)
	MkdirTemp     func(string, string) (string, error)
	RemoveAll     func(string) error
	WriteFile     func(string, []byte, os.FileMode) error
	RunGitInDir   func(string, string, ...string) (string, error)
	RunGit        func(string, ...string) (string, error)

	SourcePR   string
	TargetRepo string
	Base       string
	Title      string
	Body       string
	Draft      bool
	JSON       bool
}

// NewCmdSync creates the pr sync command.
func NewCmdSync(f *cmdutil.Factory, runF func(*SyncOptions) error) *cobra.Command {
	opts := &SyncOptions{
		IO:            f.IOStreams,
		HttpClient:    f.HttpClient,
		GetPR:         api.GetPullRequest,
		ListPRCommits: api.ListPRCommits,
		GetRepo:       api.GetRepo,
		CreatePR:      api.CreatePullRequest,
		MkdirTemp:     os.MkdirTemp,
		RemoveAll:     os.RemoveAll,
		WriteFile:     os.WriteFile,
	}

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync a PR from one repository to another",
		Long: heredoc.Doc(`
			Sync (cherry-pick) a pull request from one repository to another repository.

			This command fetches the commits from the source PR, cherry-picks them into
			a new branch in the target repository, and creates a new pull request.

			The source PR can be specified as:
			- owner/repo#number (e.g., gitcode-cli/cli#123)
			- Full URL (e.g., https://gitcode.com/gitcode-cli/cli/merge_requests/123)
		`),
		Example: heredoc.Doc(`
			# Sync PR #123 from source repo to target repo
			$ gc pr sync --source-pr owner/source-repo#123 --target-repo owner/target-repo

			# Specify target branch
			$ gc pr sync --source-pr owner/source-repo#123 \
				--target-repo owner/target-repo \
				--base release/v1.0

			# Custom title and body
			$ gc pr sync --source-pr owner/source-repo#123 \
				--target-repo owner/target-repo \
				--title "[sync] Fix login bug" \
				--body "Synced from owner/source-repo#123"

			# Create as draft
			$ gc pr sync --source-pr owner/source-repo#123 \
				--target-repo owner/target-repo \
				--draft

			# JSON output
			$ gc pr sync --source-pr owner/source-repo#123 \
				--target-repo owner/target-repo \
				--json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return syncRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.SourcePR, "source-pr", "", "Source PR (owner/repo#number or URL)")
	cmd.Flags().StringVar(&opts.TargetRepo, "target-repo", "", "Target repository (owner/repo)")
	cmd.Flags().StringVar(&opts.Base, "base", "", "Base branch in target repository (default: target default branch)")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Pull request title (default: [sync] <source title>)")
	cmd.Flags().StringVar(&opts.Body, "body", "", "Pull request body (default: inherit from source PR with sync info)")
	cmd.Flags().BoolVar(&opts.Draft, "draft", false, "Create the pull request as draft")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)
	cmd.MarkFlagRequired("source-pr")
	cmd.MarkFlagRequired("target-repo")

	return cmd
}

// PRRef represents a parsed PR reference
type PRRef struct {
	Owner  string
	Repo   string
	Number int
}

// ParsePRRef parses a PR reference string
func ParsePRRef(ref string) (*PRRef, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return nil, fmt.Errorf("source PR is required")
	}

	// Try URL format: https://gitcode.com/owner/repo/merge_requests/123
	urlPattern := regexp.MustCompile(`^https?://gitcode\.com/([^/]+)/([^/]+)/(?:merge_requests|pulls)/(\d+)(?:/[^/]*)?$`)
	if matches := urlPattern.FindStringSubmatch(ref); matches != nil {
		number, err := strconv.Atoi(matches[3])
		if err != nil {
			return nil, fmt.Errorf("invalid PR number in URL: %s", matches[3])
		}
		return &PRRef{Owner: matches[1], Repo: matches[2], Number: number}, nil
	}

	// Try short format: owner/repo#123
	shortPattern := regexp.MustCompile(`^([^/]+)/([^#]+)#(\d+)$`)
	if matches := shortPattern.FindStringSubmatch(ref); matches != nil {
		number, err := strconv.Atoi(matches[3])
		if err != nil {
			return nil, fmt.Errorf("invalid PR number: %s", matches[3])
		}
		return &PRRef{Owner: matches[1], Repo: strings.TrimSuffix(matches[2], ".git"), Number: number}, nil
	}

	return nil, fmt.Errorf("invalid PR reference format. Use owner/repo#number or https://gitcode.com/owner/repo/merge_requests/number")
}

func syncRun(opts *SyncOptions) error {
	token := cmdutil.EnvToken()
	if token == "" {
		return cmdutil.NewAuthError("not authenticated. Run: gc auth login")
	}

	// Parse source PR reference
	sourcePR, err := ParsePRRef(opts.SourcePR)
	if err != nil {
		return err
	}

	// Parse target repo
	targetOwner, targetRepo, err := cmdutil.ParseRepo(opts.TargetRepo)
	if err != nil {
		return err
	}

	// Create HTTP client
	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}
	client := api.NewClientFromHTTP(httpClient)
	client.SetToken(token, "environment")

	// Get source PR details
	pr, err := opts.GetPR(client, sourcePR.Owner, sourcePR.Repo, sourcePR.Number)
	if err != nil {
		return fmt.Errorf("failed to get source PR: %w", err)
	}

	// Get source PR commits
	commits, err := opts.ListPRCommits(client, sourcePR.Owner, sourcePR.Repo, sourcePR.Number)
	if err != nil {
		return fmt.Errorf("failed to get source PR commits: %w", err)
	}

	if len(commits) == 0 {
		return fmt.Errorf("source PR has no commits")
	}

	// Get target repository info
	targetRepoInfo, err := opts.GetRepo(client, targetOwner, targetRepo)
	if err != nil {
		return fmt.Errorf("failed to get target repository: %w", err)
	}

	// Determine base branch
	baseBranch := opts.Base
	if strings.TrimSpace(baseBranch) == "" {
		baseBranch = targetRepoInfo.DefaultBranch
	}
	if strings.TrimSpace(baseBranch) == "" {
		baseBranch = "main"
	}

	// Generate sync branch name
	syncBranch := buildSyncBranch(sourcePR.Owner, sourcePR.Repo, sourcePR.Number)

	// Build title and body
	title := opts.Title
	if strings.TrimSpace(title) == "" {
		title = fmt.Sprintf("[sync] %s", pr.Title)
	}

	body := opts.Body
	if strings.TrimSpace(body) == "" {
		body = buildSyncBody(pr, sourcePR, opts.TargetRepo)
	}

	// Create temporary work directory
	workDir, err := opts.MkdirTemp("", "gc-pr-sync-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer opts.RemoveAll(workDir)

	// Create temporary credential helper script for secure authentication
	// This avoids embedding token in URL or process arguments
	// Note: credential helper is created outside workDir to avoid clone conflict
	credHelperDir, err := opts.MkdirTemp("", "gc-cred-*")
	if err != nil {
		return fmt.Errorf("failed to create credential helper directory: %w", err)
	}
	defer opts.RemoveAll(credHelperDir)

	credHelperPath := filepath.Join(credHelperDir, "git-credential-gc")
	credHelperScript := fmt.Sprintf(`#!/bin/bash
echo "protocol=https"
echo "host=gitcode.com"
echo "username=oauth2"
echo "password=%s"
`, token)
	if err := opts.WriteFile(credHelperPath, []byte(credHelperScript), 0700); err != nil {
		return fmt.Errorf("failed to create credential helper: %w", err)
	}

	// Git commands with credential helper
	gitCmd := opts.RunGit
	if gitCmd == nil {
		gitCmd = func(credHelperPath string, args ...string) (string, error) {
			fullArgs := append([]string{"-c", "credential.helper=" + credHelperPath}, args...)
			return gitpkg.Run(fullArgs...)
		}
	}

	gitCmdInDir := opts.RunGitInDir
	if gitCmdInDir == nil {
		gitCmdInDir = func(dir string, credHelperPath string, args ...string) (string, error) {
			fullArgs := append([]string{"-C", dir, "-c", "credential.helper=" + credHelperPath}, args...)
			return gitpkg.Run(fullArgs...)
		}
	}

	runGit := func(args ...string) (string, error) {
		return gitCmd(credHelperPath, args...)
	}

	runGitInDir := func(dir string, args ...string) (string, error) {
		return gitCmdInDir(dir, credHelperPath, args...)
	}

	// Clone target repository
	if _, err := runGit("clone", repositoryGitURL(targetOwner, targetRepo), workDir); err != nil {
		return fmt.Errorf("failed to clone target repository: %w", err)
	}

	// Fetch source repository to get commits
	if _, err := runGitInDir(workDir, "remote", "add", "source", repositoryGitURL(sourcePR.Owner, sourcePR.Repo)); err != nil {
		return fmt.Errorf("failed to add source remote: %w", err)
	}
	if _, err := runGitInDir(workDir, "fetch", "source"); err != nil {
		return fmt.Errorf("failed to fetch source repository: %w", err)
	}

	// Create sync branch based on target base branch
	if _, err := runGitInDir(workDir, "checkout", "-B", syncBranch, "origin/"+baseBranch); err != nil {
		return fmt.Errorf("failed to create sync branch: %w", err)
	}

	commitsSynced, conflictError := syncCommits(runGitInDir, workDir, commits)

	result := SyncResult{
		SourcePR:      fmt.Sprintf("%s/%s#%d", sourcePR.Owner, sourcePR.Repo, sourcePR.Number),
		SourcePRURL:   pr.HTMLURL,
		TargetRepo:    opts.TargetRepo,
		TargetBranch:  baseBranch,
		SyncBranch:    syncBranch,
		CommitsSynced: commitsSynced,
		ConflictError: conflictError,
	}

	if conflictError != "" {
		return writeSyncResult(opts, result, cmdutil.NewCLIError(cmdutil.ExitConflict, fmt.Sprintf("%s. Manual resolution required.", conflictError), nil))
	}

	// Push sync branch
	if _, err := runGitInDir(workDir, "push", "--force-with-lease", "-u", "origin", syncBranch); err != nil {
		return fmt.Errorf("failed to push sync branch: %w", err)
	}

	// Create pull request
	newPR, err := opts.CreatePR(client, targetOwner, targetRepo, &api.CreatePROptions{
		Title: title,
		Body:  body,
		Head:  syncBranch,
		Base:  baseBranch,
		Draft: opts.Draft,
	})
	if err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}

	prURL := newPR.HTMLURL
	if strings.TrimSpace(prURL) == "" && newPR.Number > 0 {
		prURL = pullRequestWebURL(targetOwner, targetRepo, newPR.Number)
	}

	result.PRNumber = newPR.Number
	result.PRURL = prURL
	return writeSyncResult(opts, result, nil)
}

func writeSyncResult(opts *SyncOptions, result SyncResult, err error) error {
	if opts.JSON {
		if err != nil {
			result.ConflictError = err.Error()
		}
		if writeErr := cmdutil.WriteJSON(opts.IO.Out, result); writeErr != nil {
			return writeErr
		}
		return err
	}

	cs := opts.IO.ColorScheme()
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "%s %s\n", cs.Red("✗"), err.Error())
		fmt.Fprintf(opts.IO.Out, "Partial sync to %s (branch: %s)\n", result.TargetRepo, result.SyncBranch)
		fmt.Fprintf(opts.IO.Out, "Commits synced before conflict: %d\n", result.CommitsSynced)
		return err
	}

	fmt.Fprintf(opts.IO.Out, "%s Synced PR %s to %s\n", cs.Green("✓"), result.SourcePR, result.TargetRepo)
	fmt.Fprintf(opts.IO.Out, "  Branch: %s\n", result.SyncBranch)
	fmt.Fprintf(opts.IO.Out, "  Commits: %d\n", result.CommitsSynced)
	if result.PRNumber > 0 {
		fmt.Fprintf(opts.IO.Out, "  PR #%d: %s\n", result.PRNumber, result.PRURL)
	}
	return nil
}

func buildSyncBranch(sourceOwner, sourceRepo string, sourceNumber int) string {
	timestamp := time.Now().Format("20060102")
	return fmt.Sprintf("sync/pr-%s-%s-%d-%s", sourceOwner, sourceRepo, sourceNumber, timestamp)
}

func buildSyncBody(pr *api.PullRequest, sourcePR *PRRef, targetRepo string) string {
	var body string
	if strings.TrimSpace(pr.Body) != "" {
		body = pr.Body + "\n\n---\n\n"
	}
	body += fmt.Sprintf("Synced from [%s/%s#%d](%s) to %s.",
		sourcePR.Owner, sourcePR.Repo, sourcePR.Number,
		pr.HTMLURL, targetRepo)
	return body
}

func syncCommits(runGitInDir func(string, ...string) (string, error), workDir string, commits []api.Commit) (int, string) {
	commitsSynced := 0
	for _, commit := range commits {
		if _, err := runGitInDir(workDir, "cherry-pick", "-x", commit.SHA); err != nil {
			_, _ = runGitInDir(workDir, "cherry-pick", "--abort")
			return commitsSynced, fmt.Sprintf("cherry-pick conflict on commit %s: %s", shortSHA(commit.SHA), commit.Message)
		}
		commitsSynced++
	}
	return commitsSynced, ""
}

func shortSHA(sha string) string {
	if len(sha) <= 8 {
		return sha
	}
	return sha[:8]
}

func pullRequestWebURL(owner, repo string, number int) string {
	return fmt.Sprintf("https://gitcode.com/%s/%s/merge_requests/%d", owner, repo, number)
}

// repositoryGitURL returns a Git URL without embedded credentials
func repositoryGitURL(owner, repo string) string {
	return fmt.Sprintf("https://gitcode.com/%s/%s.git", owner, repo)
}
