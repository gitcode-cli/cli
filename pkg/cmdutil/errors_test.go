package cmdutil

import (
	"errors"
	"testing"

	"gitcode.com/gitcode-cli/cli/api"
)

func TestExitCode(t *testing.T) {
	t.Run("cli usage error", func(t *testing.T) {
		if got := ExitCode(NewUsageError("bad args")); got != ExitUsage {
			t.Fatalf("ExitCode() = %d, want %d", got, ExitUsage)
		}
	})

	t.Run("api auth error", func(t *testing.T) {
		err := &api.APIError{StatusCode: 401, Message: "unauthorized"}
		if got := ExitCode(err); got != ExitAuth {
			t.Fatalf("ExitCode() = %d, want %d", got, ExitAuth)
		}
	})

	t.Run("generic error", func(t *testing.T) {
		if got := ExitCode(errors.New("boom")); got != ExitError {
			t.Fatalf("ExitCode() = %d, want %d", got, ExitError)
		}
	})
}
