// Package run implements the actions workflow run command.
package run

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

// RunOptions configures the actions workflow run command.
type RunOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	WorkflowID string
	Ref        string
	RawFields  []string
	Fields     map[string]string
	JSON       bool
}

// NewCmdRun creates the actions workflow run command.
func NewCmdRun(f *cmdutil.Factory, runF func(*RunOptions) error) *cobra.Command {
	opts := &RunOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "run [<workflow-id> | <filename>]",
		Short: "Trigger a workflow run",
		Long: heredoc.Doc(`
			Trigger a workflow_dispatch run for a GitCode Actions workflow.

			The workflow_id can be the ID returned by 'gc actions workflow list'
			or a filename with .yml/.yaml extension (e.g. ci.yml).
		`),
		Example: heredoc.Doc(`
			# Trigger a workflow by ID
			$ gc actions workflow run wf-1 -R owner/repo --ref main

			# Trigger by filename
			$ gc actions workflow run ci.yml -R owner/repo --ref main

			# With custom inputs
			$ gc actions workflow run ci.yml -R owner/repo --ref main -f key1=value1 -f key2=value2

			# JSON output (returns workflow_run_id)
			$ gc actions workflow run ci.yml -R owner/repo --ref main --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.WorkflowID = strings.TrimSpace(args[0])
			if opts.WorkflowID == "" {
				return cmdutil.NewUsageError("workflow id is required")
			}
			opts.Fields = make(map[string]string)
			for _, f := range opts.RawFields {
				parts := strings.SplitN(f, "=", 2)
				if len(parts) != 2 || parts[0] == "" {
					return cmdutil.NewUsageError(fmt.Sprintf("invalid field %q, expected key=value", f))
				}
				if err := cmdutil.ScanContentForSecrets(parts[1]); err != nil {
					return err
				}
				opts.Fields[parts[0]] = parts[1]
			}
			if runF != nil {
				return runF(opts)
			}
			return runRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVar(&opts.Ref, "ref", "", "Branch or tag name (required)")
	cmd.Flags().StringArrayVarP(&opts.RawFields, "field", "f", nil, "Add string input in key=value format")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func runRun(opts *RunOptions) error {
	if opts.Ref == "" {
		return cmdutil.NewUsageError("--ref is required")
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

	req := &api.DispatchWorkflowRequest{
		Ref:    opts.Ref,
		Inputs: opts.Fields,
	}

	result, err := api.DispatchActionsWorkflow(client, owner, repo, opts.WorkflowID, req)
	if err != nil {
		return fmt.Errorf("failed to trigger workflow: %w", err)
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, result)
	}

	fmt.Fprintf(opts.IO.ErrOut, "Triggered workflow %s on %s. Run ID: %s\n",
		result.WorkflowID, opts.Ref, result.WorkflowRunID)
	return nil
}
