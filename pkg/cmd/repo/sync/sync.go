// Package sync implements the repo sync command.
package sync

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	gitpkg "gitcode.com/gitcode-cli/cli/git"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

var gitRun = gitpkg.Run

type SyncResult struct {
	SourceRepo    string `json:"source_repo"`
	SourceDir     string `json:"source_dir"`
	TargetRepo    string `json:"target_repo"`
	TargetDir     string `json:"target_dir"`
	BaseBranch    string `json:"base_branch"`
	SyncBranch    string `json:"sync_branch"`
	Changed       bool   `json:"changed"`
	CommitMessage string `json:"commit_message,omitempty"`
	PRNumber      int    `json:"pr_number,omitempty"`
	PRURL         string `json:"pr_url,omitempty"`
}

type SyncOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	RootDir    func() (string, error)
	Branch     func() (string, error)
	BaseRepo   func() (string, error)
	GetRepo    func(*api.Client, string, string) (*api.Repository, error)
	CreatePR   func(*api.Client, string, string, *api.CreatePROptions) (*api.PullRequest, error)
	GitRun     func(string, ...string) (string, error)
	MkdirTemp  func(string, string) (string, error)
	RemoveAll  func(string) error

	TargetRepo string
	SourceDir  string
	TargetDir  string
	Title      string
	Body       string
	Base       string
	BranchName string
	CommitMsg  string
	Draft      bool
	JSON       bool
}

