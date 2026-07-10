// Package delete implements the actions artifact delete command.
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

// DeleteResult represents the result of an artifact delete operation.
type DeleteResult struct {
	ArtifactID string `json:"artifact_id"`
	Owner      string `json:"owner"`
	Repo       string `json:"repo"`
	Action     string `json:"action"`
}

// DeleteOptions configures the actions artifact delete command.
type DeleteOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	ArtifactID string
	Yes        bool
	DryRun     bool
	JSON       bool
}

// NewCmdDelete creates the actions artifact delete command.
func NewCmdDelete(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "delete <artifact-id>",
		Short: "Delete an artifact",
		Long: heredoc.Doc(`
			Delete a workflow artifact from a GitCode repository.

			This is a destructive operation. By default it asks for confirmation
			(type the artifact id to confirm). In non-interactive mode, --yes is
			required to skip confirmation. Use --dry-run to preview without deleting.
		`),
		Example: heredoc.Doc(`
			# Delete an artifact (interactive)
			$ gc actions artifact delete <artifact-id> -R owner/repo

			# Non-interactive (requires --yes)
			$ gc actions artifact delete <artifact-id> -R owner/repo --yes

			# Preview the deletion
			$ gc actions artifact delete <artifact-id> -R owner/repo --dry-run

			# JSON output
			$ gc actions artifact delete <artifact-id> -R owner/repo --yes --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ArtifactID = args[0]
			if opts.ArtifactID == "" {
				return cmdutil.NewUsageError("artifact id is required")
			}
			if runF != nil {
				return runF(opts)
			}
			return deleteRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVar(&opts.Yes, "yes", false, "Skip confirmation (required in non-interactive mode)")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Preview the deletion without deleting")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func deleteRun(opts *DeleteOptions) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}
	client, err := cmdutil.AuthenticatedClient(httpClient)
	if err != nil {
		return err
	}

	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}
	owner, repo, err := cmdutil.ParseRepo(repository)
	if err != nil {
		return err
	}

	// Dry run
	if opts.DryRun {
		result := DeleteResult{
			ArtifactID: opts.ArtifactID,
			Owner:      owner,
			Repo:       repo,
			Action:     "dry_run",
		}
		if opts.JSON {
			return cmdutil.WriteJSON(opts.IO.Out, result)
		}
		fmt.Fprintf(opts.IO.Out, "Dry run: would delete artifact %s from %s/%s\n", opts.ArtifactID, owner, repo)
		return nil
	}

	// Confirmation gate
	if err := cmdutil.ConfirmOrAbort(cmdutil.ConfirmOptions{
		IO:       opts.IO,
		Yes:      opts.Yes,
		Expected: opts.ArtifactID,
		Prompt:   fmt.Sprintf("! This will delete artifact %s\nType the artifact id to confirm: ", cs.Bold(opts.ArtifactID)),
	}); err != nil {
		return err
	}

	// Delete
	err = api.DeleteActionsArtifact(client, owner, repo, opts.ArtifactID)
	if err != nil {
		return fmt.Errorf("failed to delete artifact: %w", err)
	}

	result := DeleteResult{
		ArtifactID: opts.ArtifactID,
		Owner:      owner,
		Repo:       repo,
		Action:     "deleted",
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, result)
	}

	fmt.Fprintf(opts.IO.Out, "%s Deleted artifact %s from %s/%s\n", cs.Red("✗"), opts.ArtifactID, owner, repo)
	return nil
}
