package download

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

func TestDownloadAssetUsesAuthorizationHeader(t *testing.T) {
	tempDir := t.TempDir()
	ioStreams, _, _, _ := iostreams.Test()
	cs := ioStreams.ColorScheme()

	var gotPath string
	var gotAuth string
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			gotPath = req.URL.Path
			if req.URL.RawQuery != "" {
				gotPath += "?" + req.URL.RawQuery
			}
			gotAuth = req.Header.Get("Authorization")
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     http.StatusText(http.StatusOK),
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader("asset-body")),
			}, nil
		}),
	}

	client := api.NewClientFromHTTP(httpClient)
	client.SetToken("test-token", "test")

	err := downloadAsset(api.ReleaseAsset{Name: "app.tar.gz"}, tempDir, httpClient, cs, io.Discard, client, "owner", "repo", "v1.0.0")
	if err != nil {
		t.Fatalf("downloadAsset() error = %v", err)
	}

	if strings.Contains(gotPath, "access_token=") {
		t.Fatalf("request path unexpectedly contains access_token query: %q", gotPath)
	}
	if gotAuth != "Bearer test-token" {
		t.Fatalf("Authorization header = %q, want %q", gotAuth, "Bearer test-token")
	}

	outputPath := filepath.Join(tempDir, "app.tar.gz")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(content) != "asset-body" {
		t.Fatalf("downloaded content = %q, want %q", string(content), "asset-body")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
