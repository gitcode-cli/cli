package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Discussion represents a GitCode organization discussion (discuss).
//
// Field shapes follow the official GitCode OpenAPI
// (gc-api-doc: GET /api/v5/orgs/{org}/discuss and /discuss/{number}).
// Nested objects that the CLI does not need to interpret (category, namespace,
// vote, vote_options, labels) are kept as raw JSON so --json output passes
// them through verbatim.
type Discussion struct {
	ID            string          `json:"id"`
	Number        int             `json:"number"`
	Title         string          `json:"title"`
	MdContent     string          `json:"md_content,omitempty"` // detail only
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
	Author        *User           `json:"author"`
	IsLock        int             `json:"is_lock"`
	IsPin         int             `json:"is_pin"`
	IsCategoryPin int             `json:"is_category_pin"`
	IsClosed      int             `json:"is_closed"`
	IsAnswered    int             `json:"is_answered"`
	CommentTotal  int             `json:"comment_total"`
	Category      json.RawMessage `json:"category,omitempty"`
	Namespace     json.RawMessage `json:"namespace,omitempty"`
	Vote          json.RawMessage `json:"vote,omitempty"`         // detail only
	VoteOptions   json.RawMessage `json:"vote_options,omitempty"` // detail only
	Labels        json.RawMessage `json:"labels,omitempty"`
}

// ListOrgDiscussionsOptions filters and paginates ListOrgDiscussions.
type ListOrgDiscussionsOptions struct {
	Page      int
	PerPage   int
	Sort      string // created (default) or comment_size
	Direction string // asc or desc (default desc)
	Search    string // filter by title/description
}

// ListOrgDiscussions lists discussions in an organization.
//
// It calls GET /api/v5/orgs/{org}/discuss with optional pagination, sort,
// direction, and search query parameters (per the official OpenAPI).
func ListOrgDiscussions(client *Client, org string, opts *ListOrgDiscussionsOptions) ([]*Discussion, error) {
	endpoint := "/orgs/" + url.PathEscape(org) + "/discuss"
	if opts != nil {
		values := url.Values{}
		if opts.Page > 0 {
			values.Set("page", itoa(opts.Page))
		}
		if opts.PerPage > 0 {
			values.Set("per_page", itoa(opts.PerPage))
		}
		if opts.Sort != "" {
			values.Set("sort", opts.Sort)
		}
		if opts.Direction != "" {
			values.Set("direction", opts.Direction)
		}
		if opts.Search != "" {
			values.Set("search", opts.Search)
		}
		if len(values) > 0 {
			endpoint += "?" + values.Encode()
		}
	}

	var discussions []*Discussion
	if err := client.Get(endpoint, &discussions); err != nil {
		return nil, fmt.Errorf("failed to list org discussions: %w", err)
	}
	return discussions, nil
}

// GetOrgDiscussion fetches a single organization discussion by number.
//
// It calls GET /api/v5/orgs/{org}/discuss/{number}.
func GetOrgDiscussion(client *Client, org string, number int) (*Discussion, error) {
	endpoint := "/orgs/" + url.PathEscape(org) + "/discuss/" + itoa(number)

	var d Discussion
	if err := client.Get(endpoint, &d); err != nil {
		return nil, fmt.Errorf("failed to get org discussion: %w", err)
	}
	return &d, nil
}
