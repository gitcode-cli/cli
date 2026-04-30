package cmdutil

import "github.com/spf13/cobra"

const FlagEnumAnnotation = "gc.enum"

// SetFlagEnum records a stable enum set for schema/export consumers.
func SetFlagEnum(cmd *cobra.Command, name string, values ...string) {
	_ = cmd.Flags().SetAnnotation(name, FlagEnumAnnotation, values)
}
