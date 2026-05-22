package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdCreate(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "create PR with title and head",
			args:    []string{"--title", "Feature", "--head", "feature-branch", "--repo", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "create draft PR",
			args:    []string{"--title", "WIP", "--head", "draft", "--draft", "--repo", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "create with base",
			args:    []string{"--title", "Feature", "--head", "feature", "--base", "develop", "--repo", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "create cross-repo PR with fork",
			args:    []string{"--title", "Feature", "--head", "feature", "--fork", "myfork/repo", "--repo", "upstream/repo"},
			wantErr: false,
		},
		{
			name:    "create with json output",
			args:    []string{"--title", "Feature", "--head", "feature", "--repo", "owner/repo", "--json"},
			wantErr: false,
		},
		{
			name:    "create with body file",
			args:    []string{"--title", "Feature", "--head", "feature", "--repo", "owner/repo", "--body-file", "body.md"},
			wantErr: false,
		},
		{
			name:    "missing title",
			args:    []string{"--head", "feature", "--repo", "owner/repo"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			// For validation tests, don't provide runF so actual validation runs
			var cmd *cobra.Command
			if tt.name == "missing title" {
				cmd = NewCmdCreate(f, nil)
			} else {
				cmd = NewCmdCreate(f, func(opts *CreateOptions) error {
					return nil
				})
			}
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateRunJSONWritesCreatedPR(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Repository: "owner/repo",
		Title:      "title",
		Body:       "body",
		Head:       "feature-branch",
		Base:       "main",
		JSON:       true,
		CreatePR: func(client *api.Client, owner, repo string, createOpts *api.CreatePROptions) (*api.PullRequest, error) {
			return &api.PullRequest{Number: 7, Title: createOpts.Title, HTMLURL: "https://gitcode.com/owner/repo/merge_requests/7"}, nil
		},
		OpenBrowser: func(url string) error { return nil },
	}

	if err := createRun(opts); err != nil {
		t.Fatalf("createRun() error = %v", err)
	}

	var got map[string]interface{}
	out := f.IOStreams.Out.(*bytes.Buffer).String()
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("JSON output did not parse: %v", err)
	}
	if got["number"] != float64(7) || got["html_url"] != "https://gitcode.com/owner/repo/merge_requests/7" {
		t.Fatalf("JSON output = %#v", got)
	}
	if strings.Contains(out, "Created PR") {
		t.Fatalf("JSON output contains text banner: %q", out)
	}
}

func TestCreateRunReadsBodyFile(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	dir := t.TempDir()
	bodyPath := filepath.Join(dir, "body.md")
	if err := os.WriteFile(bodyPath, []byte("file body\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	f := cmdutil.TestFactory()
	var createdOpts *api.CreatePROptions
	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Repository: "owner/repo",
		Title:      "title",
		BodyFile:   bodyPath,
		Head:       "feature-branch",
		Base:       "main",
		CreatePR: func(client *api.Client, owner, repo string, createOpts *api.CreatePROptions) (*api.PullRequest, error) {
			createdOpts = createOpts
			return &api.PullRequest{Number: 7, HTMLURL: "https://gitcode.com/owner/repo/merge_requests/7"}, nil
		},
		OpenBrowser: func(url string) error { return nil },
	}

	if err := createRun(opts); err != nil {
		t.Fatalf("createRun() error = %v", err)
	}
	if createdOpts == nil {
		t.Fatalf("CreatePR() was not called")
	}
	if createdOpts.Body != "file body" {
		t.Fatalf("CreatePR Body = %q", createdOpts.Body)
	}
}

func TestCreateRunReadsBodyFileFromStdin(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	f.IOStreams.In = strings.NewReader("stdin body\n")

	var createdOpts *api.CreatePROptions
	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Repository: "owner/repo",
		Title:      "title",
		BodyFile:   "-",
		Head:       "feature-branch",
		Base:       "main",
		CreatePR: func(client *api.Client, owner, repo string, createOpts *api.CreatePROptions) (*api.PullRequest, error) {
			createdOpts = createOpts
			return &api.PullRequest{Number: 7, HTMLURL: "https://gitcode.com/owner/repo/merge_requests/7"}, nil
		},
		OpenBrowser: func(url string) error { return nil },
	}

	if err := createRun(opts); err != nil {
		t.Fatalf("createRun() error = %v", err)
	}
	if createdOpts == nil {
		t.Fatalf("CreatePR() was not called")
	}
	if createdOpts.Body != "stdin body" {
		t.Fatalf("CreatePR Body = %q", createdOpts.Body)
	}
}

func TestCreateRunBodyAndBodyFileMutualExclusion(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	called := false
	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Repository: "owner/repo",
		Title:      "title",
		Body:       "body",
		BodyFile:   "body.md",
		Head:       "feature-branch",
		Base:       "main",
		CreatePR: func(client *api.Client, owner, repo string, createOpts *api.CreatePROptions) (*api.PullRequest, error) {
			called = true
			return nil, nil
		},
		OpenBrowser: func(url string) error { return nil },
	}

	err := createRun(opts)
	if err == nil {
		t.Fatal("createRun() error = nil")
	}
	if !strings.Contains(err.Error(), "cannot use both --body and --body-file") {
		t.Fatalf("createRun() error = %v", err)
	}
	if called {
		t.Fatal("CreatePR should not be called when body inputs conflict")
	}
}

func TestCreateRunJSONWithWebReturnsUsageBeforeCreate(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	called := false
	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Repository: "owner/repo",
		Title:      "title",
		Head:       "feature-branch",
		Base:       "main",
		JSON:       true,
		Web:        true,
		CreatePR: func(client *api.Client, owner, repo string, createOpts *api.CreatePROptions) (*api.PullRequest, error) {
			called = true
			return nil, nil
		},
		OpenBrowser: func(url string) error { return nil },
	}

	err := createRun(opts)
	if err == nil {
		t.Fatal("createRun() error = nil")
	}
	if got := cmdutil.ExitCode(err); got != cmdutil.ExitUsage {
		t.Fatalf("ExitCode() = %d, want %d", got, cmdutil.ExitUsage)
	}
	if called {
		t.Fatal("CreatePR should not be called when --json and --web conflict")
	}
}

func TestCreateRunJSONWithWebReturnsUsageBeforeAuth(t *testing.T) {
	f := cmdutil.TestFactory()
	called := false
	opts := &CreateOptions{
		IO: f.IOStreams,
		HttpClient: func() (*http.Client, error) {
			t.Fatal("HttpClient should not be called when --json and --web conflict")
			return nil, nil
		},
		Repository: "owner/repo",
		Title:      "title",
		Head:       "feature-branch",
		Base:       "main",
		JSON:       true,
		Web:        true,
		CreatePR: func(client *api.Client, owner, repo string, createOpts *api.CreatePROptions) (*api.PullRequest, error) {
			called = true
			return nil, nil
		},
		OpenBrowser: func(url string) error { return nil },
	}

	err := createRun(opts)
	if err == nil {
		t.Fatal("createRun() error = nil")
	}
	if got := cmdutil.ExitCode(err); got != cmdutil.ExitUsage {
		t.Fatalf("ExitCode() = %d, want %d", got, cmdutil.ExitUsage)
	}
	if called {
		t.Fatal("CreatePR should not be called when --json and --web conflict")
	}
}

func TestFillFromLastCommit(t *testing.T) {
	opts := &CreateOptions{
		ExecGitCommand: func(name string, args ...string) (string, error) {
			return "feat: add api auth cleanup\n\nBody line 1\nBody line 2\n", nil
		},
	}

	if err := fillFromLastCommit(opts); err != nil {
		t.Fatalf("fillFromLastCommit() error = %v", err)
	}

	if opts.Title != "feat: add api auth cleanup" {
		t.Fatalf("Title = %q", opts.Title)
	}
	if opts.Body != "Body line 1\nBody line 2" {
		t.Fatalf("Body = %q", opts.Body)
	}
}

func TestCreateRunFillAndWeb(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	var createdOpts *api.CreatePROptions
	var openedURL string

	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Repository: "owner/repo",
		Base:       "main",
		Fill:       true,
		Web:        true,
		Branch: func() (string, error) {
			return "feature-branch", nil
		},
		ExecGitCommand: func(name string, args ...string) (string, error) {
			commandLine := name + " " + strings.Join(args, " ")
			switch commandLine {
			case "git log -1 --pretty=%B":
				return "feat: add fill behavior\n\ncommit body", nil
			default:
				return "", nil
			}
		},
		CreatePR: func(client *api.Client, owner, repo string, createOpts *api.CreatePROptions) (*api.PullRequest, error) {
			createdOpts = createOpts
			return &api.PullRequest{
				Number:  12,
				HTMLURL: "https://gitcode.com/owner/repo/merge_requests/12",
			}, nil
		},
		OpenBrowser: func(url string) error {
			openedURL = url
			return nil
		},
	}

	if err := createRun(opts); err != nil {
		t.Fatalf("createRun() error = %v", err)
	}

	if createdOpts == nil {
		t.Fatalf("CreatePR() was not called")
	}
	if createdOpts.Title != "feat: add fill behavior" {
		t.Fatalf("CreatePR Title = %q", createdOpts.Title)
	}
	if createdOpts.Body != "commit body" {
		t.Fatalf("CreatePR Body = %q", createdOpts.Body)
	}
	if createdOpts.Head != "feature-branch" {
		t.Fatalf("CreatePR Head = %q", createdOpts.Head)
	}
	if openedURL != "https://gitcode.com/owner/repo/merge_requests/12" {
		t.Fatalf("opened URL = %q", openedURL)
	}
}

