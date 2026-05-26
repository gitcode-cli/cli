package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// TestGitCodeUpdateReleaseOptions_JSON tests JSON serialization of GitCodeUpdateReleaseOptions
func TestGitCodeUpdateReleaseOptions_JSON(t *testing.T) {
	opts := &GitCodeUpdateReleaseOptions{
		Name:          "v1.0.0",
		Body:          "Release notes",
		ReleaseStatus: "pre",
	}

	data, err := json.Marshal(opts)
	if err != nil {
		t.Fatalf("Failed to marshal GitCodeUpdateReleaseOptions: %v", err)
	}

	// Verify required fields are present
	if !strings.Contains(string(data), `"name":"v1.0.0"`) {
		t.Errorf("Expected name field in JSON: %s", string(data))
	}
	if !strings.Contains(string(data), `"body":"Release notes"`) {
		t.Errorf("Expected body field in JSON: %s", string(data))
	}
	if !strings.Contains(string(data), `"release_status":"pre"`) {
		t.Errorf("Expected release_status field in JSON: %s", string(data))
	}
}

// TestGitCodeUpdateReleaseOptions_EmptyReleaseStatus tests that empty release_status is omitted
func TestGitCodeUpdateReleaseOptions_EmptyReleaseStatus(t *testing.T) {
	opts := &GitCodeUpdateReleaseOptions{
		Name: "v1.0.0",
		Body: "Release notes",
	}

	data, err := json.Marshal(opts)
	if err != nil {
		t.Fatalf("Failed to marshal GitCodeUpdateReleaseOptions: %v", err)
	}

	// Verify release_status is omitted when empty
	if strings.Contains(string(data), "release_status") {
		t.Errorf("Expected release_status to be omitted when empty: %s", string(data))
	}
}

func TestCreateReleaseOptionsReleaseStatusJSON(t *testing.T) {
	opts := &CreateReleaseOptions{
		TagName:       "v1.0.0-rc1",
		Name:          "v1.0.0 RC1",
		Prerelease:    true,
		ReleaseStatus: "pre",
	}

	data, err := json.Marshal(opts)
	if err != nil {
		t.Fatalf("Failed to marshal CreateReleaseOptions: %v", err)
	}

	if !strings.Contains(string(data), `"release_status":"pre"`) {
		t.Fatalf("Expected release_status field in JSON: %s", string(data))
	}
}

// TestUpdateReleaseByTagDirect_PatchPath tests that UpdateReleaseByTagDirect uses correct PATCH path
func TestUpdateReleaseByTagDirect_PatchPath(t *testing.T) {
	var gotMethod string
	var gotPath string

	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		return authTestResponse(http.StatusOK, `{"tag_name":"v1.0.0","name":"Updated","body":"Notes"}`), nil
	})
	client.SetToken("test-token", "test")

	opts := &GitCodeUpdateReleaseOptions{
		Name: "Updated",
		Body: "Notes",
	}

	_, err := UpdateReleaseByTagDirect(client, "owner", "repo", "v1.0.0", opts)
	if err != nil {
		t.Fatalf("UpdateReleaseByTagDirect() error = %v", err)
	}

	if gotMethod != "PATCH" {
		t.Errorf("Expected PATCH method, got %s", gotMethod)
	}
	if gotPath != "/api/v5/repos/owner/repo/releases/v1.0.0" {
		t.Errorf("Expected path '/api/v5/repos/owner/repo/releases/v1.0.0', got %s", gotPath)
	}
}

// TestUpdateReleaseByTagDirect_TagEscaping tests that tags with slashes are properly escaped
func TestUpdateReleaseByTagDirect_TagEscaping(t *testing.T) {
	var gotURL string

	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		// Use req.URL.String() to get the encoded URL, not req.URL.Path which is decoded
		gotURL = req.URL.String()
		return authTestResponse(http.StatusOK, `{"tag_name":"release/v1.0.0","name":"Updated","body":"Notes"}`), nil
	})
	client.SetToken("test-token", "test")

	opts := &GitCodeUpdateReleaseOptions{
		Name: "Updated",
		Body: "Notes",
	}

	_, err := UpdateReleaseByTagDirect(client, "owner", "repo", "release/v1.0.0", opts)
	if err != nil {
		t.Fatalf("UpdateReleaseByTagDirect() error = %v", err)
	}

	// Verify the slash in the tag is escaped to %2F in the full URL
	expectedEscaped := "release%2Fv1.0.0"
	if !strings.Contains(gotURL, expectedEscaped) {
		t.Errorf("Expected escaped tag '%s' in URL, got %s", expectedEscaped, gotURL)
	}
}

