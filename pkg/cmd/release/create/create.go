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
	Draft      bool
	Prerelease bool
	Target     string
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

			# Create a draft release
			$ gc release create v1.0.0 -R owner/repo --draft

			# Create a prerelease
			$ gc release create v1.0.0-beta -R owner/repo --prerelease

			# Create a release in a specific repository
			$ gc release create v1.0.0 -R owner/repo
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
	cmd.Flags().BoolVarP(&opts.Draft, "draft", "d", false, "Mark as draft")
	cmd.Flags().BoolVarP(&opts.Prerelease, "prerelease", "p", false, "Mark as prerelease")
	cmd.Flags().StringVarP(&opts.Target, "target", "", "", "Target commitish")

	return cmd
}

func createRun(opts *CreateOptions) error {
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

	// Create release
	release, err := api.CreateRelease(client, owner, repo, &api.CreateReleaseOptions{
		TagName:         opts.TagName,
		Name:            opts.Title,
		Body:            opts.Notes,
		Draft:           opts.Draft,
		Prerelease:      opts.Prerelease,
		TargetCommitish: opts.Target,
	})
	if err != nil {
		return fmt.Errorf("failed to create release: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Created release %s\n", cs.Green("✓"), release.TagName)
	fmt.Fprintf(opts.IO.Out, "  %s\n", release.HTMLURL)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
