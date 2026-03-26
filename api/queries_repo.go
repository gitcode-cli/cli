package api

// Repository represents a GitCode repository
type Repository struct {
	ID          interface{} `json:"id"`
	Name        string      `json:"name"`
	FullName    string      `json:"full_name"`
	Description string      `json:"description"`
	Private     bool        `json:"private"`
	Owner       *User       `json:"owner"`
	HTMLURL     string      `json:"web_url"`
	CloneURL    string      `json:"http_url_to_repo"`
	SSHURL      string      `json:"ssh_url_to_repo"`
	DefaultBranch string    `json:"default_branch"`
	CreatedAt   FlexibleTime `json:"created_at"`
	UpdatedAt   FlexibleTime `json:"updated_at"`
	StargazersCount int     `json:"stargazers_count"`
	ForksCount  int         `json:"forks_count"`
	OpenIssuesCount int     `json:"open_issues_count"`
	Language    string      `json:"language"`
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
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Private     bool   `json:"private"`
	AutoInit    bool   `json:"auto_init,omitempty"`
	GitignoreTemplate string `json:"gitignore_template,omitempty"`
	LicenseTemplate string `json:"license_template,omitempty"`
}

// ListUserRepos lists repositories for the authenticated user
func ListUserRepos(client *Client, opts *RepoListOptions) ([]Repository, error) {
	path := "/user/repos"
	if opts != nil && opts.PerPage > 0 {
		path = path + "?per_page=" + itoa(opts.PerPage)
	}

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
	err := client.Get("/repos/"+owner+"/"+name, &repo)
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
	return client.Delete("/repos/" + owner + "/" + name)
}

// ForkRepo forks a repository
func ForkRepo(client *Client, owner, name string) (*Repository, error) {
	var repo Repository
	err := client.Post("/repos/"+owner+"/"+name+"/forks", nil, &repo)
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
	Author  string `json:"author"`
	Commits int    `json:"commits"`
}

// StatItem represents a statistics item
type StatItem struct {
	Author    string `json:"author"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
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
		path += "?" + stringsJoin(params, "&")
	}

	var stats CommitStatistics
	err := client.Get(path, &stats)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func stringsJoin(s []string, sep string) string {
	result := ""
	for i, v := range s {
		if i > 0 {
			result += sep
		}
		result += v
	}
	return result
}

func itoa(i int) string {
	if i <= 0 {
		return "30"
	}
	// simple conversion
	s := ""
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	return s
}