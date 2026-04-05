package cmdutil

import (
	"github.com/spf13/cobra"
	"gitcode.com/gitcode-cli/cli/pkg/output"
)

// AddFormatFlag adds a consistent format flag to a command.
func AddFormatFlag(cmd *cobra.Command, target *string) {
	cmd.Flags().StringVar(target, "format", "simple", "Output format (table/json/simple)")
}

// ParseFormat parses the format string to output.Format.
func ParseFormat(s string) output.Format {
	switch s {
	case "table":
		return output.FormatTable
	case "json":
		return output.FormatJSON
	case "simple":
		return output.FormatSimple
	default:
		return output.FormatSimple
	}
}
