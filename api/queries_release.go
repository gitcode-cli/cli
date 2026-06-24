package api

import (
	"net/url"
	"sort"
	"strings"
)

// Release represents a GitCode release
type Release struct {
	ID              interface{}    `json:"id"`
	TagName         string         `json:"tag_name"`
	TargetCommitish string         `json:"target_commitish"`
	Name            string         `json:"name"`
	Body            string         `json:"body"`
	Draft           bool           `json:"draft"`
	Prerelease      bool           `json:"prerelease"`
	HTMLURL         string         `json:"html_url"`
	AssetsURL       string         `json:"assets_url"`
	UploadURL       string         `json:"upload_url"`
	CreatedAt       FlexibleTime   `json:"created_at"`
	PublishedAt     *FlexibleTime  `json:"published_at"`
	Author          *User          `json:"author"`
	Assets          []ReleaseAsset `json:"assets"`
}

// ReleaseAsset represents an asset in a release
type ReleaseAsset struct {
	ID                 int          `json:"id"`
	Name               string       `json:"name"`
	Label              string       `json:"label"`
	State              string       `json:"state"`
	ContentType        string       `json:"content_type"`
	Size               int          `json:"size"`
	Downloads          int          `json:"download_count"`
	URL                string       `json:"url"`
	BrowserDownloadURL string       `json:"browser_download_url"`
	CreatedAt          FlexibleTime `json:"created_at"`
	UpdatedAt          FlexibleTime `json:"updated_at"`
}

// AssetUploadURL represents the response from getting upload URL
type AssetUploadURL struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
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
	Draft           bool   `json:"draft,omitempty"`
	Prerelease      bool   `json:"prerelease,omitempty"`
	ReleaseStatus   string `json:"release_status,omitempty"`
	TargetCommitish string `json:"target_commitish,omitempty"`
}

// UpdateReleaseOptions represents options for updating a release
type UpdateReleaseOptions struct {
	TagName         string `json:"tag_name,omitempty"`
	TargetCommitish string `json:"target_commitish,omitempty"`
	Name            string `json:"name,omitempty"`
	Body            string `json:"body,omitempty"`
	Draft           *bool  `json:"draft,omitempty"`
	Prerelease      *bool  `json:"prerelease,omitempty"`
}

// GitCodeUpdateReleaseOptions represents options for GitCode's tag-based update API
// Required fields: name, body
// Optional: release_status (pre/latest)
type GitCodeUpdateReleaseOptions struct {
	Name          string `json:"name"`                     // Required
	Body          string `json:"body"`                     // Required
	ReleaseStatus string `json:"release_status,omitempty"` // Optional: "pre" or "latest"
}

// ListReleases lists releases for a repository
func ListReleases(client *Client, owner, repo string, opts *ReleaseListOptions) ([]Release, error) {
	path := buildPath("/repos/"+owner+"/"+repo+"/releases", opts)

	var releases []Release
	err := client.Get(path, &releases)
	if err != nil {
		return nil, err
	}
	// Sort by created_at descending so --limit returns the latest releases.
	// The GitCode API does not guarantee sort order.
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].CreatedAt.Time.After(releases[j].CreatedAt.Time)
	})
	return releases, nil
}

// GetRelease fetches a release by tag name
func GetRelease(client *Client, owner, repo, tag string) (*Release, error) {
	escapedTag := url.PathEscape(tag)
	var release Release
	err := client.Get("/repos/"+owner+"/"+repo+"/releases/tags/"+escapedTag, &release)
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

// GetLatestRelease fetches the latest release for a repository
func GetLatestRelease(client *Client, owner, repo string) (*Release, error) {
	var release Release
	err := client.Get("/repos/"+owner+"/"+repo+"/releases/latest", &release)
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

// UpdateRelease updates an existing release by ID
func UpdateRelease(client *Client, owner, repo string, id int64, opts *UpdateReleaseOptions) (*Release, error) {
	var release Release
	err := client.Patch("/repos/"+owner+"/"+repo+"/releases/"+itoa64(id), opts, &release)
	if err != nil {
		return nil, err
	}
	return &release, nil
}

// UpdateReleaseByTag updates a release by tag name
func UpdateReleaseByTag(client *Client, owner, repo, tag string, opts *UpdateReleaseOptions) (*Release, error) {
	// First get the release to find its ID
	release, err := GetRelease(client, owner, repo, tag)
	if err != nil {
		return nil, err
	}

	id, err := release.GetID()
	if err != nil {
		return nil, err
	}

	return UpdateRelease(client, owner, repo, id, opts)
}

// UpdateReleaseByTagDirect updates a release using GitCode's tag-based PATCH endpoint
// This bypasses the need for release ID which GitCode API doesn't return
func UpdateReleaseByTagDirect(client *Client, owner, repo, tag string, opts *GitCodeUpdateReleaseOptions) (*Release, error) {
	// URL escape the tag for path safety (e.g., "release/v1.0.0" -> "release%2Fv1.0.0")
	escapedTag := url.PathEscape(tag)
	path := "/repos/" + owner + "/" + repo + "/releases/" + escapedTag

	var release Release
	err := client.Patch(path, opts, &release)
	if err != nil {
		return nil, err
	}
	return &release, nil
}

// DeleteRelease deletes a release by ID
func DeleteRelease(client *Client, owner, repo string, id int64) error {
	return client.Delete("/repos/" + owner + "/" + repo + "/releases/" + itoa64(id))
}

// DeleteReleaseByTag deletes a release by tag name.
// It first tries the tag-based DELETE endpoint; if the API does not support it (404),
// it falls back to the ID-based flow (get release → extract ID → delete by ID).
func DeleteReleaseByTag(client *Client, owner, repo, tag string) error {
	err := DeleteReleaseByTagDirect(client, owner, repo, tag)
	if err == nil {
		return nil
	}
	// If the tag endpoint isn't available, fall back to ID-based deletion
	if !isNotFoundError(err) {
		return err
	}

	release, err := GetRelease(client, owner, repo, tag)
	if err != nil {
		return err
	}

	id, err := release.GetID()
	if err != nil {
		return err
	}

	return DeleteRelease(client, owner, repo, id)
}

// DeleteReleaseByTagDirect deletes a release by tag name directly.
// This bypasses the need for release ID which GitCode API may not return.
func DeleteReleaseByTagDirect(client *Client, owner, repo, tag string) error {
	escapedTag := url.PathEscape(tag)
	return client.Delete("/repos/" + owner + "/" + repo + "/releases/" + escapedTag)
}

// DeleteReleaseByTagKnown is like DeleteReleaseByTag but uses a previously
// fetched release to avoid an extra GetRelease call in the fallback path.
// Callers that already have the release (e.g. for pre-confirm) should use
// this to skip the redundant API round-trip.
func DeleteReleaseByTagKnown(client *Client, owner, repo, tag string, release *Release) error {
	err := DeleteReleaseByTagDirect(client, owner, repo, tag)
	if err == nil {
		return nil
	}
	if !isNotFoundError(err) {
		return err
	}
	id, idErr := release.GetID()
	if idErr != nil {
		return idErr
	}
	return DeleteRelease(client, owner, repo, id)
}

// isNotFoundError checks whether err represents a 404 response.
// It first checks for *APIError with status 404, then falls back to
// inspecting the error string for non-*APIError edge cases where
// client.REST returns a bare fmt.Errorf for a 404 response.
func isNotFoundError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == 404
	}
	if err != nil {
		return strings.Contains(err.Error(), "404")
	}
	return false
}

