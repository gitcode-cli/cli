package api

import "time"

// Release represents a GitCode release
type Release struct {
	ID          interface{} `json:"id"`
	TagName     string      `json:"tag_name"`
	Name        string      `json:"name"`
	Body        string      `json:"body"`
	Draft       bool        `json:"draft"`
	Prerelease  bool        `json:"prerelease"`
	HTMLURL     string      `json:"html_url"`
	AssetsURL   string      `json:"assets_url"`
	UploadURL   string      `json:"upload_url"`
	CreatedAt   time.Time   `json:"created_at"`
	PublishedAt *time.Time  `json:"published_at"`
	Author      *User       `json:"author"`
	Assets      []ReleaseAsset `json:"assets"`
}

// ReleaseAsset represents an asset in a release
type ReleaseAsset struct {
	ID        interface{} `json:"id"`
	Name      string      `json:"name"`
	Label     string      `json:"label"`
	State     string      `json:"state"`
	Size      int         `json:"size"`
	Downloads int         `json:"download_count"`
	URL       string      `json:"url"`
	BrowserURL string     `json:"browser_download_url"`
	CreatedAt time.Time   `json:"created_at"`
}

// ReleaseListOptions represents options for listing releases
type ReleaseListOptions struct {
	PerPage int `url:"per_page,omitempty"`
	Page    int `url:"page,omitempty"`
}

// CreateReleaseOptions represents options for creating a release
type CreateReleaseOptions struct {
	TagName         string `json:"tag_name"`
	Name            string `json:"name,omitempty"`
	Body            string `json:"body,omitempty"`
	Draft           bool   `json:"draft"`
	Prerelease      bool   `json:"prerelease"`
	TargetCommitish string `json:"target_commitish,omitempty"`
}

// UpdateReleaseOptions represents options for updating a release
type UpdateReleaseOptions struct {
	TagName    string `json:"tag_name,omitempty"`
	Name       string `json:"name,omitempty"`
	Body       string `json:"body,omitempty"`
	Draft      *bool  `json:"draft,omitempty"`
	Prerelease *bool  `json:"prerelease,omitempty"`
}

// ListReleases lists releases for a repository
func ListReleases(client *Client, owner, repo string, opts *ReleaseListOptions) ([]Release, error) {
	path := "/repos/" + owner + "/" + repo + "/releases"
	if opts != nil && opts.PerPage > 0 {
		path = path + "?per_page=" + itoa(opts.PerPage)
	}

	var releases []Release
	err := client.Get(path, &releases)
	if err != nil {
		return nil, err
	}
	return releases, nil
}

// GetRelease fetches a release by tag name
func GetRelease(client *Client, owner, repo, tag string) (*Release, error) {
	var release Release
	err := client.Get("/repos/"+owner+"/"+repo+"/releases/tags/"+tag, &release)
	if err != nil {
		return nil, err
	}
	return &release, nil
}

// GetReleaseByID fetches a release by ID
func GetReleaseByID(client *Client, owner, repo string, id int64) (*Release, error) {
	var release Release
	err := client.Get("/repos/"+owner+"/"+repo+"/releases/"+itoa64(id), &release)
	if err != nil {
		return nil, err
	}
	return &release, nil
}

// CreateRelease creates a new release
func CreateRelease(client *Client, owner, repo string, opts *CreateReleaseOptions) (*Release, error) {
	var release Release
	err := client.Post("/repos/"+owner+"/"+repo+"/releases", opts, &release)
	if err != nil {
		return nil, err
	}
	return &release, nil
}

// UpdateRelease updates an existing release
func UpdateRelease(client *Client, owner, repo string, id int64, opts *UpdateReleaseOptions) (*Release, error) {
	var release Release
	err := client.Patch("/repos/"+owner+"/"+repo+"/releases/"+itoa64(id), opts, &release)
	if err != nil {
		return nil, err
	}
	return &release, nil
}

// DeleteRelease deletes a release
func DeleteRelease(client *Client, owner, repo string, id int64) error {
	return client.Delete("/repos/" + owner + "/" + repo + "/releases/" + itoa64(id))
}

// DeleteReleaseByTag deletes a release by tag name
func DeleteReleaseByTag(client *Client, owner, repo, tag string) error {
	// First get the release to find its ID
	release, err := GetRelease(client, owner, repo, tag)
	if err != nil {
		return err
	}

	// Extract ID
	var id int64
	switch v := release.ID.(type) {
	case float64:
		id = int64(v)
	case int64:
		id = v
	case int:
		id = int64(v)
	default:
		return ErrInvalidReleaseID
	}

	return DeleteRelease(client, owner, repo, id)
}

// ErrInvalidReleaseID is returned when release ID is invalid
var ErrInvalidReleaseID = fmtError("invalid release ID")

func fmtError(msg string) error {
	return &releaseError{msg: msg}
}

type releaseError struct {
	msg string
}

func (e *releaseError) Error() string {
	return e.msg
}

func itoa64(i int64) string {
	if i <= 0 {
		return "0"
	}
	s := ""
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	return s
}