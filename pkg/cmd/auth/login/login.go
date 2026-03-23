// Package login implements the auth login command
package login

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type LoginOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Flags
	Hostname     string
	Token        string
	WithToken    bool
	GitProtocol  string
	Web          bool
}

// NewCmdLogin creates the login command
func NewCmdLogin(f *cmdutil.Factory, runF func(*LoginOptions) error) *cobra.Command {
	opts := &LoginOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to a GitCode account",
		Long: heredoc.Doc(`
			Authenticate with a GitCode account.

			The default authentication mode is a web-based flow.
			Alternatively, use --with-token to pass a token on standard input.
		`),
		Example: heredoc.Doc(`
			# Start interactive login
			$ gc auth login

			# Login with a token from stdin
			$ echo "your-token" | gc auth login --with-token

			# Login with a token from a file
			$ cat token.txt | gc auth login --with-token

			# Login to a specific host
			$ gc auth login --hostname gitcode.com
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			// Handle --with-token
			if opts.WithToken {
				return loginWithToken(opts)
			}

			// Handle --token flag
			if opts.Token != "" {
				return loginWithTokenFlag(opts)
			}

			// Interactive login
			return loginInteractive(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Hostname, "hostname", "H", "", "The hostname of the GitCode instance to authenticate with")
	cmd.Flags().StringVarP(&opts.Token, "token", "t", "", "An authentication token for GitCode")
	cmd.Flags().BoolVar(&opts.WithToken, "with-token", false, "Read token from standard input")
	cmd.Flags().StringVarP(&opts.GitProtocol, "git-protocol", "p", "https", "The Git protocol to use for operations (https/ssh)")
	cmd.Flags().BoolVarP(&opts.Web, "web", "w", false, "Open a browser to authenticate")

	return cmd
}

func loginWithToken(opts *LoginOptions) error {
	// Read token from stdin
	reader := bufio.NewReader(os.Stdin)
	token, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read token from stdin: %w", err)
	}
	token = strings.TrimSpace(token)

	if token == "" {
		return fmt.Errorf("no token provided on stdin")
	}

	opts.Token = token
	return loginWithTokenFlag(opts)
}

func loginWithTokenFlag(opts *LoginOptions) error {
	if opts.Token == "" {
		return fmt.Errorf("no token provided")
	}

	// Set default hostname
	if opts.Hostname == "" {
		opts.Hostname = "gitcode.com"
	}

	// Verify the token
	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	user, err := api.VerifyToken(httpClient, opts.Hostname, opts.Token)
	if err != nil {
		return fmt.Errorf("failed to verify token: %w", err)
	}

	// Token is valid - in memory only, not persisted to file
	// In a real implementation, this would be stored in keyring
	fmt.Fprintf(opts.IO.Out, "%s Logged in as %s\n", opts.IO.ColorScheme().Green("✓"), user.Login)
	fmt.Fprintf(opts.IO.Out, "  Host: %s\n", opts.Hostname)
	fmt.Fprintf(opts.IO.Out, "  Git protocol: %s\n", opts.GitProtocol)
	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "Note: Token is stored in memory only for this session.\n")

	return nil
}

func loginInteractive(opts *LoginOptions) error {
	// Set default hostname
	if opts.Hostname == "" {
		opts.Hostname = "gitcode.com"
	}

	cs := opts.IO.ColorScheme()

	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "%s Authenticate with GitCode\n", cs.Bold("Tip:"))
	fmt.Fprintf(opts.IO.Out, "\n")

	// Prompt for token
	fmt.Fprintf(opts.IO.Out, "? Paste your authentication token: ")
	reader := bufio.NewReader(os.Stdin)
	token, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read token: %w", err)
	}
	token = strings.TrimSpace(token)

	if token == "" {
		return fmt.Errorf("no token provided")
	}

	opts.Token = token
	return loginWithTokenFlag(opts)
}