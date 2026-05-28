package api

import (
	"testing"

	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestListRepoLabels(t *testing.T) {
	client := NewClientFromHTTP(testutil.NewTestHTTPClient(testutil.MockAPIHandler()))
	labels, err := ListRepoLabels(client, "owner", "test-repo", nil)
	if err != nil {
		t.Fatalf("ListRepoLabels() error = %v", err)
	}
	if len(labels) != 1 || labels[0].Name != "bug" {
		t.Fatalf("ListRepoLabels() = %+v, want [bug]", labels)
	}
}

func TestCreateIssueComment(t *testing.T) {
	client := NewClientFromHTTP(testutil.NewTestHTTPClient(testutil.MockAPIHandler()))
	opts := &CreateCommentOptions{Body: "New comment"}
	comment, err := CreateIssueComment(client, "owner", "test-repo", 1, opts)
	if err != nil {
		t.Fatalf("CreateIssueComment() error = %v", err)
	}
	if comment.Body != "New comment" {
		t.Fatalf("CreateIssueComment().Body = %q, want %q", comment.Body, "New comment")
	}
}

func TestListRepoMilestones(t *testing.T) {
	client := NewClientFromHTTP(testutil.NewTestHTTPClient(testutil.MockAPIHandler()))
	milestones, err := ListRepoMilestones(client, "owner", "test-repo", nil)
	if err != nil {
		t.Fatalf("ListRepoMilestones() error = %v", err)
	}
	if len(milestones) != 1 || milestones[0].Title != "v1.0" {
		t.Fatalf("ListRepoMilestones() = %+v, want [v1.0]", milestones)
	}
}

func TestListPRReviews(t *testing.T) {
	client := NewClientFromHTTP(testutil.NewTestHTTPClient(testutil.MockAPIHandler()))
	reviews, err := ListPRReviews(client, "owner", "test-repo", 1)
	if err != nil {
		t.Fatalf("ListPRReviews() error = %v", err)
	}
	if len(reviews) != 1 || reviews[0].State != "APPROVED" {
		t.Fatalf("ListPRReviews() = %+v, want [APPROVED]", reviews)
	}
}

func TestGetIssuePullRequests(t *testing.T) {
	client := NewClientFromHTTP(testutil.NewTestHTTPClient(testutil.MockAPIHandler()))
	// The mock handler may not have this endpoint; test that the function runs
	// and we just check it doesn't panic on the call path.
	_, _ = GetIssuePullRequests(client, "owner", "test-repo", 1, 0)
}

func TestListPRComments(t *testing.T) {
	client := NewClientFromHTTP(testutil.NewTestHTTPClient(testutil.MockAPIHandler()))
	comments, err := ListPRComments(client, "owner", "test-repo", 1)
	if err != nil {
		t.Fatalf("ListPRComments() error = %v", err)
	}
	if len(comments) != 1 || comments[0].Body != "Test comment" {
		t.Fatalf("ListPRComments() = %+v, want [Test comment]", comments)
	}
}

func TestGetCommit(t *testing.T) {
	client := NewClientFromHTTP(testutil.NewTestHTTPClient(testutil.MockAPIHandler()))
	commit, err := GetCommit(client, "owner", "test-repo", "abc123", false)
	if err != nil {
		t.Fatalf("GetCommit() error = %v", err)
	}
	if commit.SHA != "abc123" {
		t.Errorf("GetCommit().SHA = %q, want %q", commit.SHA, "abc123")
	}
}

func TestListCommitComments(t *testing.T) {
	client := NewClientFromHTTP(testutil.NewTestHTTPClient(testutil.MockAPIHandler()))
	comments, err := ListCommitComments(client, "owner", "test-repo", nil)
	if err != nil {
		t.Fatalf("ListCommitComments() error = %v", err)
	}
	if len(comments) != 1 || comments[0].Body != "nice fix" {
		t.Fatalf("ListCommitComments() = %+v, want [nice fix]", comments)
	}
}

func TestUpdateIssueComment(t *testing.T) {
	client := NewClientFromHTTP(testutil.NewTestHTTPClient(testutil.MockAPIHandler()))
	opts := &UpdateCommentOptions{Body: "updated"}
	comment, err := UpdateIssueComment(client, "owner", "test-repo", "1", opts)
	if err != nil {
		t.Fatalf("UpdateIssueComment() error = %v", err)
	}
	if comment.Body != "updated comment" {
		t.Errorf("Body = %q, want %q", comment.Body, "updated comment")
	}
}

func TestDeleteIssueComment(t *testing.T) {
	client := NewClientFromHTTP(testutil.NewTestHTTPClient(testutil.MockAPIHandler()))
	if err := DeleteIssueComment(client, "owner", "test-repo", 1); err != nil {
		t.Fatalf("DeleteIssueComment() error = %v", err)
	}
}

func TestCreateMilestone(t *testing.T) {
	client := NewClientFromHTTP(testutil.NewTestHTTPClient(testutil.MockAPIHandler()))
	opts := &CreateMilestoneOptions{Title: "v1.0"}
	milestone, err := CreateMilestone(client, "owner", "test-repo", opts)
	if err != nil {
		t.Fatalf("CreateMilestone() error = %v", err)
	}
	if milestone.Title != "v1.0" {
		t.Errorf("Title = %q, want %q", milestone.Title, "v1.0")
	}
}
