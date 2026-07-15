// Package create implements the milestone create command
package create

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type CreateOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	Repository string
	Title      string

	// Flags
	Description string
	DueDate     string
	JSON        bool
}

// NewCmdCreate creates the create command
func NewCmdCreate(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "create <title>",
		Short: "Create a milestone",
		Long: heredoc.Doc(`
			Create a new milestone in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Create a milestone (due date defaults to 30 days from now)
			$ gc milestone create "v1.0" -R owner/repo

			# Create with specific due date
			$ gc milestone create "v1.0" --due-date "2024-12-31" -R owner/repo

			# Create with description
			$ gc milestone create "v2.0" --description "Next release" --due-date "2025-01-31" -R owner/repo

			# Output as JSON
			$ gc milestone create "v1.0" -R owner/repo --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Title = args[0]

			if runF != nil {
				return runF(opts)
			}
			return createRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Description, "description", "d", "", "Description")
	cmd.Flags().StringVarP(&opts.DueDate, "due-date", "D", "", "Due date (YYYY-MM-DD)")
	cmd.Flags().BoolVar(&opts.JSON, "json", false, "Output as JSON")

	return cmd
}

func createRun(opts *CreateOptions) error {
	if opts.Description != "" {
		if err := cmdutil.ScanContentForSecrets(opts.Description); err != nil {
			return err
		}
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

	// Get repository
	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}
	owner, repo, err := parseRepo(repository)
	if err != nil {
		return err
	}

	// Parse due date (required by GitCode API)
	var dueOn string
	if opts.DueDate != "" {
		_, err := time.Parse("2006-01-02", opts.DueDate)
		if err != nil {
			return cmdutil.NewUsageError("invalid due date format, use YYYY-MM-DD")
		}
		dueOn = opts.DueDate
	} else {
		// Default to 30 days from now
		dueOn = time.Now().AddDate(0, 0, 30).Format("2006-01-02")
	}

	// Create milestone
	ms, err := api.CreateMilestone(client, owner, repo, &api.CreateMilestoneOptions{
		Title:       opts.Title,
		Description: opts.Description,
		DueOn:       dueOn,
	})
	if err != nil {
		return fmt.Errorf("failed to create milestone: %w", err)
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, ms)
	}

	fmt.Fprintf(opts.IO.Out, "%s Created milestone #%d %s in %s/%s\n", cs.Green("✓"), ms.Number, cs.Bold(ms.Title), owner, repo)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
