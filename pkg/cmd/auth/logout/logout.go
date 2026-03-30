// Package logout implements the auth logout command
package logout

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/internal/config"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type LogoutOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	Config     func() (config.Config, error)

	// Flags
	Hostname string
	Username string
}

// NewCmdLogout creates the logout command
func NewCmdLogout(f *cmdutil.Factory, runF func(*LogoutOptions) error) *cobra.Command {
	opts := &LogoutOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Config:     f.Config,
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

	cmd.Flags().StringVarP(&opts.Hostname, "hostname", "H", "", "The hostname of the GitCode instance to log out of")
	cmd.Flags().StringVarP(&opts.Username, "username", "u", "", "The username to log out of")

	return cmd
}

func logoutRun(opts *LogoutOptions) error {
	// Set default hostname
	cfg, err := opts.Config()
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	authCfg := cfg.Authentication()
	if opts.Hostname == "" {
		opts.Hostname, _ = authCfg.DefaultHost()
	}

	cs := opts.IO.ColorScheme()
	_, source := authCfg.ActiveToken(opts.Hostname)
	if err := authCfg.Logout(opts.Hostname, opts.Username); err != nil {
		return fmt.Errorf("failed to clear stored authentication: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Cleared stored authentication for %s\n", cs.Green("✓"), opts.Hostname)
	if source == "GC_TOKEN" || source == "GITCODE_TOKEN" {
		fmt.Fprintf(opts.IO.Out, "\n")
		fmt.Fprintf(opts.IO.Out, "Environment token is still active. Unset it manually to fully log out:\n")
		fmt.Fprintf(opts.IO.Out, "  unset GC_TOKEN\n")
		fmt.Fprintf(opts.IO.Out, "  unset GITCODE_TOKEN\n")
	}

	return nil
}
