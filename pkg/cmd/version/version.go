// Package version implements the version command
package version

import (
	"fmt"

	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// VersionInfo represents the JSON output for version command
type VersionInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit,omitempty"`
	Built   string `json:"built,omitempty"`
	URL     string `json:"url,omitempty"`
}

// NewCmdVersion creates the version command
func NewCmdVersion(ver, commit, date string) *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print gc version",
		Long:  "Print the version information of gc.",
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()

			if jsonOutput {
				info := VersionInfo{
					Version: ver,
					Commit:  commit,
					Built:   date,
					URL:     "https://gitcode.com/gitcode-cli/cli",
				}
				cmdutil.WriteJSON(out, info)
				return
			}

			fmt.Fprintf(out, "gc version %s\n", ver)
			fmt.Fprintf(out, "  commit: %s\n", commit)
			fmt.Fprintf(out, "  built:  %s\n", date)
			fmt.Fprintln(out, "https://gitcode.com/gitcode-cli/cli")
		},
	}

	cmdutil.AddJSONFlag(cmd, &jsonOutput)

	return cmd
}
