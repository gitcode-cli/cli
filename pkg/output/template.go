package output

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

// TemplatePrinter outputs data using Go templates
type TemplatePrinter struct {
	opts     *Options
	template *template.Template
}

// NewTemplatePrinter creates a printer with custom template
func NewTemplatePrinter(templateStr string, opts *Options) (*TemplatePrinter, error) {
	tmpl, err := template.New("output").Funcs(templateFuncs()).Parse(templateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid template: %w", err)
	}
	return &TemplatePrinter{opts: opts, template: tmpl}, nil
}

// PrintIssues prints issues using template
func (p *TemplatePrinter) PrintIssues(w io.Writer, issues interface{}) error {
	return p.executeTemplate(w, issues)
}

// PrintPRs prints PRs using template
func (p *TemplatePrinter) PrintPRs(w io.Writer, prs interface{}) error {
	return p.executeTemplate(w, prs)
}

// PrintRepos prints repos using template
func (p *TemplatePrinter) PrintRepos(w io.Writer, repos interface{}) error {
	return p.executeTemplate(w, repos)
}

// PrintReleases prints releases using template
func (p *TemplatePrinter) PrintReleases(w io.Writer, releases interface{}) error {
	return p.executeTemplate(w, releases)
}

// PrintOne prints single item using template
func (p *TemplatePrinter) PrintOne(w io.Writer, item interface{}) error {
	return p.executeTemplate(w, item)
}

func (p *TemplatePrinter) executeTemplate(w io.Writer, data interface{}) error {
	var buf bytes.Buffer
	if err := p.template.Execute(&buf, data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}
	_, err := w.Write(buf.Bytes())
	return err
}

// templateFuncs returns common template functions
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"upper":    upper,
		"lower":    lower,
		"trunc":    truncateStr,
		"json":     toJSON,
	}
}

func upper(s string) string {
	if len(s) > 50 {
		return s[:50]
	}
	return s
}

func lower(s string) string {
	result := make([]byte, len(s))
	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			result[i] = byte(c + 32)
		} else {
			result[i] = byte(c)
		}
	}
	return string(result)
}

func truncateStr(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func toJSON(v interface{}) string {
	return fmt.Sprintf("%+v", v)
}