// Package download implements the actions artifact download command.
package download

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// DownloadOptions configures the actions artifact download command.
type DownloadOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	ArtifactID string
	Output     string
}

// NewCmdDownload creates the actions artifact download command.
func NewCmdDownload(f *cmdutil.Factory, runF func(*DownloadOptions) error) *cobra.Command {
	opts := &DownloadOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "download <artifact-id>",
		Short: "Download an artifact",
		Long: heredoc.Doc(`
			Download a workflow artifact as a ZIP archive.

			The artifact id is the ` + "`id`" + ` returned by ` + "`gc actions artifact list`" + `.
			The endpoint returns a ZIP archive (binary); use --output to save it,
			or redirect stdout (e.g. > artifact.zip). On an interactive terminal
			without --output it refuses to dump the binary archive.
		`),
		Example: heredoc.Doc(`
			# Save artifact to a file
			$ gc actions artifact download <artifact-id> -R owner/repo --output artifact.zip

			# Stream to a file via redirect
			$ gc actions artifact download <artifact-id> -R owner/repo > artifact.zip
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ArtifactID = strings.TrimSpace(args[0])
			if opts.ArtifactID == "" {
				return cmdutil.NewUsageError("artifact id is required")
			}
			if runF != nil {
				return runF(opts)
			}
			return downloadRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Output, "output", "o", "", "Write artifact ZIP to FILE (default: stdout)")

	return cmd
}

func downloadRun(opts *DownloadOptions) error {
	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}

	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}
	owner, repo, err := cmdutil.ParseRepo(repository)
	if err != nil {
		return err
	}

	// The download returns a binary ZIP archive; refuse to dump it to an
	// interactive terminal. Piped/redirected stdout or --output are supported.
	if strings.TrimSpace(opts.Output) == "" && opts.IO.IsStdoutTTY() {
		return cmdutil.NewUsageError("artifact download is a binary ZIP archive; use --output FILE or redirect stdout (e.g. > artifact.zip)")
	}

	zipBytes, err := api.DownloadActionsArtifact(client, owner, repo, opts.ArtifactID)
	if err != nil {
		return fmt.Errorf("failed to download artifact: %w", err)
	}

	if strings.TrimSpace(opts.Output) != "" {
		if err := os.WriteFile(opts.Output, zipBytes, 0o644); err != nil {
			return fmt.Errorf("failed to write artifact file: %w", err)
		}
		fmt.Fprintf(opts.IO.ErrOut, "Saved artifact to %s (%d bytes)\n", opts.Output, len(zipBytes))
		return nil
	}

	if _, err := opts.IO.Out.Write(zipBytes); err != nil {
		return fmt.Errorf("failed to write artifact output: %w", err)
	}
	return nil
}
