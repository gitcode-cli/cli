// Package view implements the actions job view command.
package view

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// msTimestampThreshold separates second from millisecond timestamps.
// Values at or above it are treated as milliseconds and divided by 1000
// (the Actions v8 API returns milliseconds). Mirrors run/view.formatTime.
const msTimestampThreshold = 100_000_000_000 // 1e11

// ViewOptions configures the actions job view command.
type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	RunID      string
	JobID      string

	JSON bool
}

// NewCmdView creates the actions job view command.
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "view <run-id> <job-id>",
		Short: "View a workflow job",
		Long: heredoc.Doc(`
			View the detail of a single workflow job, including its steps.

			Both the run id (workflow_run_id from ` + "`gc actions run list`" + `) and the
			job id (from ` + "`gc actions job list`" + `) are required, matching the API
			path /actions/runs/{run_id}/jobs/{job_id}. Use --json for a faithful,
			machine-readable copy of the API response.
		`),
		Example: heredoc.Doc(`
			# View a workflow job
			$ gc actions job view <run-id> <job-id> -R owner/repo

			# Faithful JSON output (all API fields preserved)
			$ gc actions job view <run-id> <job-id> -R owner/repo --json
		`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.RunID = strings.TrimSpace(args[0])
			opts.JobID = strings.TrimSpace(args[1])
			if opts.RunID == "" || opts.JobID == "" {
				return cmdutil.NewUsageError("run id and job id are required")
			}
			if runF != nil {
				return runF(opts)
			}
			return viewRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func viewRun(opts *ViewOptions) error {
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

	job, raw, err := api.GetActionsJob(client, owner, repo, opts.RunID, opts.JobID)
	if err != nil {
		return fmt.Errorf("failed to get workflow job: %w", err)
	}

	if opts.JSON {
		// Output the raw API response verbatim for full fidelity (preserves
		// deep step execution fields that the typed struct elides).
		if _, err := opts.IO.Out.Write(raw); err != nil {
			return fmt.Errorf("failed to write JSON output: %w", err)
		}
		_, err := fmt.Fprintln(opts.IO.Out)
		return err
	}

	return printJob(opts, job)
}

func printJob(opts *ViewOptions, j *api.WorkflowRunJob) error {
	cs := opts.IO.ColorScheme()
	out := opts.IO.Out

	fmt.Fprintf(out, "%s  %s\n", orDash(jobName(j)), statusLabel(cs, j.Status))
	fmt.Fprintf(out, "  job id:          %s\n", orDash(j.ID))
	fmt.Fprintf(out, "  identifier:      %s\n", orDash(j.Identifier))
	fmt.Fprintf(out, "  name:            %s\n", orDash(j.Name))
	fmt.Fprintf(out, "  status:          %s\n", statusLabel(cs, j.Status))
	fmt.Fprintf(out, "  sequence:        %d\n", j.Sequence)
	fmt.Fprintf(out, "  job type:        %s\n", orDash(j.JobType))
	fmt.Fprintf(out, "  resource:        %s\n", orDash(j.Resource))
	fmt.Fprintf(out, "  exec id:         %s\n", orDash(j.ExecID))
	fmt.Fprintf(out, "  last dispatch:   %s\n", orDash(j.LastDispatchID))
	if j.Condition != "" {
		fmt.Fprintf(out, "  condition:       %s\n", j.Condition)
	}
	if len(j.DependsOn) > 0 {
		fmt.Fprintf(out, "  depends on:      %s\n", strings.Join(j.DependsOn, ", "))
	}
	fmt.Fprintf(out, "  started:         %s\n", formatTime(j.StartTime))
	fmt.Fprintf(out, "  ended:           %s\n", formatTime(j.EndTime))
	// execute_cost_time's unit is undocumented (real data shows ms for a
	// ~77s job); start/end already convey duration, so it is omitted from
	// the human view and preserved raw in --json.

	if len(j.Steps) == 0 {
		return nil
	}
	fmt.Fprintf(out, "\nSteps:\n")
	for _, s := range j.Steps {
		fmt.Fprintf(out, "  - %s  %s  task: %s  started: %s  ended: %s\n",
			orDash(stepName(s)), statusLabel(cs, s.Status), orDash(s.Task),
			formatTime(s.StartTime), formatTime(s.EndTime))
	}
	return nil
}

func statusLabel(cs *iostreams.ColorScheme, status string) string {
	if cs == nil || status == "" {
		return status
	}
	switch status {
	case "COMPLETED":
		return cs.Green(status)
	case "FAILED":
		return cs.Red(status)
	case "RUNNING":
		return cs.Yellow(status)
	case "CANCELED", "IGNORED", "PAUSED", "SUSPEND":
		return cs.Gray(status)
	default:
		return status
	}
}

func jobName(j *api.WorkflowRunJob) string {
	if j.Name != "" {
		return j.Name
	}
	return j.Identifier
}

func stepName(s api.WorkflowRunStep) string {
	if s.Name != "" {
		return s.Name
	}
	return s.Identifier
}

func orDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}

func formatTime(t int64) string {
	if t <= 0 {
		return "-"
	}
	secs := t
	if t >= msTimestampThreshold {
		secs = t / 1000
	}
	return time.Unix(secs, 0).UTC().Format(time.RFC3339)
}
