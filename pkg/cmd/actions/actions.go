// Package actions implements the actions command.
package actions

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	artifactcmd "gitcode.com/gitcode-cli/cli/pkg/cmd/actions/artifact"
	jobcmd "gitcode.com/gitcode-cli/cli/pkg/cmd/actions/job"
	runcmd "gitcode.com/gitcode-cli/cli/pkg/cmd/actions/run"
	runnergroupcmd "gitcode.com/gitcode-cli/cli/pkg/cmd/actions/runner-group"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdActions creates the actions command.
func NewCmdActions(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "actions <command>",
		Short: "Manage GitCode Actions (pipeline runs and workflow jobs)",
		Long: heredoc.Doc(`
			Work with GitCode Actions: inspect pipeline runs, workflow jobs and
			job logs.

			GitCode Actions exposes pipeline run records and workflow jobs through
			the Actions v8 API. This command group provides read-only inspection
			of CI run status.
		`),
		Example: heredoc.Doc(`
			# List recent pipeline runs
			$ gc actions run list -R owner/repo

			# Filter runs by status
			$ gc actions run list -R owner/repo --status FAILED

			# Output runs as JSON
			$ gc actions run list -R owner/repo --json

			# List jobs of a pipeline run
			$ gc actions job list <run-id> -R owner/repo
		`),
		Annotations: map[string]string{
			"IsCore": "true",
		},
	}

	cmd.AddCommand(runcmd.NewCmdRun(f))
	cmd.AddCommand(jobcmd.NewCmdJob(f))
	cmd.AddCommand(artifactcmd.NewCmdArtifact(f))
	cmd.AddCommand(runnergroupcmd.NewCmdRunnerGroup(f))

	return cmd
}
