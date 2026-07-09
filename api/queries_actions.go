package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// ActionsActor represents the actor who triggered a workflow run.
type ActionsActor struct {
	ID       string `json:"id"`
	ObjectID string `json:"object_id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
}

// WorkflowRun represents a single pipeline run record.
type WorkflowRun struct {
	WorkflowRunID string        `json:"workflow_run_id"`
	WorkflowID    string        `json:"workflow_id"`
	WorkflowName  string        `json:"workflow_name"`
	FilePath      string        `json:"file_path"`
	Title         string        `json:"title"`
	Status        string        `json:"status"`
	Event         string        `json:"event"`
	RunNumber     int           `json:"run_number"`
	HeadBranch    string        `json:"head_branch"`
	HeadSHA       string        `json:"head_sha"`
	Actor         *ActionsActor `json:"actor"`
	StartTime     int64         `json:"start_time"`
	EndTime       int64         `json:"end_time"`
	PauseTime     int64         `json:"pause_time"`
}

// WorkflowRunsResponse represents the response from listing workflow runs.
type WorkflowRunsResponse struct {
	TotalCount   int           `json:"total_count"`
	WorkflowRuns []WorkflowRun `json:"workflow_runs"`
}

// ActionsListRunsOptions represents the filter options for listing workflow runs.
//
// The fields mirror the query parameters accepted by the GitCode Actions v8 API.
// The optional access_token query parameter is intentionally omitted: the CLI
// authenticates through the standard Bearer header instead of exposing the
// token in the request URL.
type ActionsListRunsOptions struct {
	Event         string
	Status        string
	Branch        string
	Executor      string
	PullRequestID string
	WorkflowID    string
	WorkflowName  string
	PerPage       int
	Page          int
	StartTime     int64
	EndTime       int64
}

// ListActionsRuns lists all pipeline run records for a repository.
//
// It calls the GitCode Actions v8 endpoint
// GET /api/v8/repos/{owner}/{repo}/actions/runs. Unlike most queries that use
// the default v5 prefix, the Actions API lives under v8, so the request is
// issued via RawREST with a full /api/v8 path which the client preserves
// verbatim (see Client.rawURL).
func ListActionsRuns(client *Client, owner, repo string, opts *ActionsListRunsOptions) (*WorkflowRunsResponse, error) {
	endpoint := "/api/v8/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/actions/runs"
	if opts != nil {
		values := url.Values{}
		if opts.Event != "" {
			values.Set("event", opts.Event)
		}
		if opts.Status != "" {
			values.Set("status", opts.Status)
		}
		if opts.Branch != "" {
			values.Set("branch", opts.Branch)
		}
		if opts.Executor != "" {
			values.Set("executor", opts.Executor)
		}
		if opts.PullRequestID != "" {
			values.Set("pull_request_id", opts.PullRequestID)
		}
		if opts.WorkflowID != "" {
			values.Set("workflow_id", opts.WorkflowID)
		}
		if opts.WorkflowName != "" {
			values.Set("workflow_name", opts.WorkflowName)
		}
		if opts.PerPage > 0 {
			values.Set("per_page", itoa(opts.PerPage))
		}
		if opts.Page > 0 {
			values.Set("page", itoa(opts.Page))
		}
		if opts.StartTime > 0 {
			values.Set("startTime", itoa64(opts.StartTime))
		}
		if opts.EndTime > 0 {
			values.Set("endTime", itoa64(opts.EndTime))
		}
		if len(values) > 0 {
			endpoint += "?" + values.Encode()
		}
	}

	resp, err := client.RawREST("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var result WorkflowRunsResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse workflow runs response: %w", err)
	}
	return &result, nil
}
