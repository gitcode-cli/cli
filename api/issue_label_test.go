package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestRemoveIssueLabelUsesLabelEndpointsAndVerifiesState(t *testing.T) {
	var requests []string

	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		requests = append(requests, req.Method+" "+req.URL.Path)

		switch len(requests) {
		case 1:
			if req.Method != http.MethodGet || req.URL.Path != "/api/v5/repos/owner/repo/issues/7" {
				t.Fatalf("unexpected request 1: %s %s", req.Method, req.URL.Path)
			}
			return authTestResponse(http.StatusOK, `{"number":"7","title":"issue","labels":[{"name":"bug"},{"name":"risk/high"}]}`), nil
		case 2:
			if req.Method != http.MethodPut || req.URL.Path != "/api/v5/repos/owner/repo/issues/7/labels" {
				t.Fatalf("unexpected request 2: %s %s", req.Method, req.URL.Path)
			}
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("failed to read request body: %v", err)
			}
			var payload struct {
				Labels []string `json:"labels"`
			}
			if err := json.Unmarshal(body, &payload); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if got := payload.Labels; len(got) != 1 || got[0] != "risk/high" {
				t.Fatalf("labels = %v, want [risk/high]", got)
			}
			return authTestResponse(http.StatusOK, `[{"name":"risk/high"}]`), nil
		case 3:
			if req.Method != http.MethodGet || req.URL.Path != "/api/v5/repos/owner/repo/issues/7" {
				t.Fatalf("unexpected request 3: %s %s", req.Method, req.URL.Path)
			}
			return authTestResponse(http.StatusOK, `{"number":"7","title":"issue","labels":[{"name":"risk/high"}]}`), nil
		default:
			t.Fatalf("unexpected extra request %d: %s", len(requests), requests[len(requests)-1])
			return nil, nil
		}
	})
	client.SetToken("test-token", "test")

	if err := RemoveIssueLabel(client, "owner", "repo", 7, "bug"); err != nil {
		t.Fatalf("RemoveIssueLabel() error = %v", err)
	}
	if len(requests) != 3 {
		t.Fatalf("request count = %d, want 3", len(requests))
	}
}

func TestRemoveIssueLabelClearsAllLabelsWhenRemovingLastLabel(t *testing.T) {
	var requests []string

	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		requests = append(requests, req.Method+" "+req.URL.Path)

		switch len(requests) {
		case 1:
			return authTestResponse(http.StatusOK, `{"number":"7","title":"issue","labels":[{"name":"bug"}]}`), nil
		case 2:
			if req.Method != http.MethodDelete || req.URL.Path != "/api/v5/repos/owner/repo/issues/7/labels" {
				t.Fatalf("unexpected clear-label request: %s %s", req.Method, req.URL.Path)
			}
			return authTestResponse(http.StatusNoContent, ``), nil
		case 3:
			return authTestResponse(http.StatusOK, `{"number":"7","title":"issue","labels":[]}`), nil
		default:
			t.Fatalf("unexpected extra request %d: %s", len(requests), requests[len(requests)-1])
			return nil, nil
		}
	})
	client.SetToken("test-token", "test")

	if err := RemoveIssueLabel(client, "owner", "repo", 7, "bug"); err != nil {
		t.Fatalf("RemoveIssueLabel() error = %v", err)
	}
}

func TestRemoveIssueLabelErrorsWhenLabelStillPresent(t *testing.T) {
	var requests int

	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		requests++
		switch requests {
		case 1:
			return authTestResponse(http.StatusOK, `{"number":"7","title":"issue","labels":[{"name":"bug"}]}`), nil
		case 2:
			return authTestResponse(http.StatusNoContent, ``), nil
		case 3:
			return authTestResponse(http.StatusOK, `{"number":"7","title":"issue","labels":[{"name":"bug"}]}`), nil
		default:
			t.Fatalf("unexpected request %d", requests)
			return nil, nil
		}
	})
	client.SetToken("test-token", "test")

	err := RemoveIssueLabel(client, "owner", "repo", 7, "bug")
	if err == nil {
		t.Fatal("expected RemoveIssueLabel() to return an error")
	}
	if !strings.Contains(err.Error(), "still present") {
		t.Fatalf("unexpected error: %v", err)
	}
}
