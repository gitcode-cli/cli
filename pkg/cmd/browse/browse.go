// Package browse implements the browse command.
package browse

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/browser"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/config"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// BrowseOptions configures the browse command.
type BrowseOptions struct {
	IO       *iostreams.IOStreams
	Config   func() (config.Config, error)
	BaseRepo func() (string, error)
	Browser  func(string) error

	Repository string
	Arg        string

	Branch    string
	Commit    string
	PR        bool
	NoBrowser bool
}

// NewCmdBrowse creates the browse command.
func NewCmdBrowse(f *cmdutil.Factory, runF func(*BrowseOptions) error) *cobra.Command {
	opts := &BrowseOptions{
		IO:       f.IOStreams,
		Config:   f.Config,
		BaseRepo: f.BaseRepo,
		Browser:  browser.Open,
	}

	cmd := &cobra.Command{
		Use:   "browse [<number> | <path>]",
		Short: "Open a GitCode repository in the browser",
		Long: heredoc.Doc(`
			Open a GitCode repository page in your default browser.

			With no arguments, opens the repository home page. With a number,
			opens the issue or pull request with that number (use --pr to
			specify pull request). With a path, opens the corresponding
			repository page (e.g. "wiki", "releases", "issues").

			Use --branch to open a branch page, --commit to open a commit page.
			Use --no-browser to only print the URL without opening a browser
			(useful for scripts and AI agents).
		`),
		Example: heredoc.Doc(`
			# Open the repository home page
			$ gc browse -R owner/repo

			# Open issue #42
			$ gc browse 42 -R owner/repo

			# Open pull request #42
			$ gc browse 42 --pr -R owner/repo

			# Open the releases page
			$ gc browse releases -R owner/repo

			# Open a specific branch
			$ gc browse --branch feature -R owner/repo

			# Open a specific commit
			$ gc browse --commit abc123 -R owner/repo

			# Print URL without opening browser
			$ gc browse --no-browser -R owner/repo
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Arg = args[0]
			}
			if runF != nil {
				return runF(opts)
			}
			return browseRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVar(&opts.Branch, "branch", "", "Open a specific branch")
	cmd.Flags().StringVar(&opts.Commit, "commit", "", "Open a specific commit")
	cmd.Flags().BoolVar(&opts.PR, "pr", false, "Open as pull request (when passing a number)")
	cmd.Flags().BoolVar(&opts.NoBrowser, "no-browser", false, "Print URL only, do not open browser")

	return cmd
}

func browseRun(opts *BrowseOptions) error {
	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return fmt.Errorf("failed to resolve repository: %w", err)
	}
	owner, repo, err := cmdutil.ParseRepo(repository)
	if err != nil {
		return fmt.Errorf("failed to parse repository: %w", err)
	}

	host, err := resolveHost(opts.Config)
	if err != nil {
		return fmt.Errorf("failed to resolve host: %w", err)
	}

	webURL := buildURL(host, owner, repo, opts)

	if opts.NoBrowser || !opts.IO.IsStdoutTTY() || opts.IO.NoInteractive() {
		fmt.Fprintf(opts.IO.Out, "%s\n", webURL)
		return nil
	}

	fmt.Fprintf(opts.IO.ErrOut, "Opening %s in your browser.\n", webURL)
	return opts.Browser(webURL)
}

func resolveHost(configFn func() (config.Config, error)) (string, error) {
	if configFn == nil {
		return "gitcode.com", nil
	}
	cfg, err := configFn()
	if err != nil {
		return "", fmt.Errorf("failed to read config: %w", err)
	}
	host, _ := cfg.Authentication().DefaultHost()
	if host == "" {
		host = "gitcode.com"
	}
	if host == "api.gitcode.com" {
		host = "gitcode.com"
	}
	return host, nil
}

func buildURL(host, owner, repo string, opts *BrowseOptions) string {
	base := fmt.Sprintf("https://%s/%s/%s",
		host, url.PathEscape(owner), url.PathEscape(repo))

	if opts.Commit != "" {
		return base + "/commit/" + url.PathEscape(opts.Commit)
	}
	if opts.Branch != "" {
		return base + "/tree/" + opts.Branch
	}
	if opts.Arg == "" {
		return base
	}
	if isNumeric(opts.Arg) {
		section := "issues"
		if opts.PR {
			section = "pulls"
		}
		return base + "/" + section + "/" + opts.Arg
	}
	return base + "/" + opts.Arg
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
