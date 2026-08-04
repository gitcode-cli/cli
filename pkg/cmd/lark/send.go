package lark

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/larkcli"
)

type sendOptions struct {
	IO      *iostreams.IOStreams
	LarkRun larkcli.RunFunc

	ChatID   string
	UserID   string
	ToSelf   bool
	Text     string
	Markdown string
	File     string
	BodyFile string
	As       string
	DryRun   bool
	JSON     bool
}

func newCmdSend(f *cmdutil.Factory, runF func(*sendOptions) error) *cobra.Command {
	opts := &sendOptions{
		IO:      f.IOStreams,
		LarkRun: larkcli.DefaultRun,
	}

	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send a message to a Feishu/Lark chat or direct message",
		Long: heredoc.Doc(`
			Send a message to a Feishu/Lark group, chat, or direct message by
			delegating to "lark-cli im +messages-send".

			Target resolution (--chat-id / --user-id / --to-self are mutually
			exclusive; if none is given, the default chat from
			"gc lark config set --default-chat" or GC_LARK_DEFAULT_CHAT_ID is
			used):
			  --chat-id oc_xxx   send to a group/chat
			  --user-id ou_xxx   send a direct message to a person
			  --to-self          send a direct message to the current lark-cli
			                     user (resolves their open_id from "lark-cli auth
			                     status"); defaults to the bot identity, which is
			                     the common "notify myself" pattern that needs no
			                     group membership.

			Feishu OAuth credentials stay in lark-cli's OS keychain and are never
			seen by gc. Message text (from --text, --markdown, or --body-file) is
			scanned for the current GC_TOKEN/GITCODE_TOKEN value before sending.
		`),
		Example: heredoc.Doc(`
			# Notify yourself (bot -> you, no group needed)
			$ gc lark send --to-self --text "deploy done"

			# Send plain text to a specific group
			$ gc lark send --chat-id oc_xxx --text "deploy done"

			# Send a direct message to a person by open_id
			$ gc lark send --user-id ou_xxx --as bot --text "hi"

			# Use the configured default group
			$ gc lark config set --default-chat oc_xxx
			$ gc lark send --text "build green"

			# Send markdown
			$ gc lark send --chat-id oc_xxx --markdown "## Release\n- v1.2.3 shipped"

			# Pipe a notification body from stdin
			$ echo "ci passed" | gc lark send --to-self --body-file -

			# Preview the lark-cli call without sending
			$ gc lark send --to-self --text "hi" --dry-run
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return sendRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.ChatID, "chat-id", "", "Feishu/Lark chat id (oc_xxx); send to a group/chat")
	cmd.Flags().StringVar(&opts.UserID, "user-id", "", "Feishu/Lark user open_id (ou_xxx); send a direct message")
	cmd.Flags().BoolVar(&opts.ToSelf, "to-self", false, "Send a direct message to the current lark-cli user (resolves open_id from auth status)")
	cmd.Flags().StringVar(&opts.Text, "text", "", "Plain text message")
	cmd.Flags().StringVar(&opts.Markdown, "markdown", "", "Markdown message")
	cmd.Flags().StringVar(&opts.File, "file", "", "File path to send as an attachment (forwarded to lark-cli)")
	cmd.Flags().StringVar(&opts.BodyFile, "body-file", "", "Read message text from file (use - for stdin)")
	cmd.Flags().StringVar(&opts.As, "as", "", "Identity: bot | user (default: lark-cli default, or bot when --to-self)")
	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Preview the lark-cli call without sending")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	cmd.MarkFlagsMutuallyExclusive("chat-id", "user-id", "to-self")

	return cmd
}

// sendTarget is the resolved destination: exactly one of chatID/userID is set.
type sendTarget struct {
	chatID string
	userID string
}

func sendRun(opts *sendOptions) error {
	if opts.LarkRun == nil {
		opts.LarkRun = larkcli.DefaultRun
	}

	target, err := resolveTarget(opts)
	if err != nil {
		return err
	}

	content, contentFlag, err := resolveContent(opts)
	if err != nil {
		return err
	}

	if contentFlag != "file" && content != "" {
		if err := cmdutil.ScanContentForSecrets(content); err != nil {
			return err
		}
	}

	as := opts.As
	if opts.ToSelf && as == "" {
		as = "bot"
	}

	larkArgs := buildSendArgs(target, contentFlag, content, as)

	if opts.DryRun {
		return printDryRun(opts, target, larkArgs, content)
	}

	res, err := larkcli.JSONResult(opts.LarkRun, larkArgs)
	if err != nil {
		return err
	}

	cs := opts.IO.ColorScheme()

	if res.ExitCode != 0 {
		return classifyLarkSendError(res)
	}

	env, _ := parseEnvelope(res.Stdout)
	if !env.OK {
		var msg string
		if env.Error != nil && env.Error.Message != "" {
			msg = env.Error.Message
		} else {
			msg = strings.TrimSpace(string(res.Stdout))
		}
		if msg == "" {
			msg = "unexpected empty response"
		}
		return classifyAuthError(msg)
	}

	if opts.JSON {
		chatID := env.chatID()
		if chatID == "" {
			chatID = target.chatID
		}
		return cmdutil.WriteJSON(opts.IO.Out, map[string]interface{}{
			"ok":       true,
			"chat_id":  chatID,
			"user_id":  target.userID,
			"identity": env.Identity,
			"data":     env.Data,
			"meta":     env.Meta,
		})
	}

	fmt.Fprintf(opts.IO.Out, "%s Sent to %s\n", cs.Green("✓"), target.describe())
	return nil
}

// resolveTarget picks the single target from --chat-id/--user-id/--to-self,
// falling back to the configured default chat when none is given.
func resolveTarget(opts *sendOptions) (sendTarget, error) {
	if opts.ToSelf {
		openID, err := resolveSelfOpenID(opts.LarkRun)
		if err != nil {
			return sendTarget{}, fmt.Errorf("failed to resolve current user: %w", err)
		}
		if openID == "" {
			return sendTarget{}, fmt.Errorf("could not determine current lark-cli user open_id; run: lark-cli auth login")
		}
		return sendTarget{userID: openID}, nil
	}
	if opts.UserID != "" {
		return sendTarget{userID: strings.TrimSpace(opts.UserID)}, nil
	}
	if opts.ChatID != "" {
		return sendTarget{chatID: strings.TrimSpace(opts.ChatID)}, nil
	}
	if dc := larkcli.DefaultChatID(); dc != "" {
		return sendTarget{chatID: dc}, nil
	}
	return sendTarget{}, cmdutil.NewUsageError("a target is required. Use --to-self, --user-id, or --chat-id (or configure a default: gc lark config set --default-chat <oc_xxx>)")
}

// resolveContent picks the single content source and returns (content, flagName).
// For --file, content is empty (the path is forwarded) and flagName is "file".
func resolveContent(opts *sendOptions) (string, string, error) {
	set := countSet(opts.Text != "", opts.Markdown != "", opts.File != "", opts.BodyFile != "")
	if set == 0 {
		return "", "", cmdutil.NewUsageError("one of --text, --markdown, --file, or --body-file is required")
	}
	if set > 1 {
		return "", "", cmdutil.NewUsageError("only one of --text, --markdown, --file, --body-file may be set")
	}

	switch {
	case opts.BodyFile != "":
		text, err := readBodyFile(opts.BodyFile, opts.IO.In)
		if err != nil {
			return "", "", err
		}
		return text, "text", nil
	case opts.Text != "":
		return opts.Text, "text", nil
	case opts.Markdown != "":
		return opts.Markdown, "markdown", nil
	case opts.File != "":
		return opts.File, "file", nil
	}
	return "", "", cmdutil.NewUsageError("no content source set")
}

// readBodyFile resolves message text from --body-file, reusing cmdutil helpers
// (which apply BOM/encoding handling, the Windows PowerShell lossy-stdin guard,
// and secret scanning). sendRun also scans the resolved content as defense in
// depth; the redundant pass is harmless (a substring containment check).
func readBodyFile(flagValue string, stdin io.Reader) (string, error) {
	if flagValue == "-" {
		text, err := cmdutil.ReadTextFromFlag(stdin, "--body-file")
		if err != nil {
			return "", fmt.Errorf("failed to read body from stdin: %w", err)
		}
		return text, nil
	}
	text, err := cmdutil.ReadTextFile(flagValue)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", flagValue, err)
	}
	return text, nil
}

func buildSendArgs(target sendTarget, contentFlag, content, as string) []string {
	args := []string{"im", "+messages-send"}
	if target.userID != "" {
		args = append(args, "--user-id", target.userID)
	} else {
		args = append(args, "--chat-id", target.chatID)
	}
	if as != "" {
		args = append(args, "--as", as)
	}
	switch contentFlag {
	case "text":
		args = append(args, "--text", content)
	case "markdown":
		args = append(args, "--markdown", content)
	case "file":
		args = append(args, "--file", content)
	}
	return args
}

func printDryRun(opts *sendOptions, target sendTarget, larkArgs []string, content string) error {
	displayCmd := "lark-cli " + joinShellArgs(append(larkArgs, "--json"))
	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, map[string]interface{}{
			"dry_run":  true,
			"target":   target.describe(),
			"command":  displayCmd,
			"identity": opts.As,
			"preview":  content,
		})
	}
	fmt.Fprintf(opts.IO.Out, "Dry run: would invoke lark-cli:\n  %s\n", displayCmd)
	if content != "" {
		preview := content
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		fmt.Fprintf(opts.IO.Out, "  message: %s\n", preview)
	}
	return nil
}

// authStatusPayload models the relevant fields of "lark-cli auth status --json"
// for resolving the current user's open_id (--to-self). The open id field is
// emitted as "openId" by lark-cli (real-command verified); "open_id" is also
// accepted defensively in case the envelope casing changes.
type authStatusPayload struct {
	Identities struct {
		User struct {
			OpenID    string `json:"openId"`
			OpenIDAlt string `json:"open_id"`
		} `json:"user"`
	} `json:"identities"`
}

func resolveSelfOpenID(run larkcli.RunFunc) (string, error) {
	res, err := larkcli.JSONResult(run, []string{"auth", "status"})
	if err != nil {
		return "", err
	}
	if res.ExitCode != 0 {
		// auth status itself failed → login/credential problem → exit 4.
		return "", cmdutil.NewAuthError(fmt.Sprintf("lark-cli auth status failed (exit %d); run: lark-cli auth login", res.ExitCode))
	}
	var as authStatusPayload
	if err := json.Unmarshal(res.Stdout, &as); err != nil {
		return "", fmt.Errorf("failed to parse lark-cli auth status: %w", err)
	}
	if id := as.Identities.User.OpenID; id != "" {
		return id, nil
	}
	return as.Identities.User.OpenIDAlt, nil
}

func extractErrorMessage(stderr []byte) string {
	if len(stderr) == 0 {
		return ""
	}
	if env, err := parseEnvelope(stderr); err == nil && env.Error != nil {
		return env.Error.Message
	}
	return ""
}

// classifyLarkSendError builds a human-readable error from a non-zero lark-cli
// send result and upgrades it to an auth/permission error (exit 4) when the
// message indicates a login, scope, or permission problem.
func classifyLarkSendError(res *larkcli.Result) error {
	msg := extractErrorMessage(res.Stderr)
	if msg == "" {
		msg = strings.TrimSpace(string(res.Stderr))
	}
	if msg == "" {
		msg = fmt.Sprintf("lark-cli send failed (exit %d)", res.ExitCode)
	} else {
		msg = fmt.Sprintf("lark-cli send failed (exit %d): %s", res.ExitCode, msg)
	}
	return classifyAuthError(msg)
}

// classifyAuthError returns an auth error (exit code 4) when the message reads
// like a login/scope/permission failure, otherwise a generic error (exit 1).
// Matching is heuristic on the message because lark-cli reports error.type as
// "api" for most failures, including auth-related ones.
func classifyAuthError(msg string) error {
	lower := strings.ToLower(msg)
	for _, ind := range []string{"scope", "auth", "login", "logged in", "authority", "permission", "availability", "out of the chat"} {
		if strings.Contains(lower, ind) {
			return cmdutil.NewAuthError(msg)
		}
	}
	return fmt.Errorf("%s", msg)
}

func countSet(flags ...bool) int {
	n := 0
	for _, f := range flags {
		if f {
			n++
		}
	}
	return n
}

// describe returns a human-readable target label.
func (t sendTarget) describe() string {
	if t.userID != "" {
		return "user " + t.userID
	}
	return "chat " + t.chatID
}

// chatID returns the chat_id from the lark-cli response envelope if present,
// else the resolved target chat id (for P2P the envelope carries the chat_id).
func (e larkEnvelope) chatID() string {
	if len(e.Data) > 0 {
		var d struct {
			ChatID string `json:"chat_id"`
		}
		if json.Unmarshal(e.Data, &d) == nil && d.ChatID != "" {
			return d.ChatID
		}
	}
	return ""
}

// joinShellArgs joins args into a single shell-readable line, quoting any
// argument that contains whitespace or shell metacharacters. It is for
// dry-run DISPLAY only; the real exec passes args as separate elements.
func joinShellArgs(args []string) string {
	var b strings.Builder
	for i, a := range args {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(shellQuote(a))
	}
	return b.String()
}

func shellQuote(s string) string {
	if s == "" {
		return `""`
	}
	if !strings.ContainsAny(s, " \t\n\"'\\$`") {
		return s
	}
	escaped := strings.ReplaceAll(s, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `"`, `\"`)
	return `"` + escaped + `"`
}
