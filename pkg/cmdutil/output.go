package cmdutil

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/spf13/cobra"
)

// AddJSONFlag adds a consistent JSON output flag to a command.
func AddJSONFlag(cmd *cobra.Command, target *bool) {
	cmd.Flags().BoolVar(target, "json", false, "Output as JSON")
}

// WriteJSON writes indented JSON to the target writer.
//
// A nil slice is normalized to an empty slice so list commands emit `[]`
// instead of `null` when there are no results, keeping --json output stable
// and consumable by scripts and agents (see spec/foundations/agent-friendly-cli.md).
//
// Only the top-level value is normalized. A nil value with no concrete type, a
// pointer to a nil slice (e.g. *[]T), and nil slice fields nested inside a
// struct are left untouched and still encode as `null`; commands that rely on
// such fields (e.g. pr view --json emitting "comments": null) are unaffected.
func WriteJSON(w io.Writer, value interface{}) error {
	value = normalizeNilSlice(value)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		return fmt.Errorf("failed to encode JSON output: %w", err)
	}
	return nil
}

// normalizeNilSlice returns a non-nil empty slice of the same type when value is
// a top-level nil slice, so it marshals as `[]` instead of `null`. Pointers to
// slices and nil slice fields nested in structs are returned unchanged (they
// still encode as `null`); see WriteJSON for the rationale.
func normalizeNilSlice(value interface{}) interface{} {
	if value == nil {
		return value
	}
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Slice && rv.IsNil() {
		return reflect.MakeSlice(rv.Type(), 0, 0).Interface()
	}
	return value
}
