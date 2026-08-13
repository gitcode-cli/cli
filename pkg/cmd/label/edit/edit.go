// Package edit implements the label edit command
package edit

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type EditOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	Repository string
	Name       string

	// Flags
	NewName string
	Color   string
	JSON    bool
}

// NewCmdEdit creates the label edit command
func NewCmdEdit(f *cmdutil.Factory, runF func(*EditOptions) error) *cobra.Command {
	opts := &EditOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "edit <name>",
		Short: "Edit a label",
		Long: heredoc.Doc(`
			Update a label's name and/or color in a GitCode repository.

			The original label name is passed as the argument. Use --new-name
			and/or --color to specify changes; at least one is required.

			Note: the GitCode API does not support updating a label's
			description via this endpoint, so --description is not provided.
		`),
		Example: heredoc.Doc(`
			# Rename a label
			$ gc label edit bug --new-name "bug-fix" -R owner/repo

			# Change a label's color
			$ gc label edit bug --color "#ff0000" -R owner/repo

			# Output as JSON
			$ gc label edit bug --color "#ff0000" --json -R owner/repo
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			if runF != nil {
				return runF(opts)
			}
			return editRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVar(&opts.NewName, "new-name", "", "New name for the label")
	cmd.Flags().StringVarP(&opts.Color, "color", "c", "", "New color in hex format (e.g., #ff0000)")
	cmd.Flags().BoolVar(&opts.JSON, "json", false, "Output as JSON")

	return cmd
}

func editRun(opts *EditOptions) error {
	if opts.NewName == "" && opts.Color == "" {
		return cmdutil.NewUsageError("at least one of --new-name or --color is required")
	}

	colorPattern := regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	if opts.Color != "" && !colorPattern.MatchString(opts.Color) {
		return cmdutil.NewUsageError("--color must be in hex format (e.g., #ff0000)")
	}

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
	owner, repo, err := parseRepo(repository)
	if err != nil {
		return err
	}

	label, err := api.UpdateLabel(client, owner, repo, opts.Name, &api.UpdateLabelOptions{
		NewName: opts.NewName,
		Color:   opts.Color,
	})
	if err != nil {
		return fmt.Errorf("failed to update label: %w", err)
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, label)
	}

	name := label.Name
	if opts.NewName != "" {
		name = opts.NewName
	}
	fmt.Fprintf(opts.IO.Out, "%s Updated label %s (ID: %s) in %s/%s\n", cs.Green("✓"), cs.Bold(name), cmdutil.FormatAPIID(label.ID), owner, repo)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
