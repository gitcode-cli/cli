// Package delete implements the repo delete command
package delete

import (
	"fmt"
	"net/http"

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
	Repository string

	// Flags
	Yes    bool
	DryRun bool
}

// NewCmdDelete creates the delete command
func NewCmdDelete(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "delete <repository>",
		Short: "Delete a repository",
		Long: heredoc.Doc(`
			Delete a GitCode repository.

			WARNING: This action is irreversible. Use with caution.
		`),
		Example: heredoc.Doc(`
			# Delete a repository (with confirmation)
			$ gc repo delete owner/repo

			# Skip confirmation
			$ gc repo delete owner/repo --yes
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Repository = args[0]

			if runF != nil {
				return runF(opts)
			}
			return deleteRun(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Yes, "yes", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Preview the deletion without deleting the repository")

	return cmd
}

func deleteRun(opts *DeleteOptions) error {
	cs := opts.IO.ColorScheme()

	// Parse repository
	owner, name, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	// Confirm deletion
	if opts.DryRun {
		fmt.Fprintf(opts.IO.Out, "Dry run: would delete repository %s\n", opts.Repository)
		return nil
	}

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := cmdutil.EnvToken()
	if token == "" {
		return cmdutil.NewAuthError("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "environment")

	if err := cmdutil.ConfirmOrAbort(cmdutil.ConfirmOptions{
		IO:       opts.IO,
		Yes:      opts.Yes,
		Expected: name,
		Prompt:   fmt.Sprintf("! This will delete %s permanently.\nType the repository name to confirm: ", cs.Bold(opts.Repository)),
	}); err != nil {
		return err
	}

	// Delete repo
	err = api.DeleteRepo(client, owner, name)
	if err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Deleted repository %s\n", cs.Red("✗"), opts.Repository)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
