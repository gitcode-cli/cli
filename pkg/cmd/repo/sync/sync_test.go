package sync

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdSync(t *testing.T) {
	f := cmdutil.TestFactory()
	cmd := NewCmdSync(f, func(opts *SyncOptions) error {
		return nil
	})

	if cmd == nil {
		t.Fatal("NewCmdSync returned nil")
	}
	if cmd.Use != "sync" {
		t.Fatalf("cmd.Use = %q", cmd.Use)
	}
}

func TestBuildSyncBranch(t *testing.T) {
	got := buildSyncBranch("owner/repo", "feature/test_branch", "docs/api")
	want := "sync/owner-repo/feature-test-branch/docs-api"
	if got != want {
		t.Fatalf("buildSyncBranch() = %q, want %q", got, want)
	}
}

func TestValidateTargetDir(t *testing.T) {
	if _, err := validateTargetDir("."); err == nil {
		t.Fatal("expected error for root target dir")
	}
	if _, err := validateTargetDir(".git"); err == nil {
		t.Fatal("expected error for repository metadata directory")
	}
	if _, err := validateTargetDir(".git/hooks"); err == nil {
		t.Fatal("expected error for repository metadata subdirectory")
	}
	got, err := validateTargetDir("sync/contracts")
	if err != nil {
		t.Fatalf("validateTargetDir() error = %v", err)
	}
	if got != "sync/contracts" {
		t.Fatalf("validateTargetDir() = %q", got)
	}
}

func TestResolveSourceDir(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "docs")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatal(err)
	}

	resolved, display, err := resolveSourceDir(root, source)
	if err != nil {
		t.Fatalf("resolveSourceDir() error = %v", err)
	}
	if resolved != source {
		t.Fatalf("resolveSourceDir() resolved = %q, want %q", resolved, source)
	}
	if display != "docs" {
		t.Fatalf("resolveSourceDir() display = %q, want %q", display, "docs")
	}
}

func TestResolveSourceDirRejectsOutsideRoot(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()

	if _, _, err := resolveSourceDir(root, outside); err == nil {
		t.Fatal("expected error for source directory outside repository root")
	}
}

