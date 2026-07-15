// Package upload implements the release upload command
package upload

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type UploadOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	TagName string
	Files   []string

	// Flags
	Repository string
	Label      string
	JSON       bool
}

type uploadResult struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        int    `json:"size"`
	ContentType string `json:"content_type"`
}

// NewCmdUpload creates the upload command
func NewCmdUpload(f *cmdutil.Factory, runF func(*UploadOptions) error) *cobra.Command {
	opts := &UploadOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "upload <tag> <file>...",
		Short: "Upload assets to a release",
		Long: heredoc.Doc(`
			Upload files as release assets.

			You can upload multiple files at once.
		`),
		Example: heredoc.Doc(`
			# Upload a single file
			$ gc release upload v1.0.0 app.zip -R owner/repo

			# Upload multiple files
			$ gc release upload v1.0.0 app.zip checksum.txt -R owner/repo

			# Upload to a specific repository
			$ gc release upload v1.0.0 app.zip -R owner/repo

			# Label is currently unsupported by the GitCode API
			$ gc release upload v1.0.0 app.zip -R owner/repo --label "linux-amd64"

			# Output as JSON
			$ gc release upload v1.0.0 app.zip -R owner/repo --json
		`),
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.TagName = args[0]
			opts.Files = args[1:]

			if runF != nil {
				return runF(opts)
			}
			return uploadRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Label, "label", "l", "", "Asset label")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func uploadRun(opts *UploadOptions) error {
	cs := opts.IO.ColorScheme()

	if opts.Label != "" {
		return cmdutil.NewUsageError("--label is not supported by the current GitCode release upload API")
	}

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}
	client, err := cmdutil.AuthenticatedClient(httpClient)
	if err != nil {
		return err
	}

	// Get repository
	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}
	owner, repo, err := parseRepo(repository)
	if err != nil {
		return err
	}

	results := make([]uploadResult, 0, len(opts.Files))

	// Upload each file using two-step process
	for _, file := range opts.Files {
		result, err := uploadFile(client, owner, repo, opts.TagName, file, opts.Label, cs, opts.IO.Out, !opts.JSON)
		if err != nil {
			return err
		}
		results = append(results, result)
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, results)
	}

	return nil
}

func uploadFile(client *api.Client, owner, repo, tag, filePath, label string, cs *iostreams.ColorScheme, out io.Writer, writeText bool) (uploadResult, error) {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return uploadResult{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		return uploadResult{}, fmt.Errorf("failed to read file: %w", err)
	}

	// Get filename
	filename := filepath.Base(filePath)

	// Detect content type
	contentType := detectContentType(filename, content)

	// Upload using two-step process
	err = api.UploadReleaseAssetByTag(client, owner, repo, tag, filename, content, contentType)
	if err != nil {
		return uploadResult{}, fmt.Errorf("failed to upload %s: %w", filename, err)
	}

	result := uploadResult{
		Name:        filename,
		Path:        filePath,
		Size:        len(content),
		ContentType: contentType,
	}

	if writeText {
		fmt.Fprintf(out, "%s Uploaded %s (%s)\n", cs.Green("✓"), filename, formatSize(len(content)))
	}

	return result, nil
}

func detectContentType(filename string, content []byte) string {
	// Try to detect from file extension
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".zip":
		return "application/zip"
	case ".tar":
		return "application/x-tar"
	case ".gz", ".tgz":
		return "application/gzip"
	case ".bz2":
		return "application/x-bzip2"
	case ".xz":
		return "application/x-xz"
	case ".deb":
		return "application/vnd.debian.binary-package"
	case ".rpm":
		return "application/x-rpm"
	case ".dmg":
		return "application/x-apple-diskimage"
	case ".exe":
		return "application/vnd.microsoft.portable-executable"
	case ".msi":
		return "application/x-msi"
	case ".apk":
		return "application/vnd.android.package-archive"
	case ".pdf":
		return "application/pdf"
	case ".txt", ".md":
		return "text/plain"
	case ".json":
		return "application/json"
	case ".yaml", ".yml":
		return "application/x-yaml"
	case ".xml":
		return "application/xml"
	}

	// Try to detect from content
	ct := http.DetectContentType(content)
	if ct != "application/octet-stream" {
		return ct
	}

	// Use mime.TypeByExtension
	if ct := mime.TypeByExtension(ext); ct != "" {
		return ct
	}

	return "application/octet-stream"
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
