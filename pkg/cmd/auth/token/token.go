// Package token implements the auth token command
package token

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
	"github.com/gitcode-com/gitcode-cli/pkg/iostreams"
)

type TokenOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Flags
	Hostname string
}

// NewCmdToken creates the token command
func NewCmdToken(f *cmdutil.Factory, runF func(*TokenOptions) error) *cobra.Command {
	opts := &TokenOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "token",
		Short: "Print an authentication token",
		Long: heredoc.Doc(`
			Print the authentication token for the current session.

			The token is printed to stdout for use in scripts or piping.
			If no token is found, returns an error.
		`),
		Example: heredoc.Doc(`
			$ gc auth token
			gc_xxxxxxxxxxxx

			# Use token in a script
			$ TOKEN=$(gc auth token)

			# Specify hostname
			$ gc auth token --hostname gitcode.com
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return tokenRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Hostname, "hostname", "H", "", "The hostname to get the token for")

	return cmd
}

func tokenRun(opts *TokenOptions) error {
	// Set default hostname
	if opts.Hostname == "" {
		opts.Hostname = "gitcode.com"
	}

	// Check for token from environment
	token := os.Getenv("GC_TOKEN")
	if token == "" {
		token = os.Getenv("GITCODE_TOKEN")
	}

	if token == "" {
		return fmt.Errorf("no authentication token found")
	}

	// Print token to stdout (for piping)
	fmt.Fprintln(opts.IO.Out, token)

	return nil
}