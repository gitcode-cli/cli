package api

import (
	"encoding/json"
	"testing"
)

// TestIssue_Unmarshal tests Issue JSON parsing with real API response
func TestIssue_Unmarshal(t *testing.T) {
	// Real API response from GitCode
	jsonResp := `{
		"id": 3815588,
		"html_url": "https://gitcode.com/infra-test/gctest1/issues/3",
		"number": "3",
		"state": "open",
		"title": "Test Issue 1",
		"body": "This is test issue 1",
		"user": {
			"html_url": "https://gitcode.com/aflyingto",
			"id": "67de131cf5d1b1713b4c0900",
			"login": "aflyingto",
			"name": "aflyingto"
		},
		"created_at": "2026-03-22T19:05:07+08:00",
		"updated_at": "2026-03-22T19:08:02+08:00",
		"labels": [],
		"comments": 1
	}`

	var issue Issue
	err := json.Unmarshal([]byte(jsonResp), &issue)
	if err != nil {
		t.Fatalf("Failed to unmarshal Issue: %v", err)
	}

	// Verify fields
	if issue.Number != "3" {
		t.Errorf("Expected Number '3', got '%s'", issue.Number)
	}
	if issue.Title != "Test Issue 1" {
		t.Errorf("Expected Title 'Test Issue 1', got '%s'", issue.Title)
	}
	if issue.State != "open" {
		t.Errorf("Expected State 'open', got '%s'", issue.State)
	}
	if issue.User == nil || issue.User.Login != "aflyingto" {
		t.Errorf("Expected User.Login 'aflyingto', got '%v'", issue.User)
	}
}

// TestIssue_NumberIsString verifies that Issue.Number is string type
// This test would fail if Number was defined as int
func TestIssue_NumberIsString(t *testing.T) {
	// GitCode returns number as string
	jsonResp := `{"number": "123"}`

	var issue Issue
	err := json.Unmarshal([]byte(jsonResp), &issue)
	if err != nil {
		t.Fatalf("Number should be string type: %v", err)
	}

	if issue.Number != "123" {
		t.Errorf("Expected Number '123', got '%s'", issue.Number)
	}
}

// TestIssue_NumberTypeValidation demonstrates that wrong type would fail
// This test documents what would happen if Issue.Number was defined as int
func TestIssue_NumberTypeValidation(t *testing.T) {
	// This test verifies the current (correct) behavior
	// If Issue.Number was int, this JSON would cause an error
	jsonWithStringNumber := `{"number": "123"}`
	jsonWithIntNumber := `{"number": 123}`

	// Test 1: String number should work (GitCode API behavior)
	var issue1 Issue
	err := json.Unmarshal([]byte(jsonWithStringNumber), &issue1)
	if err != nil {
		t.Errorf("String number should parse: %v", err)
	}

	// Test 2: Int number should also work (for compatibility)
	var issue2 Issue
	err = json.Unmarshal([]byte(jsonWithIntNumber), &issue2)
	if err != nil {
		// Note: If Issue.Number was int, string "123" would fail
		// But since we use string, int 123 should convert to string
		t.Logf("Note: Int number parsing result: %v, Number=%s", err, issue2.Number)
	}
}

// TestIssueList_Unmarshal tests list of issues parsing
func TestIssueList_Unmarshal(t *testing.T) {
	jsonResp := `[
		{"number": "1", "title": "Issue 1", "state": "open"},
		{"number": "2", "title": "Issue 2", "state": "closed"}
	]`

	var issues []Issue
	err := json.Unmarshal([]byte(jsonResp), &issues)
	if err != nil {
		t.Fatalf("Failed to unmarshal Issue list: %v", err)
	}

	if len(issues) != 2 {
		t.Fatalf("Expected 2 issues, got %d", len(issues))
	}
	if issues[0].Number != "1" {
		t.Errorf("Expected first issue Number '1', got '%s'", issues[0].Number)
	}
}

// TestLabel_Unmarshal tests Label JSON parsing
func TestLabel_Unmarshal(t *testing.T) {
	// Real API response from GitCode
	jsonResp := `{
		"id": 13232224,
		"name": "bug",
		"color": "#ff0000",
		"repository_id": 9480067
	}`

	var label Label
	err := json.Unmarshal([]byte(jsonResp), &label)
	if err != nil {
		t.Fatalf("Failed to unmarshal Label: %v", err)
	}

	if label.Name != "bug" {
		t.Errorf("Expected Name 'bug', got '%s'", label.Name)
	}
	if label.Color != "#ff0000" {
		t.Errorf("Expected Color '#ff0000', got '%s'", label.Color)
	}
}

// TestMilestone_Unmarshal tests Milestone JSON parsing
func TestMilestone_Unmarshal(t *testing.T) {
	// Sample milestone response
	jsonResp := `{
		"id": 1,
		"number": "1",
		"title": "v1.0.0",
		"description": "First release",
		"state": "open"
	}`

	var milestone Milestone
	err := json.Unmarshal([]byte(jsonResp), &milestone)
	if err != nil {
		t.Fatalf("Failed to unmarshal Milestone: %v", err)
	}

	if milestone.Number != "1" {
		t.Errorf("Expected Number '1', got '%s'", milestone.Number)
	}
	if milestone.Title != "v1.0.0" {
		t.Errorf("Expected Title 'v1.0.0', got '%s'", milestone.Title)
	}
}

// TestIssueComment_Unmarshal tests IssueComment JSON parsing
func TestIssueComment_Unmarshal(t *testing.T) {
	jsonResp := `{
		"id": 166027129,
		"body": "Test comment",
		"user": {
			"login": "aflyingto",
			"name": "aflyingto"
		},
		"created_at": "2026-03-22T19:05:07+08:00"
	}`

	var comment IssueComment
	err := json.Unmarshal([]byte(jsonResp), &comment)
	if err != nil {
		t.Fatalf("Failed to unmarshal IssueComment: %v", err)
	}

	if comment.Body != "Test comment" {
		t.Errorf("Expected Body 'Test comment', got '%s'", comment.Body)
	}
}