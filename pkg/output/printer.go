package output

import (
	"io"
)

// Format defines the output format type
type Format string

const (
	FormatTable  Format = "table"
	FormatJSON   Format = "json"
	FormatSimple Format = "simple"
)

// Options configures the output behavior
type Options struct {
	Format     Format
	NoColor    bool
	TimeFormat string
	Columns    []string // optional: specific columns to display
}

// Printer defines the interface for output formatting
type Printer interface {
	// PrintIssues formats and prints a list of issues
	PrintIssues(w io.Writer, issues interface{}) error

	// PrintPRs formats and prints a list of pull requests
	PrintPRs(w io.Writer, prs interface{}) error

	// PrintRepos formats and prints a list of repositories
	PrintRepos(w io.Writer, repos interface{}) error

	// PrintReleases formats and prints a list of releases
	PrintReleases(w io.Writer, releases interface{}) error

	// PrintOne prints a single item in detail view
	PrintOne(w io.Writer, item interface{}) error
}

// NewPrinter creates a printer based on the output format
func NewPrinter(opts *Options) Printer {
	switch opts.Format {
	case FormatJSON:
		return &JSONPrinter{opts: opts}
	case FormatTable:
		return &TablePrinter{opts: opts}
	default:
		return &SimplePrinter{opts: opts}
	}
}