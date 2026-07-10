// Package log implements the actions job log command.
package log

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// LogOptions configures the actions job log command.
type LogOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	RunID      string
	JobID      string
	Output     string
}

// NewCmdLog creates the actions job log command.
func NewCmdLog(f *cmdutil.Factory, runF func(*LogOptions) error) *cobra.Command {
	opts := &LogOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "log <run-id> <job-id>",
		Short: "Download a workflow job log",
		Long: heredoc.Doc(`
			Download the log of a single workflow job.

			Both the run id (workflow_run_id from ` + "`gc actions run list`" + `) and the
			job id (from ` + "`gc actions job list`" + `) are required, matching the API
			path /actions/runs/{run_id}/jobs/{job_id}/download_log. The endpoint
			returns a ZIP archive of the job's step logs (binary, not plain text),
			so prefer --output to save and unzip; writing to stdout streams the raw
			archive bytes (use a redirect, e.g. > job-log.zip).
		`),
		Example: heredoc.Doc(`
			# Save the job log archive and unzip it
			$ gc actions job log <run-id> <job-id> -R owner/repo --output job-log.zip
			$ unzip job-log.zip

			# Stream the raw archive to a file via redirect
			$ gc actions job log <run-id> <job-id> -R owner/repo > job-log.zip
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
			return logRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Output, "output", "o", "", "Write log to FILE (default: stdout)")

	return cmd
}

func logRun(opts *LogOptions) error {
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

	// The download_log endpoint returns a binary ZIP archive; refuse to dump
	// it to an interactive terminal (would garble the screen). Piped/redirected
	// stdout (non-TTY) or --output are the supported targets.
	if strings.TrimSpace(opts.Output) == "" && opts.IO.IsStdoutTTY() {
		return cmdutil.NewUsageError("job log is a binary archive; use --output FILE or redirect stdout (e.g. `> job-log.zip`)")
	}

	logBytes, err := api.GetActionsJobLog(client, owner, repo, opts.RunID, opts.JobID)
	if err != nil {
		return fmt.Errorf("failed to download job log: %w", err)
	}

	if strings.TrimSpace(opts.Output) != "" {
		if err := os.WriteFile(opts.Output, logBytes, 0o644); err != nil {
			return fmt.Errorf("failed to write log file: %w", err)
		}
		fmt.Fprintf(opts.IO.ErrOut, "Saved log to %s (%d bytes)\n", opts.Output, len(logBytes))
		return nil
	}

	if _, err := opts.IO.Out.Write(logBytes); err != nil {
		return fmt.Errorf("failed to write log output: %w", err)
	}
	return nil
}
