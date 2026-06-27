// Package status implements the auth status command
package status

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/config"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type StatusOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	Config     func() (config.Config, error)

	// Flags
	Hostname    string
	HostnameSet bool
	ShowToken   bool
	JSON        bool
}

// AuthStatus represents the JSON output structure for auth status
type AuthStatus struct {
	Hostname    string `json:"hostname"`
	LoggedIn    bool   `json:"logged_in"`
	Username    string `json:"username,omitempty"`
	TokenSource string `json:"token_source,omitempty"`
	TokenValid  bool   `json:"token_valid,omitempty"`
	GitProtocol string `json:"git_protocol,omitempty"`
	Token       string `json:"token,omitempty"`
}

// NewCmdStatus creates the status command
func NewCmdStatus(f *cmdutil.Factory, runF func(*StatusOptions) error) *cobra.Command {
	opts := &StatusOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "View authentication status",
		Long: heredoc.Doc(`
			View information about your authentication status.

			When no hostname is specified, checks for token from:
			1. GC_TOKEN environment variable
			2. GITCODE_TOKEN environment variable
			3. Stored credentials from local config (~/.config/gc/auth.json)
		`),
		Example: heredoc.Doc(`
			# Check authentication status
			$ gc auth status
			gitcode.com
			  ✓ Logged in as username (GC_TOKEN)
			  ✓ Git operations protocol: https

			# Output as JSON
			$ gc auth status --json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.HostnameSet = cmd.Flags().Changed("hostname")
			if runF != nil {
				return runF(opts)
			}
			return statusRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Hostname, "hostname", "H", "", "Check a specific hostname")
	cmd.Flags().BoolVar(&opts.ShowToken, "show-token", false, "Display the full auth token")
	cmd.Flags().BoolVar(&opts.JSON, "json", false, "Output as JSON")

	return cmd
}

func statusRun(opts *StatusOptions) error {
	// Set default hostname
	cfg, err := opts.Config()
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	authCfg := cfg.Authentication()
	if opts.Hostname == "" {
		opts.Hostname, _ = authCfg.DefaultHost()
	}
	opts.Hostname, err = config.NormalizeTrustedHost(opts.Hostname)
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()

	token, tokenSource := authCfg.ActiveToken(opts.Hostname)
	if opts.HostnameSet || opts.Hostname != "gitcode.com" {
		token, tokenSource = authCfg.StoredToken(opts.Hostname)
	}

	// Build status for JSON output
	status := AuthStatus{
		Hostname: opts.Hostname,
	}

	if token == "" {
		status.LoggedIn = false

		if opts.JSON {
			return cmdutil.WriteJSON(opts.IO.Out, status)
		}

		fmt.Fprintf(opts.IO.Out, "%s\n", opts.Hostname)
		fmt.Fprintf(opts.IO.Out, "  %s Not logged in\n", cs.Red("✗"))
		fmt.Fprintf(opts.IO.Out, "\n")
		fmt.Fprintf(opts.IO.Out, "To authenticate, run: gc auth login\n")
		return nil
	}

	status.LoggedIn = true
	status.TokenSource = tokenSource

	// Verify token
	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	user, err := api.VerifyToken(httpClient, opts.Hostname, token)
	if err != nil {
		status.TokenValid = false

		if opts.JSON {
			return cmdutil.WriteJSON(opts.IO.Out, status)
		}

		fmt.Fprintf(opts.IO.Out, "%s\n", opts.Hostname)
		fmt.Fprintf(opts.IO.Out, "  %s Token is invalid or expired\n", cs.Red("✗"))
		fmt.Fprintf(opts.IO.Out, "\n")
		fmt.Fprintf(opts.IO.Out, "To re-authenticate, run: gc auth login\n")
		return nil
	}

	status.TokenValid = true
	status.Username = user.Login
	status.GitProtocol = cfg.GitProtocol(opts.Hostname).Value

	if opts.ShowToken {
		status.Token = token
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, status)
	}

	// Display logged in status
	fmt.Fprintf(opts.IO.Out, "%s\n", opts.Hostname)
	fmt.Fprintf(opts.IO.Out, "  %s Logged in as %s (%s)\n", cs.Green("✓"), user.Login, tokenSource)
	fmt.Fprintf(opts.IO.Out, "  %s Git operations protocol: %s\n", cs.Green("✓"), cfg.GitProtocol(opts.Hostname).Value)

	if opts.ShowToken {
		fmt.Fprintf(opts.IO.ErrOut, "%s Warning: displaying authentication token. Do not share this output.\n", cs.Yellow("!"))
		fmt.Fprintf(opts.IO.Out, "  %s Token: %s\n", cs.Green("✓"), token)
	}

	return nil
}
