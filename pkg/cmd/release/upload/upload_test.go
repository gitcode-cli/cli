package upload

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestUploadRun(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	io, _, out, _ := testutil.NewTestIOStreams()
	filePath := filepath.Join(t.TempDir(), "asset.txt")
	if err := os.WriteFile(filePath, []byte("hello"), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	client := testutil.NewTestHTTPClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/v5/repos/owner/repo/releases/v1.0.0/upload_url"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"url":"https://uploads.gitcode.test/upload/asset.txt","headers":{"X-Test":"1"}}`))
		case r.Method == http.MethodPut && r.URL.Path == "/upload/asset.txt":
			w.WriteHeader(http.StatusCreated)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))

	err := uploadRun(&UploadOptions{
		IO:         io,
		HttpClient: func() (*http.Client, error) { return client, nil },
		Repository: "owner/repo",
		TagName:    "v1.0.0",
		Files:      []string{filePath},
	})
	if err != nil {
		t.Fatalf("uploadRun() error = %v", err)
	}

	if !strings.Contains(out.String(), "Uploaded asset.txt") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}
