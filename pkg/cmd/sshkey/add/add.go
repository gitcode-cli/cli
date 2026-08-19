// Package add implements the ssh-key add command.
package add

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// AddOptions configures the ssh-key add command.
type AddOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	Title      string
	KeyFile    string
	JSON       bool
}

// NewCmdAdd creates the ssh-key add command.
func NewCmdAdd(f *cmdutil.Factory, runF func(*AddOptions) error) *cobra.Command {
	opts := &AddOptions{IO: f.IOStreams, HttpClient: f.HttpClient}
	cmd := &cobra.Command{
		Use:     "add",
		Short:   "Add an SSH public key",
		Long:    heredoc.Doc("Add an SSH public key from an explicitly provided file."),
		Example: heredoc.Doc("$ gc ssh-key add --title laptop --key ~/.ssh/id_ed25519.pub\n$ gc ssh-key add --title laptop --key key.pub --json"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return addRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.Title, "title", "t", "", "Title for the SSH key")
	cmd.Flags().StringVarP(&opts.KeyFile, "key", "k", "", "Path to the SSH public key file")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)
	_ = cmd.MarkFlagRequired("title")
	_ = cmd.MarkFlagRequired("key")
	return cmd
}

func addRun(opts *AddOptions) error {
	title := strings.TrimSpace(opts.Title)
	if title == "" {
		return cmdutil.NewUsageError("--title cannot be empty")
	}
	keyText, err := cmdutil.ReadTextFile(opts.KeyFile)
	if err != nil {
		return fmt.Errorf("failed to read SSH public key file: %w", err)
	}
	keyText = strings.TrimSpace(keyText)
	if err := validatePublicKey(keyText); err != nil {
		return err
	}
	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}
	key, err := api.CreateSSHKey(client, &api.CreateSSHKeyOptions{Title: title, Key: keyText})
	if err != nil {
		return fmt.Errorf("failed to add SSH key: %w", err)
	}
	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, key)
	}
	fmt.Fprintf(opts.IO.Out, "Added SSH key %s (ID: %d).\n", cmdutil.SanitizeTerminalText(key.Title), key.ID)
	return nil
}

func validatePublicKey(keyText string) error {
	if keyText == "" {
		return cmdutil.NewUsageError("SSH public key file is empty")
	}
	if strings.Contains(keyText, "PRIVATE KEY") {
		return cmdutil.NewUsageError("--key must reference a public key, not a private key")
	}
	if strings.ContainsAny(keyText, "\r\n") {
		return cmdutil.NewUsageError("SSH public key must contain exactly one line; expected a single OpenSSH public key")
	}
	key, _, options, rest, err := ssh.ParseAuthorizedKey([]byte(keyText))
	fields := strings.Fields(keyText)
	if err != nil || len(fields) < 2 || fields[0] != key.Type() || len(options) != 0 ||
		len(strings.TrimSpace(string(rest))) != 0 || !supportedKeyType(key.Type()) {
		return cmdutil.NewUsageError("SSH public key must be a single OpenSSH public key line")
	}
	return nil
}

func supportedKeyType(keyType string) bool {
	switch keyType {
	case ssh.KeyAlgoRSA, ssh.KeyAlgoDSA, ssh.KeyAlgoED25519,
		ssh.KeyAlgoECDSA256, ssh.KeyAlgoECDSA384, ssh.KeyAlgoECDSA521,
		ssh.KeyAlgoSKED25519, ssh.KeyAlgoSKECDSA256:
		return true
	default:
		return false
	}
}
