package api

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

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
