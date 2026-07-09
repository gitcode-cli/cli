// Package create implements the release create command
package create

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type CreateOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	TagName string

	// Flags
	Repository string
	Title      string
	Notes      string
	NotesFile  string
	Draft      bool
	Prerelease bool
	Target     string
	JSON       bool
}

// NewCmdCreate creates the create command
func NewCmdCreate(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "create <tag>",
		Short: "Create a new release",
		Long: heredoc.Doc(`
			Create a new release for a repository.

			A release is associated with a git tag and can include release notes,
			binary assets, and other metadata.
		`),
		Example: heredoc.Doc(`
			# Create a release
			$ gc release create v1.0.0 -R owner/repo

			# Create a release with title and notes
			$ gc release create v1.0.0 -R owner/repo --title "Version 1.0" --notes "Release notes"

			# Create a release with notes from file
			$ gc release create v1.0.0 -R owner/repo --notes-file CHANGELOG.md

			# Create a draft release
			$ gc release create v1.0.0 -R owner/repo --draft

			# Create a prerelease
			$ gc release create v1.0.0-beta -R owner/repo --prerelease

			# Create a release with JSON output
			$ gc release create v1.0.0 -R owner/repo --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.TagName = args[0]

			if runF != nil {
				return runF(opts)
			}
			return createRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "Release title")
	cmd.Flags().StringVarP(&opts.Notes, "notes", "n", "", "Release notes")
	cmd.Flags().StringVarP(&opts.NotesFile, "notes-file", "F", "", "Read release notes from file")
	cmd.Flags().BoolVarP(&opts.Draft, "draft", "d", false, "Mark as draft (currently unsupported by GitCode release create API)")
	cmd.Flags().BoolVarP(&opts.Prerelease, "prerelease", "p", false, "Mark as prerelease")
	cmd.Flags().StringVarP(&opts.Target, "target", "", "", "Target commitish")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func createRun(opts *CreateOptions) error {
	cs := opts.IO.ColorScheme()

	// Validate mutual exclusion of --notes and --notes-file
	if opts.Notes != "" && opts.NotesFile != "" {
		return cmdutil.NewUsageError("cannot use both --notes and --notes-file")
	}
	if opts.Draft {
		return cmdutil.NewUsageError("--draft is not supported by GitCode release create API; create the release without --draft or manage draft state in the web UI")
	}

	// Get release body from notes or file
	body := opts.Notes
	if body != "" {
		if err := cmdutil.ScanContentForSecrets(body); err != nil {
			return err
		}
	}
	if opts.NotesFile != "" {
		content, err := cmdutil.ReadTextFile(opts.NotesFile)
		if err != nil {
			return fmt.Errorf("failed to read notes file: %w", err)
		}
		body = content
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
	owner, repo, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	releaseStatus := ""
	if opts.Prerelease {
		releaseStatus = "pre"
	}

	// Create release
	release, err := api.CreateRelease(client, owner, repo, &api.CreateReleaseOptions{
		TagName:         opts.TagName,
		Name:            opts.Title,
		Body:            body,
		Draft:           opts.Draft,
		Prerelease:      opts.Prerelease,
		ReleaseStatus:   releaseStatus,
		TargetCommitish: opts.Target,
	})
	if err != nil {
		return fmt.Errorf("failed to create release: %w", err)
	}
	if opts.Prerelease {
		if fetched, err := api.GetRelease(client, owner, repo, opts.TagName); err == nil {
			release = fetched
		}
		if !release.Prerelease {
			return fmt.Errorf("release %s was created, but GitCode API did not mark it as prerelease", opts.TagName)
		}
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, release)
	}
	fmt.Fprintf(opts.IO.Out, "%s Created release %s\n", cs.Green("✓"), release.TagName)
	fmt.Fprintf(opts.IO.Out, "  %s\n", release.HTMLURL)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
