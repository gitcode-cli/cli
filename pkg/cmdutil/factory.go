// Package cmdutil provides utilities for command implementation
package cmdutil

import (
	"fmt"
	"net/http"

	"gitcode.com/gitcode-cli/cli/api"
	gitpkg "gitcode.com/gitcode-cli/cli/git"
	"gitcode.com/gitcode-cli/cli/internal/config"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// Factory provides dependencies for commands
type Factory struct {
	IOStreams  *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	Config     func() (config.Config, error)
	BaseRepo   func() (string, error)
	Branch     func() (string, error)
}

// NewFactory creates a new Factory with default settings
func NewFactory() *Factory {
	return &Factory{
		IOStreams: iostreams.System(),
		HttpClient: func() (*http.Client, error) {
			return api.DefaultHTTPClient(), nil
		},
		Config: func() (config.Config, error) {
			return config.New(), nil
		},
		BaseRepo: func() (string, error) {
			if !gitpkg.IsRepo() {
				return "", fmt.Errorf("not in a git repository")
			}
			repo, err := gitpkg.CurrentRepo()
			if err != nil {
				return "", err
			}
			return repo.String(), nil
		},
		Branch: func() (string, error) {
			if !gitpkg.IsRepo() {
				return "", fmt.Errorf("not in a git repository")
			}
			branch, err := gitpkg.CurrentBranch()
			if err != nil {
				return "", err
			}
			if branch == "" || branch == "HEAD" {
				return "", fmt.Errorf("could not determine current branch")
			}
			return branch, nil
		},
	}
}

// TestFactory creates a Factory for testing
func TestFactory() *Factory {
	io, _, _, _ := iostreams.Test()
	return &Factory{
		IOStreams: io,
		HttpClient: func() (*http.Client, error) {
			return api.DefaultHTTPClient(), nil
		},
		Config: func() (config.Config, error) {
			return config.New(), nil
		},
		BaseRepo: func() (string, error) {
			return "owner/repo", nil
		},
		Branch: func() (string, error) {
			return "main", nil
		},
	}
}
