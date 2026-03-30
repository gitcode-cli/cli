package cmdutil

import (
	"fmt"

	gitpkg "gitcode.com/gitcode-cli/cli/git"
)

// ResolveRepo returns the explicit repository when provided, otherwise tries
// to infer it from the current git repository.
func ResolveRepo(repo string, baseRepo func() (string, error)) (string, error) {
	if repo != "" {
		return repo, nil
	}
	if baseRepo == nil {
		return "", fmt.Errorf("no repository specified. Use -R owner/repo")
	}

	detectedRepo, err := baseRepo()
	if err != nil {
		return "", fmt.Errorf("no repository specified and could not determine current repository: %w", err)
	}
	if detectedRepo == "" {
		return "", fmt.Errorf("no repository specified and could not determine current repository")
	}

	return detectedRepo, nil
}

// ParseRepo parses a repository reference and returns the owner and repository
// name. It supports owner/repo, HTTPS URLs, and SSH URLs.
func ParseRepo(repo string) (string, string, error) {
	if repo == "" {
		return "", "", fmt.Errorf("no repository specified. Use -R owner/repo")
	}

	parsedRepo, err := gitpkg.ParseRepo(repo)
	if err != nil {
		return "", "", fmt.Errorf("invalid repository format: %s", repo)
	}

	return parsedRepo.Owner, parsedRepo.Name, nil
}
