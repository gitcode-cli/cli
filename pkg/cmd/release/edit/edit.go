// Package edit implements the release edit command
package edit

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type EditOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

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
}

// NewCmdEdit creates the edit command
func NewCmdEdit(f *cmdutil.Factory, runF func(*EditOptions) error) *cobra.Command {
	opts := &EditOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "edit <tag>",
		Short: "Edit a release",
		Long: heredoc.Doc(`
			Edit an existing release.

			You can update the release title, notes, draft status, prerelease status,
			or target branch.
		`),
		Example: heredoc.Doc(`
			# Edit release title and notes
			$ gc release edit v1.0.0 -R owner/repo --title "New Title" --notes "New notes"

			# Mark release as prerelease
			$ gc release edit v1.0.0 -R owner/repo --prerelease

			# Read notes from a file
			$ gc release edit v1.0.0 -R owner/repo --notes-file RELEASE_NOTES.md
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

	return cmd
}

func editRun(opts *EditOptions) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := cmdutil.EnvToken()
	if token == "" {
		return cmdutil.NewAuthError("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "active")

	// Get repository
	owner, repo, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	// Build update options
	updateOpts := &api.UpdateReleaseOptions{}

	if opts.Title != "" {
		updateOpts.Name = opts.Title
	}

	if opts.Notes != "" {
		updateOpts.Body = opts.Notes
	}

	if opts.NotesFile != "" {
		content, err := os.ReadFile(opts.NotesFile)
		if err != nil {
			return fmt.Errorf("failed to read notes file: %w", err)
		}
		updateOpts.Body = string(content)
	}

	if opts.Draft != "" {
		draft, err := parseBool(opts.Draft)
		if err != nil {
			return fmt.Errorf("invalid draft value: %s", opts.Draft)
		}
		updateOpts.Draft = &draft
	}

	if opts.Prerelease != "" {
		prerelease, err := parseBool(opts.Prerelease)
		if err != nil {
			return fmt.Errorf("invalid prerelease value: %s", opts.Prerelease)
		}
		updateOpts.Prerelease = &prerelease
	}

	if opts.Target != "" {
		updateOpts.TargetCommitish = opts.Target
	}

	// Update release
	release, err := api.UpdateReleaseByTag(client, owner, repo, opts.TagName, updateOpts)
	if err != nil {
		if err == api.ErrNoReleaseID {
			return fmt.Errorf("failed to update release: %w; GitCode currently omits release IDs in release lookup responses", err)
		}
		return fmt.Errorf("failed to update release: %w", err)
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
		return false, fmt.Errorf("invalid boolean value: %s", s)
	}
}
