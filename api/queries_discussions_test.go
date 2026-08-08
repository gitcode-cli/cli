package api

import (
	"net/http"
	"testing"
)

func TestListOrgDiscussionsBuildsPathAndParams(t *testing.T) {
	var gotPath, gotQuery, gotMethod string
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		gotQuery = req.URL.RawQuery
		gotMethod = req.Method
		return authTestResponse(http.StatusOK, `[]`), nil
	})

	discussions, err := ListOrgDiscussions(client, "my-org", &ListOrgDiscussionsOptions{
		Page:      2,
		PerPage:   50,
		Sort:      "comment_size",
		Direction: "asc",
		Search:    "release",
	})
	if err != nil {
		t.Fatalf("ListOrgDiscussions() error = %v", err)
	}
	if len(discussions) != 0 {
		t.Fatalf("len(discussions) = %d, want 0 for empty array", len(discussions))
	}
	if gotMethod != http.MethodGet {
		t.Fatalf("method = %q, want GET", gotMethod)
	}
	// Path must target the v5 orgs discuss endpoint (per gc-api-doc).
	if gotPath != "/api/v5/orgs/my-org/discuss" {
		t.Fatalf("path = %q, want /api/v5/orgs/my-org/discuss", gotPath)
	}
	// Query params must match the official OpenAPI names.
	for _, want := range []string{"page=2", "per_page=50", "sort=comment_size", "direction=asc", "search=release"} {
		if !containsParam(gotQuery, want) {
			t.Fatalf("query = %q, want to contain %q", gotQuery, want)
		}
	}
}

func TestListOrgDiscussionsNoParamsOmitsQuery(t *testing.T) {
	var gotQuery string
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotQuery = req.URL.RawQuery
		return authTestResponse(http.StatusOK, `[]`), nil
	})

	if _, err := ListOrgDiscussions(client, "my-org", nil); err != nil {
		t.Fatalf("ListOrgDiscussions() error = %v", err)
	}
	if gotQuery != "" {
		t.Fatalf("query = %q, want empty when no options", gotQuery)
	}
}

func TestListOrgDiscussionsParsesResponse(t *testing.T) {
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		return authTestResponse(http.StatusOK, `[{"id":"abc","number":7,"title":"release plan","comment_total":3,"is_closed":0}]`), nil
	})

	discussions, err := ListOrgDiscussions(client, "my-org", &ListOrgDiscussionsOptions{PerPage: 20})
	if err != nil {
		t.Fatalf("ListOrgDiscussions() error = %v", err)
	}
	if len(discussions) != 1 {
		t.Fatalf("len = %d, want 1", len(discussions))
	}
	d := discussions[0]
	if d.ID != "abc" || d.Number != 7 || d.Title != "release plan" || d.CommentTotal != 3 {
		t.Fatalf("discussion = %+v, want id=abc number=7 title=release plan comment_total=3", d)
	}
}

func TestGetOrgDiscussionBuildsPath(t *testing.T) {
	var gotPath, gotMethod string
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		gotMethod = req.Method
		return authTestResponse(http.StatusOK, `{"id":"abc","number":42,"title":"q","md_content":"body"}`), nil
	})

	d, err := GetOrgDiscussion(client, "my-org", 42)
	if err != nil {
		t.Fatalf("GetOrgDiscussion() error = %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Fatalf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/api/v5/orgs/my-org/discuss/42" {
		t.Fatalf("path = %q, want /api/v5/orgs/my-org/discuss/42", gotPath)
	}
	if d.ID != "abc" || d.Number != 42 || d.Title != "q" || d.MdContent != "body" {
		t.Fatalf("discussion = %+v, want id=abc number=42 title=q md_content=body", d)
	}
}

func TestGetOrgDiscussionNotFoundReturnsError(t *testing.T) {
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		return authTestResponse(http.StatusNotFound, `{"message":"not found"}`), nil
	})

	_, err := GetOrgDiscussion(client, "my-org", 999)
	if err == nil {
		t.Fatal("GetOrgDiscussion() error = nil, want error for 404")
	}
}

func TestListRepoDiscussionsParsesResponse(t *testing.T) {
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		return authTestResponse(http.StatusOK, `[{"id":"a","number":7,"title":"release plan","comment_total":3,"is_closed":0}]`), nil
	})

	discussions, err := ListRepoDiscussions(client, "owner", "repo", &ListOrgDiscussionsOptions{PerPage: 20})
	if err != nil {
		t.Fatalf("ListRepoDiscussions() error = %v", err)
	}
	if len(discussions) != 1 {
		t.Fatalf("len = %d, want 1", len(discussions))
	}
	d := discussions[0]
	if d.ID != "a" || d.Number != 7 || d.Title != "release plan" || d.CommentTotal != 3 {
		t.Fatalf("discussion = %+v, want id=a number=7 title=release plan comment_total=3", d)
	}
}

func TestListRepoDiscussionsBuildsPathAndParams(t *testing.T) {
	var gotPath, gotQuery, gotMethod string
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		gotQuery = req.URL.RawQuery
		gotMethod = req.Method
		return authTestResponse(http.StatusOK, `[]`), nil
	})

	if _, err := ListRepoDiscussions(client, "owner", "repo", &ListOrgDiscussionsOptions{
		Page: 2, PerPage: 50, Sort: "comment_size", Direction: "asc", Search: "release",
	}); err != nil {
		t.Fatalf("ListRepoDiscussions() error = %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Fatalf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/api/v5/repos/owner/repo/discuss" {
		t.Fatalf("path = %q, want /api/v5/repos/owner/repo/discuss", gotPath)
	}
	for _, want := range []string{"page=2", "per_page=50", "sort=comment_size", "direction=asc", "search=release"} {
		if !containsParam(gotQuery, want) {
			t.Fatalf("query = %q, want to contain %q", gotQuery, want)
		}
	}
}

func TestGetRepoDiscussionBuildsPath(t *testing.T) {
	var gotPath, gotMethod string
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		gotMethod = req.Method
		return authTestResponse(http.StatusOK, `{"id":"a","number":42,"title":"q","md_content":"body"}`), nil
	})

	d, err := GetRepoDiscussion(client, "owner", "repo", 42)
	if err != nil {
		t.Fatalf("GetRepoDiscussion() error = %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Fatalf("method = %q, want GET", gotMethod)
	}
	if gotPath != "/api/v5/repos/owner/repo/discuss/42" {
		t.Fatalf("path = %q, want /api/v5/repos/owner/repo/discuss/42", gotPath)
	}
	if d.ID != "a" || d.Number != 42 || d.Title != "q" || d.MdContent != "body" {
		t.Fatalf("discussion = %+v, want id=a number=42 title=q md_content=body", d)
	}
}

// containsParam checks that the url-encoded query contains the given k=v pair.
func containsParam(query, want string) bool {
	for _, p := range splitQuery(query) {
		if p == want {
			return true
		}
	}
	return false
}

func splitQuery(query string) []string {
	if query == "" {
		return nil
	}
	var out []string
	start := 0
	for i := 0; i < len(query); i++ {
		if query[i] == '&' {
			out = append(out, query[start:i])
			start = i + 1
		}
	}
	out = append(out, query[start:])
	return out
}
