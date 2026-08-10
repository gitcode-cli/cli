package config

import (
	"bytes"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestSetAndGetUpdateMode(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	f := cmdutil.TestFactory()

	set := NewCmdConfig(f)
	setOut := &bytes.Buffer{}
	set.SetOut(setOut)
	set.SetArgs([]string{"set", "update.mode", "notify"})
	if err := set.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(setOut.String(), "update.mode=notify") {
		t.Fatalf("unexpected set output: %s", setOut.String())
	}

	get := NewCmdConfig(f)
	getOut := &bytes.Buffer{}
	get.SetOut(getOut)
	get.SetArgs([]string{"get", "update.mode"})
	if err := get.Execute(); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(getOut.String()) != "notify" {
		t.Fatalf("get output = %q, want notify", getOut.String())
	}
}

func TestUpdateModeDefaultsToAuto(t *testing.T) {
	t.Setenv("GC_CONFIG_DIR", t.TempDir())
	cmd := NewCmdConfig(cmdutil.TestFactory())
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"get", "update.mode"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(out.String()) != "auto" {
		t.Fatalf("output = %q, want auto", out.String())
	}
}

func TestRejectsInvalidUpdateMode(t *testing.T) {
	cmd := NewCmdConfig(cmdutil.TestFactory())
	cmd.SetArgs([]string{"set", "update.mode", "sometimes"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "auto, notify, or off") {
		t.Fatalf("error = %v", err)
	}
}
