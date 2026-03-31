// Package delete implements the label delete command
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
	Name       string

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
		Use:   "delete <name>",
		Short: "Delete a label",
		Long: heredoc.Doc(`
			Delete a label from a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Delete a label
			$ gc label delete old-label

			# Skip confirmation
			$ gc label delete old-label --yes
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]

			if runF != nil {
				return runF(opts)
			}
			return deleteRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVarP(&opts.Yes, "yes", "y", false, "Skip confirmation")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Preview the deletion without deleting the label")

	return cmd
}

func deleteRun(opts *DeleteOptions) error {
	cs := opts.IO.ColorScheme()

	// Get repository
	owner, repo, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	// Confirm deletion
	if opts.DryRun {
		fmt.Fprintf(opts.IO.Out, "Dry run: would delete label %s from %s/%s\n", opts.Name, owner, repo)
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
		Expected: opts.Name,
		Prompt:   fmt.Sprintf("! This will delete label %s\nType the label name to confirm: ", cs.Bold(opts.Name)),
	}); err != nil {
		return err
	}

	// Delete label
	err = api.DeleteLabel(client, owner, repo, opts.Name)
	if err != nil {
		return fmt.Errorf("failed to delete label: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Deleted label %s from %s/%s\n", cs.Red("✗"), opts.Name, owner, repo)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
