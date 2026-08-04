package lark

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/larkcli"
)

// newCmdAuth creates the "gc lark auth" command group.
func newCmdAuth(f *cmdutil.Factory, runF func(*authStatusOptions) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Show lark-cli login status",
		Long: heredoc.Doc(`
			Inspect the current lark-cli login status and granted scopes by
			delegating to "lark-cli auth status --json". Feishu OAuth setup
			(config init / auth login) is performed directly with lark-cli and
			keeps credentials in the OS keychain.
		`),
		Args: cobra.NoArgs,
	}

	cmd.AddCommand(newCmdAuthStatus(f, runF))
	return cmd
}

type authStatusOptions struct {
	IO      *iostreams.IOStreams
	LarkRun larkcli.RunFunc
	JSON    bool
}

func newCmdAuthStatus(f *cmdutil.Factory, runF func(*authStatusOptions) error) *cobra.Command {
	opts := &authStatusOptions{
		IO:      f.IOStreams,
		LarkRun: larkcli.DefaultRun,
	}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show lark-cli login status",
		Example: heredoc.Doc(`
			# Show login status
			$ gc lark auth status

			# JSON output for scripts/agents
			$ gc lark auth status --json
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return authStatusRun(opts)
		},
	}

	cmdutil.AddJSONFlag(cmd, &opts.JSON)
	return cmd
}

func authStatusRun(opts *authStatusOptions) error {
	if opts.LarkRun == nil {
		opts.LarkRun = larkcli.DefaultRun
	}

	res, err := larkcli.JSONResult(opts.LarkRun, []string{"auth", "status"})
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		if msg := extractErrorMessage(res.Stderr); msg != "" {
			return fmt.Errorf("lark-cli auth status failed (exit %d): %s", res.ExitCode, msg)
		}
		return fmt.Errorf("lark-cli auth status failed (exit %d): %s", res.ExitCode, strings.TrimSpace(string(res.Stderr)))
	}

	if opts.JSON {
		// Pass lark-cli's structured status through unchanged.
		_, _ = fmt.Fprint(opts.IO.Out, string(res.Stdout))
		return nil
	}

	env, _ := parseEnvelope(res.Stderr)
	if env.Identity == "" {
		env, _ = parseEnvelope(res.Stdout)
	}
	cs := opts.IO.ColorScheme()
	if env.Identity != "" {
		fmt.Fprintf(opts.IO.Out, "%s lark-cli logged in as %s\n", cs.Green("✓"), env.Identity)
	} else {
		fmt.Fprintf(opts.IO.Out, "%s lark-cli login status unknown; raw output:\n", cs.Yellow("!"))
		fmt.Fprint(opts.IO.Out, string(res.Stdout))
	}
	return nil
}
