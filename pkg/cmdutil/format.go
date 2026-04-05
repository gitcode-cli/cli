package cmdutil

import (
	"github.com/spf13/cobra"
	"gitcode.com/gitcode-cli/cli/pkg/output"
)

// AddFormatFlag adds a consistent format flag to a command.
func AddFormatFlag(cmd *cobra.Command, target *string) {
	cmd.Flags().StringVar(target, "format", "simple", "Output format (table/json/simple)")
}

// AddTimeFormatFlag adds a time format flag to a command.
func AddTimeFormatFlag(cmd *cobra.Command, target *string) {
	cmd.Flags().StringVar(target, "time-format", "absolute", "Time format (relative/absolute)")
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

// ParseTimeFormat parses the time format string.
func ParseTimeFormat(s string) output.TimeFormat {
	switch s {
	case "relative":
		return output.TimeFormatRelative
	case "absolute":
		return output.TimeFormatAbsolute
	default:
		return output.TimeFormatAbsolute
	}
}