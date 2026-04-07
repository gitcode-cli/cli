package output

import (
	"fmt"
	"io"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// PRListOptions configures PR list text output.
type PRListOptions struct {
	Format Format
	Color  *iostreams.ColorScheme
}

// PRListPrinter renders pull request lists.
type PRListPrinter struct {
	opts PRListOptions
}

// NewPRListPrinter validates and returns a PR printer.
func NewPRListPrinter(opts PRListOptions) (*PRListPrinter, error) {
	if opts.Format == "" {
		opts.Format = FormatSimple
	}
	return &PRListPrinter{opts: opts}, nil
}

// Print renders PRs according to the configured format.
func (p *PRListPrinter) Print(w io.Writer, prs []api.PullRequest) error {
	switch p.opts.Format {
	case FormatTable:
		return p.printTable(w, prs)
	default:
		return p.printSimple(w, prs)
	}
}

func (p *PRListPrinter) printSimple(w io.Writer, prs []api.PullRequest) error {
	maxNumWidth := 0
	stateWidth := len("merged")
	for _, pr := range prs {
		if width := len(fmt.Sprintf("#%d", pr.Number)); width > maxNumWidth {
			maxNumWidth = width
		}
		if width := len(p.stateText(pr)); width > stateWidth {
			stateWidth = width
		}
	}

	for _, pr := range prs {
		fmt.Fprintf(
			w,
			"%-*s %s %s\n",
			maxNumWidth,
			fmt.Sprintf("#%d", pr.Number),
			p.stateLabel(pr, stateWidth),
			pr.Title,
		)
	}
	return nil
}

func (p *PRListPrinter) printTable(w io.Writer, prs []api.PullRequest) error {
	maxNumWidth := len("NUMBER")
	maxStateWidth := len("STATE")
	maxAuthorWidth := len("AUTHOR")
	maxReviewWidth := len("REVIEW")

	for _, pr := range prs {
		if width := len(fmt.Sprintf("#%d", pr.Number)); width > maxNumWidth {
			maxNumWidth = width
		}
		if width := len(p.stateText(pr)); width > maxStateWidth {
			maxStateWidth = width
		}
		if width := len(prAuthor(pr)); width > maxAuthorWidth {
			maxAuthorWidth = width
		}
		if width := len(prReviewStatus(pr)); width > maxReviewWidth {
			maxReviewWidth = width
		}
	}

	fmt.Fprintf(w, "%-*s  %-*s  %-*s  %-*s  %s\n", maxNumWidth, "NUMBER", maxStateWidth, "STATE", maxAuthorWidth, "AUTHOR", maxReviewWidth, "REVIEW", "TITLE")
	for _, pr := range prs {
		fmt.Fprintf(
			w,
			"%-*s  %-*s  %-*s  %-*s  %s\n",
			maxNumWidth,
			fmt.Sprintf("#%d", pr.Number),
			maxStateWidth,
			p.stateLabel(pr, maxStateWidth),
			maxAuthorWidth,
			prAuthor(pr),
			maxReviewWidth,
			prReviewStatus(pr),
			pr.Title,
		)
	}

	return nil
}

func (p *PRListPrinter) stateLabel(pr api.PullRequest, width int) string {
	label := fmt.Sprintf("%-*s", width, p.stateText(pr))
	if p.opts.Color == nil {
		return label
	}

	switch p.stateText(pr) {
	case "merged":
		return p.opts.Color.Magenta(label)
	case "closed":
		return p.opts.Color.Red(label)
	case "draft":
		return p.opts.Color.Gray(label)
	default:
		return p.opts.Color.Green(label)
	}
}

func (p *PRListPrinter) stateText(pr api.PullRequest) string {
	switch {
	case pr.Merged || pr.State == "merged":
		return "merged"
	case pr.State == "closed":
		return "closed"
	case pr.Draft:
		return "draft"
	default:
		return "open"
	}
}

func prAuthor(pr api.PullRequest) string {
	if pr.User == nil || pr.User.Login == "" {
		return "-"
	}
	return pr.User.Login
}

func prReviewStatus(pr api.PullRequest) string {
	switch {
	case pr.Draft:
		return "draft"
	case pr.Merged || pr.State == "merged":
		return "merged"
	case len(pr.Reviewers) > 0:
		return "requested"
	case pr.State == "closed":
		return "closed"
	default:
		return "pending"
	}
}
