package cmdutil

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteJSONNilSliceEmitsEmptyArray(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  string
	}{
		{
			name:  "nil int slice emits []",
			value: func() interface{} { var s []int; return s }(),
			want:  "[]",
		},
		{
			name:  "nil struct slice emits []",
			value: func() interface{} { var s []struct{ N int }; return s }(),
			want:  "[]",
		},
		{
			name:  "empty int slice emits []",
			value: []int{},
			want:  "[]",
		},
		{
			name:  "non-empty slice emits array",
			value: []int{1, 2, 3},
			want:  "[\n  1,\n  2,\n  3\n]",
		},
		{
			name:  "struct emits object unchanged",
			value: struct{ A int }{A: 1},
			want:  "{\n  \"A\": 1\n}",
		},
		{
			name:  "nil interface still emits null",
			value: nil,
			want:  "null",
		},
		{
			// normalizeNilSlice only inspects a top-level slice; a pointer to a
			// nil slice has Kind==Ptr and is returned unchanged, so the nil
			// slice it points to still marshals as `null`. This locks the
			// intended scope boundary - no current call site passes *[]T.
			name:  "pointer to nil slice still emits null (scope boundary)",
			value: func() interface{} { var s []int; return &s }(),
			want:  "null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := WriteJSON(&buf, tt.value); err != nil {
				t.Fatalf("WriteJSON() error = %v", err)
			}
			got := strings.TrimSpace(buf.String())
			if got != tt.want {
				t.Errorf("WriteJSON(%v) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}