// NewCmdSync creates the repo sync command.
func NewCmdSync(f *cmdutil.Factory, runF func(*SyncOptions) error) *cobra.Command {
	opts := &SyncOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		RootDir:    gitpkg.RootDir,
		Branch:     f.Branch,
		BaseRepo:   f.BaseRepo,
		GetRepo:    api.GetRepo,
		CreatePR:   api.CreatePullRequest,
		GitRun:     gitpkg.RunInDir,
		MkdirTemp:  os.MkdirTemp,
		RemoveAll:  os.RemoveAll,
	}

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync a local directory into another repository and create a PR",
		Long: heredoc.Doc(`
			Sync a directory from the current local repository into a target repository path.

			The command clones the target repository into a temporary directory, copies the
			source directory contents into the requested target directory, commits the change,
			pushes a sync branch, and creates a pull request.
		`),
		Example: heredoc.Doc(`
			# Sync local docs/api into another repo's sync/api directory
			$ gc repo sync \
			  --target-repo infra-test/target-repo \
			  --source-dir docs/api \
			  --target-dir sync/api

			# Override base branch and PR title
			$ gc repo sync \
			  --target-repo infra-test/target-repo \
			  --source-dir pkg/contracts \
			  --target-dir mirror/contracts \
			  --base main \
			  --title "sync: update contracts"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return syncRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.TargetRepo, "target-repo", "", "Target repository (owner/repo)")
	cmd.Flags().StringVar(&opts.SourceDir, "source-dir", "", "Source directory in the current repository")
	cmd.Flags().StringVar(&opts.TargetDir, "target-dir", "", "Target directory inside the target repository")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Pull request title")
	cmd.Flags().StringVar(&opts.Body, "body", "", "Pull request body")
	cmd.Flags().StringVar(&opts.Base, "base", "", "Base branch in the target repository (default: target default branch)")
	cmd.Flags().StringVar(&opts.BranchName, "branch", "", "Sync branch name in the target repository")
	cmd.Flags().StringVar(&opts.CommitMsg, "commit-message", "", "Commit message for the sync commit")
	cmd.Flags().BoolVar(&opts.Draft, "draft", false, "Create the target pull request as draft")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)
	cmd.MarkFlagRequired("target-repo")
	cmd.MarkFlagRequired("source-dir")
	cmd.MarkFlagRequired("target-dir")

	return cmd
}

func syncRun(opts *SyncOptions) error {
	token := cmdutil.EnvToken()
	if token == "" {
		return cmdutil.NewAuthError("not authenticated. Run: gc auth login")
	}

	rootDir, err := opts.RootDir()
	if err != nil {
		return fmt.Errorf("repo sync must be run inside a git repository: %w", err)
	}

	sourceRepo, err := opts.BaseRepo()
	if err != nil {
		return fmt.Errorf("failed to resolve current repository: %w", err)
	}

	currentBranch, err := opts.Branch()
	if err != nil {
		return fmt.Errorf("failed to resolve current branch: %w", err)
	}

	sourcePath, err := resolveSourceDir(rootDir, opts.SourceDir)
	if err != nil {
		return err
	}

	targetDir, err := validateTargetDir(opts.TargetDir)
	if err != nil {
		return err
	}

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}
	client := api.NewClientFromHTTP(httpClient)
	client.SetToken(token, "environment")

	targetOwner, targetRepo, err := cmdutil.ParseRepo(opts.TargetRepo)
	if err != nil {
		return err
	}

	repo, err := opts.GetRepo(client, targetOwner, targetRepo)
	if err != nil {
		return fmt.Errorf("failed to get target repository: %w", err)
	}

	baseBranch := opts.Base
	if strings.TrimSpace(baseBranch) == "" {
		baseBranch = repo.DefaultBranch
	}
	if strings.TrimSpace(baseBranch) == "" {
		baseBranch = "main"
	}

	syncBranch := opts.BranchName
	if strings.TrimSpace(syncBranch) == "" {
		syncBranch = buildSyncBranch(sourceRepo, currentBranch, opts.SourceDir)
	}

	title := opts.Title
	if strings.TrimSpace(title) == "" {
		title = fmt.Sprintf("sync: %s -> %s/%s", cleanPath(opts.SourceDir), targetRepo, cleanPath(targetDir))
	}

	commitMsg := opts.CommitMsg
	if strings.TrimSpace(commitMsg) == "" {
		commitMsg = title
	}

	body := opts.Body
	if strings.TrimSpace(body) == "" {
		body = defaultPRBody(sourceRepo, currentBranch, opts.SourceDir, opts.TargetRepo, targetDir)
	}

	workDir, err := opts.MkdirTemp("", "gc-repo-sync-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer opts.RemoveAll(workDir)

	authURL := authenticatedGitURL(targetOwner, targetRepo, token)
	if _, err := gitRun("clone", authURL, workDir); err != nil {
		return fmt.Errorf("failed to clone target repository: %w", err)
	}

	if _, err := opts.GitRun(workDir, "checkout", "-B", syncBranch, "origin/"+baseBranch); err != nil {
		return fmt.Errorf("failed to prepare sync branch: %w", err)
	}

	targetPath := filepath.Join(workDir, filepath.FromSlash(targetDir))
	if err := replaceDirContents(sourcePath, targetPath); err != nil {
		return fmt.Errorf("failed to sync directory contents: %w", err)
	}

	status, err := opts.GitRun(workDir, "status", "--porcelain")
	if err != nil {
		return fmt.Errorf("failed to inspect target repository changes: %w", err)
	}

	result := SyncResult{
		SourceRepo:    sourceRepo,
		SourceDir:     cleanPath(opts.SourceDir),
		TargetRepo:    opts.TargetRepo,
		TargetDir:     cleanPath(targetDir),
		BaseBranch:    baseBranch,
		SyncBranch:    syncBranch,
		Changed:       strings.TrimSpace(status) != "",
		CommitMessage: commitMsg,
	}

	if !result.Changed {
		return writeSyncResult(opts, result)
	}

	if _, err := opts.GitRun(workDir, "add", "--all", "--", targetDir); err != nil {
		return fmt.Errorf("failed to stage synced changes: %w", err)
	}
	if _, err := opts.GitRun(workDir, "commit", "-m", commitMsg); err != nil {
		return fmt.Errorf("failed to create sync commit: %w", err)
	}
	if _, err := opts.GitRun(workDir, "push", "--force-with-lease", "-u", "origin", syncBranch); err != nil {
		return fmt.Errorf("failed to push sync branch: %w", err)
	}

	pr, err := opts.CreatePR(client, targetOwner, targetRepo, &api.CreatePROptions{
		Title: title,
		Body:  body,
		Head:  syncBranch,
		Base:  baseBranch,
		Draft: opts.Draft,
	})
	if err != nil {
		return fmt.Errorf("failed to create sync pull request: %w", err)
	}

	prURL := pr.HTMLURL
	if strings.TrimSpace(prURL) == "" && pr.Number > 0 {
		prURL = fmt.Sprintf("https://gitcode.com/%s/%s/merge_requests/%d", targetOwner, targetRepo, pr.Number)
	}
	result.PRNumber = pr.Number
	result.PRURL = prURL
	return writeSyncResult(opts, result)
}

func writeSyncResult(opts *SyncOptions, result SyncResult) error {
	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, result)
	}
	cs := opts.IO.ColorScheme()
	if !result.Changed {
		fmt.Fprintf(opts.IO.Out, "%s No changes to sync for %s -> %s/%s\n", cs.Gray("-"), result.SourceDir, result.TargetRepo, result.TargetDir)
		return nil
	}
	fmt.Fprintf(opts.IO.Out, "%s Synced %s to %s/%s\n", cs.Green("✓"), result.SourceDir, result.TargetRepo, result.TargetDir)
	fmt.Fprintf(opts.IO.Out, "  Branch: %s\n", result.SyncBranch)
	if result.PRNumber > 0 {
		fmt.Fprintf(opts.IO.Out, "  PR #%d: %s\n", result.PRNumber, result.PRURL)
	}
	return nil
}

func resolveSourceDir(rootDir, sourceDir string) (string, error) {
	if strings.TrimSpace(sourceDir) == "" {
		return "", fmt.Errorf("source directory is required")
	}
	sourcePath := sourceDir
	if !filepath.IsAbs(sourcePath) {
		sourcePath = filepath.Join(rootDir, filepath.FromSlash(sourceDir))
	}
	info, err := os.Stat(sourcePath)
	if err != nil {
		return "", fmt.Errorf("failed to access source directory: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("source path must be a directory: %s", sourceDir)
	}
	return sourcePath, nil
}

func validateTargetDir(targetDir string) (string, error) {
	cleaned := cleanPath(targetDir)
	if cleaned == "" || cleaned == "." {
		return "", fmt.Errorf("target directory must be a non-root directory")
	}
	if strings.HasPrefix(cleaned, "../") || cleaned == ".." {
		return "", fmt.Errorf("target directory must stay inside the target repository")
	}
	return cleaned, nil
}

func replaceDirContents(sourcePath, targetPath string) error {
	if err := os.RemoveAll(targetPath); err != nil {
		return err
	}
	if err := os.MkdirAll(targetPath, 0o755); err != nil {
		return err
	}
	return copyDirContents(sourcePath, targetPath)
}

func copyDirContents(sourcePath, targetPath string) error {
	return filepath.WalkDir(sourcePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		dst := filepath.Join(targetPath, rel)
		if d.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		return copyFile(path, dst)
	})
}

func copyFile(sourcePath, targetPath string) error {
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return err
	}
	src, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer src.Close()

	info, err := src.Stat()
	if err != nil {
		return err
	}

	dst, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func buildSyncBranch(sourceRepo, currentBranch, sourceDir string) string {
	repo := strings.ReplaceAll(sourceRepo, "/", "-")
	branch := sanitizeBranchSegment(currentBranch)
	dir := sanitizeBranchSegment(cleanPath(sourceDir))
	return fmt.Sprintf("sync/%s/%s/%s", repo, branch, dir)
}

func sanitizeBranchSegment(value string) string {
	value = strings.ToLower(value)
	replacer := strings.NewReplacer("/", "-", "\\", "-", " ", "-", "_", "-")
	value = replacer.Replace(value)
	value = strings.Trim(value, "-.")
	if value == "" {
		return "sync"
	}
	return value
}

func cleanPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	return filepath.ToSlash(filepath.Clean(path))
}

func defaultPRBody(sourceRepo, currentBranch, sourceDir, targetRepo, targetDir string) string {
	return fmt.Sprintf(
		"Sync `%s/%s` from `%s` branch `%s` into `%s/%s`.",
		sourceRepo,
		cleanPath(sourceDir),
		sourceRepo,
		currentBranch,
		targetRepo,
		cleanPath(targetDir),
	)
}

func authenticatedGitURL(owner, repo, token string) string {
	return fmt.Sprintf("https://oauth2:%s@gitcode.com/%s/%s.git", token, owner, repo)
}
