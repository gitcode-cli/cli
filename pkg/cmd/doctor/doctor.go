// Package doctor implements local GitCode CLI diagnostics.
package doctor

import (
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/doctor/install"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdDoctor creates the doctor command.
func NewCmdDoctor(f *cmdutil.Factory, version, commit, built string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor <command>",
		Short: "Diagnose the local GitCode CLI installation",
		Long:  "Inspect the local GitCode CLI installation without authentication or network access.",
	}
	cmd.AddCommand(install.NewCmdInstall(f, version, commit, built))
	return cmd
}
