package cmdutil

import (
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

func TestConfirmOrAbort_YesFlagSkipsConfirmation(t *testing.T) {
	opts := ConfirmOptions{
		Yes:      true,
		Expected: "123",
	}
	if err := ConfirmOrAbort(opts); err != nil {
		t.Fatalf("ConfirmOrAbort() with --yes = %v, want nil", err)
	}
}

func TestConfirmOrAbort_YesFlagWorksWithNilIO(t *testing.T) {
	opts := ConfirmOptions{
		IO:       nil,
		Yes:      true,
		Expected: "123",
	}
	if err := ConfirmOrAbort(opts); err != nil {
		t.Fatalf("ConfirmOrAbort() with --yes and nil IO = %v, want nil", err)
	}
}

func TestConfirmOrAbort_MissingExpectedValue(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := ConfirmOptions{
		IO:       io,
		Expected: "",
	}
	err := ConfirmOrAbort(opts)
	if err == nil {
		t.Fatal("ConfirmOrAbort() with empty Expected = nil, want error")
	}
	if _, ok := err.(*CLIError); !ok {
		t.Fatalf("ConfirmOrAbort() error type = %T, want *CLIError", err)
	}
	if ExitCode(err) != ExitUsage {
		t.Fatalf("ExitCode() = %d, want %d", ExitCode(err), ExitUsage)
	}
}

func TestConfirmOrAbort_NonInteractiveReturnsUsageError(t *testing.T) {
	// Create non-interactive IO (CanPrompt returns false)
	io, _, _, _ := iostreams.Test()
	opts := ConfirmOptions{
		IO:       io,
		Expected: "123",
		Prompt:   "Type the number to confirm: ",
	}
	err := ConfirmOrAbort(opts)
	if err == nil {
		t.Fatal("ConfirmOrAbort() in non-interactive mode = nil, want error")
	}
	if _, ok := err.(*CLIError); !ok {
		t.Fatalf("ConfirmOrAbort() error type = %T, want *CLIError", err)
	}
	if !strings.Contains(err.Error(), "--yes") {
		t.Fatalf("ConfirmOrAbort() error = %q, want mention of --yes", err.Error())
	}
}

func TestConfirmOrAbort_NilIOReturnsUsageError(t *testing.T) {
	opts := ConfirmOptions{
		IO:       nil,
		Expected: "123",
		Prompt:   "Type the number: ",
	}
	err := ConfirmOrAbort(opts)
	if err == nil {
		t.Fatal("ConfirmOrAbort() with nil IO = nil, want error")
	}
	if _, ok := err.(*CLIError); !ok {
		t.Fatalf("ConfirmOrAbort() error type = %T, want *CLIError", err)
	}
	if !strings.Contains(err.Error(), "--yes") {
		t.Fatalf("ConfirmOrAbort() error = %q, want mention of --yes", err.Error())
	}
}
