// Package list implements the actions job list command.
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

// ListOptions configures the actions job list command.
type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	RunID      string

	JSON   bool
	Format string
}

// NewCmdList creates the actions job list command.
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "list <run-id>",
		Short: "List workflow jobs of a pipeline run",
		Long: heredoc.Doc(`
			List the workflow jobs of a single pipeline run.

			The run id is the workflow_run_id returned by ` + "`gc actions run list`" + `.
			Use --json for machine-readable output.
		`),
		Example: heredoc.Doc(`
			# List jobs of a pipeline run
			$ gc actions job list <run-id> -R owner/repo

			# Render as a table
			$ gc actions job list <run-id> -R owner/repo --format table

			# Output as JSON
			$ gc actions job list <run-id> -R owner/repo --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.RunID = args[0]
			if opts.RunID == "" {
				return cmdutil.NewUsageError("run id is required")
			}
			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)
	cmdutil.AddFormatFlag(cmd, &opts.Format)

	return cmd
}

func listRun(opts *ListOptions) error {
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

	resp, err := api.ListActionsRunJobs(client, owner, repo, opts.RunID)
	if err != nil {
		return fmt.Errorf("failed to list workflow jobs: %w", err)
	}
	jobs := resp.Jobs
	if jobs == nil {
		jobs = []api.WorkflowRunJob{}
	}

	if len(jobs) == 0 {
		if format == output.FormatJSON {
			return cmdutil.WriteJSON(opts.IO.Out, jobs)
		}
		fmt.Fprintf(opts.IO.Out, "No workflow jobs found\n")
		return nil
	}

	if format == output.FormatJSON {
		return cmdutil.WriteJSON(opts.IO.Out, jobs)
	}

	printer, err := output.NewWorkflowJobListPrinter(output.WorkflowJobListOptions{
		Format: format,
		Color:  opts.IO.ColorScheme(),
	})
	if err != nil {
		return err
	}
	return printer.Print(opts.IO.Out, jobs)
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
