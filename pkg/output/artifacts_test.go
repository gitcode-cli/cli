package output

import (
	"strings"
	"testing"
)

func TestSizeLabel(t *testing.T) {
	tests := []struct {
		name  string
		bytes int64
		want  string
	}{
		{name: "zero", bytes: 0, want: "0 B"},
		{name: "bytes", bytes: 512, want: "512 B"},
		{name: "kib boundary", bytes: 1024, want: "1.0 KiB"},
		{name: "kib", bytes: 1536, want: "1.5 KiB"},
		{name: "mib boundary", bytes: 1048576, want: "1.0 MiB"},
		{name: "mib", bytes: 1572864, want: "1.5 MiB"},
		{name: "gib boundary", bytes: 1073741824, want: "1.0 GiB"},
		{name: "gib", bytes: 1610612736, want: "1.5 GiB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sizeLabel(tt.bytes); got != tt.want {
				t.Errorf("sizeLabel(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestFormatMsTimeString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: "-"},
		{name: "zero", in: "0", want: "-"},
		{name: "negative", in: "-1", want: "-"},
		{name: "invalid", in: "not-a-number", want: "-"},
		{name: "seconds", in: "1783500745", want: "2026-"},
		{name: "milliseconds", in: "1783500745000", want: "2026-"},
		{name: "ms equals sec", in: "same", want: "skip"}, // special: see below
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "ms equals sec" {
				// 1783500745 s and 1783500745000 ms must render identically (same instant).
				gotSec := formatMsTimeString("1783500745")
				gotMs := formatMsTimeString("1783500745000")
				if gotSec != gotMs {
					t.Fatalf("formatMsTimeString(sec)=%q != formatMsTimeString(ms)=%q", gotSec, gotMs)
				}
				if !strings.HasPrefix(gotSec, "2026-") {
					t.Fatalf("formatMsTimeString = %q, want 2026- prefix", gotSec)
				}
				return
			}
			got := formatMsTimeString(tt.in)
			if tt.want == "-" {
				if got != "-" {
					t.Errorf("formatMsTimeString(%q) = %q, want -", tt.in, got)
				}
				return
			}
			if !strings.HasPrefix(got, tt.want) {
				t.Errorf("formatMsTimeString(%q) = %q, want prefix %q", tt.in, got, tt.want)
			}
		})
	}
}
