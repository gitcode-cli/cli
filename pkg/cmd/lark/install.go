package lark

import (
	"fmt"
	"io"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/larkcli"
)

type installOptions struct {
	IO        *iostreams.IOStreams
	installFn func(stdout, stderr io.Writer) error
}

func newCmdInstall(f *cmdutil.Factory, runF func(*installOptions) error) *cobra.Command {
	opts := &installOptions{
		IO: f.IOStreams,
	}

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the official lark-cli tool",
		Long: heredoc.Doc(`
			Bootstrap the official Feishu/Lark CLI by running
			"npx @larksuite/cli@latest install". After install, run
			"gc lark auth status" (or "lark-cli config init" + "lark-cli auth login")
			to complete Feishu OAuth setup.
		`),
		Example: heredoc.Doc(`
			# Install lark-cli
			$ gc lark install
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.installFn = defaultInstall
			if runF != nil {
				return runF(opts)
			}
			return installRun(opts)
		},
	}

	return cmd
}

func installRun(opts *installOptions) error {
	cs := opts.IO.ColorScheme()

	fmt.Fprintf(opts.IO.ErrOut, "%s Installing lark-cli via npx @larksuite/cli@latest install...\n", cs.Blue("i"))
	if err := opts.installFn(opts.IO.Out, opts.IO.ErrOut); err != nil {
		return err
	}

	// Verify the install landed on PATH (npx wrapper may need a new shell).
	if bin := larkcli.FindLarkCLI(); bin == "" {
		fmt.Fprintf(opts.IO.ErrOut, "%s lark-cli installed but not found on PATH of the current shell. Open a new terminal or set %s to the binary path.\n", cs.Yellow("!"), larkcli.EnvBin)
	} else {
		fmt.Fprintf(opts.IO.Out, "%s lark-cli ready: %s\n", cs.Green("✓"), bin)
	}
	return nil
}

func defaultInstall(stdout, stderr io.Writer) error {
	installer := &larkcli.Installer{
		Stdout: stdout,
		Stderr: stderr,
		Stdin:  os.Stdin,
	}
	return installer.Install()
}
