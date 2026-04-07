package output

import (
	"fmt"
	"io"

	"gitcode.com/gitcode-cli/cli/api"
)

// RepoListOptions configures repo list text output.
type RepoListOptions struct {
	Format Format
}

// RepoListPrinter renders repository lists.
type RepoListPrinter struct {
	opts RepoListOptions
}

// NewRepoListPrinter validates and returns a repository printer.
func NewRepoListPrinter(opts RepoListOptions) (*RepoListPrinter, error) {
	if opts.Format == "" {
		opts.Format = FormatSimple
	}
	return &RepoListPrinter{opts: opts}, nil
}

// Print renders repositories according to the configured format.
func (p *RepoListPrinter) Print(w io.Writer, repos []api.Repository) error {
	switch p.opts.Format {
	case FormatTable:
		return p.printTable(w, repos)
	default:
		return p.printSimple(w, repos)
	}
}

func (p *RepoListPrinter) printSimple(w io.Writer, repos []api.Repository) error {
	for _, repo := range repos {
		fmt.Fprintf(w, "%s  %s  %s\n", repoDisplayName(repo), repoVisibility(repo), repo.Description)
	}
	return nil
}

func (p *RepoListPrinter) printTable(w io.Writer, repos []api.Repository) error {
	maxNameWidth := len("NAME")
	maxVisibilityWidth := len("VISIBILITY")
	maxLanguageWidth := len("LANGUAGE")

	for _, repo := range repos {
		if width := len(repoDisplayName(repo)); width > maxNameWidth {
			maxNameWidth = width
		}
		if width := len(repoVisibility(repo)); width > maxVisibilityWidth {
			maxVisibilityWidth = width
		}
		if width := len(repoLanguage(repo)); width > maxLanguageWidth {
			maxLanguageWidth = width
		}
	}

	fmt.Fprintf(w, "%-*s  %-*s  %-*s  %s\n", maxNameWidth, "NAME", maxVisibilityWidth, "VISIBILITY", maxLanguageWidth, "LANGUAGE", "DESCRIPTION")
	for _, repo := range repos {
		fmt.Fprintf(
			w,
			"%-*s  %-*s  %-*s  %s\n",
			maxNameWidth,
			repoDisplayName(repo),
			maxVisibilityWidth,
			repoVisibility(repo),
			maxLanguageWidth,
			repoLanguage(repo),
			repo.Description,
		)
	}

	return nil
}

func repoDisplayName(repo api.Repository) string {
	if repo.FullName != "" {
		return repo.FullName
	}
	if repo.Name != "" {
		return repo.Name
	}
	return "-"
}

func repoVisibility(repo api.Repository) string {
	if repo.Private {
		return "private"
	}
	return "public"
}

func repoLanguage(repo api.Repository) string {
	if repo.Language == "" {
		return "-"
	}
	return repo.Language
}
