package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestListActionsRunsBuildsV8Query(t *testing.T) {
	var gotPath string
	var gotAuth string
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		if req.URL.RawQuery != "" {
			gotPath += "?" + req.URL.RawQuery
		}
		gotAuth = req.Header.Get("Authorization")
		return authTestResponse(http.StatusOK, `{"total_count":0,"workflow_runs":[]}`), nil
	})
	client.SetToken("test-token", "test")

	resp, err := ListActionsRuns(client, "owner", "repo", &ActionsListRunsOptions{
		Event:         "Push",
		Status:        "FAILED",
		Branch:        "main",
		Executor:      "dev",
		PullRequestID: "42",
		WorkflowID:    "wf-1",
		WorkflowName:  "ci",
		PerPage:       50,
		Page:          2,
		StartTime:     1700000000,
		EndTime:       1700001000,
	})
	if err != nil {
		t.Fatalf("ListActionsRuns() error = %v", err)
	}
	if resp.TotalCount != 0 {
		t.Fatalf("TotalCount = %d, want 0", resp.TotalCount)
	}

	assertNoAccessTokenQuery(t, gotPath)

	wantPrefix := "/api/v8/repos/owner/repo/actions/runs?"
	if !strings.HasPrefix(gotPath, wantPrefix) {
		t.Fatalf("request path = %q, want prefix %q", gotPath, wantPrefix)
	}
	if gotAuth != "Bearer test-token" {
		t.Fatalf("Authorization = %q, want Bearer test-token", gotAuth)
	}

	parsed, err := url.Parse(gotPath)
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}
	q := parsed.Query()
	for _, key := range []string{"event", "status", "branch", "executor", "pull_request_id", "workflow_id", "workflow_name", "per_page", "page", "startTime", "endTime"} {
		if _, ok := q[key]; !ok {
			t.Fatalf("query missing %s in %s", key, q.Encode())
		}
	}
	if q.Get("status") != "FAILED" {
		t.Fatalf("status param = %q, want FAILED", q.Get("status"))
	}
	if q.Get("per_page") != "50" {
		t.Fatalf("per_page param = %q, want 50", q.Get("per_page"))
	}
}

func TestListActionsRunsNoOptions(t *testing.T) {
	var gotPath string
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotPath = req.URL.Path
		if req.URL.RawQuery != "" {
			gotPath += "?" + req.URL.RawQuery
		}
		return authTestResponse(http.StatusOK, `{"total_count":0,"workflow_runs":[]}`), nil
	})

	if _, err := ListActionsRuns(client, "owner", "repo", nil); err != nil {
		t.Fatalf("ListActionsRuns() error = %v", err)
	}

	want := "/api/v8/repos/owner/repo/actions/runs"
	if gotPath != want {
		t.Fatalf("request path = %q, want %q (no query)", gotPath, want)
	}
}

func TestListActionsRunsParsesResponse(t *testing.T) {
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		return authTestResponse(http.StatusOK, `{
			"total_count": 2,
			"workflow_runs": [
				{
					"workflow_run_id": "run-1",
					"workflow_id": "wf-1",
					"workflow_name": "CI",
					"file_path": ".gitcode/workflows/ci.yml",
					"title": "run CI",
					"status": "COMPLETED",
					"event": "Push",
					"run_number": 7,
					"head_branch": "main",
					"head_sha": "abc123",
					"actor": {"id": "1", "object_id": "u1", "login": "dev", "name": "Dev"},
					"start_time": 1700000000,
					"end_time": 1700000100,
					"pause_time": 0
				}
			]
		}`), nil
	})

	resp, err := ListActionsRuns(client, "owner", "repo", &ActionsListRunsOptions{Status: "COMPLETED"})
	if err != nil {
		t.Fatalf("ListActionsRuns() error = %v", err)
	}
	if resp.TotalCount != 2 {
		t.Fatalf("TotalCount = %d, want 2", resp.TotalCount)
	}
	if len(resp.WorkflowRuns) != 1 {
		t.Fatalf("len(WorkflowRuns) = %d, want 1", len(resp.WorkflowRuns))
	}
	run := resp.WorkflowRuns[0]
	if run.WorkflowRunID != "run-1" {
		t.Fatalf("WorkflowRunID = %q, want run-1", run.WorkflowRunID)
	}
	if run.RunNumber != 7 {
		t.Fatalf("RunNumber = %d, want 7", run.RunNumber)
	}
	if run.Actor == nil || run.Actor.Login != "dev" {
		t.Fatalf("Actor.Login = %v, want dev", run.Actor)
	}
	if run.StartTime != 1700000000 {
		t.Fatalf("StartTime = %d, want 1700000000", run.StartTime)
	}
}

func TestListActionsRunsEmptyRuns(t *testing.T) {
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		return authTestResponse(http.StatusOK, `{"total_count":3,"workflow_runs":null}`), nil
	})

	resp, err := ListActionsRuns(client, "owner", "repo", nil)
	if err != nil {
		t.Fatalf("ListActionsRuns() error = %v", err)
	}
	if resp.TotalCount != 3 {
		t.Fatalf("TotalCount = %d, want 3", resp.TotalCount)
	}
	if resp.WorkflowRuns != nil {
		t.Fatalf("WorkflowRuns = %v, want nil", resp.WorkflowRuns)
	}
}

func TestListActionsRunsError(t *testing.T) {
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		return authTestResponse(http.StatusBadRequest, `{"message":"bad request"}`), nil
	})

	_, err := ListActionsRuns(client, "owner", "repo", nil)
	if err == nil {
		t.Fatal("ListActionsRuns() error = nil, want error")
	}
	var apiErr *APIError
	if err := json.Unmarshal([]byte(err.Error()), &apiErr); err == nil {
		t.Fatalf("error should be *APIError, got parseable json")
	}
	if !strings.Contains(err.Error(), "400") {
		t.Fatalf("error = %q, want to contain 400", err.Error())
	}
}