// GetID extracts the release ID as int64
func (r *Release) GetID() (int64, error) {
	if r.ID == nil {
		return 0, ErrNoReleaseID
	}

	switch v := r.ID.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case string:
		// Try to parse string as number
		var id int64
		for _, c := range v {
			if c >= '0' && c <= '9' {
				id = id*10 + int64(c-'0')
			} else {
				return 0, ErrNoReleaseID
			}
		}
		return id, nil
	default:
		return 0, ErrNoReleaseID
	}
}

// ListReleaseAssets lists assets for a release
func ListReleaseAssets(client *Client, owner, repo string, releaseID int64) ([]ReleaseAsset, error) {
	var assets []ReleaseAsset
	err := client.Get("/repos/"+owner+"/"+repo+"/releases/"+itoa64(releaseID)+"/assets", &assets)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

// GetReleaseAsset fetches a single release asset
func GetReleaseAsset(client *Client, owner, repo string, assetID int64) (*ReleaseAsset, error) {
	var asset ReleaseAsset
	err := client.Get("/repos/"+owner+"/"+repo+"/releases/assets/"+itoa64(assetID), &asset)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// DeleteReleaseAsset deletes a release asset
func DeleteReleaseAsset(client *Client, owner, repo string, assetID int64) error {
	return client.Delete("/repos/" + owner + "/" + repo + "/releases/assets/" + itoa64(assetID))
}

// GetReleaseUploadURL fetches the upload URL for a release asset
func GetReleaseUploadURL(client *Client, owner, repo, tag, filename string) (*AssetUploadURL, error) {
	escapedTag := url.PathEscape(tag)
	values := url.Values{}
	values.Set("file_name", filename)
	path := "/repos/" + owner + "/" + repo + "/releases/" + escapedTag + "/upload_url?" + values.Encode()

	var result AssetUploadURL
	err := client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UploadReleaseAsset uploads a file as a release asset
func UploadReleaseAsset(client *Client, owner, repo string, releaseID int64, filename string, content []byte, contentType string) (*ReleaseAsset, error) {
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return client.UploadAsset("/repos/"+owner+"/"+repo+"/releases/"+itoa64(releaseID)+"/assets", filename, content, contentType)
}

// UploadReleaseAssetByTag uploads a file to a release by tag name using two-step process
func UploadReleaseAssetByTag(client *Client, owner, repo, tag, filename string, content []byte, contentType string) error {
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Step 1: Get upload URL and headers
	uploadInfo, err := GetReleaseUploadURL(client, owner, repo, tag, filename)
	if err != nil {
		return fmtError("failed to get upload URL: " + err.Error())
	}

	if uploadInfo.URL == "" {
		return fmtError("upload URL is empty")
	}

	// Step 2: Upload file to the returned URL with headers
	return client.UploadToURL(uploadInfo.URL, filename, content, contentType, uploadInfo.Headers)
}

// ErrNoReleaseID is returned when the GitCode API omits release IDs.
var ErrNoReleaseID = fmtError("release id was not returned by GitCode API")

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

func buildPath(base string, opts *ReleaseListOptions) string {
	if opts == nil {
		return base
	}

	params := ""
	if opts.PerPage > 0 {
		params = "?per_page=" + itoa64(int64(opts.PerPage))
	}
	if opts.Page > 0 {
		if params != "" {
			params += "&page=" + itoa64(int64(opts.Page))
		} else {
			params = "?page=" + itoa64(int64(opts.Page))
		}
	}
	return base + params
}
