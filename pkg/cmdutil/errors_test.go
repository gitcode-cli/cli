package cmdutil

import (
	"errors"
	"testing"

	"gitcode.com/gitcode-cli/cli/api"
)

func TestExitCode(t *testing.T) {
	t.Run("success returns ExitSuccess", func(t *testing.T) {
		if got := ExitCode(nil); got != ExitSuccess {
			t.Fatalf("ExitCode(nil) = %d, want %d", got, ExitSuccess)
		}
	})

	t.Run("cli usage error", func(t *testing.T) {
		if got := ExitCode(NewUsageError("bad args")); got != ExitUsage {
			t.Fatalf("ExitCode() = %d, want %d", got, ExitUsage)
		}
	})

	t.Run("cli auth error", func(t *testing.T) {
		if got := ExitCode(NewAuthError("not authenticated")); got != ExitAuth {
			t.Fatalf("ExitCode() = %d, want %d", got, ExitAuth)
		}
	})

	t.Run("cli not found error", func(t *testing.T) {
		if got := ExitCode(NewNotFoundError("resource not found", nil)); got != ExitNotFound {
			t.Fatalf("ExitCode() = %d, want %d", got, ExitNotFound)
		}
	})

	t.Run("cli conflict error", func(t *testing.T) {
		if got := ExitCode(NewConflictError("conflict")); got != ExitConflict {
			t.Fatalf("ExitCode() = %d, want %d", got, ExitConflict)
		}
	})

	t.Run("api auth error", func(t *testing.T) {
		err := &api.APIError{StatusCode: 401, Message: "unauthorized"}
		if got := ExitCode(err); got != ExitAuth {
			t.Fatalf("ExitCode() = %d, want %d", got, ExitAuth)
		}
	})

	t.Run("api embedded not found error code", func(t *testing.T) {
		err := &api.APIError{StatusCode: 400, ErrorCode: 404, ErrorMessage: "404 Not Found Commit"}
		if got := ExitCode(err); got != ExitNotFound {
			t.Fatalf("ExitCode() = %d, want %d", got, ExitNotFound)
		}
	})

	t.Run("generic error", func(t *testing.T) {
		if got := ExitCode(errors.New("boom")); got != ExitError {
			t.Fatalf("ExitCode() = %d, want %d", got, ExitError)
		}
	})

	// Cobra usage errors should return ExitUsage (2)
	t.Run("cobra missing args error", func(t *testing.T) {
		err := errors.New("accepts 1 arg(s), received 0")
		if got := ExitCode(err); got != ExitUsage {
			t.Fatalf("ExitCode(cobra missing args) = %d, want %d", got, ExitUsage)
		}
	})

	t.Run("cobra too many args error", func(t *testing.T) {
		err := errors.New("accepts at most 1 arg(s), received 2")
		if got := ExitCode(err); got != ExitUsage {
			t.Fatalf("ExitCode(cobra too many args) = %d, want %d", got, ExitUsage)
		}
	})

	t.Run("cobra minimum args error", func(t *testing.T) {
		err := errors.New("requires at least 2 arg(s), only received 1")
		if got := ExitCode(err); got != ExitUsage {
			t.Fatalf("ExitCode(cobra minimum args) = %d, want %d", got, ExitUsage)
		}
	})

	t.Run("cobra required flag error", func(t *testing.T) {
		err := errors.New("required flag(s) \"target-repo\" not set")
		if got := ExitCode(err); got != ExitUsage {
			t.Fatalf("ExitCode(cobra required flag) = %d, want %d", got, ExitUsage)
		}
	})

	t.Run("cobra unknown command error", func(t *testing.T) {
		err := errors.New("unknown command \"foo\" for \"gc\"")
		if got := ExitCode(err); got != ExitUsage {
			t.Fatalf("ExitCode(cobra unknown command) = %d, want %d", got, ExitUsage)
		}
	})

	t.Run("cobra unknown flag error", func(t *testing.T) {
		err := errors.New("unknown flag: --unknown-flag")
		if got := ExitCode(err); got != ExitUsage {
			t.Fatalf("ExitCode(cobra unknown flag) = %d, want %d", got, ExitUsage)
		}
	})

	t.Run("cobra unknown shorthand flag error", func(t *testing.T) {
		err := errors.New("unknown shorthand flag: -f in `-f'")
		if got := ExitCode(err); got != ExitUsage {
			t.Fatalf("ExitCode(cobra unknown shorthand flag) = %d, want %d", got, ExitUsage)
		}
	})

	t.Run("cobra range args error", func(t *testing.T) {
		err := errors.New("accepts between 1 and 3 arg(s), received 4")
		if got := ExitCode(err); got != ExitUsage {
			t.Fatalf("ExitCode(cobra range args) = %d, want %d", got, ExitUsage)
		}
	})
}

func TestFormatAPIID(t *testing.T) {
	cases := map[string]struct {
		input interface{}
		want  string
	}{
		"string":  {input: "123", want: "123"},
		"float64": {input: 123.0, want: "123"},
		"int":     {input: 123, want: "123"},
		"nil":     {input: nil, want: ""},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := FormatAPIID(tc.input); got != tc.want {
				t.Fatalf("FormatAPIID(%v) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
