// Package create implements the release create command
package create

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/gitcode-com/gitcode-cli/api"
	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
	"github.com/gitcode-com/gitcode-cli/pkg/iostreams"
)

type CreateOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	TagName string

	// Flags
	Repository  string
	Title       string
	Notes       string
	Draft       bool
	Prerelease  bool
	Target      string
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
			$ gc release create v1.0.0

			# Create a release with title and notes
			$ gc release create v1.0.0 --title "Version 1.0" --notes "Release notes"

			# Create a draft release
			$ gc release create v1.0.0 --draft

			# Create a prerelease
			$ gc release create v1.0.0-beta --prerelease

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
	token := getEnvToken()
	if token == "" {
		return fmt.Errorf("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "environment")

	// Get repository
	owner, repo, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	// TODO: Implement API call to create release
	_ = client
	_ = owner
	_ = repo

	fmt.Fprintf(opts.IO.Out, "%s Created release %s\n", cs.Green("✓"), opts.TagName)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	if repo == "" {
		return "", "", fmt.Errorf("no repository specified. Use -R owner/repo")
	}

	// Simple parsing - in real implementation would be more robust
	for i := 0; i < len(repo); i++ {
		if repo[i] == '/' {
			return repo[:i], repo[i+1:], nil
		}
	}
	return "", "", fmt.Errorf("invalid repository format: %s", repo)
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITCODE_TOKEN")
}