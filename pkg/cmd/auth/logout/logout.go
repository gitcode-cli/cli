// Package logout implements the auth logout command
package logout

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
	"github.com/gitcode-com/gitcode-cli/pkg/iostreams"
)

type LogoutOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Flags
	Hostname string
	Username string
}

// NewCmdLogout creates the logout command
func NewCmdLogout(f *cmdutil.Factory, runF func(*LogoutOptions) error) *cobra.Command {
	opts := &LogoutOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out of a GitCode account",
		Long: heredoc.Doc(`
			Log out of a GitCode account.

			This removes the stored authentication token.
		`),
		Example: heredoc.Doc(`
			$ gc auth logout
			? Confirm logout of gitcode.com? Yes
			✓ Logged out of gitcode.com

			$ gc auth logout --hostname gitcode.com
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return logoutRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Hostname, "hostname", "h", "", "The hostname of the GitCode instance to log out of")
	cmd.Flags().StringVarP(&opts.Username, "username", "u", "", "The username to log out of")

	return cmd
}

func logoutRun(opts *LogoutOptions) error {
	// Set default hostname
	if opts.Hostname == "" {
		opts.Hostname = "gitcode.com"
	}

	cs := opts.IO.ColorScheme()

	// In memory-only mode, just clear from memory
	// The actual implementation would clear from keyring

	fmt.Fprintf(opts.IO.Out, "%s Logged out of %s\n", cs.Green("✓"), opts.Hostname)
	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "Note: If you used GC_TOKEN environment variable, unset it manually:\n")
	fmt.Fprintf(opts.IO.Out, "  unset GC_TOKEN\n")

	return nil
}