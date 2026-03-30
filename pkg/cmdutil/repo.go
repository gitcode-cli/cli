package cmdutil

import "fmt"

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
