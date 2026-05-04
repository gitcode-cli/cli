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
