// Package view implements the actions run view command.
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

// ViewOptions configures the actions run view command.
type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	RunID      string

	JSON bool
}

// NewCmdView creates the actions run view command.
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "view <run-id>",
		Short: "View a pipeline run",
		Long: heredoc.Doc(`
			View the detail of a single pipeline (workflow) run, including its
			stages and jobs.

			The run id is the workflow_run_id returned by ` + "`gc actions run list`" + `.
			Use --json for a faithful, machine-readable copy of the API response.
		`),
		Example: heredoc.Doc(`
			# View a pipeline run
			$ gc actions run view <run-id> -R owner/repo

			# Faithful JSON output (all API fields preserved)
			$ gc actions run view <run-id> -R owner/repo --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.RunID = strings.TrimSpace(args[0])
			if opts.RunID == "" {
				return cmdutil.NewUsageError("run id is required")
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

	detail, raw, err := api.GetActionsRun(client, owner, repo, opts.RunID)
	if err != nil {
		return fmt.Errorf("failed to get pipeline run: %w", err)
	}

	if opts.JSON {
		// Output the raw API response verbatim for full fidelity (preserves
		// deep stage/job/step execution fields that the typed struct elides).
		if _, err := opts.IO.Out.Write(raw); err != nil {
			return fmt.Errorf("failed to write JSON output: %w", err)
		}
		_, err := fmt.Fprintln(opts.IO.Out)
		return err
	}

	return printDetail(opts, detail)
}

func printDetail(opts *ViewOptions, d *api.WorkflowRunDetail) error {
	cs := opts.IO.ColorScheme()
	out := opts.IO.Out

	fmt.Fprintf(out, "%s #%d  %s\n", d.WorkflowName, d.RunNumber, statusLabel(cs, d.Status))
	fmt.Fprintf(out, "  run id:            %s\n", orDash(d.WorkflowRunID))
	fmt.Fprintf(out, "  workflow id:       %s\n", orDash(d.WorkflowID))
	fmt.Fprintf(out, "  title:             %s\n", orDash(d.Title))
	fmt.Fprintf(out, "  file path:         %s\n", orDash(d.FilePath))
	fmt.Fprintf(out, "  status:            %s\n", statusLabel(cs, d.Status))
	fmt.Fprintf(out, "  event:             %s\n", orDash(d.Event))
	fmt.Fprintf(out, "  branch:            %s\n", orDash(d.HeadBranch))
	fmt.Fprintf(out, "  head sha:          %s\n", orDash(d.HeadSHA))
	fmt.Fprintf(out, "  actor:             %s\n", actorText(d.Actor))
	fmt.Fprintf(out, "  in default branch: %v\n", d.ExistInDefaultBranch)
	fmt.Fprintf(out, "  started:           %s\n", formatTime(d.StartTime))
	fmt.Fprintf(out, "  ended:             %s\n", formatTime(d.EndTime))
	if d.PauseTime > 0 {
		fmt.Fprintf(out, "  paused:            %s\n", formatTime(d.PauseTime))
	}

	if len(d.Stages) == 0 {
		return nil
	}
	fmt.Fprintf(out, "\nStages:\n")
	for _, stage := range d.Stages {
		fmt.Fprintf(out, "  - %s  %s  jobs: %d  started: %s  ended: %s\n",
			orDash(stage.Name), statusLabel(cs, stage.Status), len(stage.Jobs),
			formatTime(stage.StartTime), formatTime(stage.EndTime))
		for _, job := range stage.Jobs {
			fmt.Fprintf(out, "      - %s  %s  steps: %d  started: %s  ended: %s\n",
				orDash(job.Name), statusLabel(cs, job.Status), len(job.Steps),
				formatTime(job.StartTime), formatTime(job.EndTime))
		}
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

func actorText(actor *api.ActionsActor) string {
	if actor == nil {
		return "-"
	}
	if actor.Name != "" && actor.Login != "" {
		return fmt.Sprintf("%s (%s)", actor.Login, actor.Name)
	}
	if actor.Login != "" {
		return actor.Login
	}
	return orDash(actor.Name)
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
	return time.Unix(t, 0).UTC().Format(time.RFC3339)
}
