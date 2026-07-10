package output

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// msTimestampThreshold separates second from millisecond timestamps.
// The Actions v8 API returns millisecond timestamps as strings; values at or
// above this threshold (as parsed int64) are treated as ms and divided by 1000.
const msTimestampThreshold = 100_000_000_000 // 1e11

// ArtifactListOptions configures artifact list text output.
type ArtifactListOptions struct {
	Format Format
	Color  *iostreams.ColorScheme
}

// ArtifactListPrinter renders artifact lists.
type ArtifactListPrinter struct {
	opts ArtifactListOptions
}

// NewArtifactListPrinter validates and returns an artifact printer.
func NewArtifactListPrinter(opts ArtifactListOptions) (*ArtifactListPrinter, error) {
	if opts.Format == "" {
		opts.Format = FormatSimple
	}
	return &ArtifactListPrinter{opts: opts}, nil
}

// Print renders artifacts according to the configured format.
func (p *ArtifactListPrinter) Print(w io.Writer, artifacts []api.Artifact) error {
	switch p.opts.Format {
	case FormatTable:
		return p.printTable(w, artifacts)
	default:
		return p.printSimple(w, artifacts)
	}
}

func (p *ArtifactListPrinter) printSimple(w io.Writer, artifacts []api.Artifact) error {
	maxNameWidth := len("NAME")
	maxSizeWidth := len("SIZE")
	for _, a := range artifacts {
		if width := len(a.Name); width > maxNameWidth {
			maxNameWidth = width
		}
		if width := len(sizeLabel(a.SizeBytes)); width > maxSizeWidth {
			maxSizeWidth = width
		}
	}
	for _, a := range artifacts {
		fmt.Fprintf(w, "%-*s  %-*s  %s\n", maxNameWidth, a.Name, maxSizeWidth, sizeLabel(a.SizeBytes), formatMsTimeString(a.CreatedAt))
	}
	return nil
}

func (p *ArtifactListPrinter) printTable(w io.Writer, artifacts []api.Artifact) error {
	maxNameWidth := len("NAME")
	maxIDWidth := len("ID")
	maxSizeWidth := len("SIZE")
	for _, a := range artifacts {
		if width := len(a.Name); width > maxNameWidth {
			maxNameWidth = width
		}
		if width := len(a.ID); width > maxIDWidth {
			maxIDWidth = width
		}
		if width := len(sizeLabel(a.SizeBytes)); width > maxSizeWidth {
			maxSizeWidth = width
		}
	}
	fmt.Fprintf(w, "%-*s  %-*s  %-*s  %-20s  %s\n", maxNameWidth, "NAME", maxIDWidth, "ID", maxSizeWidth, "SIZE", "CREATED", "EXPIRES")
	for _, a := range artifacts {
		fmt.Fprintf(w, "%-*s  %-*s  %-*s  %-20s  %s\n", maxNameWidth, a.Name, maxIDWidth, a.ID, maxSizeWidth, sizeLabel(a.SizeBytes), formatMsTimeString(a.CreatedAt), formatMsTimeString(a.ExpiresAt))
	}
	return nil
}

// sizeLabel formats a byte count as a human-readable string.
func sizeLabel(bytes int64) string {
	const (
		kiB = 1024
		miB = 1024 * 1024
		giB = 1024 * 1024 * 1024
	)
	switch {
	case bytes >= giB:
		return fmt.Sprintf("%.1f GiB", float64(bytes)/float64(giB))
	case bytes >= miB:
		return fmt.Sprintf("%.1f MiB", float64(bytes)/float64(miB))
	case bytes >= kiB:
		return fmt.Sprintf("%.1f KiB", float64(bytes)/float64(kiB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// formatMsTimeString parses a string timestamp (ms, as returned by the Actions
// v8 API) and formats it as RFC3339; returns "-" for empty/invalid values.
func formatMsTimeString(s string) string {
	if s == "" {
		return "-"
	}
	t, err := strconv.ParseInt(s, 10, 64)
	if err != nil || t <= 0 {
		return "-"
	}
	secs := t
	if t >= msTimestampThreshold {
		secs = t / 1000
	}
	return time.Unix(secs, 0).UTC().Format(time.RFC3339)
}
