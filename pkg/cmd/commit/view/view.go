// Package view implements the commit view command
package view

import (
	"fmt"
	"net/http"
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
	BaseRepo   func() (string, error)

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
		BaseRepo:   f.BaseRepo,
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

	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
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

	// Get commit
	commit, err := api.GetCommit(client, owner, repo, opts.SHA, opts.ShowDiff)
	if err != nil {
		return cmdutil.WrapNotFound(err, "commit %s not found in %s/%s", opts.SHA, owner, repo)
	}

	// Open in browser if --web flag is set
	if opts.Web {
		fmt.Fprintf(opts.IO.Out, "Opening %s in your browser.\n", commit.HTMLURL)
		if err := browser.Open(commit.HTMLURL); err != nil {
			if opts.IO.IsStdoutTTY() {
				return err
			}
			fmt.Fprintf(opts.IO.ErrOut, "Failed to open browser: %v\n", err)
		}
		return nil
	}

	// Output as JSON
	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, commit)
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
	return cmdutil.ParseRepo(repo)
}
