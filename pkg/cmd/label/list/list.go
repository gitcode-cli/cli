// Package list implements the label list command
package list

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

type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	Repository string

	// Flags
	Limit      int
	Page       int
	PerPage    int
	LimitSet   bool
	PerPageSet bool
	JSON       bool
}

// NewCmdList creates the list command
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List labels",
		Long: heredoc.Doc(`
			List labels in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# List labels
			$ gc label list -R owner/repo

			# List labels with pagination
			$ gc label list -R owner/repo --limit 50 --page 2

			# List labels as JSON
			$ gc label list -R owner/repo --json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			opts.LimitSet = cmd.Flags().Changed("limit")
			opts.PerPageSet = cmd.Flags().Changed("per-page")
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of labels to list")
	cmd.Flags().IntVar(&opts.Page, "page", 1, "Page number")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "API page size (default: --limit)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func listRun(opts *ListOptions) error {
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

	if opts.Limit <= 0 {
		return cmdutil.NewUsageError("--limit must be greater than 0")
	}
	if opts.Page < 0 {
		return cmdutil.NewUsageError("--page must be greater than or equal to 0")
	}
	if opts.PerPage < 0 {
		return cmdutil.NewUsageError("--per-page must be greater than or equal to 0")
	}

	// List labels
	labels, err := api.ListRepoLabels(client, owner, repo, &api.LabelListOptions{
		PerPage: resolvePerPage(opts),
		Page:    opts.Page,
	})
	if err != nil {
		return fmt.Errorf("failed to list labels: %w", err)
	}
	labels = trimLabels(labels, opts)

	// Output
	if len(labels) == 0 {
		if opts.JSON {
			return cmdutil.WriteJSON(opts.IO.Out, labels)
		}
		fmt.Fprintf(opts.IO.Out, "No labels found\n")
		return nil
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, labels)
	}

	fmt.Fprintf(opts.IO.Out, "\n")
	for _, label := range labels {
		color := label.Color
		if !strings.HasPrefix(color, "#") {
			color = "#" + color
		}
		_ = color // TODO: use color in output formatting
		fmt.Fprintf(opts.IO.Out, "%s  %s\n", cs.Bold(label.Name), label.Description)
	}
	fmt.Fprintf(opts.IO.Out, "\n")

	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

// resolvePerPage returns the API page size: the explicit --per-page when set,
// otherwise --limit. Mirrors the pr list / issue list pattern.
func resolvePerPage(opts *ListOptions) int {
	if opts.PerPageSet && opts.PerPage > 0 {
		return opts.PerPage
	}
	return opts.Limit
}

// trimLabels truncates to --limit when both --limit and --per-page were set.
func trimLabels(labels []api.Label, opts *ListOptions) []api.Label {
	if opts.PerPageSet && opts.LimitSet && len(labels) > opts.Limit {
		return labels[:opts.Limit]
	}
	return labels
}
