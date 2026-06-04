// Package check implements the `precommit check` command.
package check

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	gitpkg "gitcode.com/gitcode-cli/cli/git"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/precommit"
)

// CheckOptions holds dependencies and flags for the check command.
type CheckOptions struct {
	IO      *iostreams.IOStreams
	GitRoot func() (string, error)
	Runner  precommit.CommandRunner

	Run       bool
	NoInstall bool
	Yes       bool
	JSON      bool
}

// NewCmdCheck creates the `precommit check` command.
func NewCmdCheck(f *cmdutil.Factory, runF func(*CheckOptions) error) *cobra.Command {
	opts := &CheckOptions{
		IO:      f.IOStreams,
		GitRoot: gitpkg.RootDir,
		Runner:  precommit.NewExecRunner(),
	}

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check pre-commit configuration and local environment before committing",
		Long: heredoc.Doc(`
			Check whether the current repository configures pre-commit and whether the
			local environment is ready to run it before committing code.

			The command:
			  1. Detects a .pre-commit-config.yaml (or .yml) in the repository root.
			  2. Verifies the pre-commit tool is installed.
			  3. Verifies the git pre-commit hook is initialized.
			  4. Optionally runs the hooks with --run.

			When something is missing it auto-installs/initializes in an interactive
			terminal. In a non-interactive (non-TTY) environment, pass --yes to allow
			environment changes, or --no-install to only diagnose. Cross-platform
			(Windows, Linux x86/arm, macOS).
		`),
		Example: heredoc.Doc(`
			# Verify the environment is ready
			$ gc precommit check

			# Verify and actually run the hooks
			$ gc precommit check --run

			# Only diagnose, never modify the environment
			$ gc precommit check --no-install

			# Machine-readable output
			$ gc precommit check --json
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return checkRun(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Run, "run", false, "Run pre-commit hooks (pre-commit run --all-files) after verifying")
	cmd.Flags().BoolVar(&opts.NoInstall, "no-install", false, "Only diagnose; never install the tool or hook")
	cmd.Flags().BoolVarP(&opts.Yes, "yes", "y", false, "Allow environment changes (install/init) in non-interactive mode")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	// --no-install (never mutate) and --yes (authorize mutation) express
	// opposite intents; reject passing both rather than silently letting one win.
	cmd.MarkFlagsMutuallyExclusive("no-install", "yes")

	return cmd
}

func checkRun(opts *CheckOptions) error {
	root, err := opts.GitRoot()
	if err != nil {
		return cmdutil.NewCLIError(cmdutil.ExitError, "not in a git repository", err)
	}

	allowInstall := !opts.NoInstall && (opts.IO.CanPrompt() || opts.Yes)

	res, err := precommit.Check(opts.Runner, precommit.Options{
		Root:         root,
		AllowInstall: allowInstall,
		Run:          opts.Run,
	})
	if err != nil {
		return cmdutil.NewCLIError(cmdutil.ExitError, "pre-commit check failed", err)
	}

	if opts.JSON {
		if writeErr := cmdutil.WriteJSON(opts.IO.Out, res); writeErr != nil {
			return writeErr
		}
	} else {
		printResult(opts, res, allowInstall)
	}

	if !res.OK {
		return cmdutil.NewCLIError(cmdutil.ExitError, "pre-commit environment is not ready", nil)
	}
	return nil
}

func printResult(opts *CheckOptions, res precommit.Result, allowInstall bool) {
	cs := opts.IO.ColorScheme()
	out := opts.IO.Out

	if !res.ConfigFound {
		fmt.Fprintf(out, "%s No pre-commit configuration found; nothing to check.\n", cs.Green("✓"))
		return
	}

	mark := func(ok bool) string {
		if ok {
			return cs.Green("✓")
		}
		return cs.Red("✗")
	}

	fmt.Fprintf(out, "%s pre-commit configuration found\n", mark(res.ConfigFound))
	if res.ToolInstalled {
		fmt.Fprintf(out, "%s pre-commit tool installed (%s)\n", mark(true), res.ToolVersion)
	} else {
		fmt.Fprintf(out, "%s pre-commit tool not installed\n", mark(false))
	}
	fmt.Fprintf(out, "%s git hook initialized\n", mark(res.HookInstalled))

	for _, a := range res.ActionsTaken {
		fmt.Fprintf(out, "  - %s\n", a)
	}

	switch res.RunResult {
	case "passed":
		fmt.Fprintf(out, "%s pre-commit run passed\n", mark(true))
	case "failed":
		fmt.Fprintf(out, "%s pre-commit run failed\n", mark(false))
	default:
		// No --run was requested; nothing to report for the run step.
	}

	if res.OK {
		return
	}

	// The environment is ready but the hooks themselves failed: that is a
	// check failure, not an unready environment. Keep the two cases distinct so
	// users don't chase the wrong problem.
	envReady := res.ConfigFound && res.ToolInstalled && res.HookInstalled
	if envReady && res.RunResult == "failed" {
		fmt.Fprintf(opts.IO.ErrOut, "\npre-commit checks failed.\n")
		if res.RunOutput != "" {
			fmt.Fprintf(opts.IO.ErrOut, "%s\n", res.RunOutput)
		}
		return
	}

	fmt.Fprintf(opts.IO.ErrOut, "\nEnvironment not ready.\n")
	// Auto-install was wanted (not --no-install) but couldn't run: non-TTY without --yes.
	if !allowInstall && !opts.NoInstall && (!res.ToolInstalled || !res.HookInstalled) {
		fmt.Fprintf(opts.IO.ErrOut, "Re-run in a terminal, or pass --yes to auto-install/initialize.\n")
	}
	if !res.ToolInstalled {
		fmt.Fprintf(opts.IO.ErrOut, "Install pre-commit, e.g.: pipx install pre-commit (or pip install --user pre-commit).\n")
	} else if !res.HookInstalled {
		fmt.Fprintf(opts.IO.ErrOut, "Initialize hooks: pre-commit install\n")
	}
}
