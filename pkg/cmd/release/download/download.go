// Package download implements the release download command
package download

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type DownloadOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	TagName string
	Assets  []string

	// Flags
	Repository string
	Output     string
	All        bool
}

// NewCmdDownload creates the download command
func NewCmdDownload(f *cmdutil.Factory, runF func(*DownloadOptions) error) *cobra.Command {
	opts := &DownloadOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "download <tag> [asset]",
		Short: "Download release assets",
		Long: heredoc.Doc(`
			Download assets from a release.

			Without an asset name, downloads all assets.
			With --all, downloads source archives as well.
		`),
		Example: heredoc.Doc(`
			# Download all assets from latest release
			$ gc release download -R owner/repo

			# Download all assets from a specific release
			$ gc release download v1.0.0 -R owner/repo

			# Download a specific asset
			$ gc release download v1.0.0 app.zip -R owner/repo

			# Download to a specific directory
			$ gc release download v1.0.0 -R owner/repo --output ./downloads/
		`),
		Args: cobra.MaximumNArgs(64),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.TagName = args[0]
			}
			if len(args) > 1 {
				opts.Assets = args[1:]
			}

			if runF != nil {
				return runF(opts)
			}
			return downloadRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Output, "output", "o", ".", "Output directory")
	cmd.Flags().BoolVarP(&opts.All, "all", "A", false, "Download source archives as well")

	return cmd
}

func downloadRun(opts *DownloadOptions) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := getEnvToken()
	if token == "" {
		return cmdutil.NewAuthError("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "environment")

	// Get repository
	owner, repo, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	// Get release
	var release *api.Release
	if opts.TagName == "" {
		release, err = api.GetLatestRelease(client, owner, repo)
		if err != nil {
			return fmt.Errorf("failed to get latest release: %w", err)
		}
	} else {
		release, err = api.GetRelease(client, owner, repo, opts.TagName)
		if err != nil {
			return fmt.Errorf("failed to get release: %w", err)
		}
	}

	// Create output directory
	if err := os.MkdirAll(opts.Output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Filter assets
	assets := release.Assets
	if len(opts.Assets) > 0 {
		assets = filterAssets(assets, opts.Assets)
	}

	// Filter out source archives unless --all is set
	if !opts.All {
		assets = filterSourceArchives(assets)
	}

	if len(assets) == 0 {
		fmt.Fprintf(opts.IO.Out, "No assets to download\n")
		return nil
	}

	// Download each asset
	for _, asset := range assets {
		err := downloadAsset(asset, opts.Output, httpClient, cs, opts.IO.Out, client, owner, repo, release.TagName)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadAsset(asset api.ReleaseAsset, outputDir string, httpClient *http.Client, cs *iostreams.ColorScheme, out io.Writer, client *api.Client, owner, repo, tag string) error {
	// Create output file
	outputPath := filepath.Join(outputDir, asset.Name)
	fmt.Fprintf(out, "%s Downloading %s...\n", cs.Blue("⬇"), asset.Name)

	downloadURL := assetDownloadURL(asset, owner, repo, tag)

	// Create request
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	if token := client.Token(); token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	// Execute request
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	// Create file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy content
	written, err := io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Fprintf(out, "%s Downloaded %s (%s)\n", cs.Green("✓"), asset.Name, formatSize(int(written)))
	return nil
}

func assetDownloadURL(asset api.ReleaseAsset, owner, repo, tag string) string {
	if isSourceArchiveAsset(asset) {
		return asset.BrowserDownloadURL
	}

	return fmt.Sprintf("https://api.gitcode.com/api/v5/repos/%s/%s/releases/%s/attach_files/%s/download",
		owner, repo, tag, url.PathEscape(asset.Name))
}

func isSourceArchiveAsset(asset api.ReleaseAsset) bool {
	if strings.TrimSpace(asset.BrowserDownloadURL) == "" {
		return false
	}

	return strings.Contains(asset.BrowserDownloadURL, "/archive/refs/heads/")
}

func filterAssets(assets []api.ReleaseAsset, names []string) []api.ReleaseAsset {
	var result []api.ReleaseAsset
	for _, asset := range assets {
		for _, name := range names {
			if asset.Name == name {
				result = append(result, asset)
				break
			}
		}
	}
	return result
}

func filterSourceArchives(assets []api.ReleaseAsset) []api.ReleaseAsset {
	var result []api.ReleaseAsset
	for _, asset := range assets {
		// Skip source archives (zip, tar.gz, etc. with tag name)
		if strings.HasSuffix(asset.Name, ".zip") ||
			strings.HasSuffix(asset.Name, ".tar.gz") ||
			strings.HasSuffix(asset.Name, ".tar.bz2") ||
			strings.HasSuffix(asset.Name, ".tar") {
			// Check if it's a source archive (name contains tag name)
			if strings.Contains(asset.Name, "v") && strings.Contains(asset.Name, ".") {
				continue
			}
		}
		result = append(result, asset)
	}
	return result
}

func formatSize(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := unit, 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	if token := os.Getenv("GITCODE_TOKEN"); token != "" {
		return token
	}
	return cmdutil.EnvToken()
}
