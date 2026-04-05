package output

import (
	"encoding/json"
	"io"
)

// JSONPrinter outputs data in JSON format
type JSONPrinter struct {
	opts *Options
}

// PrintIssues prints issues as JSON
func (p *JSONPrinter) PrintIssues(w io.Writer, issues interface{}) error {
	return writeJSON(w, issues)
}

// PrintPRs prints pull requests as JSON
func (p *JSONPrinter) PrintPRs(w io.Writer, prs interface{}) error {
	return writeJSON(w, prs)
}

// PrintRepos prints repositories as JSON
func (p *JSONPrinter) PrintRepos(w io.Writer, repos interface{}) error {
	return writeJSON(w, repos)
}

// PrintReleases prints releases as JSON
func (p *JSONPrinter) PrintReleases(w io.Writer, releases interface{}) error {
	return writeJSON(w, releases)
}

// PrintOne prints a single item as JSON
func (p *JSONPrinter) PrintOne(w io.Writer, item interface{}) error {
	return writeJSON(w, item)
}

// writeJSON writes indented JSON to the writer
func writeJSON(w io.Writer, value interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		return err
	}
	return nil
}