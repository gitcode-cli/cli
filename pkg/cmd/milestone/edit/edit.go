// Package edit implements the milestone edit command
package edit

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

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

	Repository string
	Number     int

	// Flags
	Title           string
	Description     string
	DescriptionFile string
	State           string
	DueDate         string
	JSON            bool
	Yes             bool // Skip confirmation for state changes
}

// NewCmdEdit creates the edit command
func NewCmdEdit(f *cmdutil.Factory, runF func(*EditOptions) error) *cobra.Command {
	opts := &EditOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "edit <number>",
		Short: "Edit a milestone",
		Long: heredoc.Doc(`
			Edit an existing milestone in a GitCode repository.

			You can update the title, description, state, and due date.
		`),
		Example: heredoc.Doc(`
			# Edit milestone title
			$ gc milestone edit 5 --title "New Title" -R owner/repo

			# Edit milestone description
			$ gc milestone edit 5 --description "Updated description" -R owner/repo

			# Edit milestone description from file
			$ gc milestone edit 5 --description-file milestone-desc.md -R owner/repo

			# Close a milestone
			$ gc milestone edit 5 --state closed -R owner/repo

			# Reopen a milestone
			$ gc milestone edit 5 --state open -R owner/repo

			# Edit due date
			$ gc milestone edit 5 --due-date "2024-12-31" -R owner/repo

			# Edit with JSON output
			$ gc milestone edit 5 --title "New Title" --json -R owner/repo

			# Combine multiple edits
			$ gc milestone edit 5 --title "v2.0" --description "Next release" --due-date "2025-01-31" -R owner/repo
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return cmdutil.NewUsageError(fmt.Sprintf("invalid milestone number: %s", args[0]))
			}
			opts.Number = number

			if runF != nil {
				return runF(opts)
			}
			return editRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "New milestone title")
	cmd.Flags().StringVarP(&opts.Description, "description", "d", "", "New description")
	cmd.Flags().StringVarP(&opts.DescriptionFile, "description-file", "F", "", "Read description from file")
	cmd.Flags().StringVarP(&opts.State, "state", "s", "", "Milestone state (open/closed)")
	cmdutil.SetFlagEnum(cmd, "state", "open", "closed")
	cmd.Flags().StringVarP(&opts.DueDate, "due-date", "D", "", "Due date (YYYY-MM-DD)")
	cmd.Flags().BoolVarP(&opts.Yes, "yes", "y", false, "Skip confirmation for state changes")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func editRun(opts *EditOptions) error {
	cs := opts.IO.ColorScheme()

	// Read description from file if specified
	description := opts.Description
	if opts.DescriptionFile != "" {
		content, err := os.ReadFile(opts.DescriptionFile)
		if err != nil {
			return fmt.Errorf("failed to read description file: %w", err)
		}
		description = string(content)
	}

	// Validate at least one edit option is provided
	if opts.Title == "" && description == "" && opts.State == "" && opts.DueDate == "" {
		return cmdutil.NewUsageError("at least one edit option is required (e.g., --title, --description, --description-file, --state, --due-date)")
	}

	// Validate state value
	if opts.State != "" && opts.State != "open" && opts.State != "closed" {
		return cmdutil.NewUsageError(fmt.Sprintf("invalid state value '%s': must be 'open' or 'closed'", opts.State))
	}

	// Validate due-date format
	if opts.DueDate != "" {
		_, err := time.Parse("2006-01-02", opts.DueDate)
		if err != nil {
			return cmdutil.NewUsageError(fmt.Sprintf("invalid due date format '%s': use YYYY-MM-DD", opts.DueDate))
		}
	}

	// Confirm state changes (close)
	if opts.State == "closed" {
		if err := cmdutil.ConfirmOrAbort(cmdutil.ConfirmOptions{
			IO:       opts.IO,
			Yes:      opts.Yes,
			Expected: "closed",
			Prompt:   fmt.Sprintf("This will close milestone #%d. Type 'closed' to confirm: ", opts.Number),
		}); err != nil {
			return err
		}
	}

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

	// Build update options
	updateOpts := &api.UpdateMilestoneOptions{
		Title:       opts.Title,
		Description: description,
		State:       opts.State,
		DueOn:       opts.DueDate,
	}

	// If only partial fields are provided, fetch current milestone to preserve others
	// GitCode API may require certain fields (like due_on) to be non-blank
	if opts.Title == "" || opts.DueDate == "" {
		current, err := api.GetMilestone(client, owner, repo, opts.Number)
		if err != nil {
			return fmt.Errorf("failed to fetch current milestone: %w", err)
		}
		// Preserve current values for fields not being updated
		if opts.Title == "" && current.Title != "" {
			updateOpts.Title = current.Title
		}
		if description == "" && current.Description != "" {
			updateOpts.Description = current.Description
		}
		if opts.DueDate == "" && current.DueOn != "" {
			updateOpts.DueOn = current.DueOn
		}
		if opts.State == "" && current.State != "" {
			updateOpts.State = current.State
		}
	}

	ms, err := api.UpdateMilestone(client, owner, repo, opts.Number, updateOpts)
	if err != nil {
		return fmt.Errorf("failed to update milestone: %w", err)
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, ms)
	}

	fmt.Fprintf(opts.IO.Out, "%s Updated milestone #%d in %s/%s\n", cs.Green("✓"), ms.Number, owner, repo)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
