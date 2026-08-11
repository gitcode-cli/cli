// Package list implements the repo branch list command.
package list

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/output"
)

// ListOptions configures the repo branch list command.
type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	Limit      int
	Page       int
	JSON       bool
	Format     string
}

// NewCmdList creates the repo branch list command.
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List branches in a repository",
		Long: heredoc.Doc(`
			List all branches in a GitCode repository.

			Use --json for machine-readable output.
		`),
		Example: heredoc.Doc(`
			# List branches
			$ gc repo branch list -R owner/repo

			# JSON output
			$ gc repo branch list -R owner/repo --json
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of branches to list")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)
	cmdutil.AddFormatFlag(cmd, &opts.Format)

	return cmd
}

func listRun(opts *ListOptions) error {
	if opts.Limit <= 0 {
		return cmdutil.NewUsageError("--limit must be greater than 0")
	}
	if opts.Page < 0 {
		return cmdutil.NewUsageError("--page must be greater than or equal to 0")
	}

	format, err := resolveOutputFormat(opts.JSON, opts.Format)
	if err != nil {
		return err
	}

	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
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

	branches, err := api.ListBranches(client, owner, repo, &api.BranchListOptions{
		PerPage: opts.Limit,
		Page:    opts.Page,
	})
	if err != nil {
		return fmt.Errorf("failed to list branches: %w", err)
	}

	if len(branches) == 0 {
		if format == output.FormatJSON {
			return cmdutil.WriteJSON(opts.IO.Out, branches)
		}
		fmt.Fprintln(opts.IO.Out, "No branches found.")
		return nil
	}

	if format == output.FormatJSON {
		return cmdutil.WriteJSON(opts.IO.Out, branches)
	}

	cs := opts.IO.ColorScheme()
	for _, b := range branches {
		marker := ""
		if b.DefaultBranch {
			marker = cs.Green(" (default)")
		}
		fmt.Fprintf(opts.IO.Out, "%s%s\n", b.Name, marker)
	}
	return nil
}

func resolveOutputFormat(jsonFlag bool, raw string) (output.Format, error) {
	format, err := output.ParseFormat(raw)
	if err != nil {
		return "", cmdutil.NewUsageError(err.Error())
	}
	if jsonFlag {
		if raw != "" && format != output.FormatJSON {
			return "", cmdutil.NewUsageError("--json cannot be combined with --format unless --format json")
		}
		return output.FormatJSON, nil
	}
	return format, nil
}
