package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
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

func TestListUserReposBuildsQuery(t *testing.T) {
	var gotPath string
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		if req.URL.RawQuery != "" {
			gotPath += "?" + req.URL.RawQuery
		}
		return authTestResponse(http.StatusOK, `[]`), nil
	})

	_, err := ListUserRepos(client, &RepoListOptions{
		Visibility:  "private",
		Affiliation: "owner",
		Type:        "member",
		Sort:        "pushed",
		Direction:   "desc",
		PerPage:     50,
		Page:        2,
	})
	if err != nil {
		t.Fatalf("ListUserRepos() error = %v", err)
	}

	assertRepoListRequest(t, gotPath, "/api/v5/user/repos", map[string]string{
		"visibility":  "private",
		"affiliation": "owner",
		"type":        "member",
		"sort":        "pushed",
		"direction":   "desc",
		"per_page":    "50",
		"page":        "2",
	})
}

func TestListOrgReposBuildsQuery(t *testing.T) {
	var gotPath string
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		if req.URL.RawQuery != "" {
			gotPath += "?" + req.URL.RawQuery
		}
		return authTestResponse(http.StatusOK, `[]`), nil
	})

	_, err := ListOrgRepos(client, "infra-test", &RepoListOptions{
		Visibility: "public",
		PerPage:    25,
	})
	if err != nil {
		t.Fatalf("ListOrgRepos() error = %v", err)
	}

	assertRepoListRequest(t, gotPath, "/api/v5/orgs/infra-test/repos", map[string]string{
		"visibility": "public",
		"per_page":   "25",
	})
}

func assertRepoListRequest(t *testing.T, gotPath, wantPath string, wantQuery map[string]string) {
	t.Helper()

	if gotPath == "" {
		t.Fatal("request path was empty")
	}
	if !strings.HasPrefix(gotPath, wantPath) {
		t.Fatalf("request path = %q, want prefix %q", gotPath, wantPath)
	}

	rawQuery := ""
	if len(gotPath) > len(wantPath) {
		rawQuery = strings.TrimPrefix(gotPath[len(wantPath):], "?")
	}
	query, err := url.ParseQuery(rawQuery)
	if err != nil {
		t.Fatalf("url.ParseQuery() error = %v", err)
	}
	for key, want := range wantQuery {
		if got := query.Get(key); got != want {
			t.Fatalf("query[%q] = %q, want %q", key, got, want)
		}
	}
	if len(query) != len(wantQuery) {
		t.Fatalf("query = %#v, want %#v", query, wantQuery)
	}
}
