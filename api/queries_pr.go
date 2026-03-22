package api

import "time"

// PullRequest represents a GitCode pull request
type PullRequest struct {
	ID          interface{} `json:"id"`
	Number      int         `json:"number"`
	Title       string      `json:"title"`
	Body        string      `json:"body"`
	State       string      `json:"state"`
	HTMLURL     string      `json:"html_url"`
	DiffURL     string      `json:"diff_url"`
	PatchURL    string      `json:"patch_url"`
	User        *User       `json:"user"`
	Head        *PRBranch   `json:"head"`
	Base        *PRBranch   `json:"base"`
	Merged      bool        `json:"merged"`
	MergedAt    *time.Time  `json:"merged_at"`
	Mergeable   *bool       `json:"mergeable"`
	MergeState  string      `json:"mergeable_state"`
	Draft       bool        `json:"draft"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	ClosedAt    *time.Time  `json:"closed_at"`
	Comments    int         `json:"comments"`
	Commits     int         `json:"commits"`
	Additions   int         `json:"additions"`
	Deletions   int         `json:"deletions"`
	ChangedFiles int        `json:"changed_files"`
	Labels      []*Label    `json:"labels"`
	Assignees   []*User     `json:"assignees"`
	Reviewers   []*User     `json:"requested_reviewers"`
}

// PRBranch represents a branch in a PR
type PRBranch struct {
	Label string      `json:"label"`
	Ref   string      `json:"ref"`
	SHA   string      `json:"sha"`
	Repo  *Repository `json:"repo"`
}

// PRComment represents a comment on a PR
type PRComment struct {
	ID        interface{} `json:"id"`
	Body      string      `json:"body"`
	User      *User       `json:"user"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Path      string      `json:"path"`
	Position  int         `json:"position"`
}

// PRReview represents a review on a PR
type PRReview struct {
	ID          interface{} `json:"id"`
	User        *User       `json:"user"`
	Body        string      `json:"body"`
	State       string      `json:"state"`
	SubmittedAt time.Time   `json:"submitted_at"`
}

// PRListOptions represents options for listing PRs
type PRListOptions struct {
	State     string `url:"state,omitempty"`
	Head      string `url:"head,omitempty"`
	Base      string `url:"base,omitempty"`
	Sort      string `url:"sort,omitempty"`
	Direction string `url:"direction,omitempty"`
	PerPage   int    `url:"per_page,omitempty"`
	Page      int    `url:"page,omitempty"`
}

// CreatePROptions represents options for creating a PR
type CreatePROptions struct {
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
	Head  string `json:"head"`
	Base  string `json:"base"`
	Draft bool   `json:"draft,omitempty"`
}

// UpdatePROptions represents options for updating a PR
type UpdatePROptions struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
	State string `json:"state,omitempty"`
	Base  string `json:"base,omitempty"`
	Draft *bool  `json:"draft,omitempty"`
}

// CreatePRCommentOptions represents options for creating a PR comment
type CreatePRCommentOptions struct {
	Body     string `json:"body"`
	Path     string `json:"path,omitempty"`
	Position int    `json:"position,omitempty"`
}

// CreatePRReviewOptions represents options for creating a PR review
type CreatePRReviewOptions struct {
	Body  string   `json:"body,omitempty"`
	Event string   `json:"event"`
	Comments []PRReviewComment `json:"comments,omitempty"`
}

// PRReviewComment represents a comment in a review
type PRReviewComment struct {
	Path     string `json:"path"`
	Position int    `json:"position"`
	Body     string `json:"body"`
}

// MergePROptions represents options for merging a PR
type MergePROptions struct {
	CommitTitle   string `json:"commit_title,omitempty"`
	CommitMessage string `json:"commit_message,omitempty"`
	SHA           string `json:"sha,omitempty"`
	MergeMethod   string `json:"merge_method,omitempty"`
}

