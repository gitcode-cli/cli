package lark

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/larkcli"
)

type doctorOptions struct {
	IO      *iostreams.IOStreams
	LarkRun larkcli.RunFunc
	JSON    bool
}

func newCmdDoctor(f *cmdutil.Factory, runF func(*doctorOptions) error) *cobra.Command {
	opts := &doctorOptions{
		IO:      f.IOStreams,
		LarkRun: larkcli.DefaultRun,
	}

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check lark-cli installation and login readiness",
		Long: heredoc.Doc(`
			Run a quick health check of the lark-cli integration: whether
			lark-cli is installed and whether the current Feishu login is ready.
			Use this before automating Feishu notifications.
		`),
		Example: heredoc.Doc(`
			# Diagnose lark-cli setup
			$ gc lark doctor

			# JSON output for scripts/agents
			$ gc lark doctor --json
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return doctorRun(opts)
		},
	}

	cmdutil.AddJSONFlag(cmd, &opts.JSON)
	return cmd
}

func doctorRun(opts *doctorOptions) error {
	if opts.LarkRun == nil {
		opts.LarkRun = larkcli.DefaultRun
	}

	bin := larkcli.FindLarkCLI()
	installed := bin != ""

	// Login readiness: delegate to lark-cli auth status. A missing binary or
	// non-zero exit means not ready. We never surface raw credentials.
	loginReady := false
	loginDetail := ""
	if installed {
		res, err := larkcli.JSONResult(opts.LarkRun, []string{"auth", "status"})
		if err == nil && res.ExitCode == 0 {
			if env, perr := parseEnvelope(res.Stderr); perr == nil && env.Identity != "" {
				loginReady = true
				loginDetail = env.Identity
			} else if env, perr := parseEnvelope(res.Stdout); perr == nil && env.Identity != "" {
				loginReady = true
				loginDetail = env.Identity
			}
		} else if err == nil && res.ExitCode != 0 {
			if msg := extractErrorMessage(res.Stderr); msg != "" {
				loginDetail = msg
			}
		}
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, map[string]interface{}{
			"installed":    installed,
			"bin":          bin,
			"login_ready":  loginReady,
			"identity":     loginDetail,
			"default_chat": larkcli.DefaultChatID(),
		})
	}

	cs := opts.IO.ColorScheme()
	if installed {
		fmt.Fprintf(opts.IO.Out, "%s lark-cli installed: %s\n", cs.Green("✓"), bin)
	} else {
		fmt.Fprintf(opts.IO.Out, "%s lark-cli not installed; run: gc lark install\n", cs.Red("✗"))
	}
	if installed {
		if loginReady {
			fmt.Fprintf(opts.IO.Out, "%s login ready (identity: %s)\n", cs.Green("✓"), loginDetail)
		} else {
			fmt.Fprintf(opts.IO.Out, "%s login not ready; run: lark-cli config init && lark-cli auth login\n", cs.Yellow("!"))
		}
	}
	if dc := larkcli.DefaultChatID(); dc != "" {
		fmt.Fprintf(opts.IO.Out, "%s default chat: %s\n", cs.Blue("i"), dc)
	} else {
		fmt.Fprintf(opts.IO.Out, "%s no default chat configured (set with: gc lark config set --default-chat <oc_xxx>)\n", cs.Blue("i"))
	}
	return nil
}
