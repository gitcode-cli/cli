// Package sync implements the pr sync command.
package sync

import (
	"bytes"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

func TestParsePRRef(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *PRRef
		wantErr bool
	}{
		{
			name:  "short format",
			input: "owner/repo#123",
			want:  &PRRef{Owner: "owner", Repo: "repo", Number: 123},
		},
		{
			name:  "short format with .git suffix",
			input: "owner/repo.git#123",
			want:  &PRRef{Owner: "owner", Repo: "repo", Number: 123},
		},
		{
			name:  "merge request URL format",
			input: "https://gitcode.com/owner/repo/merge_requests/123",
			want:  &PRRef{Owner: "owner", Repo: "repo", Number: 123},
		},
		{
			name:  "pulls URL format remains supported",
			input: "https://gitcode.com/owner/repo/pulls/123",
			want:  &PRRef{Owner: "owner", Repo: "repo", Number: 123},
		},
		{
			name:  "merge request URL format with trailing slash",
			input: "https://gitcode.com/owner/repo/merge_requests/123/",
			want:  &PRRef{Owner: "owner", Repo: "repo", Number: 123},
		},
		{
			name:  "merge request URL format with additional path",
			input: "https://gitcode.com/owner/repo/merge_requests/123/commits",
			want:  &PRRef{Owner: "owner", Repo: "repo", Number: 123},
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid format",
			input:   "invalid-format",
			wantErr: true,
		},
		{
			name:    "missing number",
			input:   "owner/repo#",
			wantErr: true,
		},
		{
			name:    "invalid number",
			input:   "owner/repo#abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePRRef(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePRRef(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Number != tt.want.Number {
				t.Errorf("ParsePRRef(%q) Number = %v, want %v", tt.input, got.Number, tt.want.Number)
			}
			if !tt.wantErr && got.Owner != tt.want.Owner {
				t.Errorf("ParsePRRef(%q) Owner = %v, want %v", tt.input, got.Owner, tt.want.Owner)
			}
			if !tt.wantErr && got.Repo != tt.want.Repo {
				t.Errorf("ParsePRRef(%q) Repo = %v, want %v", tt.input, got.Repo, tt.want.Repo)
			}
		})
	}
}

func TestBuildSyncBranch(t *testing.T) {
	branch := buildSyncBranch("owner", "repo", 123)
	// Should match pattern: sync/pr-owner-repo-123-YYYYMMDD
	if !bytes.HasPrefix([]byte(branch), []byte("sync/pr-owner-repo-123-")) {
		t.Errorf("buildSyncBranch() = %q, expected prefix sync/pr-owner-repo-123-", branch)
	}
}

func TestBuildSyncBody(t *testing.T) {
	pr := &api.PullRequest{
		Title:   "Test PR",
		Body:    "Original body",
		HTMLURL: "https://gitcode.com/owner/repo/merge_requests/123",
	}
	sourcePR := &PRRef{Owner: "source-owner", Repo: "source-repo", Number: 123}
	targetRepo := "target-owner/target-repo"

	body := buildSyncBody(pr, sourcePR, targetRepo)

	// Should contain original body
	if !bytes.Contains([]byte(body), []byte("Original body")) {
		t.Errorf("buildSyncBody() missing original body")
	}
	// Should contain sync info
	if !bytes.Contains([]byte(body), []byte("Synced from")) {
		t.Errorf("buildSyncBody() missing sync info")
	}
	if !bytes.Contains([]byte(body), []byte("source-owner/source-repo#123")) {
		t.Errorf("buildSyncBody() missing source PR reference")
	}
}

func TestNewCmdSync(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	f := &cmdutil.Factory{
		IOStreams: io,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{}, nil
		},
	}

	cmd := NewCmdSync(f, nil)
	if cmd == nil {
		t.Fatal("NewCmdSync returned nil")
	}
	if cmd.Use != "sync" {
		t.Errorf("NewCmdSync().Use = %q, want sync", cmd.Use)
	}

	// Check required flags
	requiredFlags := []string{"source-pr", "target-repo"}
	for _, flag := range requiredFlags {
		if cmd.Flags().Lookup(flag) == nil {
			t.Errorf("NewCmdSync() missing required flag %q", flag)
		}
	}

	// Check optional flags
	optionalFlags := []string{"base", "title", "body", "draft", "json"}
	for _, flag := range optionalFlags {
		if cmd.Flags().Lookup(flag) == nil {
			t.Errorf("NewCmdSync() missing optional flag %q", flag)
		}
	}
}