// ListPullRequests lists pull requests for a repository
func ListPullRequests(client *Client, owner, repo string, opts *PRListOptions) ([]PullRequest, error) {
	path := "/repos/" + owner + "/" + repo + "/pulls"
	if opts != nil && opts.PerPage > 0 {
		path = path + "?per_page=" + itoa(opts.PerPage)
		if opts.State != "" {
			path = path + "&state=" + opts.State
		}
	}

	var prs []PullRequest
	err := client.Get(path, &prs)
	if err != nil {
		return nil, err
	}
	return prs, nil
}

// GetPullRequest fetches a PR by number
func GetPullRequest(client *Client, owner, repo string, number int) (*PullRequest, error) {
	var pr PullRequest
	err := client.Get("/repos/"+owner+"/"+repo+"/pulls/"+itoa(number), &pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// CreatePullRequest creates a new PR
func CreatePullRequest(client *Client, owner, repo string, opts *CreatePROptions) (*PullRequest, error) {
	var pr PullRequest
	err := client.Post("/repos/"+owner+"/"+repo+"/pulls", opts, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// UpdatePullRequest updates an existing PR
func UpdatePullRequest(client *Client, owner, repo string, number int, opts *UpdatePROptions) (*PullRequest, error) {
	var pr PullRequest
	err := client.Patch("/repos/"+owner+"/"+repo+"/pulls/"+itoa(number), opts, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// ClosePullRequest closes a PR
func ClosePullRequest(client *Client, owner, repo string, number int) (*PullRequest, error) {
	return UpdatePullRequest(client, owner, repo, number, &UpdatePROptions{State: "closed"})
}

// ReopenPullRequest reopens a closed PR
func ReopenPullRequest(client *Client, owner, repo string, number int) (*PullRequest, error) {
	return UpdatePullRequest(client, owner, repo, number, &UpdatePROptions{State: "open"})
}

// MergePullRequest merges a PR
func MergePullRequest(client *Client, owner, repo string, number int, opts *MergePROptions) (*PullRequest, error) {
	var pr PullRequest
	err := client.Put("/repos/"+owner+"/"+repo+"/pulls/"+itoa(number)+"/merge", opts, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// ListPRComments lists comments on a PR
func ListPRComments(client *Client, owner, repo string, number int) ([]PRComment, error) {
	var comments []PRComment
	err := client.Get("/repos/"+owner+"/"+repo+"/pulls/"+itoa(number)+"/comments", &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// CreatePRComment creates a comment on a PR
func CreatePRComment(client *Client, owner, repo string, number int, opts *CreatePRCommentOptions) (*PRComment, error) {
	var comment PRComment
	err := client.Post("/repos/"+owner+"/"+repo+"/pulls/"+itoa(number)+"/comments", opts, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// CreatePRReview creates a review on a PR
func CreatePRReview(client *Client, owner, repo string, number int, opts *CreatePRReviewOptions) (*PRReview, error) {
	var review PRReview
	err := client.Post("/repos/"+owner+"/"+repo+"/pulls/"+itoa(number)+"/reviews", opts, &review)
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// ListPRReviews lists reviews on a PR
func ListPRReviews(client *Client, owner, repo string, number int) ([]PRReview, error) {
	var reviews []PRReview
	err := client.Get("/repos/"+owner+"/"+repo+"/pulls/"+itoa(number)+"/reviews", &reviews)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

// ListPRCommits lists commits in a PR
func ListPRCommits(client *Client, owner, repo string, number int) ([]Commit, error) {
	var commits []Commit
	err := client.Get("/repos/"+owner+"/"+repo+"/pulls/"+itoa(number)+"/commits", &commits)
	if err != nil {
		return nil, err
	}
	return commits, nil
}

// Commit represents a Git commit
type Commit struct {
	SHA       string    `json:"sha"`
	Message   string    `json:"message"`
	Author    *User     `json:"author"`
	Committer *User     `json:"committer"`
}

// GetPRDiff gets the diff of a PR
func GetPRDiff(client *Client, owner, repo string, number int) (string, error) {
	reqURL := "/repos/" + owner + "/" + repo + "/pulls/" + itoa(number)
	var result struct {
		Diff string `json:"diff"`
	}
	err := client.Get(reqURL, &result)
	if err != nil {
		return "", err
	}
	return result.Diff, nil
}