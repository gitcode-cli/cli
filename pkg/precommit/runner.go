// Package precommit detects and manages pre-commit configuration and environment.
package precommit

import "os/exec"

// CommandRunner runs external commands. It is abstracted so tests can avoid
// invoking real pre-commit / python binaries.
type CommandRunner interface {
	// Look reports whether an executable named name exists on PATH.
	Look(name string) bool
	// Run executes name with args. If dir != "", it is the working directory.
	// It returns the combined stdout+stderr output and any execution error.
	Run(dir, name string, args ...string) (string, error)
}

type execRunner struct{}

// NewExecRunner returns a CommandRunner backed by os/exec.
func NewExecRunner() CommandRunner { return execRunner{} }

func (execRunner) Look(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func (execRunner) Run(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}
