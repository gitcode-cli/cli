// Package version implements the version command
package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCmdVersion creates the version command
func NewCmdVersion(ver, commit, date string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print gc version",
		Long:  "Print the version information of gc.",
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "gc version %s\n", ver)
			fmt.Fprintf(out, "  commit: %s\n", commit)
			fmt.Fprintf(out, "  built:  %s\n", date)
			fmt.Fprintln(out, "https://gitcode.com/gitcode-cli/cli")
		},
	}

	return cmd
}