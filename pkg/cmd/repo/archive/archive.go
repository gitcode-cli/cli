// Package archive implements the repo archive and unarchive commands.
package archive

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

const (
	statusActive   = 0
	statusArchived = 2
)

// Options holds the configuration for the archive/unarchive commands.
type Options struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string

	// Flags
	Yes    bool
	DryRun bool
	JSON   bool
}

// Result represents the JSON output for archive/unarchive.
type Result struct {
	Repository string `json:"repository"`
	Action     string `json:"action"`
}

// NewCmdArchive creates the repo archive command.
func NewCmdArchive(f *cmdutil.Factory, runF func(*Options) error) *cobra.Command {
	return newCmd(f, runF, "archive", "Archive a repository", statusArchived, "archived", true)
}

// NewCmdUnarchive creates the repo unarchive command.
func NewCmdUnarchive(f *cmdutil.Factory, runF func(*Options) error) *cobra.Command {
	return newCmd(f, runF, "unarchive", "Unarchive a repository", statusActive, "unarchived", false)
}

func newCmd(f *cmdutil.Factory, runF func(*Options) error, verb, short string, status int, pastTense string, irreversible bool) *cobra.Command {
	opts := &Options{IO: f.IOStreams, HttpClient: f.HttpClient}

	long := fmt.Sprintf("%s a GitCode repository.", strings.ToUpper(verb[:1])+verb[1:])
	if irreversible {
		long += " Archived repositories become read-only."
	}
	long += "\n\n\t\tNon-interactive mode: requires --yes to skip confirmation."

	cmd := &cobra.Command{
		Use:   verb + " <repository>",
		Short: short,
		Long:  heredoc.Doc(long),
		Example: heredoc.Doc(fmt.Sprintf(`
			# %s a repository
			$ gc repo %s owner/repo --yes
		`, strings.ToUpper(verb[:1])+verb[1:], verb)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Repository = args[0]
			if runF != nil {
				return runF(opts)
			}
			return run(opts, status, verb, pastTense)
		},
	}

	cmd.Flags().BoolVarP(&opts.Yes, "yes", "y", false, "Skip confirmation")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Preview without "+verb+"ing the repository")
	cmd.Flags().BoolVar(&opts.JSON, "json", false, "Output as JSON")

	return cmd
}

func run(opts *Options, status int, verb, pastTense string) error {
	cs := opts.IO.ColorScheme()

	owner, name, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	if opts.DryRun {
		result := Result{Repository: opts.Repository, Action: "dry_run"}
		if opts.JSON {
			return cmdutil.WriteJSON(opts.IO.Out, result)
		}
		fmt.Fprintf(opts.IO.Out, "Dry run: would %s repository %s\n", verb, opts.Repository)
		return nil
	}

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}
	client, err := cmdutil.AuthenticatedClient(httpClient)
	if err != nil {
		return err
	}

	if err := cmdutil.ConfirmOrAbort(cmdutil.ConfirmOptions{
		IO:       opts.IO,
		Yes:      opts.Yes,
		Expected: name,
		Prompt:   fmt.Sprintf("! This will %s %s\nType the repository name to confirm: ", verb, cs.Bold(opts.Repository)),
	}); err != nil {
		return err
	}

	if err := api.SetRepoStatus(client, owner, name, status); err != nil {
		return fmt.Errorf("failed to %s repository: %w", verb, err)
	}

	result := Result{Repository: opts.Repository, Action: pastTense}
	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, result)
	}

	fmt.Fprintf(opts.IO.Out, "%s %s repository %s\n", cs.Green("✓"), strings.ToUpper(pastTense[:1])+pastTense[1:], opts.Repository)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
