package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestListPullRequestsBuildsQuery(t *testing.T) {
	var gotPath string
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		if req.URL.RawQuery != "" {
			gotPath += "?" + req.URL.RawQuery
		}
		return authTestResponse(http.StatusOK, `[]`), nil
	})

	_, err := ListPullRequests(client, "owner", "repo", &PRListOptions{
		State:     "open",
		Head:      "feature/login",
		Base:      "main",
		Sort:      "updated",
		Direction: "asc",
		PerPage:   50,
		Page:      2,
	})
	if err != nil {
		t.Fatalf("ListPullRequests() error = %v", err)
	}

	assertPRListRequest(t, gotPath, "/api/v5/repos/owner/repo/pulls", map[string]string{
		"state":     "open",
		"head":      "feature/login",
		"base":      "main",
		"sort":      "updated",
		"direction": "asc",
		"per_page":  "50",
		"page":      "2",
	})
}

func TestUpdatePullRequestUsesFormEncoding(t *testing.T) {
	draft := false
	closeRelated := true

	var gotPath string
	var gotAuth string
	var gotContentType string
	var gotBody string

	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		if req.URL.RawQuery != "" {
			gotPath += "?" + req.URL.RawQuery
		}
		gotAuth = req.Header.Get("Authorization")
		gotContentType = req.Header.Get("Content-Type")
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		gotBody = string(body)
		return authTestResponse(http.StatusOK, `{"number":123,"title":"updated"}`), nil
	})
	client.SetToken("test-token", "test")

	_, err := UpdatePullRequest(client, "owner", "repo", 123, &UpdatePROptions{
		Title:             "updated",
		Body:              "body text",
		State:             "open",
		StateEvent:        "reopen",
		Base:              "main",
		Draft:             &draft,
		MilestoneNumber:   5,
		Labels:            []string{"type/feature", "risk/medium"},
		CloseRelatedIssue: &closeRelated,
	})
	if err != nil {
		t.Fatalf("UpdatePullRequest() error = %v", err)
	}

	assertNoAccessTokenQuery(t, gotPath)
	if gotPath != "/api/v5/repos/owner/repo/pulls/123" {
		t.Fatalf("request path = %q, want %q", gotPath, "/api/v5/repos/owner/repo/pulls/123")
	}
	if gotAuth != "Bearer test-token" {
		t.Fatalf("Authorization header = %q, want %q", gotAuth, "Bearer test-token")
	}
	if gotContentType != "application/x-www-form-urlencoded" {
		t.Fatalf("Content-Type = %q, want %q", gotContentType, "application/x-www-form-urlencoded")
	}

	expectedPairs := []string{
		"title=updated",
		"body=body+text",
		"state=open",
		"state_event=reopen",
		"base=main",
		"draft=false",
		"milestone_number=5",
		"labels%5B%5D=type%2Ffeature",
		"labels%5B%5D=risk%2Fmedium",
		"close_related_issue=true",
	}
	for _, pair := range expectedPairs {
		if !strings.Contains(gotBody, pair) {
			t.Fatalf("request body %q does not contain %q", gotBody, pair)
		}
	}
}

func TestBuildPRUpdateFormValuesNilOptions(t *testing.T) {
	formValues := buildPRUpdateFormValues(nil)
	if len(formValues) != 0 {
		t.Fatalf("expected empty form values, got %v", formValues)
	}
}

func TestPullRequestUnmarshal(t *testing.T) {
	jsonResp := `{
		"id": 8483763,
		"number": 95,
		"title": "feat: complete issue output contracts and view details",
		"body": "body",
		"state": "open",
		"html_url": "https://gitcode.com/gitcode-cli/cli/merge_requests/95",
		"draft": false,
		"created_at": "2026-04-07T10:20:21+08:00",
		"updated_at": "2026-04-07T11:08:27+08:00",
		"user": {"login": "aflyingto"},
		"labels": [{"name": "type/feature"}],
		"requested_reviewers": [{"login": "reviewer1"}]
	}`

	var pr PullRequest
	if err := json.Unmarshal([]byte(jsonResp), &pr); err != nil {
		t.Fatalf("Failed to unmarshal PullRequest: %v", err)
	}

	if pr.Number != 95 {
		t.Fatalf("Expected Number 95, got %d", pr.Number)
	}
	if pr.User == nil || pr.User.Login != "aflyingto" {
		t.Fatalf("Expected User.Login aflyingto, got %#v", pr.User)
	}
	if len(pr.Labels) != 1 || pr.Labels[0].Name != "type/feature" {
		t.Fatalf("Expected labels to include type/feature, got %#v", pr.Labels)
	}
	if len(pr.Reviewers) != 1 || pr.Reviewers[0].Login != "reviewer1" {
		t.Fatalf("Expected requested_reviewers to include reviewer1, got %#v", pr.Reviewers)
	}
}

