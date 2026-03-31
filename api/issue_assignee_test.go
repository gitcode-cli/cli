package api

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestCreateIssueUsesAssigneeIDs(t *testing.T) {
	var gotBody url.Values

	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}
		gotBody, err = url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("ParseQuery() error = %v", err)
		}
		return authTestResponse(http.StatusOK, `{"number":"1","title":"created"}`), nil
	})
	client.SetToken("test-token", "test")

	_, err := CreateIssue(client, "owner", "repo", &CreateIssueOptions{
		Title:       "created",
		AssigneeIDs: []string{"4744798"},
	})
	if err != nil {
		t.Fatalf("CreateIssue() error = %v", err)
	}

	if gotBody.Get("assignee_id") != "4744798" {
		t.Fatalf("assignee_id = %q, want %q", gotBody.Get("assignee_id"), "4744798")
	}
	if got := gotBody["assignee_ids[]"]; len(got) != 1 || got[0] != "4744798" {
		t.Fatalf("assignee_ids[] = %v, want [%q]", got, "4744798")
	}
}

func TestUpdateIssueUsesAssigneeIDs(t *testing.T) {
	var gotBody url.Values

	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}
		gotBody, err = url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("ParseQuery() error = %v", err)
		}
		return authTestResponse(http.StatusOK, `{"number":"1","title":"updated"}`), nil
	})
	client.SetToken("test-token", "test")

	_, err := UpdateIssue(client, "owner", "repo", 1, &UpdateIssueOptions{
		AssigneeIDs: []string{"4744798"},
	})
	if err != nil {
		t.Fatalf("UpdateIssue() error = %v", err)
	}

	if gotBody.Get("assignee_id") != "4744798" {
		t.Fatalf("assignee_id = %q, want %q", gotBody.Get("assignee_id"), "4744798")
	}
	if got := gotBody["assignee_ids[]"]; len(got) != 1 || got[0] != "4744798" {
		t.Fatalf("assignee_ids[] = %v, want [%q]", got, "4744798")
	}
}

func TestResolveUserIDs(t *testing.T) {
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/api/v5/users/alice":
			return authTestResponse(http.StatusOK, `{"id":"101","login":"alice"}`), nil
		case "/api/v5/users/bob":
			return authTestResponse(http.StatusOK, `{"id":202,"login":"bob"}`), nil
		default:
			return authTestResponse(http.StatusNotFound, `{"message":"not found"}`), nil
		}
	})
	client.SetToken("test-token", "test")

	got, err := ResolveUserIDs(client, []string{"alice", "bob"})
	if err != nil {
		t.Fatalf("ResolveUserIDs() error = %v", err)
	}
	if strings.Join(got, ",") != "101,202" {
		t.Fatalf("ResolveUserIDs() = %v, want %v", got, []string{"101", "202"})
	}
}
