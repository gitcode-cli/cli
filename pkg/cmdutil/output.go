package cmdutil

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// AddJSONFlag adds a consistent JSON output flag to a command.
func AddJSONFlag(cmd *cobra.Command, target *bool) {
	cmd.Flags().BoolVar(target, "json", false, "Output as JSON")
}

// WriteJSON writes indented JSON to the target writer.
func WriteJSON(w io.Writer, value interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		return fmt.Errorf("failed to encode JSON output: %w", err)
	}
	return nil
}
