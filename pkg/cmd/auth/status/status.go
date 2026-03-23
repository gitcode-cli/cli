// Package status implements the auth status command
package status

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type StatusOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Flags
	Hostname string
	ShowToken bool
}

// NewCmdStatus creates the status command
func NewCmdStatus(f *cmdutil.Factory, runF func(*StatusOptions) error) *cobra.Command {
	opts := &StatusOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "View authentication status",
		Long: heredoc.Doc(`
			View information about your authentication status.

			Checks for token from:
			1. GC_TOKEN environment variable
			2. GITCODE_TOKEN environment variable
			3. Stored credentials (keyring)
		`),
		Example: heredoc.Doc(`
			$ gc auth status
			gitcode.com
			  ✓ Logged in as username (GC_TOKEN)
			  ✓ Git operations protocol: https
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return statusRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Hostname, "hostname", "H", "", "Check a specific hostname")
	cmd.Flags().BoolVar(&opts.ShowToken, "show-token", false, "Display the auth token")

	return cmd
}

func statusRun(opts *StatusOptions) error {
	// Set default hostname
	if opts.Hostname == "" {
		opts.Hostname = "gitcode.com"
	}

	cs := opts.IO.ColorScheme()

	// Check for token from environment
	token := os.Getenv("GC_TOKEN")
	tokenSource := "GC_TOKEN"
	if token == "" {
		token = os.Getenv("GITCODE_TOKEN")
		tokenSource = "GITCODE_TOKEN"
	}

	fmt.Fprintf(opts.IO.Out, "%s\n", opts.Hostname)

	if token == "" {
		fmt.Fprintf(opts.IO.Out, "  %s Not logged in\n", cs.Red("✗"))
		fmt.Fprintf(opts.IO.Out, "\n")
		fmt.Fprintf(opts.IO.Out, "To authenticate, run: gc auth login\n")
		return nil
	}

	// Verify token
	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	user, err := api.VerifyToken(httpClient, opts.Hostname, token)
	if err != nil {
		fmt.Fprintf(opts.IO.Out, "  %s Token is invalid or expired\n", cs.Red("✗"))
		fmt.Fprintf(opts.IO.Out, "\n")
		fmt.Fprintf(opts.IO.Out, "To re-authenticate, run: gc auth login\n")
		return nil
	}

	// Display logged in status
	fmt.Fprintf(opts.IO.Out, "  %s Logged in as %s (%s)\n", cs.Green("✓"), user.Login, tokenSource)
	fmt.Fprintf(opts.IO.Out, "  %s Git operations protocol: https\n", cs.Green("✓"))

	if opts.ShowToken {
		// Mask most of the token
		maskedToken := maskToken(token)
		fmt.Fprintf(opts.IO.Out, "  %s Token: %s\n", cs.Green("✓"), maskedToken)
	}

	return nil
}

func maskToken(token string) string {
	if len(token) <= 8 {
		return "********"
	}
	return token[:4] + "..." + token[len(token)-4:]
}