// Package version_test tests the version command
package version

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewCmdVersion(t *testing.T) {
	tests := []struct {
		name         string
		version      string
		commit       string
		date         string
		wantContains []string
	}{
		{
			name:         "with valid version info",
			version:      "v0.2.8",
			commit:       "abc1234",
			date:         "2026-03-24",
			wantContains: []string{"gc version v0.2.8", "commit: abc1234", "built:  2026-03-24", "https://gitcode.com/gitcode-cli/cli"},
		},
		{
			name:         "with dev version",
			version:      "dev",
			commit:       "none",
			date:         "unknown",
			wantContains: []string{"gc version dev", "commit: none", "built:  unknown"},
		},
		{
			name:         "with empty version info",
			version:      "",
			commit:       "",
			date:         "",
			wantContains: []string{"gc version", "commit:", "built:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCmdVersion(tt.version, tt.commit, tt.date)

			if cmd.Use != "version" {
				t.Errorf("NewCmdVersion().Use = %q, want %q", cmd.Use, "version")
			}

			if cmd.Short != "Print gc version" {
				t.Errorf("NewCmdVersion().Short = %q, want %q", cmd.Short, "Print gc version")
			}

			// Capture output
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.Execute()

			output := buf.String()
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("version output missing %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestVersionOutput(t *testing.T) {
	tests := []struct {
		name    string
		version string
		commit  string
		date    string
	}{
		{"release version", "v1.0.0", "a1b2c3d", "2026-01-15"},
		{"dev version", "dev", "none", "unknown"},
		{"snapshot version", "v1.0.0-dirty", "xyz789", "2026-03-24"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCmdVersion(tt.version, tt.commit, tt.date)

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			err := cmd.Execute()
			if err != nil {
				t.Errorf("cmd.Execute() error = %v", err)
			}

			output := buf.String()

			// Check that output contains expected parts
			if !strings.Contains(output, "gc version") {
				t.Error("output missing 'gc version'")
			}
			if !strings.Contains(output, "commit:") {
				t.Error("output missing 'commit:'")
			}
			if !strings.Contains(output, "built:") {
				t.Error("output missing 'built:'")
			}
			if !strings.Contains(output, "https://gitcode.com/gitcode-cli/cli") {
				t.Error("output missing project URL")
			}
		})
	}
}