// TestUpdateReleaseByTagDirect_RequestBody tests that request body contains correct fields
func TestUpdateReleaseByTagDirect_RequestBody(t *testing.T) {
	var gotBody string

	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		gotBody = string(body)
		return authTestResponse(http.StatusOK, `{"tag_name":"v1.0.0","name":"Updated Title","body":"Updated Notes"}`), nil
	})
	client.SetToken("test-token", "test")

	opts := &GitCodeUpdateReleaseOptions{
		Name:          "Updated Title",
		Body:          "Updated Notes",
		ReleaseStatus: "pre",
	}

	_, err := UpdateReleaseByTagDirect(client, "owner", "repo", "v1.0.0", opts)
	if err != nil {
		t.Fatalf("UpdateReleaseByTagDirect() error = %v", err)
	}

	// Verify request body contains required fields
	if !strings.Contains(gotBody, `"name":"Updated Title"`) {
		t.Errorf("Expected name in request body: %s", gotBody)
	}
	if !strings.Contains(gotBody, `"body":"Updated Notes"`) {
		t.Errorf("Expected body in request body: %s", gotBody)
	}
	if !strings.Contains(gotBody, `"release_status":"pre"`) {
		t.Errorf("Expected release_status in request body: %s", gotBody)
	}
}

// TestUpdateReleaseByTagDirect_ReturnsRelease tests that response is parsed correctly
func TestUpdateReleaseByTagDirect_ReturnsRelease(t *testing.T) {
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		return authTestResponse(http.StatusOK, `{
			"tag_name": "v1.0.0",
			"name": "Updated Title",
			"body": "Updated Notes",
			"prerelease": true,
			"html_url": "https://gitcode.com/owner/repo/releases/v1.0.0"
		}`), nil
	})
	client.SetToken("test-token", "test")

	opts := &GitCodeUpdateReleaseOptions{
		Name:          "Updated Title",
		Body:          "Updated Notes",
		ReleaseStatus: "pre",
	}

	release, err := UpdateReleaseByTagDirect(client, "owner", "repo", "v1.0.0", opts)
	if err != nil {
		t.Fatalf("UpdateReleaseByTagDirect() error = %v", err)
	}

	if release.TagName != "v1.0.0" {
		t.Errorf("Expected TagName 'v1.0.0', got '%s'", release.TagName)
	}
	if release.Name != "Updated Title" {
		t.Errorf("Expected Name 'Updated Title', got '%s'", release.Name)
	}
	if release.Body != "Updated Notes" {
		t.Errorf("Expected Body 'Updated Notes', got '%s'", release.Body)
	}
	if !release.Prerelease {
		t.Errorf("Expected Prerelease true, got false")
	}
}

// TestUpdateReleaseByTagDirect_Error tests error handling
func TestUpdateReleaseByTagDirect_Error(t *testing.T) {
	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		return authTestResponse(http.StatusNotFound, `{"error_message":"release not found"}`), nil
	})
	client.SetToken("test-token", "test")

	opts := &GitCodeUpdateReleaseOptions{
		Name: "Updated",
		Body: "Notes",
	}

	_, err := UpdateReleaseByTagDirect(client, "owner", "repo", "nonexistent", opts)
	if err == nil {
		t.Fatal("Expected error for nonexistent release")
	}
}

// TestGetRelease_TagEscaping tests that GetRelease properly escapes tags with slashes
func TestGetRelease_TagEscaping(t *testing.T) {
	var gotURL string

	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotURL = req.URL.String()
		return authTestResponse(http.StatusOK, `{"tag_name":"release/v1.0.0","name":"Test Release","body":"Notes"}`), nil
	})
	client.SetToken("test-token", "test")

	_, err := GetRelease(client, "owner", "repo", "release/v1.0.0")
	if err != nil {
		t.Fatalf("GetRelease() error = %v", err)
	}

	// Verify the slash in the tag is escaped to %2F in the full URL
	expectedEscaped := "release%2Fv1.0.0"
	if !strings.Contains(gotURL, expectedEscaped) {
		t.Errorf("Expected escaped tag '%s' in URL, got %s", expectedEscaped, gotURL)
	}
}

// TestGetReleaseUploadURL_TagEscaping tests that GetReleaseUploadURL properly escapes tags with slashes
func TestGetReleaseUploadURL_TagEscaping(t *testing.T) {
	var gotURL string

	client := newAuthTestClient(func(req *http.Request) (*http.Response, error) {
		gotURL = req.URL.String()
		return authTestResponse(http.StatusOK, `{"url":"https://uploads.example.test/file","headers":{}}`), nil
	})
	client.SetToken("test-token", "test")

	_, err := GetReleaseUploadURL(client, "owner", "repo", "release/v1.0.0", "test.txt")
	if err != nil {
		t.Fatalf("GetReleaseUploadURL() error = %v", err)
	}

	// Verify the slash in the tag is escaped to %2F in the full URL
	expectedEscaped := "release%2Fv1.0.0"
	if !strings.Contains(gotURL, expectedEscaped) {
		t.Errorf("Expected escaped tag '%s' in URL, got %s", expectedEscaped, gotURL)
	}
}
