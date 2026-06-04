// Package precommit implements the precommit command group.
package precommit

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/precommit/check"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdPrecommit creates the precommit command group.
func NewCmdPrecommit(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "precommit <command>",
		Short: "Work with pre-commit configuration and hooks",
		Long: heredoc.Doc(`
			Detect and verify pre-commit configuration and the local pre-commit
			environment before committing code.
		`),
		Example: heredoc.Doc(`
			# Verify the environment is ready
			$ gc precommit check

			# Verify and run the hooks
			$ gc precommit check --run
		`),
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "precommit",
		},
	}

	cmd.AddCommand(check.NewCmdCheck(f, nil))

	return cmd
}
