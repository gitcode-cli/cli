// Package view implements the release view command
package view

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/browser"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	TagName string

	// Flags
	Repository string
	Web        bool
}

// NewCmdView creates the view command
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "view <tag>",
		Short: "View release details",
		Long: heredoc.Doc(`
			View details of a release.
		`),
		Example: heredoc.Doc(`
			# View a release
			$ gc release view v1.0.0

			# View in browser
			$ gc release view v1.0.0 --web
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.TagName = args[0]

			if runF != nil {
				return runF(opts)
			}
			return viewRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVarP(&opts.Web, "web", "w", false, "Open in browser")

	return cmd
}

func viewRun(opts *ViewOptions) error {
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

	// Get release
	release, err := api.GetRelease(client, owner, repo, opts.TagName)
	if err != nil {
		return fmt.Errorf("failed to get release: %w", err)
	}

	// Open in browser if requested
	if opts.Web {
		fmt.Fprintf(opts.IO.Out, "Opening %s in your browser.\n", release.HTMLURL)
		return browser.Open(release.HTMLURL)
	}

	// Output
	title := release.TagName
	if release.Name != "" {
		title = release.Name
	}

	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "%s\n", cs.Bold(title))
	fmt.Fprintf(opts.IO.Out, "  Tag: %s\n", release.TagName)

	// Status
	if release.Draft {
		fmt.Fprintf(opts.IO.Out, "  Status: %s\n", cs.Gray("draft"))
	} else if release.Prerelease {
		fmt.Fprintf(opts.IO.Out, "  Status: %s\n", cs.Yellow("pre-release"))
	} else {
		fmt.Fprintf(opts.IO.Out, "  Status: %s\n", cs.Green("published"))
	}

	// Dates
	if release.PublishedAt != nil {
		fmt.Fprintf(opts.IO.Out, "  Published: %s\n", release.PublishedAt.Format("2006-01-02 15:04:05"))
	}
	if !release.CreatedAt.IsZero() {
		fmt.Fprintf(opts.IO.Out, "  Created: %s\n", release.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	// URL
	fmt.Fprintf(opts.IO.Out, "  URL: %s\n", release.HTMLURL)

	// Body
	if release.Body != "" {
		fmt.Fprintf(opts.IO.Out, "\n")
		fmt.Fprintf(opts.IO.Out, "%s\n", release.Body)
	}

	// Assets
	if len(release.Assets) > 0 {
		fmt.Fprintf(opts.IO.Out, "\n%s\n", cs.Bold("Assets:"))
		for _, asset := range release.Assets {
			fmt.Fprintf(opts.IO.Out, "  %s (%d bytes, %d downloads)\n", asset.Name, asset.Size, asset.Downloads)
		}
	}

	fmt.Fprintf(opts.IO.Out, "\n")
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITCODE_TOKEN")
}
