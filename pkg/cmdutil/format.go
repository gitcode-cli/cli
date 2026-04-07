package cmdutil

import "github.com/spf13/cobra"

// AddFormatFlag adds a consistent format flag to a command.
func AddFormatFlag(cmd *cobra.Command, target *string) {
	cmd.Flags().StringVar(target, "format", "", "Output format (json/simple/table)")
}

// AddTimeFormatFlag adds a consistent time-format flag to a command.
func AddTimeFormatFlag(cmd *cobra.Command, target *string) {
	cmd.Flags().StringVar(target, "time-format", "", "Time format (absolute/relative)")
}

// AddTemplateFlag adds a consistent template flag to a command.
func AddTemplateFlag(cmd *cobra.Command, target *string) {
	cmd.Flags().StringVar(target, "template", "", "Format output using a Go template")
}
