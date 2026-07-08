package edit

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestEditRunWithNotesFile(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	origCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(origCwd) })
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	streams, _, out, _ := testutil.NewTestIOStreams()
	if err := os.WriteFile("notes.md", []byte("updated notes"), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v5/repos/owner/repo/releases/tags/v1.0.0":
			// Return existing release with name and body
			_, _ = w.Write([]byte(`{"tag_name":"v1.0.0","name":"v1.0.0","body":"original notes"}`))
		case r.Method == http.MethodPatch && r.URL.Path == "/api/v5/repos/owner/repo/releases/v1.0.0":
			// Verify request body contains required fields
			var opts api.GitCodeUpdateReleaseOptions
			body, _ := io.ReadAll(r.Body)
			if err := json.Unmarshal(body, &opts); err != nil {
				t.Fatalf("Failed to unmarshal request body: %v", err)
			}
			if opts.Name != "v1.0.0" {
				t.Errorf("Expected name 'v1.0.0', got '%s'", opts.Name)
			}
			if opts.Body != "updated notes" {
				t.Errorf("Expected body 'updated notes', got '%s'", opts.Body)
			}
			_, _ = w.Write([]byte(`{"tag_name":"v1.0.0","name":"v1.0.0","body":"updated notes","html_url":"https://gitcode.com/owner/repo/releases/v1.0.0"}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))

	err = editRun(&EditOptions{
		IO:         streams,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		TagName:    "v1.0.0",
		NotesFile:  "notes.md",
	})
	if err != nil {
		t.Fatalf("editRun() error = %v", err)
	}

	if !strings.Contains(out.String(), "Updated release v1.0.0") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestEditRunWithTitleOnly(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	streams, _, out, _ := testutil.NewTestIOStreams()

	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v5/repos/owner/repo/releases/tags/v1.0.0":
			_, _ = w.Write([]byte(`{"tag_name":"v1.0.0","name":"Old Title","body":"original notes"}`))
		case r.Method == http.MethodPatch && r.URL.Path == "/api/v5/repos/owner/repo/releases/v1.0.0":
			// Verify request body contains existing body and new title
			var opts api.GitCodeUpdateReleaseOptions
			body, _ := io.ReadAll(r.Body)
			if err := json.Unmarshal(body, &opts); err != nil {
				t.Fatalf("Failed to unmarshal request body: %v", err)
			}
			if opts.Name != "New Title" {
				t.Errorf("Expected name 'New Title', got '%s'", opts.Name)
			}
			if opts.Body != "original notes" {
				t.Errorf("Expected body to be preserved 'original notes', got '%s'", opts.Body)
			}
			_, _ = w.Write([]byte(`{"tag_name":"v1.0.0","name":"New Title","body":"original notes"}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))

	err := editRun(&EditOptions{
		IO:         streams,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		TagName:    "v1.0.0",
		Title:      "New Title",
	})
	if err != nil {
		t.Fatalf("editRun() error = %v", err)
	}

	if !strings.Contains(out.String(), "Updated release New Title") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestEditRunWithPrereleaseTrue(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	streams, _, out, _ := testutil.NewTestIOStreams()

	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v5/repos/owner/repo/releases/tags/v1.0.0":
			_, _ = w.Write([]byte(`{"tag_name":"v1.0.0","name":"v1.0.0","body":"notes"}`))
		case r.Method == http.MethodPatch && r.URL.Path == "/api/v5/repos/owner/repo/releases/v1.0.0":
			// Verify release_status is set to "pre"
			var opts api.GitCodeUpdateReleaseOptions
			body, _ := io.ReadAll(r.Body)
			if err := json.Unmarshal(body, &opts); err != nil {
				t.Fatalf("Failed to unmarshal request body: %v", err)
			}
			if opts.ReleaseStatus != "pre" {
				t.Errorf("Expected release_status 'pre', got '%s'", opts.ReleaseStatus)
			}
			_, _ = w.Write([]byte(`{"tag_name":"v1.0.0","name":"v1.0.0","body":"notes","prerelease":true}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))

	err := editRun(&EditOptions{
		IO:         streams,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		TagName:    "v1.0.0",
		Prerelease: "true",
	})
	if err != nil {
		t.Fatalf("editRun() error = %v", err)
	}

	if !strings.Contains(out.String(), "Updated release") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestEditRunWithPrereleaseFalse(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	streams, _, out, _ := testutil.NewTestIOStreams()

	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v5/repos/owner/repo/releases/tags/v1.0.0":
			_, _ = w.Write([]byte(`{"tag_name":"v1.0.0","name":"v1.0.0","body":"notes"}`))
		case r.Method == http.MethodPatch && r.URL.Path == "/api/v5/repos/owner/repo/releases/v1.0.0":
			// Verify release_status is set to "latest"
			var opts api.GitCodeUpdateReleaseOptions
			body, _ := io.ReadAll(r.Body)
			if err := json.Unmarshal(body, &opts); err != nil {
				t.Fatalf("Failed to unmarshal request body: %v", err)
			}
			if opts.ReleaseStatus != "latest" {
				t.Errorf("Expected release_status 'latest', got '%s'", opts.ReleaseStatus)
			}
			_, _ = w.Write([]byte(`{"tag_name":"v1.0.0","name":"v1.0.0","body":"notes","prerelease":false}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))

	err := editRun(&EditOptions{
		IO:         streams,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		TagName:    "v1.0.0",
		Prerelease: "false",
	})
	if err != nil {
		t.Fatalf("editRun() error = %v", err)
	}

	if !strings.Contains(out.String(), "Updated release") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestEditRunWithDraftFlagShowsWarning(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	streams, _, out, errOut := testutil.NewTestIOStreams()

	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v5/repos/owner/repo/releases/tags/v1.0.0":
			_, _ = w.Write([]byte(`{"tag_name":"v1.0.0","name":"v1.0.0","body":"notes"}`))
		case r.Method == http.MethodPatch && r.URL.Path == "/api/v5/repos/owner/repo/releases/v1.0.0":
			_, _ = w.Write([]byte(`{"tag_name":"v1.0.0","name":"New Title","body":"notes"}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))

	err := editRun(&EditOptions{
		IO:         streams,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		TagName:    "v1.0.0",
		Title:      "New Title",
		Draft:      "true",
	})
	if err != nil {
		t.Fatalf("editRun() error = %v", err)
	}

	// Verify warning is output
	if !strings.Contains(errOut.String(), "warning: --draft is not supported") {
		t.Fatalf("expected warning about --draft, got errOut: %q", errOut.String())
	}
	// Verify command still succeeds
	if !strings.Contains(out.String(), "Updated release") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestEditRunWithTargetFlagShowsWarning(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	streams, _, out, errOut := testutil.NewTestIOStreams()

	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v5/repos/owner/repo/releases/tags/v1.0.0":
			_, _ = w.Write([]byte(`{"tag_name":"v1.0.0","name":"v1.0.0","body":"notes"}`))
		case r.Method == http.MethodPatch && r.URL.Path == "/api/v5/repos/owner/repo/releases/v1.0.0":
			_, _ = w.Write([]byte(`{"tag_name":"v1.0.0","name":"New Title","body":"notes"}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))

	err := editRun(&EditOptions{
		IO:         streams,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		TagName:    "v1.0.0",
		Title:      "New Title",
		Target:     "main",
	})
	if err != nil {
		t.Fatalf("editRun() error = %v", err)
	}

	// Verify warning is output
	if !strings.Contains(errOut.String(), "warning: --target is not supported") {
		t.Fatalf("expected warning about --target, got errOut: %q", errOut.String())
	}
	// Verify command still succeeds
	if !strings.Contains(out.String(), "Updated release") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestEditRunWithSlashTag(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	streams, _, out, _ := testutil.NewTestIOStreams()

	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Note: http.Request.URL.Path decodes %2F to /
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v5/repos/owner/repo/releases/tags/release/v1.0.0":
			_, _ = w.Write([]byte(`{"tag_name":"release/v1.0.0","name":"release/v1.0.0","body":"notes"}`))
		case r.Method == http.MethodPatch && r.URL.Path == "/api/v5/repos/owner/repo/releases/release/v1.0.0":
			// Verify the raw URL contains escaped %2F
			rawURL := r.URL.String()
			if !strings.Contains(rawURL, "release%2Fv1.0.0") {
				expectedEscaped := "release%2Fv1.0.0"
				t.Errorf("Expected escaped tag '%s' in URL, got %s", expectedEscaped, rawURL)
			}
			_, _ = w.Write([]byte(`{"tag_name":"release/v1.0.0","name":"New Title","body":"notes"}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))

	err := editRun(&EditOptions{
		IO:         streams,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		TagName:    "release/v1.0.0",
		Title:      "New Title",
	})
	if err != nil {
		t.Fatalf("editRun() error = %v", err)
	}

	if !strings.Contains(out.String(), "Updated release") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestEditRunWithNoChanges(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	streams, _, _, _ := testutil.NewTestIOStreams()

	err := editRun(&EditOptions{
		IO:         streams,
		HttpClient: func() (*http.Client, error) { return nil, nil },
		Repository: "owner/repo",
		TagName:    "v1.0.0",
	})
	if err == nil {
		t.Fatal("Expected error for no changes")
	}
	if !strings.Contains(err.Error(), "no changes specified") {
		t.Fatalf("Expected 'no changes specified' error, got: %v", err)
	}
}
