package larkcli

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"strings"
	"testing"
)

func TestInstaller_BootstrapCommandArgs(t *testing.T) {
	i := &Installer{}
	runner, args, err := i.bootstrapCommand()
	if err != nil {
		t.Fatalf("bootstrapCommand err = %v", err)
	}
	if runner == "" {
		t.Error("runner binary should not be empty")
	}
	wantPkg := "@larksuite/cli@latest"
	found := false
	for _, a := range args {
		if a == wantPkg {
			found = true
		}
	}
	if !found {
		t.Errorf("args = %v, want to contain %s", args, wantPkg)
	}
	if args[len(args)-1] != "install" {
		t.Errorf("last arg = %q, want install", args[len(args)-1])
	}
}

func TestInstaller_ExecFnInvoked(t *testing.T) {
	var gotName string
	var gotArgs []string
	i := &Installer{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		execFn: func(name string, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
			gotName = name
			gotArgs = args
			return nil
		},
	}
	if err := i.Install(); err != nil {
		t.Fatalf("Install err = %v", err)
	}
	if gotName == "" {
		t.Error("execFn not invoked with a runner name")
	}
	if len(gotArgs) == 0 {
		t.Error("execFn not invoked with args")
	}
}

func TestInstaller_ExitErrorWrapped(t *testing.T) {
	i := &Installer{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		execFn: func(name string, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
			return &exec.ExitError{}
		},
	}
	err := i.Install()
	if err == nil {
		t.Fatal("Install err = nil, want error for exit failure")
	}
	if !strings.Contains(err.Error(), "lark-cli installer exited") {
		t.Errorf("err = %v, want installer exit message", err)
	}
}

func TestInstaller_GenericErrorWrapped(t *testing.T) {
	i := &Installer{
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
		execFn: func(name string, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
			return errors.New("boom")
		},
	}
	err := i.Install()
	if err == nil {
		t.Fatal("Install err = nil, want error")
	}
	if !strings.Contains(err.Error(), "failed to install lark-cli") {
		t.Errorf("err = %v, want install failure message", err)
	}
}
