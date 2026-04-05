package output

import (
	"fmt"
	"io"
)

// SimplePrinter outputs data in simple text format
type SimplePrinter struct {
	opts *Options
}

// PrintIssues prints issues in simple format
func (p *SimplePrinter) PrintIssues(w io.Writer, issues interface{}) error {
	issueList, ok := issues.([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid issue list type")
	}

	for _, issue := range issueList {
		number := fmt.Sprintf("#%v", issue["number"])
		state := fmt.Sprintf("%v", issue["state"])
		title := fmt.Sprintf("%v", issue["title"])

		fmt.Fprintf(w, "%-8s %s  %s\n", number, state, title)
	}

	return nil
}

// PrintPRs prints pull requests in simple format
func (p *SimplePrinter) PrintPRs(w io.Writer, prs interface{}) error {
	prList, ok := prs.([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid PR list type")
	}

	for _, pr := range prList {
		number := fmt.Sprintf("#%v", pr["number"])
		state := fmt.Sprintf("%v", pr["state"])
		title := fmt.Sprintf("%v", pr["title"])

		fmt.Fprintf(w, "%-8s %s  %s\n", number, state, title)
	}

	return nil
}

// PrintRepos prints repositories in simple format
func (p *SimplePrinter) PrintRepos(w io.Writer, repos interface{}) error {
	repoList, ok := repos.([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid repo list type")
	}

	for _, repo := range repoList {
		name := fmt.Sprintf("%v", repo["name"])
		desc := fmt.Sprintf("%v", repo["description"])

		fmt.Fprintf(w, "%s  %s\n", name, desc)
	}

	return nil
}

// PrintReleases prints releases in simple format
func (p *SimplePrinter) PrintReleases(w io.Writer, releases interface{}) error {
	relList, ok := releases.([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid release list type")
	}

	for _, rel := range relList {
		tag := fmt.Sprintf("%v", rel["tag"])
		name := fmt.Sprintf("%v", rel["name"])

		fmt.Fprintf(w, "%s  %s\n", tag, name)
	}

	return nil
}

// PrintOne prints a single item in simple format
func (p *SimplePrinter) PrintOne(w io.Writer, item interface{}) error {
	data, ok := item.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid item type")
	}

	for key, value := range data {
		fmt.Fprintf(w, "%s: %v\n", key, value)
	}

	return nil
}