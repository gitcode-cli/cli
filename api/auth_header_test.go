package api

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestListPRCommentsUsesAuthorizationHeader(t *testing.T) {
	var gotPath string
	var gotAuth string

	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		if req.URL.RawQuery != "" {
			gotPath += "?" + req.URL.RawQuery
		}
		gotAuth = req.Header.Get("Authorization")
		return authTestResponse(http.StatusOK, `[]`), nil
	})
	client.SetToken("test-token", "test")

	_, err := ListPRComments(client, "owner", "repo", 123)
	if err != nil {
		t.Fatalf("ListPRComments() error = %v", err)
	}

	assertNoAccessTokenQuery(t, gotPath)
	if gotAuth != "Bearer test-token" {
		t.Fatalf("Authorization header = %q, want %q", gotAuth, "Bearer test-token")
	}
}

func TestUpdateIssueUsesAuthorizationHeader(t *testing.T) {
	var gotPath string
	var gotAuth string

	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		if req.URL.RawQuery != "" {
			gotPath += "?" + req.URL.RawQuery
		}
		gotAuth = req.Header.Get("Authorization")
		return authTestResponse(http.StatusOK, `{"number":"1","title":"updated"}`), nil
	})
	client.SetToken("test-token", "test")

	_, err := UpdateIssue(client, "owner", "repo", 1, &UpdateIssueOptions{
		Title: "updated",
	})
	if err != nil {
		t.Fatalf("UpdateIssue() error = %v", err)
	}

	assertNoAccessTokenQuery(t, gotPath)
	if gotAuth != "Bearer test-token" {
		t.Fatalf("Authorization header = %q, want %q", gotAuth, "Bearer test-token")
	}
}

func TestGetReleaseUploadURLUsesAuthorizationHeader(t *testing.T) {
	var gotPath string
	var gotAuth string

	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		if req.URL.RawQuery != "" {
			gotPath += "?" + req.URL.RawQuery
		}
		gotAuth = req.Header.Get("Authorization")
		return authTestResponse(http.StatusOK, `{"url":"https://uploads.example.test/file","headers":{}}`), nil
	})
	client.SetToken("test-token", "test")

	_, err := GetReleaseUploadURL(client, "owner", "repo", "v1.0.0", "app bundle+linux.tar.gz")
	if err != nil {
		t.Fatalf("GetReleaseUploadURL() error = %v", err)
	}

	assertNoAccessTokenQuery(t, gotPath)
	if !strings.Contains(gotPath, "file_name=app+bundle%2Blinux.tar.gz") {
		t.Fatalf("request path = %q, want encoded file_name query", gotPath)
	}
	if gotAuth != "Bearer test-token" {
		t.Fatalf("Authorization header = %q, want %q", gotAuth, "Bearer test-token")
	}
}

func TestReviewPRUsesAuthorizationHeaderAndReviewEndpoint(t *testing.T) {
	var gotPath string
	var gotAuth string

	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		if req.URL.RawQuery != "" {
			gotPath += "?" + req.URL.RawQuery
		}
		gotAuth = req.Header.Get("Authorization")
		return authTestResponse(http.StatusNoContent, ``), nil
	})
	client.SetToken("test-token", "test")

	err := ReviewPR(client, "owner", "repo", 123, &ReviewPROptions{})
	if err != nil {
		t.Fatalf("ReviewPR() error = %v", err)
	}

	assertNoAccessTokenQuery(t, gotPath)
	if gotPath != "/api/v5/repos/owner/repo/pulls/123/review" {
		t.Fatalf("request path = %q, want %q", gotPath, "/api/v5/repos/owner/repo/pulls/123/review")
	}
	if gotAuth != "Bearer test-token" {
		t.Fatalf("Authorization header = %q, want %q", gotAuth, "Bearer test-token")
	}
}

func TestReviewPRReturnsErrorMessageField(t *testing.T) {
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		return authTestResponse(http.StatusBadRequest, `{"error_code_name":"UN_KNOW","error_message":"403 Forbidden - You don't have the authority to approval this merge request."}`), nil
	})
	client.SetToken("test-token", "test")

	err := ReviewPR(client, "owner", "repo", 123, &ReviewPROptions{})
	if err == nil {
		t.Fatal("expected ReviewPR() to return an error")
	}
	if !strings.Contains(err.Error(), "You don't have the authority") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func newAuthTestClient(fn func(*http.Request) (*http.Response, error)) *Client {
	return NewClientFromHTTP(&http.Client{
		Transport: roundTripFunc(fn),
	})
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func authTestResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func assertNoAccessTokenQuery(t *testing.T, path string) {
	t.Helper()
	if strings.Contains(path, "access_token=") {
		t.Fatalf("request path unexpectedly contains access_token query: %q", path)
	}
}
