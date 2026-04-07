package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/template"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// IssueListOptions configures issue list text output.
type IssueListOptions struct {
	Format     Format
	TimeFormat TimeFormat
	Template   string
	Color      *iostreams.ColorScheme
}

// IssueListPrinter renders issues for text and template outputs.
type IssueListPrinter struct {
	opts IssueListOptions
	tmpl *template.Template
}

// NewIssueListPrinter validates and returns a printer for issue lists.
func NewIssueListPrinter(opts IssueListOptions) (*IssueListPrinter, error) {
	if opts.Format == "" {
		opts.Format = FormatSimple
	}
	if opts.TimeFormat == "" {
		opts.TimeFormat = TimeFormatAbsolute
	}

	var tmpl *template.Template
	if strings.TrimSpace(opts.Template) != "" {
		parsed, err := template.New("issues").Funcs(templateFuncs(opts.TimeFormat)).Parse(opts.Template)
		if err != nil {
			return nil, fmt.Errorf("invalid template: %w", err)
		}
		tmpl = parsed
	}

	return &IssueListPrinter{opts: opts, tmpl: tmpl}, nil
}

// Print renders an issue list according to the printer options.
func (p *IssueListPrinter) Print(w io.Writer, issues []api.Issue) error {
	if p.tmpl != nil {
		return p.tmpl.Execute(w, issues)
	}

	switch p.opts.Format {
	case FormatTable:
		return p.printTable(w, issues)
	default:
		return p.printSimple(w, issues)
	}
}

func (p *IssueListPrinter) printSimple(w io.Writer, issues []api.Issue) error {
	maxNumWidth := 0
	for _, issue := range issues {
		if width := len(fmt.Sprintf("#%s", issue.Number)); width > maxNumWidth {
			maxNumWidth = width
		}
	}

	for _, issue := range issues {
		fmt.Fprintf(
			w,
			"%-*s %-8s %-16s %s\n",
			maxNumWidth,
			fmt.Sprintf("#%s", issue.Number),
			p.stateLabel(issue.State),
			FormatFlexibleTime(issue.UpdatedAt, p.opts.TimeFormat),
			issue.Title,
		)
	}

	return nil
}

func (p *IssueListPrinter) printTable(w io.Writer, issues []api.Issue) error {
	maxNumWidth := len("NUMBER")
	maxStateWidth := len("STATE")
	maxUpdatedWidth := len("UPDATED")

	for _, issue := range issues {
		if width := len(fmt.Sprintf("#%s", issue.Number)); width > maxNumWidth {
			maxNumWidth = width
		}
		if width := len(issue.State); width > maxStateWidth {
			maxStateWidth = width
		}
		if width := len(FormatFlexibleTime(issue.UpdatedAt, p.opts.TimeFormat)); width > maxUpdatedWidth {
			maxUpdatedWidth = width
		}
	}

	fmt.Fprintf(w, "%-*s  %-*s  %-*s  %s\n", maxNumWidth, "NUMBER", maxStateWidth, "STATE", maxUpdatedWidth, "UPDATED", "TITLE")
	for _, issue := range issues {
		fmt.Fprintf(
			w,
			"%-*s  %-*s  %-*s  %s\n",
			maxNumWidth,
			fmt.Sprintf("#%s", issue.Number),
			maxStateWidth,
			p.stateLabel(issue.State),
			maxUpdatedWidth,
			FormatFlexibleTime(issue.UpdatedAt, p.opts.TimeFormat),
			issue.Title,
		)
	}

	return nil
}

func (p *IssueListPrinter) stateLabel(state string) string {
	if p.opts.Color == nil {
		return state
	}
	switch state {
	case "closed":
		return p.opts.Color.Red("closed")
	case "open":
		return p.opts.Color.Green("open")
	default:
		return state
	}
}

func templateFuncs(timeFormat TimeFormat) template.FuncMap {
	return template.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"trunc": truncate,
		"json":  toJSON,
		"time": func(v interface{}) string {
			switch t := v.(type) {
			case api.FlexibleTime:
				return FormatFlexibleTime(t, timeFormat)
			default:
				return fmt.Sprintf("%v", v)
			}
		},
	}
}

func truncate(value string, max int) string {
	if max <= 0 || len(value) <= max {
		return value
	}
	if max <= 3 {
		return value[:max]
	}
	return value[:max-3] + "..."
}

func toJSON(v interface{}) (string, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return "", err
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}
