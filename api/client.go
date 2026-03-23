// Package api provides GitCode API client
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultHost is the default GitCode host
	DefaultHost = "api.gitcode.com"
	// DefaultAPIVersion is the default API version
	DefaultAPIVersion = "v5"
)

// Client is a GitCode API client
type Client struct {
	httpClient  *http.Client
	host        string
	token       string
	tokenSource string
}

// NewClient creates a new API client
func NewClient(httpClient *http.Client, host, token string) *Client {
	return &Client{
		httpClient: httpClient,
		host:       host,
		token:      token,
	}
}

// NewClientFromHTTP creates a client from an existing HTTP client
func NewClientFromHTTP(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
		host:       DefaultHost,
	}
}

// SetToken sets the authentication token
func (c *Client) SetToken(token, source string) {
	c.token = token
	c.tokenSource = source
}

// SetHost sets the API host
func (c *Client) SetHost(host string) {
	c.host = host
}

// Host returns the current host
func (c *Client) Host() string {
	return c.host
}

// Token returns the current token
func (c *Client) Token() string {
	return c.token
}

// REST performs a REST API call
func (c *Client) REST(method, path string, body interface{}, response interface{}) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	reqURL := fmt.Sprintf("https://%s/api/%s%s", c.host, DefaultAPIVersion, path)
	req, err := http.NewRequest(method, reqURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Add authentication
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err == nil && apiErr.Message != "" {
			return &apiErr
		}
		return fmt.Errorf("API error: %s", resp.Status)
	}

	// Parse response
	if response != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, response); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

// Get performs a GET request
func (c *Client) Get(path string, response interface{}) error {
	return c.REST("GET", path, nil, response)
}

// Post performs a POST request
func (c *Client) Post(path string, body interface{}, response interface{}) error {
	return c.REST("POST", path, body, response)
}

// Put performs a PUT request
func (c *Client) Put(path string, body interface{}, response interface{}) error {
	return c.REST("PUT", path, body, response)
}

// Patch performs a PATCH request
func (c *Client) Patch(path string, body interface{}, response interface{}) error {
	return c.REST("PATCH", path, body, response)
}

// Delete performs a DELETE request
func (c *Client) Delete(path string) error {
	return c.REST("DELETE", path, nil, nil)
}

// UploadToURL uploads a file to an external URL with custom headers
func (c *Client) UploadToURL(uploadURL, filename string, content []byte, contentType string, headers map[string]string) error {
	req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("failed to create upload request: %w", err)
	}

	// Set Content-Type
	req.Header.Set("Content-Type", contentType)

	// Set custom headers from API response
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		return fmt.Errorf("upload failed: %s - %s", resp.Status, string(respBody))
	}

	return nil
}

// UploadAsset uploads a file to a release
func (c *Client) UploadAsset(path, filename string, content []byte, contentType string) (*ReleaseAsset, error) {
	reqURL := fmt.Sprintf("https://%s/api/%s%s?name=%s", c.host, DefaultAPIVersion, path, url.QueryEscape(filename))
	req, err := http.NewRequest("POST", reqURL, bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", contentType)

	// Add authentication
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err == nil && apiErr.Message != "" {
			return nil, &apiErr
		}
		return nil, fmt.Errorf("upload failed: %s", resp.Status)
	}

	// Parse response
	var asset ReleaseAsset
	if err := json.Unmarshal(respBody, &asset); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &asset, nil
}

// APIError represents an API error
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	ErrorName  string `json:"error"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("API error (%d)", e.StatusCode)
}

// DefaultHTTPClient returns the default HTTP client
func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

// BuildURL builds a URL with path parameters
func BuildURL(base string, params map[string]string) string {
	for k, v := range params {
		base = replacePath(base, k, v)
	}
	return base
}

func replacePath(urlStr, key, value string) string {
	placeholder := "{" + key + "}"
	return replaceAll(urlStr, placeholder, url.PathEscape(value))
}

func replaceAll(s, old, new string) string {
	result := ""
	for {
		idx := findString(s, old)
		if idx == -1 {
			return result + s
		}
		result += s[:idx] + new
		s = s[idx+len(old):]
	}
}

func findString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}