package output

import (
	"fmt"
	"io"
	"strings"
)

// TablePrinter outputs data in table format
type TablePrinter struct {
	opts *Options
}

// PrintIssues prints issues in table format
func (p *TablePrinter) PrintIssues(w io.Writer, issues interface{}) error {
	// TODO: implement with actual issue type
	issueList, ok := issues.([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid issue list type")
	}

	if len(issueList) == 0 {
		fmt.Fprintln(w, "No issues found")
		return nil
	}

	// Header
	p.printRow(w, []string{"NUMBER", "STATE", "TITLE", "AUTHOR", "LABELS"})
	p.printSeparator(w, 5)

	// Rows
	for _, issue := range issueList {
		number := fmt.Sprintf("#%v", issue["number"])
		state := fmt.Sprintf("%v", issue["state"])
		title := truncate(fmt.Sprintf("%v", issue["title"]), 40)
		author := fmt.Sprintf("%v", issue["author"])
		labels := formatLabels(issue["labels"])

		p.printRow(w, []string{number, state, title, author, labels})
	}

	return nil
}

// PrintPRs prints pull requests in table format
func (p *TablePrinter) PrintPRs(w io.Writer, prs interface{}) error {
	prList, ok := prs.([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid PR list type")
	}

	if len(prList) == 0 {
		fmt.Fprintln(w, "No pull requests found")
		return nil
	}

	p.printRow(w, []string{"NUMBER", "STATE", "TITLE", "AUTHOR", "REVIEW"})
	p.printSeparator(w, 5)

	for _, pr := range prList {
		number := fmt.Sprintf("#%v", pr["number"])
		state := fmt.Sprintf("%v", pr["state"])
		title := truncate(fmt.Sprintf("%v", pr["title"]), 40)
		author := fmt.Sprintf("%v", pr["author"])
		review := fmt.Sprintf("%v", pr["review_status"])

		p.printRow(w, []string{number, state, title, author, review})
	}

	return nil
}

// PrintRepos prints repositories in table format
func (p *TablePrinter) PrintRepos(w io.Writer, repos interface{}) error {
	repoList, ok := repos.([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid repo list type")
	}

	if len(repoList) == 0 {
		fmt.Fprintln(w, "No repositories found")
		return nil
	}

	p.printRow(w, []string{"NAME", "VISIBILITY", "DESCRIPTION", "LANGUAGE"})
	p.printSeparator(w, 4)

	for _, repo := range repoList {
		name := fmt.Sprintf("%v", repo["name"])
		visibility := fmt.Sprintf("%v", repo["visibility"])
		desc := truncate(fmt.Sprintf("%v", repo["description"]), 30)
		language := fmt.Sprintf("%v", repo["language"])

		p.printRow(w, []string{name, visibility, desc, language})
	}

	return nil
}

// PrintReleases prints releases in table format
func (p *TablePrinter) PrintReleases(w io.Writer, releases interface{}) error {
	relList, ok := releases.([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid release list type")
	}

	if len(relList) == 0 {
		fmt.Fprintln(w, "No releases found")
		return nil
	}

	p.printRow(w, []string{"TAG", "NAME", "TYPE", "CREATED"})
	p.printSeparator(w, 4)

	for _, rel := range relList {
		tag := fmt.Sprintf("%v", rel["tag"])
		name := truncate(fmt.Sprintf("%v", rel["name"]), 30)
		typ := fmt.Sprintf("%v", rel["type"])
		created := fmt.Sprintf("%v", rel["created_at"])

		p.printRow(w, []string{tag, name, typ, created})
	}

	return nil
}

// PrintOne prints a single item as key-value table
func (p *TablePrinter) PrintOne(w io.Writer, item interface{}) error {
	data, ok := item.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid item type")
	}

	for key, value := range data {
		fmt.Fprintf(w, "%s: %v\n", strings.ToUpper(key), value)
	}

	return nil
}

func (p *TablePrinter) printRow(w io.Writer, cols []string) {
	fmt.Fprintf(w, "│")
	for _, col := range cols {
		fmt.Fprintf(w, " %-*s │", 12, truncate(col, 12))
	}
	fmt.Fprintln(w)
}

func (p *TablePrinter) printSeparator(w io.Writer, count int) {
	fmt.Fprintf(w, "├")
	for i := 0; i < count; i++ {
		fmt.Fprintf(w, "────────────┼")
	}
	fmt.Fprintf(w, "────────────┤\n")
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func formatLabels(labels interface{}) string {
	if labels == nil {
		return ""
	}
	labelList, ok := labels.([]string)
	if !ok {
		return fmt.Sprintf("%v", labels)
	}
	return strings.Join(labelList, ", ")
}