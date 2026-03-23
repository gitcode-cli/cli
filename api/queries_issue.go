package api

import (
	"fmt"
	"net/url"
	"time"
)

// Issue represents a GitCode issue
type Issue struct {
	ID          interface{} `json:"id"`
	Number      string      `json:"number"`
	Title       string      `json:"title"`
	Body        string      `json:"body"`
	State       string      `json:"state"`
	HTMLURL     string      `json:"html_url"`
	User        *User       `json:"user"`
	Assignees   []*User     `json:"assignees"`
	Labels      []*Label    `json:"labels"`
	Milestone   *Milestone  `json:"milestone"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	ClosedAt    *time.Time  `json:"closed_at"`
	Comments    int         `json:"comments"`
}

// Label represents a GitCode label
type Label struct {
	ID          interface{} `json:"id"`
	Name        string      `json:"name"`
	Color       string      `json:"color"`
	Description string      `json:"description"`
}

// Milestone represents a GitCode milestone
type Milestone struct {
	ID          interface{} `json:"id"`
	Number      string      `json:"number"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	State       string      `json:"state"`
	DueOn       *time.Time  `json:"due_on"`
}

// IssueComment represents a comment on an issue
type IssueComment struct {
	ID        interface{} `json:"id"`
	Body      string      `json:"body"`
	User      *User       `json:"user"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// IssueListOptions represents options for listing issues
type IssueListOptions struct {
	State     string `url:"state,omitempty"`
	Labels    string `url:"labels,omitempty"`
	Sort      string `url:"sort,omitempty"`
	Direction string `url:"direction,omitempty"`
	Since     string `url:"since,omitempty"`
	PerPage   int    `url:"per_page,omitempty"`
	Page      int    `url:"page,omitempty"`
}

// CreateIssueOptions represents options for creating an issue
type CreateIssueOptions struct {
	Title     string   `json:"title"`
	Body      string   `json:"body,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	Milestone int      `json:"milestone,omitempty"`
}

