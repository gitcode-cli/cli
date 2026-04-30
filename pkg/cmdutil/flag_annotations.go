package cmdutil

import (
	"fmt"

	"github.com/spf13/cobra"
)

const FlagEnumAnnotation = "gc.enum"

// SetFlagEnum records a stable enum set for schema/export consumers.
func SetFlagEnum(cmd *cobra.Command, name string, values ...string) {
	if err := cmd.Flags().SetAnnotation(name, FlagEnumAnnotation, values); err != nil {
		panic(fmt.Sprintf("failed to annotate flag %q: %v", name, err))
	}
}