func TestCreateRunFillPreservesExplicitTitleAndBody(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	var createdOpts *api.CreatePROptions

	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Repository: "owner/repo",
		Title:      "explicit title",
		Body:       "explicit body",
		Head:       "feature-branch",
		Fill:       true,
		Branch: func() (string, error) {
			return "ignored", nil
		},
		ExecGitCommand: func(name string, args ...string) (string, error) {
			return "commit title\n\ncommit body", nil
		},
		CreatePR: func(client *api.Client, owner, repo string, createOpts *api.CreatePROptions) (*api.PullRequest, error) {
			createdOpts = createOpts
			return &api.PullRequest{
				Number:  1,
				HTMLURL: "https://gitcode.com/owner/repo/merge_requests/1",
			}, nil
		},
		OpenBrowser: func(url string) error { return nil },
	}

	if err := createRun(opts); err != nil {
		t.Fatalf("createRun() error = %v", err)
	}

	if createdOpts.Title != "explicit title" || createdOpts.Body != "explicit body" {
		t.Fatalf("explicit values were overwritten: %+v", createdOpts)
	}
}

func TestCreateRunFillPreservesBodyFile(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	dir := t.TempDir()
	bodyPath := filepath.Join(dir, "body.md")
	if err := os.WriteFile(bodyPath, []byte("file body\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	f := cmdutil.TestFactory()
	var createdOpts *api.CreatePROptions

	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Repository: "owner/repo",
		Title:      "explicit title",
		BodyFile:   bodyPath,
		Head:       "feature-branch",
		Fill:       true,
		Branch: func() (string, error) {
			return "ignored", nil
		},
		ExecGitCommand: func(name string, args ...string) (string, error) {
			return "commit title\n\ncommit body", nil
		},
		CreatePR: func(client *api.Client, owner, repo string, createOpts *api.CreatePROptions) (*api.PullRequest, error) {
			createdOpts = createOpts
			return &api.PullRequest{
				Number:  1,
				HTMLURL: "https://gitcode.com/owner/repo/merge_requests/1",
			}, nil
		},
		OpenBrowser: func(url string) error { return nil },
	}

	if err := createRun(opts); err != nil {
		t.Fatalf("createRun() error = %v", err)
	}
	if createdOpts.Title != "explicit title" || createdOpts.Body != "file body" {
		t.Fatalf("explicit values were overwritten: %+v", createdOpts)
	}
}

func TestCreateRunUsesFactoryBranch(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	var createdOpts *api.CreatePROptions

	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Repository: "owner/repo",
		Title:      "title",
		Body:       "body",
		Base:       "main",
		Branch: func() (string, error) {
			return "feature/from-factory", nil
		},
		ExecGitCommand: func(name string, args ...string) (string, error) {
			t.Fatalf("ExecGitCommand should not be used for branch detection")
			return "", nil
		},
		CreatePR: func(client *api.Client, owner, repo string, createOpts *api.CreatePROptions) (*api.PullRequest, error) {
			createdOpts = createOpts
			return &api.PullRequest{Number: 7, HTMLURL: "https://gitcode.com/owner/repo/merge_requests/7"}, nil
		},
		OpenBrowser: func(url string) error { return nil },
	}

	if err := createRun(opts); err != nil {
		t.Fatalf("createRun() error = %v", err)
	}
	if createdOpts == nil {
		t.Fatalf("CreatePR() was not called")
	}
	if createdOpts.Head != "feature/from-factory" {
		t.Fatalf("CreatePR Head = %q", createdOpts.Head)
	}
}

func TestCreateRunBranchError(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Repository: "owner/repo",
		Title:      "title",
		Body:       "body",
		Branch: func() (string, error) {
			return "", fmt.Errorf("not in a git repository")
		},
		ExecGitCommand: func(name string, args ...string) (string, error) {
			t.Fatalf("ExecGitCommand should not be used for branch detection")
			return "", nil
		},
		CreatePR:    api.CreatePullRequest,
		OpenBrowser: func(url string) error { return nil },
	}

	err := createRun(opts)
	if err == nil {
		t.Fatalf("createRun() error = nil")
	}
	if !strings.Contains(err.Error(), "could not determine current branch") {
		t.Fatalf("createRun() error = %v", err)
	}
}
