package api

import (
	"fmt"
	"net/url"
)

// SearchRepoResult represents a repository search result item.
type SearchRepoResult struct {
	ID           interface{} `json:"id"`
	FullName     string      `json:"full_name"`
	HumanName    string      `json:"human_name"`
	URL          string      `json:"url"`
	Namespace    interface{} `json:"namespace"`
	Path         string      `json:"path"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Status       interface{} `json:"status"`
	SSHURLToRepo string      `json:"ssh_url_to_repo"`
}

// SearchIssueResult represents an issue search result item.
type SearchIssueResult struct {
	ID         interface{} `json:"id"`
	HTMLURL    string      `json:"html_url"`
	Number     interface{} `json:"number"`
	State      string      `json:"state"`
	Title      string      `json:"title"`
	Body       string      `json:"body"`
	Repository interface{} `json:"repository"`
	CreatedAt  string      `json:"created_at"`
	UpdatedAt  string      `json:"updated_at"`
}

// SearchOptions represents common search query parameters.
type SearchOptions struct {
	Q       string
	Sort    string
	Order   string
	PerPage int
	Page    int
}

// SearchRepos searches for repositories.
//
// It calls GET /api/v5/search/repositories.
func SearchRepos(client *Client, opts *SearchOptions) ([]SearchRepoResult, error) {
	path := "/search/repositories" + buildSearchQuery(opts)
	var results []SearchRepoResult
	if err := client.Get(path, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// SearchIssues searches for issues.
//
// It calls GET /api/v5/search/issues.
func SearchIssues(client *Client, opts *SearchOptions, repo, state string) ([]SearchIssueResult, error) {
	q := buildSearchQuery(opts)
	if repo != "" {
		q += fmt.Sprintf("&repo=%s", url.QueryEscape(repo))
	}
	if state != "" {
		q += fmt.Sprintf("&state=%s", url.QueryEscape(state))
	}
	path := "/search/issues" + q
	var results []SearchIssueResult
	if err := client.Get(path, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// SearchUsers searches for users.
//
// It calls GET /api/v5/search/users.
func SearchUsers(client *Client, opts *SearchOptions) ([]User, error) {
	path := "/search/users" + buildSearchQuery(opts)
	var results []User
	if err := client.Get(path, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func buildSearchQuery(opts *SearchOptions) string {
	if opts == nil {
		return ""
	}
	q := newQueryBuilder().
		Set("q", opts.Q).
		Set("sort", opts.Sort).
		Set("order", opts.Order).
		SetInt("per_page", opts.PerPage).
		SetInt("page", opts.Page).
		String()
	return q
}