// UpdateIssueOptions represents options for updating an issue
type UpdateIssueOptions struct {
	Repo      string   `json:"repo,omitempty"`
	Title     string   `json:"title,omitempty"`
	Body      string   `json:"body,omitempty"`
	State     string   `json:"state,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	Milestone int      `json:"milestone,omitempty"`
}

// CreateCommentOptions represents options for creating a comment
type CreateCommentOptions struct {
	Body string `json:"body"`
}

// ListRepoIssues lists issues for a repository
func ListRepoIssues(client *Client, owner, repo string, opts *IssueListOptions) ([]Issue, error) {
	path := "/repos/" + owner + "/" + repo + "/issues"
	if opts != nil && opts.PerPage > 0 {
		path = path + "?per_page=" + itoa(opts.PerPage)
		if opts.State != "" {
			path = path + "&state=" + opts.State
		}
	}

	var issues []Issue
	err := client.Get(path, &issues)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

// GetIssue fetches an issue by number
func GetIssue(client *Client, owner, repo string, number int) (*Issue, error) {
	var issue Issue
	err := client.Get("/repos/"+owner+"/"+repo+"/issues/"+itoa(number), &issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

// CreateIssue creates a new issue
func CreateIssue(client *Client, owner, repo string, opts *CreateIssueOptions) (*Issue, error) {
	var issue Issue
	err := client.Post("/repos/"+owner+"/"+repo+"/issues", opts, &issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

// UpdateIssue updates an existing issue
// GitCode API: PATCH /repos/:owner/issues/:number with repo param in body
func UpdateIssue(client *Client, owner, repo string, number int, opts *UpdateIssueOptions) (*Issue, error) {
	token := client.Token()
	path := "/repos/" + owner + "/issues/" + itoa(number)
	if token != "" {
		path += "?access_token=" + token
	}

	// Ensure repo is set in opts
	if opts.Repo == "" {
		opts.Repo = repo
	}

	// Use form data for GitCode API compatibility
	formValues := url.Values{}
	formValues.Set("repo", opts.Repo)
	if opts.Title != "" {
		formValues.Set("title", opts.Title)
	}
	if opts.Body != "" {
		formValues.Set("body", opts.Body)
	}
	if opts.State != "" {
		formValues.Set("state", opts.State)
	}
	for _, label := range opts.Labels {
		formValues.Add("labels[]", label)
	}

	var issue Issue
	err := client.PatchForm(path, formValues, &issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

// CloseIssue closes an issue
func CloseIssue(client *Client, owner, repo string, number int) (*Issue, error) {
	// Get current issue to preserve its title
	issue, err := GetIssue(client, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}
	return UpdateIssue(client, owner, repo, number, &UpdateIssueOptions{
		Repo:  repo,
		State: "close",
		Title: issue.Title,
	})
}

// ReopenIssue reopens a closed issue
func ReopenIssue(client *Client, owner, repo string, number int) (*Issue, error) {
	// Get current issue to preserve its title
	issue, err := GetIssue(client, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}
	return UpdateIssue(client, owner, repo, number, &UpdateIssueOptions{
		Repo:  repo,
		State: "reopen",
		Title: issue.Title,
	})
}

// ListIssueComments lists comments on an issue
func ListIssueComments(client *Client, owner, repo string, number int) ([]IssueComment, error) {
	var comments []IssueComment
	err := client.Get("/repos/"+owner+"/"+repo+"/issues/"+itoa(number)+"/comments", &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// CreateIssueComment creates a comment on an issue
func CreateIssueComment(client *Client, owner, repo string, number int, opts *CreateCommentOptions) (*IssueComment, error) {
	var comment IssueComment
	err := client.Post("/repos/"+owner+"/"+repo+"/issues/"+itoa(number)+"/comments", opts, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// DeleteIssueComment deletes a comment on an issue
func DeleteIssueComment(client *Client, owner, repo string, commentID int64) error {
	return client.Delete("/repos/" + owner + "/" + repo + "/issues/comments/" + itoa(int(commentID)))
}

// ListRepoLabels lists labels for a repository
func ListRepoLabels(client *Client, owner, repo string) ([]Label, error) {
	var labels []Label
	err := client.Get("/repos/"+owner+"/"+repo+"/labels", &labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

// AddIssueLabels adds labels to an issue by updating the issue
func AddIssueLabels(client *Client, owner, repo string, number int, labels []string) ([]*Label, error) {
	// Get current issue to preserve existing labels
	issue, err := GetIssue(client, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	// Merge existing labels with new ones
	existingLabels := make([]string, len(issue.Labels))
	for i, l := range issue.Labels {
		existingLabels[i] = l.Name
	}

	// Add new labels (avoid duplicates)
	labelMap := make(map[string]bool)
	for _, l := range existingLabels {
		labelMap[l] = true
	}
	for _, l := range labels {
		labelMap[l] = true
	}

	allLabels := make([]string, 0, len(labelMap))
	for l := range labelMap {
		allLabels = append(allLabels, l)
	}

	// Update issue with new labels
	updated, err := UpdateIssue(client, owner, repo, number, &UpdateIssueOptions{
		Repo:   repo,
		Title:  issue.Title,
		Labels: allLabels,
	})
	if err != nil {
		return nil, err
	}

	return updated.Labels, nil
}

// RemoveIssueLabel removes a label from an issue by updating the issue
func RemoveIssueLabel(client *Client, owner, repo string, number int, label string) error {
	// Get current issue
	issue, err := GetIssue(client, owner, repo, number)
	if err != nil {
		return fmt.Errorf("failed to get issue: %w", err)
	}

	// Filter out the label to remove
	labels := make([]string, 0)
	for _, l := range issue.Labels {
		if l.Name != label {
			labels = append(labels, l.Name)
		}
	}

	// Update issue with remaining labels
	_, err = UpdateIssue(client, owner, repo, number, &UpdateIssueOptions{
		Repo:   repo,
		Title:  issue.Title,
		Labels: labels,
	})
	return err
}

// ListRepoMilestones lists milestones for a repository
func ListRepoMilestones(client *Client, owner, repo string) ([]Milestone, error) {
	var milestones []Milestone
	err := client.Get("/repos/"+owner+"/"+repo+"/milestones", &milestones)
	if err != nil {
		return nil, err
	}
	return milestones, nil
}