package review

import (
	"net/http"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestReviewRun_Approve(t *testing.T) {
	io, _, out, _ := testutil.NewTestIOStreams()
	restoreToken := testutil.SetTestToken()
	defer restoreToken()

	reviewCalled := false
	commentCalled := false

	err := reviewRun(&ReviewOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return &http.Client{}, nil },
		Repository: "owner/repo",
		Number:     123,
		Approve:    true,
		ReviewPR: func(client *api.Client, owner, repo string, number int, opts *api.ReviewPROptions) error {
			reviewCalled = true
			if owner != "owner" || repo != "repo" || number != 123 {
				t.Fatalf("unexpected review args: %s/%s #%d", owner, repo, number)
			}
			if opts == nil || opts.Force {
				t.Fatalf("unexpected review opts: %+v", opts)
			}
			return nil
		},
		CreatePRComment: func(client *api.Client, owner, repo string, number int, opts *api.CreatePRCommentOptions) (*api.PRComment, error) {
			commentCalled = true
			return &api.PRComment{Body: opts.Body}, nil
		},
	})
	if err != nil {
		t.Fatalf("reviewRun() unexpected error: %v", err)
	}
	if !reviewCalled {
		t.Fatal("expected approve to call ReviewPR")
	}
	if commentCalled {
		t.Fatal("did not expect approve without comment to call CreatePRComment")
	}
	if !strings.Contains(out.String(), "approved PR #123") {
		t.Fatalf("expected approval output, got %q", out.String())
	}
}

func TestReviewRun_ApproveWithComment(t *testing.T) {
	io, _, out, _ := testutil.NewTestIOStreams()
	restoreToken := testutil.SetTestToken()
	defer restoreToken()

	reviewCalled := false
	commentCalled := false

	err := reviewRun(&ReviewOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return &http.Client{}, nil },
		Repository: "owner/repo",
		Number:     123,
		Approve:    true,
		Comment:    "LGTM",
		ReviewPR: func(client *api.Client, owner, repo string, number int, opts *api.ReviewPROptions) error {
			reviewCalled = true
			return nil
		},
		CreatePRComment: func(client *api.Client, owner, repo string, number int, opts *api.CreatePRCommentOptions) (*api.PRComment, error) {
			commentCalled = true
			if opts.Body != "LGTM" {
				t.Fatalf("unexpected comment body %q", opts.Body)
			}
			return &api.PRComment{Body: opts.Body}, nil
		},
	})
	if err != nil {
		t.Fatalf("reviewRun() unexpected error: %v", err)
	}
	if !reviewCalled || !commentCalled {
		t.Fatalf("expected approve with comment to call both comment and review, got review=%v comment=%v", reviewCalled, commentCalled)
	}
	if !strings.Contains(out.String(), "LGTM") {
		t.Fatalf("expected comment echoed in output, got %q", out.String())
	}
}

func TestReviewRun_RequestUnsupported(t *testing.T) {
	io, _, _, _ := testutil.NewTestIOStreams()
	restoreToken := testutil.SetTestToken()
	defer restoreToken()

	err := reviewRun(&ReviewOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return &http.Client{}, nil },
		Repository: "owner/repo",
		Number:     123,
		Request:    true,
		ReviewPR: func(client *api.Client, owner, repo string, number int, opts *api.ReviewPROptions) error {
			t.Fatal("did not expect request to call ReviewPR")
			return nil
		},
		CreatePRComment: func(client *api.Client, owner, repo string, number int, opts *api.CreatePRCommentOptions) (*api.PRComment, error) {
			t.Fatal("did not expect request to call CreatePRComment")
			return nil, nil
		},
	})
	if err == nil {
		t.Fatal("expected error for unsupported request changes")
	}
	if !strings.Contains(err.Error(), "not supported by the current GitCode API") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewCmdReview(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "approve PR",
			args:    []string{"123", "--approve"},
			wantErr: false,
		},
		{
			name:    "request changes",
			args:    []string{"123", "--request"},
			wantErr: false,
		},
		{
			name:    "review with comment",
			args:    []string{"123", "--comment", "Looks good"},
			wantErr: false,
		},
		{
			name:    "approve with body",
			args:    []string{"123", "--approve", "--comment", "LGTM"},
			wantErr: false,
		},
		{
			name:    "no PR number",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCmdReview(cmdutil.TestFactory(), func(opts *ReviewOptions) error {
				return nil
			})
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
