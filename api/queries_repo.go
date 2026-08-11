package api

import (
	"net/url"
	"strings"
)

// Repository represents a GitCode repository
type Repository struct {
	ID              interface{}  `json:"id"`
	Name            string       `json:"name"`
	FullName        string       `json:"full_name"`
	Description     string       `json:"description"`
	Private         bool         `json:"private"`
	Owner           *User        `json:"owner"`
	HTMLURL         string       `json:"web_url"`
	CloneURL        string       `json:"http_url_to_repo"`
	SSHURL          string       `json:"ssh_url_to_repo"`
	DefaultBranch   string       `json:"default_branch"`
	CreatedAt       FlexibleTime `json:"created_at"`
	UpdatedAt       FlexibleTime `json:"updated_at"`
	StargazersCount int          `json:"stargazers_count"`
	ForksCount      int          `json:"forks_count"`
	OpenIssuesCount int          `json:"open_issues_count"`
	Language        string       `json:"language"`
}

// RepoListOptions represents options for listing repositories
type RepoListOptions struct {
	Visibility  string `url:"visibility,omitempty"`
	Affiliation string `url:"affiliation,omitempty"`
	Type        string `url:"type,omitempty"`
	Sort        string `url:"sort,omitempty"`
	Direction   string `url:"direction,omitempty"`
	PerPage     int    `url:"per_page,omitempty"`
	Page        int    `url:"page,omitempty"`
}

// CreateRepoOptions represents options for creating a repository
type CreateRepoOptions struct {
	Name              string `json:"name"`
	Description       string `json:"description,omitempty"`
	Private           bool   `json:"private"`
	AutoInit          bool   `json:"auto_init,omitempty"`
	GitignoreTemplate string `json:"gitignore_template,omitempty"`
	LicenseTemplate   string `json:"license_template,omitempty"`
}

