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

// WorkflowRunStep represents a single step within a workflow job.
type WorkflowRunStep struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Task           string `json:"task"`
	Identifier     string `json:"identifier"`
	Status         string `json:"status"`
	Sequence       int    `json:"sequence"`
	JobRunID       string `json:"job_run_id"`
	LastDispatchID string `json:"last_dispatch_id"`
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time"`
}

// WorkflowRunJob represents a job within a workflow run stage.
type WorkflowRunJob struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Identifier      string            `json:"identifier"`
	Status          string            `json:"status"`
	Sequence        int               `json:"sequence"`
	JobType         string            `json:"job_type"`
	Resource        string            `json:"resource"`
	Condition       string            `json:"condition"`
	IsSelect        bool              `json:"is_select"`
	DependsOn       []string          `json:"depends_on"`
	StartTime       int64             `json:"start_time"`
	EndTime         int64             `json:"end_time"`
	ExecuteCostTime int64             `json:"execute_cost_time"`
	ExecID          string            `json:"exec_id"`
	LastDispatchID  string            `json:"last_dispatch_id"`
	Steps           []WorkflowRunStep `json:"steps"`
}

// WorkflowRunStage represents a stage within a workflow run.
type WorkflowRunStage struct {
	ID         string           `json:"id"`
	Category   string           `json:"category"`
	Name       string           `json:"name"`
	Identifier string           `json:"identifier"`
	Status     string           `json:"status"`
	Sequence   int              `json:"sequence"`
	RunAlways  bool             `json:"run_always"`
	FailFast   bool             `json:"fail_fast"`
	IsSelect   bool             `json:"is_select"`
	DependsOn  []string         `json:"depends_on"`
	StartTime  int64            `json:"start_time"`
	EndTime    int64            `json:"end_time"`
	PauseTime  int64            `json:"pause_time"`
	Jobs       []WorkflowRunJob `json:"jobs"`
}

// WorkflowRunDetail represents the detail of a single pipeline run.
//
// It embeds WorkflowRun for the base run fields and adds the detail-only
// fields (exist_in_default_branch and stages). Only the fields needed for the
// human-facing view are modeled; deep or nullable execution fields (pre/post
// hooks, parallel, condition, inputs, env, ...) are left unmodeled and
// preserved verbatim through the raw response returned alongside this struct
// for faithful --json output.
type WorkflowRunDetail struct {
	WorkflowRun
	ExistInDefaultBranch bool               `json:"exist_in_default_branch"`
	Stages               []WorkflowRunStage `json:"stages"`
}

// GetActionsRun fetches the detail of a single pipeline run.
//
// It calls GET /api/v8/repos/{owner}/{repo}/actions/runs/{run_id}. It returns
// both a typed detail (for the human-facing view) and the raw response body
// (for faithful --json output that preserves every API field verbatim,
// including the deep stage/job/step execution fields that the typed struct
// does not model).
func GetActionsRun(client *Client, owner, repo, runID string) (*WorkflowRunDetail, json.RawMessage, error) {
	endpoint := "/api/v8/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/actions/runs/" + url.PathEscape(runID)

	resp, err := client.RawREST("GET", endpoint, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var detail WorkflowRunDetail
	if err := json.Unmarshal(resp.Body, &detail); err != nil {
		return nil, nil, fmt.Errorf("failed to parse workflow run detail response: %w", err)
	}
	return &detail, resp.Body, nil
}

// WorkflowRunJobsResponse represents the response from listing the jobs of a
// workflow run.
type WorkflowRunJobsResponse struct {
	TotalCount int              `json:"total_count"`
	Jobs       []WorkflowRunJob `json:"jobs"`
}

// ListActionsRunJobs lists the jobs of a single pipeline run.
//
// It calls GET /api/v8/repos/{owner}/{repo}/actions/runs/{run_id}/jobs. The
// job items share the WorkflowRunJob type used by the run detail response
// (the job shape is identical in both endpoints).
func ListActionsRunJobs(client *Client, owner, repo, runID string) (*WorkflowRunJobsResponse, error) {
	endpoint := "/api/v8/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/actions/runs/" + url.PathEscape(runID) + "/jobs"

	resp, err := client.RawREST("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var result WorkflowRunJobsResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse workflow run jobs response: %w", err)
	}
	return &result, nil
}

// GetActionsJob fetches the detail of a single workflow job.
//
// It calls GET /api/v8/repos/{owner}/{repo}/actions/runs/{run_id}/jobs/{job_id}.
// It returns both the typed job (for the human-facing view) and the raw
// response body (for faithful --json output that preserves the deep step
// execution fields the typed struct does not model), mirroring GetActionsRun.
func GetActionsJob(client *Client, owner, repo, runID, jobID string) (*WorkflowRunJob, json.RawMessage, error) {
	endpoint := "/api/v8/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/actions/runs/" + url.PathEscape(runID) + "/jobs/" + url.PathEscape(jobID)

	resp, err := client.RawREST("GET", endpoint, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var job WorkflowRunJob
	if err := json.Unmarshal(resp.Body, &job); err != nil {
		return nil, nil, fmt.Errorf("failed to parse workflow job response: %w", err)
	}
	return &job, resp.Body, nil
}
