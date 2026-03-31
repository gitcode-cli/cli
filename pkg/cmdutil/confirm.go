package cmdutil

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// ConfirmOptions controls destructive-action confirmation prompts.
type ConfirmOptions struct {
	IO       *iostreams.IOStreams
	Yes      bool
	Expected string
	Prompt   string
}

// ConfirmOrAbort enforces a shared confirmation flow.
func ConfirmOrAbort(opts ConfirmOptions) error {
	if opts.Yes {
		return nil
	}
	if opts.Expected == "" {
		return NewUsageError("missing confirmation target")
	}
	if opts.IO == nil || !opts.IO.CanPrompt() {
		return NewUsageError("confirmation required in non-interactive mode; rerun with --yes")
	}

	if opts.Prompt != "" {
		fmt.Fprint(opts.IO.ErrOut, opts.Prompt)
	}
	reader := bufio.NewReader(opts.IO.In)
	input, err := reader.ReadString('\n')
	if err != nil && strings.TrimSpace(input) == "" {
		if err == io.EOF {
			return NewUsageError("confirmation required in non-interactive mode; rerun with --yes")
		}
		return fmt.Errorf("failed to read confirmation: %w", err)
	}
	if strings.TrimSpace(input) != opts.Expected {
		return NewUsageError("confirmation did not match expected value")
	}
	return nil
}
