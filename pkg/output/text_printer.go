package output

import (
	"fmt"
	"io"

	"gitcode.com/gitcode-cli/cli/api"
)

// SimplePrinter outputs data in simple text format
type SimplePrinter struct {
	opts *Options
}

// PrintIssues prints issues in simple format
func (p *SimplePrinter) PrintIssues(w io.Writer, issues interface{}) error {
	switch v := issues.(type) {
	case []api.Issue:
		return p.printAPIIssues(w, v)
	case []*api.Issue:
		return p.printAPIIssuePtrs(w, v)
	case []map[string]interface{}:
		return p.printMapIssues(w, v)
	default:
		return fmt.Errorf("invalid issue list type: %T", issues)
	}
}

func (p *SimplePrinter) printAPIIssues(w io.Writer, issueList []api.Issue) error {
	// Calculate max width for issue numbers
	maxNumWidth := 0
	for _, issue := range issueList {
		w := len(fmt.Sprintf("#%s", issue.Number))
		if w > maxNumWidth {
			maxNumWidth = w
		}
	}

	for _, issue := range issueList {
		fmt.Fprintf(w, "%-*s %s  %s\n", maxNumWidth, fmt.Sprintf("#%s", issue.Number), issue.State, issue.Title)
	}

	return nil
}

func (p *SimplePrinter) printAPIIssuePtrs(w io.Writer, issueList []*api.Issue) error {
	maxNumWidth := 0
	for _, issue := range issueList {
		w := len(fmt.Sprintf("#%s", issue.Number))
		if w > maxNumWidth {
			maxNumWidth = w
		}
	}

	for _, issue := range issueList {
		fmt.Fprintf(w, "%-*s %s  %s\n", maxNumWidth, fmt.Sprintf("#%s", issue.Number), issue.State, issue.Title)
	}

	return nil
}

func (p *SimplePrinter) printMapIssues(w io.Writer, issueList []map[string]interface{}) error {
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
	switch v := prs.(type) {
	case []*api.PullRequest:
		return p.printAPIDocsPRs(w, v)
	case []map[string]interface{}:
		return p.printMapPRs(w, v)
	default:
		return fmt.Errorf("invalid PR list type: %T", prs)
	}
}

func (p *SimplePrinter) printAPIDocsPRs(w io.Writer, prList []*api.PullRequest) error {
	maxNumWidth := 0
	for _, pr := range prList {
		w := len(fmt.Sprintf("#%d", pr.Number))
		if w > maxNumWidth {
			maxNumWidth = w
		}
	}

	for _, pr := range prList {
		fmt.Fprintf(w, "%-*s %s  %s\n", maxNumWidth, fmt.Sprintf("#%d", pr.Number), pr.State, pr.Title)
	}

	return nil
}

func (p *SimplePrinter) printMapPRs(w io.Writer, prList []map[string]interface{}) error {
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