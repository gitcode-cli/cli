package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// ActionsPlugin represents a plugin entry in the Actions plugin list response.
type ActionsPlugin struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

// ActionsPluginVersionEntry represents a single version entry within the
// vision_content array of a plugin detail response.
type ActionsPluginVersionEntry struct {
	Version string `json:"version"`
	Readme  string `json:"readme"`
}

// ActionsPluginDetail represents the full plugin detail response.
type ActionsPluginDetail struct {
	Name          string                      `json:"name"`
	DisplayName   string                      `json:"display_name"`
	Description   string                      `json:"description"`
	VisionContent []ActionsPluginVersionEntry `json:"vision_content"`
}

// ActionsListPluginsOptions controls pagination for ListActionsPlugins.
type ActionsListPluginsOptions struct {
	PerPage int
	Page    int
}

// webAPIHeaders returns headers required by the web-api.gitcode.com plugin
// directory endpoints, including a Referer for same-site request validation.
func webAPIHeaders() map[string]string {
	return map[string]string{
		"Referer": "https://gitcode.com/",
	}
}

// ListActionsPlugins lists official Actions plugins for a project.
//
// It calls GET https://web-api.gitcode.com/api/v2/projects/{project}/actions/plugins/all.
// The project path is the URL-encoded "owner/repo" form. The raw response body
// is returned so callers can preserve full API fields for --json output.
func ListActionsPlugins(client *Client, project string, opts *ActionsListPluginsOptions) ([]byte, error) {
	endpoint := "/api/v2/projects/" + url.PathEscape(project) + "/actions/plugins/all"
	if opts != nil {
		endpoint += newQueryBuilder().
			SetInt("per_page", opts.PerPage).
			SetInt("page", opts.Page).
			String()
	}

	resp, err := client.RawRESTToHost("GET", WebAPIHost, endpoint, nil, webAPIHeaders())
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// ViewActionsPlugin retrieves the detail of a specific Actions plugin.
//
// It calls
// GET https://web-api.gitcode.com/api/v2/projects/{project}/actions/plugins/detail?name={name}.
// The raw response body is returned so callers can preserve full API fields
// for --json output.
func ViewActionsPlugin(client *Client, project, name string) ([]byte, error) {
	endpoint := "/api/v2/projects/" + url.PathEscape(project) + "/actions/plugins/detail" +
		newQueryBuilder().Set("name", name).String()

	resp, err := client.RawRESTToHost("GET", WebAPIHost, endpoint, nil, webAPIHeaders())
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// ParseActionsPluginsList parses the raw list response into typed plugin
// entries. It handles both a plain JSON array and a paginated wrapper object
// with common field names ("plugins", "list", or "data").
func ParseActionsPluginsList(raw []byte) ([]ActionsPlugin, error) {
	var plugins []ActionsPlugin
	if err := json.Unmarshal(raw, &plugins); err == nil {
		return plugins, nil
	}

	var wrapper struct {
		Plugins []ActionsPlugin `json:"plugins"`
		List    []ActionsPlugin `json:"list"`
		Data    []ActionsPlugin `json:"data"`
	}
	if err := json.Unmarshal(raw, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse plugins list response: %w", err)
	}
	if len(wrapper.Plugins) > 0 {
		return wrapper.Plugins, nil
	}
	if len(wrapper.List) > 0 {
		return wrapper.List, nil
	}
	return wrapper.Data, nil
}

// ParseActionsPluginsListRaw parses the raw list response into raw JSON
// entries, preserving full API fields for --json output. It handles both a
// plain JSON array and a paginated wrapper object.
func ParseActionsPluginsListRaw(raw []byte) ([]json.RawMessage, error) {
	var entries []json.RawMessage
	if err := json.Unmarshal(raw, &entries); err == nil {
		return entries, nil
	}

	var wrapper struct {
		Plugins []json.RawMessage `json:"plugins"`
		List    []json.RawMessage `json:"list"`
		Data    []json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse plugins list response: %w", err)
	}
	if len(wrapper.Plugins) > 0 {
		return wrapper.Plugins, nil
	}
	if len(wrapper.List) > 0 {
		return wrapper.List, nil
	}
	return wrapper.Data, nil
}

// CountPluginsEntries returns the number of plugin entries in a raw list
// response, used to detect short pages during pagination.
func CountPluginsEntries(raw []byte) int {
	entries, err := ParseActionsPluginsListRaw(raw)
	if err != nil {
		return 0
	}
	return len(entries)
}