func TestRepositoryGitURL(t *testing.T) {
	url := repositoryGitURL("owner", "repo")
	expected := "https://gitcode.com/owner/repo.git"
	if url != expected {
		t.Errorf("repositoryGitURL() = %q, want %q", url, expected)
	}
}

func TestPullRequestWebURL(t *testing.T) {
	url := pullRequestWebURL("owner", "repo", 123)
	expected := "https://gitcode.com/owner/repo/merge_requests/123"
	if url != expected {
		t.Errorf("pullRequestWebURL() = %q, want %q", url, expected)
	}
}

func TestSyncCommitsPreservesCommitBoundaries(t *testing.T) {
	commits := []api.Commit{
		{SHA: "1111111111111111111111111111111111111111", Message: "first"},
		{SHA: "2222222222222222222222222222222222222222", Message: "second"},
	}

	var calls [][]string
	runGitInDir := func(dir string, args ...string) (string, error) {
		if dir != "/tmp/workdir" {
			t.Fatalf("runGitInDir dir = %q, want /tmp/workdir", dir)
		}
		calls = append(calls, append([]string(nil), args...))
		return "", nil
	}

	commitsSynced, conflictError := syncCommits(runGitInDir, "/tmp/workdir", commits)
	if conflictError != "" {
		t.Fatalf("syncCommits() unexpected conflict error = %q", conflictError)
	}
	if commitsSynced != 2 {
		t.Fatalf("syncCommits() commitsSynced = %d, want 2", commitsSynced)
	}

	wantCalls := [][]string{
		{"cherry-pick", "-x", commits[0].SHA},
		{"cherry-pick", "-x", commits[1].SHA},
	}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Fatalf("syncCommits() calls = %#v, want %#v", calls, wantCalls)
	}
}

func TestSyncCommitsReportsActualCountOnConflict(t *testing.T) {
	commits := []api.Commit{
		{SHA: "1111111111111111111111111111111111111111", Message: "first"},
		{SHA: "2222222222222222222222222222222222222222", Message: "second"},
	}

	var calls [][]string
	runGitInDir := func(dir string, args ...string) (string, error) {
		calls = append(calls, append([]string(nil), args...))
		if len(args) == 3 && args[0] == "cherry-pick" && args[2] == commits[1].SHA {
			return "", errors.New("conflict")
		}
		return "", nil
	}

	commitsSynced, conflictError := syncCommits(runGitInDir, "/tmp/workdir", commits)
	if commitsSynced != 1 {
		t.Fatalf("syncCommits() commitsSynced = %d, want 1", commitsSynced)
	}
	if conflictError == "" {
		t.Fatal("syncCommits() conflictError = empty, want non-empty")
	}

	wantCalls := [][]string{
		{"cherry-pick", "-x", commits[0].SHA},
		{"cherry-pick", "-x", commits[1].SHA},
		{"cherry-pick", "--abort"},
	}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Fatalf("syncCommits() calls = %#v, want %#v", calls, wantCalls)
	}
}

func TestWriteSyncResultJSONReturnsConflictExitCode(t *testing.T) {
	f := cmdutil.TestFactory()
	out := &strings.Builder{}
	f.IOStreams.Out = out

	err := writeSyncResult(&SyncOptions{
		IO:   f.IOStreams,
		JSON: true,
	}, SyncResult{
		SourcePR:      "owner/repo#1",
		TargetRepo:    "target/repo",
		SyncBranch:    "sync/test",
		CommitsSynced: 1,
	}, cmdutil.NewCLIError(cmdutil.ExitConflict, "conflict detected", nil))
	if err == nil {
		t.Fatal("writeSyncResult() error = nil, want conflict error")
	}
	if got := cmdutil.ExitCode(err); got != cmdutil.ExitConflict {
		t.Fatalf("ExitCode() = %d, want %d", got, cmdutil.ExitConflict)
	}
	if !strings.Contains(out.String(), "\"conflict_error\": \"conflict detected\"") {
		t.Fatalf("output = %q", out.String())
	}
}
