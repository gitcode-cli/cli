// Package deleteasset implements the release delete-asset command
package deleteasset

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type DeleteAssetOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	TagName      string
	AttachFileID string

	// Flags
	Repository string
	Yes        bool
	DryRun     bool
}

// NewCmdDeleteAsset creates the release delete-asset command
func NewCmdDeleteAsset(f *cmdutil.Factory, runF func(*DeleteAssetOptions) error) *cobra.Command {
	opts := &DeleteAssetOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "delete-asset <tag> <asset-id>",
		Short: "Delete a release asset",
		Long: heredoc.Doc(`
			Delete an attached file (asset) from a release.

			Non-interactive mode: requires --yes to skip confirmation.
		`),
		Example: heredoc.Doc(`
			# Delete a release asset
			$ gc release delete-asset v1.0.0 12345 -R owner/repo

			# Delete without confirmation
			$ gc release delete-asset v1.0.0 12345 -R owner/repo --yes

			# Preview the deletion
			$ gc release delete-asset v1.0.0 12345 -R owner/repo --dry-run
		`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.TagName = args[0]
			opts.AttachFileID = args[1]
			if runF != nil {
				return runF(opts)
			}
			return deleteAssetRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVarP(&opts.Yes, "yes", "y", false, "Skip confirmation")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Preview the deletion without deleting the asset")

	return cmd
}

func deleteAssetRun(opts *DeleteAssetOptions) error {
	cs := opts.IO.ColorScheme()

	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}
	owner, repo, err := parseRepo(repository)
	if err != nil {
		return err
	}

	if opts.DryRun {
		fmt.Fprintf(opts.IO.Out, "Dry run: would delete asset %s from release %s in %s/%s\n", opts.AttachFileID, opts.TagName, owner, repo)
		return nil
	}

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}
	client, err := cmdutil.AuthenticatedClient(httpClient)
	if err != nil {
		return err
	}

	if err := cmdutil.ConfirmOrAbort(cmdutil.ConfirmOptions{
		IO:       opts.IO,
		Yes:      opts.Yes,
		Expected: opts.AttachFileID,
		Prompt:   fmt.Sprintf("! This will delete asset %s from release %s\nType the asset id to confirm: ", cs.Bold(opts.AttachFileID), opts.TagName),
	}); err != nil {
		return err
	}

	if err := api.DeleteReleaseAssetByTag(client, owner, repo, opts.TagName, opts.AttachFileID); err != nil {
		return fmt.Errorf("failed to delete release asset: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Deleted asset %s from release %s\n", cs.Red("✓"), opts.AttachFileID, opts.TagName)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