func TestReplaceDirContents(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	target := filepath.Join(root, "target")
	if err := os.MkdirAll(filepath.Join(source, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "nested", "file.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "stale.txt"), []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := replaceDirContents(source, target); err != nil {
		t.Fatalf("replaceDirContents() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(target, "stale.txt")); !os.IsNotExist(err) {
		t.Fatalf("stale file should be removed, got err=%v", err)
	}
	content, err := os.ReadFile(filepath.Join(target, "nested", "file.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "hello" {
		t.Fatalf("copied content = %q", string(content))
	}
}

func TestSyncRunNoChanges(t *testing.T) {
	f := cmdutil.TestFactory()
	root := t.TempDir()
	source := filepath.Join(root, "docs")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "api.txt"), []byte("same"), 0o644); err != nil {
		t.Fatal(err)
	}

	targetRepoDir := filepath.Join(root, "target")
	if err := os.MkdirAll(filepath.Join(targetRepoDir, "mirror"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(targetRepoDir, "mirror", "api.txt"), []byte("same"), 0o644); err != nil {
		t.Fatal(err)
	}

	var cloneArgs []string
	var cloneEnv map[string]string
	prCalled := false
	opts := &SyncOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		RootDir: func() (string, error) {
			return root, nil
		},
		Branch: func() (string, error) {
			return "feature/demo", nil
		},
		BaseRepo: func() (string, error) {
			return "owner/source", nil
		},
		GetRepo: func(client *api.Client, owner, repo string) (*api.Repository, error) {
			return &api.Repository{DefaultBranch: "main"}, nil
		},
		CreatePR: func(client *api.Client, owner, repo string, opts *api.CreatePROptions) (*api.PullRequest, error) {
			prCalled = true
			return nil, nil
		},
		GitRun: func(dir string, env map[string]string, args ...string) (string, error) {
			switch strings.Join(args, " ") {
			case "checkout -B sync/owner-source/feature-demo/docs origin/main":
				return "", nil
			case "status --porcelain":
				return "", nil
			default:
				return "", nil
			}
		},
		MkdirTemp: func(dir, pattern string) (string, error) {
			return targetRepoDir, nil
		},
		RemoveAll:  func(string) error { return nil },
		TargetRepo: "infra-test/target",
		SourceDir:  "docs",
		TargetDir:  "mirror",
	}

	oldToken := os.Getenv("GC_TOKEN")
	t.Cleanup(func() { _ = os.Setenv("GC_TOKEN", oldToken) })
	_ = os.Setenv("GC_TOKEN", "token")

	originalGitRun := gitRun
	gitRun = func(env map[string]string, args ...string) (string, error) {
		cloneEnv = env
		cloneArgs = append([]string{}, args...)
		return "", nil
	}
	t.Cleanup(func() { gitRun = originalGitRun })

	if err := syncRun(opts); err != nil {
		t.Fatalf("syncRun() error = %v", err)
	}
	if prCalled {
		t.Fatal("CreatePR should not be called when there are no changes")
	}
	if len(cloneArgs) != 3 {
		t.Fatalf("unexpected clone args: %#v", cloneArgs)
	}
	if cloneEnv["GIT_CONFIG_COUNT"] != "1" || cloneEnv["GIT_CONFIG_KEY_0"] != "http.extraHeader" || cloneEnv["GIT_CONFIG_VALUE_0"] != "Authorization: Bearer token" {
		t.Fatalf("unexpected auth env: %#v", cloneEnv)
	}
	if cloneArgs[1] != "https://gitcode.com/infra-test/target.git" {
		t.Fatalf("unexpected clone URL: %#v", cloneArgs)
	}
}

func TestAuthenticatedGitEnv(t *testing.T) {
	env := authenticatedGitEnv("token")
	if env["GIT_CONFIG_COUNT"] != "1" || env["GIT_CONFIG_KEY_0"] != "http.extraHeader" {
		t.Fatalf("authenticatedGitEnv() = %#v", env)
	}
	if env["GIT_CONFIG_VALUE_0"] != "Authorization: Bearer token" {
		t.Fatalf("authenticatedGitEnv() token = %q", env["GIT_CONFIG_VALUE_0"])
	}
}

func TestSyncRunBuildsPRURLWhenCreateResponseOmitsHTMLURL(t *testing.T) {
	f := cmdutil.TestFactory()
	buf, ok := f.IOStreams.Out.(*bytes.Buffer)
	if !ok {
		t.Fatalf("output type = %T", f.IOStreams.Out)
	}

	root := t.TempDir()
	source := filepath.Join(root, "docs")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "api.txt"), []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}

	targetRepoDir := filepath.Join(root, "target")
	if err := os.MkdirAll(targetRepoDir, 0o755); err != nil {
		t.Fatal(err)
	}

	opts := &SyncOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		RootDir: func() (string, error) {
			return root, nil
		},
		Branch: func() (string, error) {
			return "feature/demo", nil
		},
		BaseRepo: func() (string, error) {
			return "owner/source", nil
		},
		GetRepo: func(client *api.Client, owner, repo string) (*api.Repository, error) {
			return &api.Repository{DefaultBranch: "main"}, nil
		},
		CreatePR: func(client *api.Client, owner, repo string, opts *api.CreatePROptions) (*api.PullRequest, error) {
			return &api.PullRequest{Number: 7}, nil
		},
		GitRun: func(dir string, env map[string]string, args ...string) (string, error) {
			switch strings.Join(args, " ") {
			case "checkout -B sync/owner-source/feature-demo/docs origin/main":
				return "", nil
			case "status --porcelain":
				return "M  mirror/api.txt\n", nil
			case "add --all -- mirror":
				return "", nil
			case "commit -m sync: owner/source -> target":
				return "", nil
			case "push --force-with-lease -u origin sync/owner-source/feature-demo/docs":
				return "", nil
			default:
				return "", nil
			}
		},
		MkdirTemp: func(dir, pattern string) (string, error) {
			return targetRepoDir, nil
		},
		RemoveAll:  func(string) error { return nil },
		TargetRepo: "infra-test/target",
		SourceDir:  "docs",
		TargetDir:  "mirror",
		CommitMsg:  "sync: owner/source -> target",
		JSON:       true,
	}

	oldToken := os.Getenv("GC_TOKEN")
	t.Cleanup(func() { _ = os.Setenv("GC_TOKEN", oldToken) })
	_ = os.Setenv("GC_TOKEN", "token")

	originalGitRun := gitRun
	gitRun = func(env map[string]string, args ...string) (string, error) { return "", nil }
	t.Cleanup(func() { gitRun = originalGitRun })

	if err := syncRun(opts); err != nil {
		t.Fatalf("syncRun() error = %v", err)
	}

	var result SyncResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v, body=%q", err, buf.String())
	}
	if result.PRURL != "https://gitcode.com/infra-test/target/merge_requests/7" {
		t.Fatalf("result.PRURL = %q", result.PRURL)
	}
}
