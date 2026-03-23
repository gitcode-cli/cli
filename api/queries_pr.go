package api

import (
	"fmt"
	"time"
)

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
	MergedAt    *string     `json:"merged_at"`
	Mergeable   *bool       `json:"mergeable"`
	MergeState  interface{} `json:"mergeable_state"`
	Draft       bool        `json:"draft"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	ClosedAt    *string     `json:"closed_at"`
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
	ID           interface{} `json:"id"`
	DiscussionID string      `json:"discussion_id"`
	Body         string      `json:"body"`
	User         *User       `json:"user"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	CommentType  string      `json:"comment_type"`
	Resolved     bool        `json:"resolved"`
	DiffFile     string      `json:"diff_file"`
	DiffPosition interface{} `json:"diff_position"`
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
	Title             string   `json:"title,omitempty"`
	Body              string   `json:"body,omitempty"`
	State             string   `json:"state,omitempty"`
	StateEvent        string   `json:"state_event,omitempty"`
	Base              string   `json:"base,omitempty"`
	Draft             *bool    `json:"draft,omitempty"`
	MilestoneNumber   int      `json:"milestone_number,omitempty"`
	Labels            []string `json:"labels,omitempty"`
	CloseRelatedIssue *bool    `json:"close_related_issue,omitempty"`
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
	// GitCode API requires at least one other field along with state_event
	// Get current PR to preserve its title
	pr, err := GetPullRequest(client, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR: %w", err)
	}
	return UpdatePullRequest(client, owner, repo, number, &UpdatePROptions{
		StateEvent: "close",
		Title:      pr.Title,
	})
}

