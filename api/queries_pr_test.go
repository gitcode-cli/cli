package api

import (
	"encoding/json"
	"testing"
)

// TestPullRequest_Unmarshal tests PullRequest JSON parsing with real API response
func TestPullRequest_Unmarshal(t *testing.T) {
	// Real API response from GitCode
	jsonResp := `{
		"id": 8434619,
		"html_url": "https://gitcode.com/infra-test/gctest1/merge_requests/2",
		"number": 2,
		"state": "open",
		"title": "Test PR 2",
		"body": "Test PR 2",
		"draft": false,
		"merged": false,
		"mergeable": false,
		"merged_at": "",
		"closed_at": "",
		"created_at": "2026-03-22T19:13:49+08:00",
		"updated_at": "2026-03-22T19:13:50+08:00",
		"user": {
			"id": "67de131cf5d1b1713b4c0900",
			"login": "aflyingto",
			"name": "aflyingto"
		},
		"head": {
			"ref": "feature-test-2",
			"sha": "f52cf5fa8383fd2b644b8104593d8e1c19aa3cc2",
			"label": "feature-test-2"
		},
		"base": {
			"ref": "main",
			"sha": "d69484cba7d5f81780d39378e72fa1f138d69c0c",
			"label": "main"
		},
		"mergeable_state": {
			"state": false,
			"conflict_passed": false
		}
	}`

	var pr PullRequest
	err := json.Unmarshal([]byte(jsonResp), &pr)
	if err != nil {
		t.Fatalf("Failed to unmarshal PullRequest: %v", err)
	}

	// Verify fields
	if pr.Number != 2 {
		t.Errorf("Expected Number 2, got %d", pr.Number)
	}
	if pr.Title != "Test PR 2" {
		t.Errorf("Expected Title 'Test PR 2', got '%s'", pr.Title)
	}
	if pr.State != "open" {
		t.Errorf("Expected State 'open', got '%s'", pr.State)
	}
	if pr.Draft != false {
		t.Errorf("Expected Draft false, got %v", pr.Draft)
	}
	if pr.Head == nil || pr.Head.Ref != "feature-test-2" {
		t.Errorf("Expected Head.Ref 'feature-test-2', got '%v'", pr.Head)
	}
	if pr.Base == nil || pr.Base.Ref != "main" {
		t.Errorf("Expected Base.Ref 'main', got '%v'", pr.Base)
	}
}

// TestPullRequest_NumberIsInt verifies that PullRequest.Number is int type
// This test would fail if Number was defined as string
func TestPullRequest_NumberIsInt(t *testing.T) {
	// GitCode returns PR number as int
	jsonResp := `{"number": 123}`

	var pr PullRequest
	err := json.Unmarshal([]byte(jsonResp), &pr)
	if err != nil {
		t.Fatalf("Number should be int type: %v", err)
	}

	if pr.Number != 123 {
		t.Errorf("Expected Number 123, got %d", pr.Number)
	}
}

// TestPullRequest_EmptyTimeFields tests handling of empty time fields
func TestPullRequest_EmptyTimeFields(t *testing.T) {
	// GitCode returns empty string for merged_at/closed_at when not merged/closed
	jsonResp := `{
		"number": 1,
		"merged_at": "",
		"closed_at": ""
	}`

	var pr PullRequest
	err := json.Unmarshal([]byte(jsonResp), &pr)
	if err != nil {
		t.Fatalf("Failed to unmarshal PullRequest with empty time fields: %v", err)
	}

	// Empty strings should be handled gracefully
	if pr.MergedAt == nil {
		t.Error("MergedAt should not be nil (even for empty string)")
	}
	if pr.ClosedAt == nil {
		t.Error("ClosedAt should not be nil (even for empty string)")
	}
}

// TestPullRequest_MergeableStateIsObject tests handling of mergeable_state as object
func TestPullRequest_MergeableStateIsObject(t *testing.T) {
	// GitCode returns mergeable_state as an object, not a string
	jsonResp := `{
		"number": 1,
		"mergeable_state": {
			"state": true,
			"conflict_passed": true,
			"branch_missing_passed": true
		}
	}`

	var pr PullRequest
	err := json.Unmarshal([]byte(jsonResp), &pr)
	if err != nil {
		t.Fatalf("Failed to unmarshal PullRequest with object mergeable_state: %v", err)
	}

	// mergeable_state should be parsed as interface{}
	if pr.MergeState == nil {
		t.Error("MergeState should not be nil")
	}
}

// TestPullRequestList_Unmarshal tests list of PRs parsing
func TestPullRequestList_Unmarshal(t *testing.T) {
	jsonResp := `[
		{"number": 1, "title": "PR 1", "state": "open"},
		{"number": 2, "title": "PR 2", "state": "closed"}
	]`

	var prs []PullRequest
	err := json.Unmarshal([]byte(jsonResp), &prs)
	if err != nil {
		t.Fatalf("Failed to unmarshal PR list: %v", err)
	}

	if len(prs) != 2 {
		t.Fatalf("Expected 2 PRs, got %d", len(prs))
	}
	if prs[0].Number != 1 {
		t.Errorf("Expected first PR Number 1, got %d", prs[0].Number)
	}
}

// TestPRBranch_Unmarshal tests PRBranch JSON parsing
func TestPRBranch_Unmarshal(t *testing.T) {
	jsonResp := `{
		"ref": "feature-branch",
		"sha": "abc123def456",
		"label": "feature-branch"
	}`

	var branch PRBranch
	err := json.Unmarshal([]byte(jsonResp), &branch)
	if err != nil {
		t.Fatalf("Failed to unmarshal PRBranch: %v", err)
	}

	if branch.Ref != "feature-branch" {
		t.Errorf("Expected Ref 'feature-branch', got '%s'", branch.Ref)
	}
	if branch.SHA != "abc123def456" {
		t.Errorf("Expected SHA 'abc123def456', got '%s'", branch.SHA)
	}
}

// TestPRComment_Unmarshal tests PRComment JSON parsing
func TestPRComment_Unmarshal(t *testing.T) {
	jsonResp := `{
		"id": 12345,
		"body": "LGTM!",
		"user": {
			"login": "reviewer"
		},
		"created_at": "2026-03-22T19:00:00+08:00"
	}`

	var comment PRComment
	err := json.Unmarshal([]byte(jsonResp), &comment)
	if err != nil {
		t.Fatalf("Failed to unmarshal PRComment: %v", err)
	}

	if comment.Body != "LGTM!" {
		t.Errorf("Expected Body 'LGTM!', got '%s'", comment.Body)
	}
}

// TestPRReview_Unmarshal tests PRReview JSON parsing
func TestPRReview_Unmarshal(t *testing.T) {
	jsonResp := `{
		"id": 999,
		"user": {
			"login": "reviewer"
		},
		"body": "Approved",
		"state": "APPROVED",
		"submitted_at": "2026-03-22T19:00:00+08:00"
	}`

	var review PRReview
	err := json.Unmarshal([]byte(jsonResp), &review)
	if err != nil {
		t.Fatalf("Failed to unmarshal PRReview: %v", err)
	}

	if review.State != "APPROVED" {
		t.Errorf("Expected State 'APPROVED', got '%s'", review.State)
	}
}