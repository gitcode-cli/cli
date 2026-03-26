// Package view implements the commit view command
package view

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

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
	Repository string
	SHA        string

	// Flags
	ShowDiff bool
	Web      bool
	JSON     bool
}

// NewCmdView creates the view command
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "view <sha>",
		Short: "View a commit",
		Long: heredoc.Doc(`
			View a commit in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# View a commit
			$ gc commit view abc123 -R owner/repo

			# View commit with diff files
			$ gc commit view abc123 -R owner/repo --show-diff

			# View commit in browser
			$ gc commit view abc123 -R owner/repo --web

			# Output as JSON
			$ gc commit view abc123 -R owner/repo --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.SHA = args[0]

			if runF != nil {
				return runF(opts)
			}
			return viewRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVar(&opts.ShowDiff, "show-diff", false, "Show diff files")
	cmd.Flags().BoolVarP(&opts.Web, "web", "w", false, "Open in browser")
	cmd.Flags().BoolVar(&opts.JSON, "json", false, "Output as JSON")

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

	// Get commit
	commit, err := api.GetCommit(client, owner, repo, opts.SHA, opts.ShowDiff)
	if err != nil {
		return fmt.Errorf("failed to get commit: %w", err)
	}

	// Open in browser if --web flag is set
	if opts.Web {
		fmt.Fprintf(opts.IO.Out, "Opening %s in your browser.\n", commit.HTMLURL)
		return browser.Open(commit.HTMLURL)
	}

	// Output as JSON
	if opts.JSON {
		data, err := json.MarshalIndent(commit, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal commit: %w", err)
		}
		fmt.Fprintf(opts.IO.Out, "%s\n", data)
		return nil
	}

	// Output
	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Commit:"), commit.SHA)
	if commit.Commit != nil && commit.Commit.Author != nil {
		fmt.Fprintf(opts.IO.Out, "%s %s <%s>\n", cs.Bold("Author:"), commit.Commit.Author.Name, commit.Commit.Author.Email)
	}
	if commit.Commit != nil {
		fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Date:"), commit.Commit.Author.Date)
		fmt.Fprintf(opts.IO.Out, "\n")
		fmt.Fprintf(opts.IO.Out, "    %s\n", strings.ReplaceAll(commit.Commit.Message, "\n", "\n    "))
	}
	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "  %s\n", commit.HTMLURL)
	fmt.Fprintf(opts.IO.Out, "\n")

	// Show stats
	if commit.Stats != nil {
		fmt.Fprintf(opts.IO.Out, "  %d files changed, %d insertions(+), %d deletions(-)\n",
			len(commit.Files),
			commit.Stats.Additions,
			commit.Stats.Deletions)
		fmt.Fprintf(opts.IO.Out, "\n")
	}

	// Show files if show-diff is set
	if opts.ShowDiff && len(commit.Files) > 0 {
		fmt.Fprintf(opts.IO.Out, "--- Files ---\n\n")
		for _, f := range commit.Files {
			status := f.Status
			if status == "" {
				status = "modified"
			}
			fmt.Fprintf(opts.IO.Out, "%s\t%s\n", status, f.Filename)
		}
		fmt.Fprintf(opts.IO.Out, "\n")
	}

	return nil
}

func parseRepo(repo string) (string, string, error) {
	if repo == "" {
		return "", "", fmt.Errorf("no repository specified. Use -R owner/repo")
	}

	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repository format: %s", repo)
	}
	return parts[0], parts[1], nil
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITCODE_TOKEN")
}