func TestClosePullRequestUsesStateAndVerifiesClosedState(t *testing.T) {
	var requests []string

	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		path := req.URL.Path
		if req.URL.RawQuery != "" {
			path += "?" + req.URL.RawQuery
		}
		requests = append(requests, req.Method+" "+path)

		switch len(requests) {
		case 1:
			if req.Method != http.MethodGet {
				t.Fatalf("request 1 method = %s, want GET", req.Method)
			}
			return authTestResponse(http.StatusOK, `{"number":123,"title":"updated","state":"open"}`), nil
		case 2:
			if req.Method != http.MethodPatch {
				t.Fatalf("request 2 method = %s, want PATCH", req.Method)
			}
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("failed to read request body: %v", err)
			}
			values, err := url.ParseQuery(string(body))
			if err != nil {
				t.Fatalf("ParseQuery() error = %v", err)
			}
			if got := values.Get("state"); got != "closed" {
				t.Fatalf("state = %q, want %q", got, "closed")
			}
			if got := values.Get("title"); got != "updated" {
				t.Fatalf("title = %q, want %q", got, "updated")
			}
			if got := values.Get("state_event"); got != "" {
				t.Fatalf("state_event = %q, want empty", got)
			}
			return authTestResponse(http.StatusOK, `{"number":123,"title":"updated","state":"closed"}`), nil
		default:
			t.Fatalf("unexpected extra request %d: %s", len(requests), requests[len(requests)-1])
			return nil, nil
		}
	})
	client.SetToken("test-token", "test")

	pr, err := ClosePullRequest(client, "owner", "repo", 123)
	if err != nil {
		t.Fatalf("ClosePullRequest() error = %v", err)
	}
	if pr.State != "closed" {
		t.Fatalf("pr.State = %q, want %q", pr.State, "closed")
	}
	if len(requests) != 2 {
		t.Fatalf("request count = %d, want 2", len(requests))
	}
}

func TestClosePullRequestErrorsWhenVerificationStillOpen(t *testing.T) {
	var requests int

	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		requests++
		switch requests {
		case 1:
			return authTestResponse(http.StatusOK, `{"number":123,"title":"updated","state":"open"}`), nil
		case 2:
			return authTestResponse(http.StatusOK, `{"number":123,"title":"updated","state":"open"}`), nil
		case 3:
			return authTestResponse(http.StatusOK, `{"number":123,"title":"updated","state":"open"}`), nil
		default:
			t.Fatalf("unexpected request %d", requests)
			return nil, nil
		}
	})
	client.SetToken("test-token", "test")

	_, err := ClosePullRequest(client, "owner", "repo", 123)
	if err == nil {
		t.Fatal("expected ClosePullRequest() to return an error")
	}
	if !strings.Contains(err.Error(), "still open") {
		t.Fatalf("unexpected error: %v", err)
	}
	if requests != 3 {
		t.Fatalf("request count = %d, want 3", requests)
	}
}

func TestPRCommentAndReviewUnmarshal(t *testing.T) {
	commentJSON := `{
		"id": 1,
		"discussion_id": "thread-1",
		"body": "Looks good",
		"user": {"login": "reviewer1"},
		"comment_type": "discussion",
		"resolved": false,
		"diff_file": "pkg/cmd/pr/view/view.go",
		"created_at": "2026-04-07T10:20:21+08:00"
	}`
	reviewJSON := `{
		"id": 2,
		"user": {"login": "reviewer2"},
		"body": "LGTM",
		"state": "approved",
		"submitted_at": "2026-04-07T10:30:21+08:00"
	}`

	var comment PRComment
	if err := json.Unmarshal([]byte(commentJSON), &comment); err != nil {
		t.Fatalf("Failed to unmarshal PRComment: %v", err)
	}
	if comment.DiscussionID != "thread-1" || comment.User == nil || comment.User.Login != "reviewer1" {
		t.Fatalf("Unexpected comment payload: %#v", comment)
	}

	var review PRReview
	if err := json.Unmarshal([]byte(reviewJSON), &review); err != nil {
		t.Fatalf("Failed to unmarshal PRReview: %v", err)
	}
	if review.State != "approved" || review.User == nil || review.User.Login != "reviewer2" {
		t.Fatalf("Unexpected review payload: %#v", review)
	}
}

func assertPRListRequest(t *testing.T, gotPath, wantPath string, wantQuery map[string]string) {
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
