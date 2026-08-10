package cmdutil

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// ConfirmTokenDisclosure gates commands that print a full authentication token.
func ConfirmTokenDisclosure(ioStreams *iostreams.IOStreams, hostname string) error {
	if ioStreams == nil || !ioStreams.CanPrompt() {
		return NewUsageError("printing authentication token requires interactive confirmation")
	}

	fmt.Fprintf(ioStreams.ErrOut, "This will print the full authentication token for %s.\n", hostname)
	fmt.Fprintf(ioStreams.ErrOut, "Type %s to continue: ", hostname)

	reader := bufio.NewReader(ioStreams.In)
	input, err := reader.ReadString('\n')
	if err != nil {
		// Surface read errors even when partial input was received, instead of
		// falling through to a misleading "did not match" (#361).
		if err == io.EOF && strings.TrimSpace(input) == "" {
			return NewUsageError("printing authentication token requires interactive confirmation")
		}
		return NewCLIError(ExitUsage, "failed to read confirmation", err)
	}
	if strings.TrimSpace(input) != hostname {
		return NewUsageError("confirmation did not match expected hostname")
	}
	return nil
}
