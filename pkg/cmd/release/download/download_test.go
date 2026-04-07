package download

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

func TestNewCmdDownloadAllowsLatestReleaseWithoutTag(t *testing.T) {
	f := cmdutil.TestFactory()
	cmd := NewCmdDownload(f, func(opts *DownloadOptions) error {
		return nil
	})
	cmd.SetArgs([]string{"-R", "owner/repo"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestDownloadRunUsesLatestReleaseWhenTagOmitted(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	tempDir := t.TempDir()
	ioStreams, _, _, _ := iostreams.Test()
	var gotPaths []string
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			path := req.URL.Path
			if req.URL.RawQuery != "" {
				path += "?" + req.URL.RawQuery
			}
			gotPaths = append(gotPaths, path)

			switch {
			case req.Method == http.MethodGet && path == "/api/v5/repos/owner/repo/releases/latest":
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
					Header:     make(http.Header),
					Body: io.NopCloser(strings.NewReader(`{
						"tag_name":"v1.0.0",
						"assets":[{"name":"app.tar.gz"}]
					}`)),
				}, nil
			case req.Method == http.MethodGet && path == "/api/v5/repos/owner/repo/releases/v1.0.0/attach_files/app.tar.gz/download":
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader("asset-body")),
				}, nil
			default:
				t.Fatalf("unexpected request: %s %s", req.Method, path)
				return nil, nil
			}
		}),
	}

	err := downloadRun(&DownloadOptions{
		IO:         ioStreams,
		HttpClient: func() (*http.Client, error) { return httpClient, nil },
		Repository: "owner/repo",
		Output:     tempDir,
	})
	if err != nil {
		t.Fatalf("downloadRun() error = %v", err)
	}
	if len(gotPaths) < 2 || gotPaths[0] != "/api/v5/repos/owner/repo/releases/latest" {
		t.Fatalf("requests = %#v, want latest release lookup first", gotPaths)
	}
}

func TestDownloadRunLatestWithAllUsesBrowserDownloadURLForSourceArchives(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	tempDir := t.TempDir()
	ioStreams, _, _, _ := iostreams.Test()
	var gotURLs []string
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			gotURLs = append(gotURLs, req.URL.String())

			switch req.URL.String() {
			case "https://api.gitcode.com/api/v5/repos/owner/repo/releases/latest":
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
					Header:     make(http.Header),
					Body: io.NopCloser(strings.NewReader(`{
						"tag_name":"v1.0.0",
						"assets":[
							{
								"name":"v1.0.0.zip",
								"browser_download_url":"https://raw.gitcode.com/owner/repo/archive/refs/heads/v1.0.0.zip"
							},
							{
								"name":"asset.txt",
								"browser_download_url":"https://api.gitcode.com/owner/repo/releases/download/v1.0.0/asset.txt"
							}
						]
					}`)),
				}, nil
			case "https://raw.gitcode.com/owner/repo/archive/refs/heads/v1.0.0.zip":
				if got := req.Header.Get("Authorization"); got != "Bearer test-token" {
					t.Fatalf("Authorization header for source archive = %q, want %q", got, "Bearer test-token")
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader("zip-body")),
				}, nil
			case "https://api.gitcode.com/api/v5/repos/owner/repo/releases/v1.0.0/attach_files/asset.txt/download":
				if got := req.Header.Get("Authorization"); got != "Bearer test-token" {
					t.Fatalf("Authorization header for uploaded asset = %q, want %q", got, "Bearer test-token")
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader("asset-body")),
				}, nil
			default:
				t.Fatalf("unexpected request: %s %s", req.Method, req.URL.String())
				return nil, nil
			}
		}),
	}

	err := downloadRun(&DownloadOptions{
		IO:         ioStreams,
		HttpClient: func() (*http.Client, error) { return httpClient, nil },
		Repository: "owner/repo",
		Output:     tempDir,
		All:        true,
	})
	if err != nil {
		t.Fatalf("downloadRun() error = %v", err)
	}

	if len(gotURLs) != 3 {
		t.Fatalf("requests = %#v, want latest lookup plus 2 asset downloads", gotURLs)
	}

	zipContent, err := os.ReadFile(filepath.Join(tempDir, "v1.0.0.zip"))
	if err != nil {
		t.Fatalf("ReadFile(v1.0.0.zip) error = %v", err)
	}
	if string(zipContent) != "zip-body" {
		t.Fatalf("zip content = %q, want %q", string(zipContent), "zip-body")
	}

	assetContent, err := os.ReadFile(filepath.Join(tempDir, "asset.txt"))
	if err != nil {
		t.Fatalf("ReadFile(asset.txt) error = %v", err)
	}
	if string(assetContent) != "asset-body" {
		t.Fatalf("asset content = %q, want %q", string(assetContent), "asset-body")
	}
}

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

func TestDownloadAssetFallsBackToAttachFilesEndpointWithoutBrowserDownloadURL(t *testing.T) {
	tempDir := t.TempDir()
	ioStreams, _, _, _ := iostreams.Test()
	cs := ioStreams.ColorScheme()

	var gotURL string
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			gotURL = req.URL.String()
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

	wantURL := "https://api.gitcode.com/api/v5/repos/owner/repo/releases/v1.0.0/attach_files/app.tar.gz/download"
	if gotURL != wantURL {
		t.Fatalf("download URL = %q, want %q", gotURL, wantURL)
	}
}

func TestDownloadAssetIgnoresNonArchiveBrowserDownloadURL(t *testing.T) {
	tempDir := t.TempDir()
	ioStreams, _, _, _ := iostreams.Test()
	cs := ioStreams.ColorScheme()

	var gotURL string
	httpClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			gotURL = req.URL.String()
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

	err := downloadAsset(api.ReleaseAsset{
		Name:               "asset.txt",
		BrowserDownloadURL: "https://api.gitcode.com/owner/repo/releases/download/v1.0.0/asset.txt",
	}, tempDir, httpClient, cs, io.Discard, client, "owner", "repo", "v1.0.0")
	if err != nil {
		t.Fatalf("downloadAsset() error = %v", err)
	}

	wantURL := "https://api.gitcode.com/api/v5/repos/owner/repo/releases/v1.0.0/attach_files/asset.txt/download"
	if gotURL != wantURL {
		t.Fatalf("download URL = %q, want %q", gotURL, wantURL)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
