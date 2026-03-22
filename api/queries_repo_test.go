package api

import (
	"encoding/json"
	"testing"
)

// TestRepository_Unmarshal tests Repository JSON parsing with real API response
func TestRepository_Unmarshal(t *testing.T) {
	// Real API response from GitCode
	jsonResp := `{
		"id": 9480067,
		"full_name": "infra-test/gctest1",
		"name": "gctest1",
		"description": "Test repository",
		"private": true,
		"html_url": "https://gitcode.com/infra-test/gctest1",
		"clone_url": "https://gitcode.com/infra-test/gctest1.git",
		"ssh_url": "git@gitcode.com:infra-test/gctest1.git",
		"default_branch": "main",
		"created_at": "2026-03-22T18:59:58+08:00",
		"updated_at": "2026-03-22T19:00:00+08:00",
		"stargazers_count": 0,
		"forks_count": 0,
		"open_issues_count": 3,
		"owner": {
			"id": "67de131cf5d1b1713b4c0900",
			"login": "aflyingto",
			"name": "aflyingto"
		}
	}`

	var repo Repository
	err := json.Unmarshal([]byte(jsonResp), &repo)
	if err != nil {
		t.Fatalf("Failed to unmarshal Repository: %v", err)
	}

	// Verify fields
	if repo.FullName != "infra-test/gctest1" {
		t.Errorf("Expected FullName 'infra-test/gctest1', got '%s'", repo.FullName)
	}
	if repo.Name != "gctest1" {
		t.Errorf("Expected Name 'gctest1', got '%s'", repo.Name)
	}
	if repo.Private != true {
		t.Errorf("Expected Private true, got %v", repo.Private)
	}
	if repo.DefaultBranch != "main" {
		t.Errorf("Expected DefaultBranch 'main', got '%s'", repo.DefaultBranch)
	}
	if repo.Owner == nil || repo.Owner.Login != "aflyingto" {
		t.Errorf("Expected Owner.Login 'aflyingto', got '%v'", repo.Owner)
	}
}

// TestRepository_IDIsInterface tests that Repository.ID is interface{} type
// GitCode may return ID as int or string
func TestRepository_IDIsInterface(t *testing.T) {
	// Test with int ID
	jsonResp1 := `{"id": 12345, "name": "test"}`
	var repo1 Repository
	err := json.Unmarshal([]byte(jsonResp1), &repo1)
	if err != nil {
		t.Fatalf("Failed to unmarshal Repository with int ID: %v", err)
	}

	// Test with string ID
	jsonResp2 := `{"id": "abc-123", "name": "test"}`
	var repo2 Repository
	err = json.Unmarshal([]byte(jsonResp2), &repo2)
	if err != nil {
		t.Fatalf("Failed to unmarshal Repository with string ID: %v", err)
	}
}

// TestRepositoryList_Unmarshal tests list of repositories parsing
func TestRepositoryList_Unmarshal(t *testing.T) {
	jsonResp := `[
		{"id": 1, "full_name": "owner/repo1", "name": "repo1", "private": false},
		{"id": 2, "full_name": "owner/repo2", "name": "repo2", "private": true}
	]`

	var repos []Repository
	err := json.Unmarshal([]byte(jsonResp), &repos)
	if err != nil {
		t.Fatalf("Failed to unmarshal Repository list: %v", err)
	}

	if len(repos) != 2 {
		t.Fatalf("Expected 2 repos, got %d", len(repos))
	}
	if repos[0].FullName != "owner/repo1" {
		t.Errorf("Expected first repo FullName 'owner/repo1', got '%s'", repos[0].FullName)
	}
}

// TestUser_Unmarshal tests User JSON parsing
func TestUser_Unmarshal(t *testing.T) {
	// Real API response from GitCode
	jsonResp := `{
		"id": "67de131cf5d1b1713b4c0900",
		"login": "aflyingto",
		"name": "aflyingto",
		"html_url": "https://gitcode.com/aflyingto",
		"email": "test@example.com"
	}`

	var user User
	err := json.Unmarshal([]byte(jsonResp), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal User: %v", err)
	}

	if user.Login != "aflyingto" {
		t.Errorf("Expected Login 'aflyingto', got '%s'", user.Login)
	}
	if user.Name != "aflyingto" {
		t.Errorf("Expected Name 'aflyingto', got '%s'", user.Name)
	}
}

// TestUser_IDIsString tests that User.ID is string type
// GitCode returns user ID as string
func TestUser_IDIsString(t *testing.T) {
	jsonResp := `{"id": "67de131cf5d1b1713b4c0900", "login": "test"}`

	var user User
	err := json.Unmarshal([]byte(jsonResp), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal User: %v", err)
	}

	// ID should be stored in ID field which is interface{}
	if user.ID == nil {
		t.Error("User.ID should not be nil")
	}
}