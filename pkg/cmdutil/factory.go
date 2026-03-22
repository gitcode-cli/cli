// Package cmdutil provides utilities for command implementation
package cmdutil

import (
	"net/http"

	"github.com/gitcode-com/gitcode-cli/internal/config"
	"github.com/gitcode-com/gitcode-cli/pkg/iostreams"
)

// Factory provides dependencies for commands
type Factory struct {
	IOStreams  *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	Config     func() (config.Config, error)
	Branch     func() (string, error)
}

// NewFactory creates a new Factory with default settings
func NewFactory() *Factory {
	return &Factory{
		IOStreams: iostreams.System(),
		HttpClient: func() (*http.Client, error) {
			return &http.Client{}, nil
		},
		Config: func() (config.Config, error) {
			return config.New(), nil
		},
		Branch: func() (string, error) {
			return "", nil // TODO: implement git branch detection
		},
	}
}