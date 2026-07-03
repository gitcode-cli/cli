// Package token implements the auth token command
package token

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/config"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// TokenInfo represents the JSON output for auth token command
type TokenInfo struct {
	Hostname string `json:"hostname"`
	Token    string `json:"token"`
	Source   string `json:"source,omitempty"`
}

type TokenOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	Config     func() (config.Config, error)

	// Flags
	Hostname    string
	HostnameSet bool
	JSON        bool
}

// tokenWarning is the security warning printed to stderr before displaying the token.
const tokenWarning = "Warning: displaying authentication token. Do not share this output."

// NewCmdToken creates the token command
func NewCmdToken(f *cmdutil.Factory, runF func(*TokenOptions) error) *cobra.Command {
	opts := &TokenOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:   "token",
		Short: "Print an authentication token",
		Long: heredoc.Doc(`
			Print the authentication token for the current session.

			The token is printed to stdout for use in scripts or piping.
			For safety, this command requires an interactive confirmation.
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
			opts.HostnameSet = cmd.Flags().Changed("hostname")
			if runF != nil {
				return runF(opts)
			}
			return tokenRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Hostname, "hostname", "H", "", "The hostname to get the token for")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func tokenRun(opts *TokenOptions) error {
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

	token, source := authCfg.ActiveToken(opts.Hostname)
	if opts.HostnameSet || opts.Hostname != "gitcode.com" {
		token, source = authCfg.StoredToken(opts.Hostname)
	}

	if token == "" {
		return cmdutil.NewAuthError("no authentication token found")
	}
	if err := cmdutil.ConfirmTokenDisclosure(opts.IO, opts.Hostname); err != nil {
		return err
	}

	fmt.Fprintln(opts.IO.ErrOut, tokenWarning)

	if opts.JSON {
		info := TokenInfo{
			Hostname: opts.Hostname,
			Token:    token,
			Source:   source,
		}
		return cmdutil.WriteJSON(opts.IO.Out, info)
	}

	// Print token to stdout (for piping)
	fmt.Fprintln(opts.IO.Out, token)

	return nil
}
