package testutil

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/gitcode-com/gitcode-cli/pkg/iostreams"
)

// NewTestIOStreams creates IOStreams for testing
func NewTestIOStreams() (*iostreams.IOStreams, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	return &iostreams.IOStreams{
		In:     in,
		Out:    out,
		ErrOut: errOut,
	}, in, out, errOut
}

// NewTestHTTPClient creates an HTTP client for testing with a mock server
func NewTestHTTPClient(handler http.Handler) *http.Client {
	server := httptest.NewServer(handler)
	return &http.Client{
		Transport: &testTransport{server: server},
	}
}

type testTransport struct {
	server *httptest.Server
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = t.server.URL[:4] // http
	req.URL.Host = t.server.URL[7:]   // remove http://
	return http.DefaultTransport.RoundTrip(req)
}

// SetTestToken sets test environment token
func SetTestToken() func() {
	os.Setenv("GC_TOKEN", "test-token")
	return func() {
		os.Unsetenv("GC_TOKEN")
	}
}

// MockAPIHandler creates a mock API handler for testing
func MockAPIHandler() http.Handler {
	mux := http.NewServeMux()

	// User endpoints
	mux.HandleFunc("/api/v5/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "1", "login": "testuser", "name": "Test User"}`))
	})

	// Repository endpoints
	mux.HandleFunc("/api/v5/repos/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "1", "name": "test-repo", "full_name": "owner/test-repo", "html_url": "https://gitcode.com/owner/test-repo"}`))
		} else if r.Method == "DELETE" {
			w.WriteHeader(http.StatusNoContent)
		}
	})

	// User repos list
	mux.HandleFunc("/api/v5/user/repos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": "1", "name": "test-repo", "full_name": "owner/test-repo", "html_url": "https://gitcode.com/owner/test-repo"}]`))
	})

	// Create repo
	mux.HandleFunc("/api/v5/user/repos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id": "1", "name": "new-repo", "full_name": "owner/new-repo", "html_url": "https://gitcode.com/owner/new-repo"}`))
		}
	})

	// Issue endpoints
	mux.HandleFunc("/api/v5/repos/owner/test-repo/issues", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id": "1", "number": 1, "title": "Test Issue", "state": "open", "html_url": "https://gitcode.com/owner/test-repo/issues/1"}]`))
		} else if r.Method == "POST" {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id": "1", "number": 1, "title": "New Issue", "state": "open", "html_url": "https://gitcode.com/owner/test-repo/issues/1"}`))
		}
	})

	mux.HandleFunc("/api/v5/repos/owner/test-repo/issues/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "1", "number": 1, "title": "Test Issue", "state": "open", "body": "Issue body", "html_url": "https://gitcode.com/owner/test-repo/issues/1"}`))
		} else if r.Method == "PATCH" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "1", "number": 1, "title": "Updated Issue", "state": "closed", "html_url": "https://gitcode.com/owner/test-repo/issues/1"}`))
		}
	})

	mux.HandleFunc("/api/v5/repos/owner/test-repo/issues/1/comments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id": "1", "body": "Test comment"}]`))
		} else if r.Method == "POST" {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id": "1", "body": "New comment"}`))
		}
	})

	// PR endpoints
	mux.HandleFunc("/api/v5/repos/owner/test-repo/pulls", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id": "1", "number": 1, "title": "Test PR", "state": "open", "html_url": "https://gitcode.com/owner/test-repo/pull/1"}]`))
		} else if r.Method == "POST" {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id": "1", "number": 1, "title": "New PR", "state": "open", "html_url": "https://gitcode.com/owner/test-repo/pull/1"}`))
		}
	})

	mux.HandleFunc("/api/v5/repos/owner/test-repo/pulls/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "1", "number": 1, "title": "Test PR", "state": "open", "body": "PR body", "head": {"ref": "feature"}, "base": {"ref": "main"}, "html_url": "https://gitcode.com/owner/test-repo/pull/1"}`))
		} else if r.Method == "PATCH" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "1", "number": 1, "title": "Updated PR", "state": "closed", "html_url": "https://gitcode.com/owner/test-repo/pull/1"}`))
		}
	})

	mux.HandleFunc("/api/v5/repos/owner/test-repo/pulls/1/merge", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "1", "number": 1, "merged": true}`))
	})

	mux.HandleFunc("/api/v5/repos/owner/test-repo/pulls/1/reviews", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id": "1", "state": "APPROVED", "body": "LGTM"}]`))
		} else if r.Method == "POST" {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id": "1", "state": "APPROVED", "body": "Approved"}`))
		}
	})

	mux.HandleFunc("/api/v5/repos/owner/test-repo/pulls/1/comments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": "1", "body": "Test comment"}]`))
	})

	// Label endpoints
	mux.HandleFunc("/api/v5/repos/owner/test-repo/labels", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id": "1", "name": "bug", "color": "ff0000"}]`))
		} else if r.Method == "POST" {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id": "1", "name": "bug", "color": "ff0000"}`))
		}
	})

	// Milestone endpoints
	mux.HandleFunc("/api/v5/repos/owner/test-repo/milestones", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id": "1", "number": 1, "title": "v1.0", "state": "open"}]`))
		} else if r.Method == "POST" {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id": "1", "number": 1, "title": "v1.0", "state": "open"}`))
		}
	})

	mux.HandleFunc("/api/v5/repos/owner/test-repo/milestones/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "1", "number": 1, "title": "v1.0", "state": "open", "description": "First release"}`))
		} else if r.Method == "DELETE" {
			w.WriteHeader(http.StatusNoContent)
		}
	})

	return mux
}