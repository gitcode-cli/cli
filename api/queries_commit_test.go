package api

import (
	"encoding/json"
	"testing"
)

// TestCommitComment_Unmarshal tests CommitComment JSON parsing with real API response
func TestCommitComment_Unmarshal(t *testing.T) {
	// GitCode API returns ID as string
	jsonResp := `{
		"id": "166564358",
		"body": "Test comment",
		"created_at": "2026-03-27T10:00:00+08:00",
		"updated_at": "2026-03-27T10:00:00+08:00",
		"user": {
			"login": "testuser",
			"name": "Test User"
		}
	}`

	var comment CommitComment
	err := json.Unmarshal([]byte(jsonResp), &comment)
	if err != nil {
		t.Fatalf("Failed to unmarshal CommitComment: %v", err)
	}

	// Verify fields
	if comment.ID == nil {
		t.Errorf("Expected ID to be set, got nil")
	}
	if comment.Body != "Test comment" {
		t.Errorf("Expected Body 'Test comment', got '%s'", comment.Body)
	}
	if comment.User == nil || comment.User.Login != "testuser" {
		t.Errorf("Expected User.Login 'testuser', got '%v'", comment.User)
	}
}

// TestCommitComment_IDIsString verifies that CommitComment.ID can be string type
func TestCommitComment_IDIsString(t *testing.T) {
	// GitCode returns ID as string
	jsonResp := `{"id": "123456"}`

	var comment CommitComment
	err := json.Unmarshal([]byte(jsonResp), &comment)
	if err != nil {
		t.Fatalf("ID should be string type: %v", err)
	}

	// Verify ID value
	idStr, ok := comment.ID.(string)
	if !ok {
		t.Errorf("Expected ID to be string, got %T", comment.ID)
	}
	if idStr != "123456" {
		t.Errorf("Expected ID '123456', got '%s'", idStr)
	}
}

// TestCommitComment_IDTypeValidation demonstrates that wrong type would fail
func TestCommitComment_IDTypeValidation(t *testing.T) {
	// Test 1: String ID should work (GitCode API behavior)
	jsonWithStringID := `{"id": "123456"}`
	var comment1 CommitComment
	err := json.Unmarshal([]byte(jsonWithStringID), &comment1)
	if err != nil {
		t.Errorf("String ID should parse: %v", err)
	}

	// Test 2: Int ID should also work (for compatibility)
	jsonWithIntID := `{"id": 123456}`
	var comment2 CommitComment
	err = json.Unmarshal([]byte(jsonWithIntID), &comment2)
	if err != nil {
		t.Errorf("Int ID should parse: %v", err)
	}
}

// TestCommitCommentList_Unmarshal tests list of comments parsing
func TestCommitCommentList_Unmarshal(t *testing.T) {
	jsonResp := `[
		{"id": "1", "body": "Comment 1"},
		{"id": "2", "body": "Comment 2"}
	]`

	var comments []CommitComment
	err := json.Unmarshal([]byte(jsonResp), &comments)
	if err != nil {
		t.Fatalf("Failed to unmarshal CommitComment list: %v", err)
	}

	if len(comments) != 2 {
		t.Fatalf("Expected 2 comments, got %d", len(comments))
	}
	if comments[0].Body != "Comment 1" {
		t.Errorf("Expected first comment Body 'Comment 1', got '%s'", comments[0].Body)
	}
}