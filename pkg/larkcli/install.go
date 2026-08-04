package larkcli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

// installPackage is the npm package spec used by the official installer.
const installPackage = "@larksuite/cli@latest"

// Installer runs the official lark-cli installer. It is split out so the
// command layer can inject a fake for tests.
type Installer struct {
	// Stdout and Stderr receive installer output.
	Stdout io.Writer
	Stderr io.Writer
	// Stdin is forwarded to the installer (interactive prompts).
	Stdin io.Reader
	// execFn is the execution seam; defaults to exec.Command(...).Run().
	execFn func(name string, args []string, stdin io.Reader, stdout, stderr io.Writer) error
}

// Install runs `npx @larksuite/cli@latest install`. On success lark-cli is
// available on PATH.
func (i *Installer) Install() error {
	if i.execFn == nil {
		i.execFn = defaultExec
	}
	runner, args, err := i.bootstrapCommand()
	if err != nil {
		return err
	}
	if err := i.execFn(runner, args, i.Stdin, i.Stdout, i.Stderr); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return fmt.Errorf("lark-cli installer exited with code %d: %w", exitErr.ExitCode(), err)
		}
		return fmt.Errorf("failed to install lark-cli: %w", err)
	}
	return nil
}

// bootstrapCommand returns the runner binary and args used to invoke the
// installer. On Windows it prefers the npm.cmd shim; elsewhere it uses npx
// directly. When npx/npm cannot be located, it falls back to "npx" and lets
// the OS report the missing-binary error to the user.
func (i *Installer) bootstrapCommand() (string, []string, error) {
	runner := "npx"
	if runtime.GOOS == "windows" {
		if p, err := exec.LookPath("npx.cmd"); err == nil {
			runner = p
		}
	} else {
		if p, err := exec.LookPath("npx"); err == nil {
			runner = p
		}
	}
	args := []string{"-y", installPackage, "install"}
	return runner, args, nil
}

func defaultExec(name string, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = sanitizedEnv(os.Environ())
	return cmd.Run()
}
