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
			fmt.Printf("gc version %s\n", ver)
			fmt.Printf("  commit: %s\n", commit)
			fmt.Printf("  built:  %s\n", date)
			fmt.Println("https://gitcode.com/gitcode-com/gitcode-cli")
		},
	}

	return cmd
}