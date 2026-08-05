// Package validate implements the actions yaml validate command.
package validate

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// ValidateOptions configures the actions yaml validate command.
type ValidateOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	File       string

	JSON bool
}

// NewCmdValidate creates the actions yaml validate command.
func NewCmdValidate(f *cmdutil.Factory, runF func(*ValidateOptions) error) *cobra.Command {
	opts := &ValidateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate workflow YAML syntax",
		Long: heredoc.Doc(`
			Validate a GitCode Actions workflow YAML file against the official
			GitCode API.

			The YAML content is base64-encoded and sent to the Actions v8
			validate endpoint. Use --file - to read the YAML from stdin.
		`),
		Example: heredoc.Doc(`
			# Validate a workflow file
			$ gc actions yaml validate --file .gitcode/workflows/ci.yml -R owner/repo

			# Validate from stdin
			$ cat ci.yml | gc actions yaml validate --file - -R owner/repo

			# Machine-readable output
			$ gc actions yaml validate --file ci.yml -R owner/repo --json
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.File == "" {
				return cmdutil.NewUsageError("--file is required (use - for stdin)")
			}
			if runF != nil {
				return runF(opts)
			}
			return validateRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVar(&opts.File, "file", "", "Workflow YAML file to validate (use - for stdin)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func validateRun(opts *ValidateOptions) error {
	content, err := readYAML(opts.File, opts.IO.In)
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

	result, raw, err := api.ValidateActionsWorkflow(client, owner, repo, content)
	if err != nil {
		return fmt.Errorf("failed to validate workflow YAML: %w", err)
	}

	if opts.JSON {
		if _, err := opts.IO.Out.Write(raw); err != nil {
			return fmt.Errorf("failed to write JSON output: %w", err)
		}
		if _, err := fmt.Fprintln(opts.IO.Out); err != nil {
			return err
		}
		if !result.Valid {
			return cmdutil.NewCLIError(cmdutil.ExitError, "workflow YAML is invalid", nil)
		}
		return nil
	}

	printResult(opts.IO, result)
	if !result.Valid {
		return cmdutil.NewCLIError(cmdutil.ExitError, "workflow YAML is invalid", nil)
	}
	return nil
}

func readYAML(file string, stdin io.Reader) ([]byte, error) {
	var content []byte
	var err error
	if file == "-" {
		content, err = io.ReadAll(stdin)
		if err != nil {
			return nil, fmt.Errorf("failed to read workflow YAML from stdin: %w", err)
		}
	} else {
		content, err = os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read workflow YAML file %s: %w", file, err)
		}
	}
	return []byte(cmdutil.DecodeUserText(content)), nil
}

func printResult(io *iostreams.IOStreams, result *api.ValidateWorkflowResponse) {
	cs := io.ColorScheme()
	if result.Valid {
		fmt.Fprintf(io.Out, "%s Workflow YAML is valid\n", cs.Green("✓"))
		return
	}

	fmt.Fprintf(io.Out, "%s Workflow YAML is invalid\n", cs.Red("✗"))
	for _, d := range result.Diagnostics {
		loc := fmt.Sprintf("line %d, column %d", d.Range.Start.Line, d.Range.Start.Column)
		if d.Range.End.Line != 0 || d.Range.End.Column != 0 {
			loc += fmt.Sprintf(" - line %d, column %d", d.Range.End.Line, d.Range.End.Column)
		}
		severity := d.Severity
		if severity == "" {
			severity = "error"
		}
		fmt.Fprintf(io.Out, "  - %s [%s] %s\n", loc, severity, d.Message)
	}
}
