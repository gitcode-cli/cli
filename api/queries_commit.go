package api

import (
	"fmt"
	"time"
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
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
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