// ListUserRepos lists repositories for the authenticated user.
func ListUserRepos(client *Client, opts *RepoListOptions) ([]Repository, error) {
	path := buildRepoListPath("/user/repos", opts)

	var repos []Repository
	err := client.Get(path, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// ListOrgRepos lists repositories for an organization.
func ListOrgRepos(client *Client, org string, opts *RepoListOptions) ([]Repository, error) {
	path := buildRepoListPath("/orgs/"+url.PathEscape(org)+"/repos", opts)

	var repos []Repository
	err := client.Get(path, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// GetRepo fetches a repository by owner/name
func GetRepo(client *Client, owner, name string) (*Repository, error) {
	var repo Repository
	err := client.Get(escapedRepoPath(owner, name), &repo)
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

// CreateRepo creates a new repository for the authenticated user
func CreateRepo(client *Client, opts *CreateRepoOptions) (*Repository, error) {
	var repo Repository
	err := client.Post("/user/repos", opts, &repo)
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

// CreateOrgRepo creates a new repository in an organization
func CreateOrgRepo(client *Client, org string, opts *CreateRepoOptions) (*Repository, error) {
	var repo Repository
	err := client.Post("/orgs/"+org+"/repos", opts, &repo)
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

// DeleteRepo deletes a repository
func DeleteRepo(client *Client, owner, name string) error {
	return client.Delete(escapedRepoPath(owner, name))
}

// UpdateRepoRequest is the request body for updating a repository.
type UpdateRepoRequest struct {
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	Homepage      string `json:"homepage,omitempty"`
	Private       *bool  `json:"private,omitempty"`
	DefaultBranch string `json:"default_branch,omitempty"`
}

// UpdateRepo updates a repository's settings.
//
// It calls PATCH /repos/{owner}/{repo}.
func UpdateRepo(client *Client, owner, repo string, req *UpdateRepoRequest) (*Repository, error) {
	var result Repository
	err := client.Patch("/repos/"+url.PathEscape(owner)+"/"+url.PathEscape(repo), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteBranch deletes a branch from a repository.
func DeleteBranch(client *Client, owner, name, branch string) error {
	return client.Delete("/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(name) + "/branches/" + url.PathEscape(branch))
}

// BranchListOptions represents the query parameters for listing branches.
type BranchListOptions struct {
	Sort      string
	Direction string
	PerPage   int
	Page      int
}

// ListBranches lists all branches in a repository.
//
// It calls GET /repos/{owner}/{repo}/branches.
func ListBranches(client *Client, owner, repo string, opts *BranchListOptions) ([]Branch, error) {
	path := "/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/branches"
	if opts != nil {
		path += newQueryBuilder().
			Set("sort", opts.Sort).
			Set("direction", opts.Direction).
			SetInt("per_page", opts.PerPage).
			SetInt("page", opts.Page).
			String()
	}
	var branches []Branch
	if err := client.Get(path, &branches); err != nil {
		return nil, err
	}
	return branches, nil
}

// CreateBranchRequest is the request body for creating a branch.
type CreateBranchRequest struct {
	Refs        string `json:"refs"`
	BranchName  string `json:"branch_name"`
	Description string `json:"description,omitempty"`
}

// CreateBranch creates a new branch in a repository.
//
// It calls POST /repos/{owner}/{repo}/branches.
func CreateBranch(client *Client, owner, repo string, req *CreateBranchRequest) (*Branch, error) {
	var branch Branch
	err := client.Post("/repos/"+url.PathEscape(owner)+"/"+url.PathEscape(repo)+"/branches", req, &branch)
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

// Branch represents a repository branch
type Branch struct {
	Name          string        `json:"name"`
	Protected     bool          `json:"protected"`
	Commit        *BranchCommit `json:"commit"`
	Description   string        `json:"description"`
	DefaultBranch bool          `json:"default_branch"`
}

// BranchCommit represents the commit a branch points to
type BranchCommit struct {
	ID        string       `json:"id"`
	ShortID   string       `json:"short_id"`
	Title     string       `json:"title"`
	Message   string       `json:"message"`
	Author    *User        `json:"author"`
	Committer *User        `json:"committer"`
	CreatedAt FlexibleTime `json:"created_at"`
}

// GetBranch fetches a branch by name
func GetBranch(client *Client, owner, repo, branch string) (*Branch, error) {
	var b Branch
	err := client.Get("/repos/"+url.PathEscape(owner)+"/"+url.PathEscape(repo)+"/branches/"+url.PathEscape(branch), &b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// ForkRepo forks a repository
func ForkRepo(client *Client, owner, name string) (*Repository, error) {
	var repo Repository
	err := client.Post(escapedRepoPath(owner, name)+"/forks", nil, &repo)
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

// CommitStatistics represents code contribution statistics
type CommitStatistics struct {
	Commits    []CommitStatItem `json:"commits"`
	Statistics []StatItem       `json:"statistics"`
	Total      int              `json:"total"`
}

// CommitStatItem represents a single commit stat item
type CommitStatItem struct {
	Author  string `json:"user_name"`
	Commits int    `json:"commit_count"`
}

// StatItem represents a statistics item
type StatItem struct {
	Author    string `json:"user_name"`
	Additions int    `json:"add_lines"`
	Deletions int    `json:"delete_lines"`
	Total     int    `json:"total"`
}

// CommitStatsOptions represents options for getting commit statistics
type CommitStatsOptions struct {
	BranchName string `url:"branch_name,omitempty"`
	Author     string `url:"author,omitempty"`
	OnlySelf   bool   `url:"only_self,omitempty"`
	Since      string `url:"since,omitempty"`
	Until      string `url:"until,omitempty"`
}

// GetCommitStatistics gets code contribution statistics for a repository
func GetCommitStatistics(client *Client, owner, repo string, opts *CommitStatsOptions) (*CommitStatistics, error) {
	path := "/" + owner + "/" + repo + "/repository/commit_statistics"

	// Build query string
	params := []string{}
	if opts != nil {
		if opts.BranchName != "" {
			params = append(params, "branch_name="+opts.BranchName)
		}
		if opts.Author != "" {
			params = append(params, "author="+opts.Author)
		}
		if opts.OnlySelf {
			params = append(params, "only_self=true")
		}
		if opts.Since != "" {
			params = append(params, "since="+opts.Since)
		}
		if opts.Until != "" {
			params = append(params, "until="+opts.Until)
		}
	}
	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}

	var stats CommitStatistics
	err := client.Get(path, &stats)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func buildRepoListPath(base string, opts *RepoListOptions) string {
	if opts == nil {
		return base
	}
	return base + newQueryBuilder().
		Set("visibility", strings.TrimSpace(opts.Visibility)).
		Set("affiliation", strings.TrimSpace(opts.Affiliation)).
		Set("type", strings.TrimSpace(opts.Type)).
		Set("sort", strings.TrimSpace(opts.Sort)).
		Set("direction", strings.TrimSpace(opts.Direction)).
		SetInt("per_page", opts.PerPage).
		SetInt("page", opts.Page).
		String()
}
