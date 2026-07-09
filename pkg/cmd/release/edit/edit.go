// Package edit implements the release edit command
package edit

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type EditOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	TagName string

	// Flags
	Repository string
	Title      string
	Notes      string
	NotesFile  string
	Draft      string
	Prerelease string
	Target     string
	JSON       bool
}

// NewCmdEdit creates the edit command
func NewCmdEdit(f *cmdutil.Factory, runF func(*EditOptions) error) *cobra.Command {
	opts := &EditOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "edit <tag>",
		Short: "Edit a release",
		Long: heredoc.Doc(`
			Edit an existing release.

			You can update the release title, notes, and prerelease status.
			--draft and --target flags are not supported by the GitCode API.
		`),
		Example: heredoc.Doc(`
			# Edit release title and notes
			$ gc release edit v1.0.0 -R owner/repo --title "New Title" --notes "New notes"

			# Mark release as prerelease (release_status=pre)
			$ gc release edit v1.0.0 -R owner/repo --prerelease true

			# Mark release as full release (release_status=latest)
			$ gc release edit v1.0.0 -R owner/repo --prerelease false

			# Read notes from a file
			$ gc release edit v1.0.0 -R owner/repo --notes-file RELEASE_NOTES.md

			# Output as JSON
			$ gc release edit v1.0.0 -R owner/repo --title "New Title" --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.TagName = args[0]

			if runF != nil {
				return runF(opts)
			}
			return editRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "Release title")
	cmd.Flags().StringVarP(&opts.Notes, "notes", "n", "", "Release notes")
	cmd.Flags().StringVarP(&opts.NotesFile, "notes-file", "F", "", "Read release notes from file")
	cmd.Flags().StringVar(&opts.Draft, "draft", "", "Mark as draft (true/false)")
	cmd.Flags().StringVar(&opts.Prerelease, "prerelease", "", "Mark as prerelease (true/false)")
	cmd.Flags().StringVarP(&opts.Target, "target", "T", "", "Target branch or commit SHA")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func editRun(opts *EditOptions) error {
	cs := opts.IO.ColorScheme()

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

	// Check unsupported flags and output warnings
	if opts.Draft != "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s warning: --draft is not supported by GitCode release edit API, ignoring\n", cs.Yellow("⚠"))
	}
	if opts.Target != "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s warning: --target is not supported by GitCode release edit API, ignoring\n", cs.Yellow("⚠"))
	}

	// Check if there's anything to update (only check supported flags)
	// This check is done before API calls to avoid unnecessary requests
	if opts.Title == "" && opts.Notes == "" && opts.NotesFile == "" && opts.Prerelease == "" {
		return cmdutil.NewUsageError("no changes specified. Use --title, --notes, --notes-file, or --prerelease to specify what to edit")
	}

	// Scan inline --notes for secrets before any API call (fail-fast)
	if opts.Notes != "" {
		if err := cmdutil.ScanContentForSecrets(opts.Notes); err != nil {
			return err
		}
	}

	// Build update options - name and body are required by GitCode API
	// First get existing release to populate required fields
	existingRelease, err := api.GetRelease(client, owner, repo, opts.TagName)
	if err != nil {
		return cmdutil.WrapNotFound(err, "release %s not found in %s/%s", opts.TagName, owner, repo)
	}

	// Build GitCodeUpdateReleaseOptions with required fields
	updateOpts := &api.GitCodeUpdateReleaseOptions{
		Name: existingRelease.Name, // Default to existing
		Body: existingRelease.Body, // Default to existing
	}

	// Override with user-provided values
	if opts.Title != "" {
		updateOpts.Name = opts.Title
	}
	if opts.Notes != "" {
		updateOpts.Body = opts.Notes
	}
	if opts.NotesFile != "" {
		content, err := cmdutil.ReadTextFile(opts.NotesFile)
		if err != nil {
			return fmt.Errorf("failed to read notes file: %w", err)
		}
		updateOpts.Body = content
	}

	// Map prerelease to release_status
	if opts.Prerelease != "" {
		prerelease, err := parseBool(opts.Prerelease)
		if err != nil {
			return cmdutil.NewUsageError(fmt.Sprintf("invalid prerelease value: %s", opts.Prerelease))
		}
		if prerelease {
			updateOpts.ReleaseStatus = "pre"
		} else {
			updateOpts.ReleaseStatus = "latest"
		}
	}

	// Update release using direct PATCH by tag
	release, err := api.UpdateReleaseByTagDirect(client, owner, repo, opts.TagName, updateOpts)
	if err != nil {
		return fmt.Errorf("failed to update release: %w", err)
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, release)
	}

	title := release.TagName
	if release.Name != "" {
		title = release.Name
	}

	fmt.Fprintf(opts.IO.Out, "%s Updated release %s\n", cs.Green("✓"), title)
	if release.HTMLURL != "" {
		fmt.Fprintf(opts.IO.Out, "  %s\n", release.HTMLURL)
	}
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

func parseBool(s string) (bool, error) {
	switch s {
	case "true", "yes", "1":
		return true, nil
	case "false", "no", "0":
		return false, nil
	default:
		return false, cmdutil.NewUsageError(fmt.Sprintf("invalid boolean value: %s", s))
	}
}
