package api

import (
	"fmt"
)

// RepositoryCommit represents a detailed commit from the commits API
type RepositoryCommit struct {
	SHA         string        `json:"sha"`
	URL         string        `json:"url"`
	HTMLURL     string        `json:"html_url"`
	CommentsURL string        `json:"comments_url"`
	Commit      *CommitInfo   `json:"commit"`
	Author      *User         `json:"author"`
	Committer   *User         `json:"committer"`
	Stats       *CommitStats  `json:"stats"`
	Files       []CommitFile  `json:"files"`
}

// CommitInfo represents the commit details
type CommitInfo struct {
	Author    *CommitAuthor `json:"author"`
	Committer *CommitAuthor `json:"committer"`
	Message   string        `json:"message"`
}

// CommitAuthor represents commit author info
type CommitAuthor struct {
	Name  string      `json:"name"`
	Email string      `json:"email"`
	Date  FlexibleTime `json:"date"`
}

// CommitStats represents commit statistics
type CommitStats struct {
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
	Total     int `json:"total"`
}

// CommitFile represents a file in a commit
type CommitFile struct {
	Filename  string `json:"filename"`
	Status    string `json:"status"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	Changes   int    `json:"changes"`
}

// GetCommit fetches a single commit by SHA
func GetCommit(client *Client, owner, repo, sha string, showDiff bool) (*RepositoryCommit, error) {
	path := "/repos/" + owner + "/" + repo + "/commits/" + sha
	if showDiff {
		path += "?show_diff=true"
	}

	var commit RepositoryCommit
	err := client.Get(path, &commit)
	if err != nil {
		return nil, err
	}
	return &commit, nil
}

// GetCommitDiff fetches the diff of a commit
func GetCommitDiff(client *Client, owner, repo, sha string) (string, error) {
	path := "/repos/" + owner + "/" + repo + "/commit/" + sha + "/diff"

	var result map[string]interface{}
	err := client.Get(path, &result)
	if err != nil {
		return "", err
	}

	// Return formatted diff
	return fmt.Sprintf("%v", result), nil
}

// GetCommitPatch fetches the patch of a commit
func GetCommitPatch(client *Client, owner, repo, sha string) (string, error) {
	path := "/repos/" + owner + "/" + repo + "/commit/" + sha + "/patch"

	var result map[string]interface{}
	err := client.Get(path, &result)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", result), nil
}

// CommitComment represents a commit comment
type CommitComment struct {
	ID        int     `json:"id"`
	Body      string  `json:"body"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	User      *User   `json:"user"`
	Target    *Commit `json:"target"`
}

// CreateCommitComment creates a comment on a commit
func CreateCommitComment(client *Client, owner, repo, sha, body string) (*CommitComment, error) {
	path := "/repos/" + owner + "/" + repo + "/commits/" + sha + "/comments"
	payload := map[string]string{"body": body}

	var comment CommitComment
	err := client.Post(path, payload, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetCommitComment fetches a single commit comment
func GetCommitComment(client *Client, owner, repo string, id int) (*CommitComment, error) {
	path := fmt.Sprintf("/repos/%s/%s/comments/%d", owner, repo, id)

	var comment CommitComment
	err := client.Get(path, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// UpdateCommitComment updates a commit comment
func UpdateCommitComment(client *Client, owner, repo string, id int, body string) (*CommitComment, error) {
	path := fmt.Sprintf("/repos/%s/%s/comments/%d", owner, repo, id)
	payload := map[string]string{"body": body}

	var comment CommitComment
	err := client.Patch(path, payload, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// ListCommitComments lists all commit comments in a repository
func ListCommitComments(client *Client, owner, repo string, opts *ListOptions) ([]CommitComment, error) {
	path := "/repos/" + owner + "/" + repo + "/comments"
	if opts != nil && opts.PerPage > 0 {
		path += fmt.Sprintf("?per_page=%d&page=%d", opts.PerPage, opts.Page)
		if opts.Order != "" {
			path += "&order=" + opts.Order
		}
	}

	var comments []CommitComment
	err := client.Get(path, &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// ListCommentsForCommit lists comments for a specific commit
func ListCommentsForCommit(client *Client, owner, repo, sha string, opts *ListOptions) ([]CommitComment, error) {
	path := "/repos/" + owner + "/" + repo + "/commits/" + sha + "/comments"
	if opts != nil && opts.PerPage > 0 {
		path += fmt.Sprintf("?per_page=%d&page=%d", opts.PerPage, opts.Page)
	}

	var comments []CommitComment
	err := client.Get(path, &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// ListOptions represents common list options
type ListOptions struct {
	Page    int
	PerPage int
	Order   string
}

// CommitCommentsBySHAOptions represents options for listing comments by SHA
type CommitCommentsBySHAOptions struct {
	Page    int
	PerPage int
}