// ReopenPullRequest reopens a closed PR
func ReopenPullRequest(client *Client, owner, repo string, number int) (*PullRequest, error) {
	// GitCode API requires at least one other field along with state_event
	// Get current PR to preserve its title
	pr, err := GetPullRequest(client, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR: %w", err)
	}
	return UpdatePullRequest(client, owner, repo, number, &UpdatePROptions{
		StateEvent: "reopen",
		Title:      pr.Title,
	})
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
	token := client.Token()
	path := "/repos/" + owner + "/" + repo + "/pulls/" + itoa(number) + "/comments"
	if token != "" {
		path += "?access_token=" + token
	}

	var comments []PRComment
	err := client.Get(path, &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// CreatePRComment creates a comment on a PR
func CreatePRComment(client *Client, owner, repo string, number int, opts *CreatePRCommentOptions) (*PRComment, error) {
	token := client.Token()
	path := "/repos/" + owner + "/" + repo + "/pulls/" + itoa(number) + "/comments"
	if token != "" {
		path += "?access_token=" + token
	}

	var comment PRComment
	err := client.Post(path, opts, &comment)
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

// ReviewPROptions represents options for reviewing a PR
type ReviewPROptions struct {
	Force bool `json:"force,omitempty"` // Force approval (admin only)
}

// ReviewPR handles PR review (approve/force pass)
func ReviewPR(client *Client, owner, repo string, number int, opts *ReviewPROptions) error {
	token := client.Token()
	path := "/repos/" + owner + "/" + repo + "/pulls/" + itoa(number) + "/review"
	if token != "" {
		path += "?access_token=" + token
	}

	return client.Post(path, opts, nil)
}

// EditPR updates a PR's information
func EditPR(client *Client, owner, repo string, number int, opts *UpdatePROptions) (*PullRequest, error) {
	token := client.Token()
	path := "/repos/" + owner + "/" + repo + "/pulls/" + itoa(number)
	if token != "" {
		path += "?access_token=" + token
	}

	var pr PullRequest
	err := client.Patch(path, opts, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// TestPROptions represents options for PR test
type TestPROptions struct {
	Force bool `json:"force,omitempty"` // Force test pass (admin only)
}

// TestPR handles PR test
func TestPR(client *Client, owner, repo string, number int, opts *TestPROptions) error {
	token := client.Token()
	path := "/repos/" + owner + "/" + repo + "/pulls/" + itoa(number) + "/test"
	if token != "" {
		path += "?access_token=" + token
	}

	return client.Post(path, opts, nil)
}

// EditPRCommentOptions represents options for editing a PR comment
type EditPRCommentOptions struct {
	Body string `json:"body"`
}

// EditPRComment edits a PR comment
func EditPRComment(client *Client, owner, repo string, commentID int, opts *EditPRCommentOptions) (*PRComment, error) {
	token := client.Token()
	path := "/repos/" + owner + "/" + repo + "/pulls/comments/" + itoa(commentID)
	if token != "" {
		path += "?access_token=" + token
	}

	var comment PRComment
	err := client.Patch(path, opts, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// ReplyPRCommentOptions represents options for replying to a PR comment
type ReplyPRCommentOptions struct {
	Body string `json:"body"`
}

// ReplyPRCommentReply represents the response from replying to a PR comment
type ReplyPRCommentReply struct {
	ID      string `json:"id"`
	NoteID  int    `json:"noteId"`
	Body    string `json:"body"`
}

// ReplyPRComment replies to a PR comment discussion
func ReplyPRComment(client *Client, owner, repo string, number int, discussionID string, opts *ReplyPRCommentOptions) (*ReplyPRCommentReply, error) {
	token := client.Token()
	path := "/repos/" + owner + "/" + repo + "/pulls/" + itoa(number) + "/discussions/" + discussionID + "/comments"
	if token != "" {
		path += "?access_token=" + token
	}

	var result ReplyPRCommentReply
	err := client.Post(path, opts, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ResolvePRCommentOptions represents options for resolving a PR comment
type ResolvePRCommentOptions struct {
	Resolved bool `json:"resolved"`
}

// ResolvePRComment updates the resolution status of a PR comment
func ResolvePRComment(client *Client, owner, repo string, number int, discussionID string, opts *ResolvePRCommentOptions) error {
	token := client.Token()
	path := "/repos/" + owner + "/" + repo + "/pulls/" + itoa(number) + "/comments/" + discussionID
	if token != "" {
		path += "?access_token=" + token
	}

	return client.Put(path, opts, nil)
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

// PRFilesResponse represents the response from PR files API
type PRFilesResponse struct {
	Code        int          `json:"code"`
	AddedLines  int          `json:"added_lines"`
	RemoveLines int          `json:"remove_lines"`
	Count       int          `json:"count"`
	DiffRefs    *PRDiffRefs  `json:"diff_refs"`
	Diffs       []*PRDiff    `json:"diffs"`
}

// PRDiffRefs represents diff references
type PRDiffRefs struct {
	BaseSHA string `json:"base_sha"`
	StartSHA string `json:"start_sha"`
	HeadSHA  string `json:"head_sha"`
}

// PRDiff represents a single file diff
type PRDiff struct {
	NewBlobID string        `json:"new_blob_id"`
	Statistic *PRStatistic  `json:"statistic"`
	Type      string        `json:"type"`
	Path      string        `json:"path"`
	OldPath   string        `json:"old_path"`
	NewPath   string        `json:"new_path"`
	View      int           `json:"view"`
	Head      *PRDiffHead   `json:"head"`
	Content   *PRDiffContent `json:"content"`
}

// PRStatistic represents file change statistics
type PRStatistic struct {
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
}

// PRDiffHead represents diff head info
type PRDiffHead struct {
	URL       string `json:"url"`
	CommitID  string `json:"commit_id"`
	Additions int    `json:"added_lines"`
	Deletions int    `json:"remove_lines"`
}

// PRDiffContent represents diff content
type PRDiffContent struct {
	Text []*PRDiffLine `json:"text"`
}

// PRDiffLine represents a single diff line
type PRDiffLine struct {
	LineContent string      `json:"line_content"`
	OldLine     interface{} `json:"old_line"` // can be string "..." or object
	NewLine     interface{} `json:"new_line"` // can be string "..." or object
	Type        string      `json:"type"`
}

// GetPRFiles gets the files and diffs of a PR
func GetPRFiles(client *Client, owner, repo string, number int) (*PRFilesResponse, error) {
	token := client.Token()
	path := "/repos/" + owner + "/" + repo + "/pulls/" + itoa(number) + "/files.json"
	if token != "" {
		path += "?access_token=" + token
	}

	var result PRFilesResponse
	err := client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}