// Package yaml implements the actions yaml command.
package yaml

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/actions/yaml/validate"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdYaml creates the actions yaml command.
func NewCmdYaml(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yaml <command>",
		Short: "Manage workflow YAML configuration",
		Long: heredoc.Doc(`
			Work with GitCode Actions workflow YAML configuration files.
		`),
		Example: heredoc.Doc(`
			# Validate a workflow YAML file
			$ gc actions yaml validate --file .gitcode/workflows/ci.yml -R owner/repo
		`),
	}

	cmd.AddCommand(validate.NewCmdValidate(f, nil))

	return cmd
}
