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
