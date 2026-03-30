package create

import (
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
		ExecGitCommand: func(name string, args ...string) (string, error) {
			commandLine := name + " " + strings.Join(args, " ")
			switch commandLine {
			case "git log -1 --pretty=%B":
				return "feat: add fill behavior\n\ncommit body", nil
			case "git rev-parse --abbrev-ref HEAD":
				return "feature-branch\n", nil
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
