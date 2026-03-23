// Package delete implements the release delete command
package delete

import (
	"bufio"
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

type DeleteOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	TagName string

	// Flags
	Repository string
	Yes        bool
}

// NewCmdDelete creates the delete command
func NewCmdDelete(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "delete <tag>",
		Short: "Delete a release",
		Long: heredoc.Doc(`
			Delete a release from a repository.

			This will delete the release but not the associated git tag.
		`),
		Example: heredoc.Doc(`
			# Delete a release
			$ gc release delete v1.0.0

			# Delete without confirmation
			$ gc release delete v1.0.0 --yes
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.TagName = args[0]

			if runF != nil {
				return runF(opts)
			}
			return deleteRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVarP(&opts.Yes, "yes", "y", false, "Skip confirmation")

	return cmd
}

func deleteRun(opts *DeleteOptions) error {
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

	// Get release for confirmation
	release, err := api.GetRelease(client, owner, repo, opts.TagName)
	if err != nil {
		return fmt.Errorf("failed to get release: %w", err)
	}

	title := release.TagName
	if release.Name != "" {
		title = release.Name
	}

	// Confirm deletion
	if !opts.Yes {
		fmt.Fprintf(opts.IO.ErrOut, "! This will delete release %s\n", cs.Bold(title))
		fmt.Fprintf(opts.IO.ErrOut, "Type the tag name to confirm: ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input != opts.TagName {
			return fmt.Errorf("confirmation did not match tag name")
		}
	}

	// Delete release
	err = api.DeleteReleaseByTag(client, owner, repo, opts.TagName)
	if err != nil {
		return fmt.Errorf("failed to delete release: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Deleted release %s\n", cs.Red("✓"), opts.TagName)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	if repo == "" {
		return "", "", fmt.Errorf("no repository specified. Use -R owner/repo")
	}

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