package cmdutil

import (
	"errors"
	"fmt"

	"gitcode.com/gitcode-cli/cli/api"
)

const (
	ExitSuccess  = 0
	ExitError    = 1
	ExitUsage    = 2
	ExitNotFound = 3
	ExitAuth     = 4
	ExitConflict = 5
)

// CLIError represents a stable CLI-facing error with an exit code.
type CLIError struct {
	Code    int
	Message string
	Cause   error
}

func (e *CLIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message == "" && e.Cause != nil {
		return e.Cause.Error()
	}
	if e.Message == "" {
		return "unknown error"
	}
	if e.Cause == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Cause)
}

func (e *CLIError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// NewCLIError creates a CLIError with a stable exit code.
func NewCLIError(code int, message string, cause error) error {
	return &CLIError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

func NewUsageError(message string) error {
	return NewCLIError(ExitUsage, message, nil)
}

func NewAuthError(message string) error {
	return NewCLIError(ExitAuth, message, nil)
}

func NewNotFoundError(message string, cause error) error {
	return NewCLIError(ExitNotFound, message, cause)
}

// ExitCode maps a command error to a stable process exit code.
func ExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}

	var cliErr *CLIError
	if errors.As(err, &cliErr) {
		if cliErr.Code > 0 {
			return cliErr.Code
		}
	}

	var apiErr *api.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case 400:
			return ExitUsage
		case 401, 403:
			return ExitAuth
		case 404:
			return ExitNotFound
		case 409:
			return ExitConflict
		default:
			return ExitError
		}
	}

	return ExitError
}
