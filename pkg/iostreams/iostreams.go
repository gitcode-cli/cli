// Package iostreams provides input/output stream management
package iostreams

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

// IOStreams holds the standard input, output, and error streams
type IOStreams struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer

	colorEnabled bool
	isTerminal   func(io.Writer) bool
	pager        string
}

// System returns IOStreams connected to standard input, output, and error
func System() *IOStreams {
	return &IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,

		colorEnabled: !isColorDisabled(),
		isTerminal:   isTerminal,
		pager:        os.Getenv("PAGER"),
	}
}

// ColorEnabled returns true if color output is enabled
func (s *IOStreams) ColorEnabled() bool {
	return s.colorEnabled && s.isTerminal(s.Out)
}

// ColorScheme returns a ColorScheme for output
func (s *IOStreams) ColorScheme() *ColorScheme {
	if !s.ColorEnabled() {
		return &ColorScheme{noColor: true}
	}
	return &ColorScheme{noColor: false}
}

// SetPager sets the pager command
func (s *IOStreams) SetPager(pager string) {
	s.pager = pager
}

// StartPager starts a pager for output
func (s *IOStreams) StartPager() error {
	if s.pager == "" {
		return nil
	}

	pagerCmd := exec.Command("sh", "-c", s.pager)
	pagerCmd.Stdout = os.Stdout
	pagerCmd.Stderr = os.Stderr

	stdin, err := pagerCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := pagerCmd.Start(); err != nil {
		return err
	}

	s.Out = stdin
	return nil
}

// StopPager stops the pager
func (s *IOStreams) StopPager() error {
	if closer, ok := s.Out.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func isColorDisabled() bool {
	return os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb"
}

func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		fi, err := f.Stat()
		if err != nil {
			return false
		}
		return (fi.Mode() & os.ModeCharDevice) != 0
	}
	return false
}

// IsStdinTTY returns true if stdin is a terminal
func (s *IOStreams) IsStdinTTY() bool {
	if f, ok := s.In.(*os.File); ok {
		fi, err := f.Stat()
		if err != nil {
			return false
		}
		return (fi.Mode() & os.ModeCharDevice) != 0
	}
	return false
}

// IsStdoutTTY returns true if stdout is a terminal
func (s *IOStreams) IsStdoutTTY() bool {
	return s.isTerminal(s.Out)
}

// IsStderrTTY returns true if stderr is a terminal
func (s *IOStreams) IsStderrTTY() bool {
	return s.isTerminal(s.ErrOut)
}

// IsInputTTY returns true if input is from a terminal
func IsInputTTY() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// CanPromptForInput returns true if we can prompt for input
func CanPromptForInput() bool {
	return IsInputTTY() && os.Getenv("GC_TEST_DISABLE_PROMPT") == ""
}

// CaptureOutput captures stdout and stderr for testing
func CaptureOutput(f func() error) (stdout, stderr string, err error) {
	oldOut := os.Stdout
	oldErr := os.Stderr

	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	err = f()

	wOut.Close()
	wErr.Close()
	os.Stdout = oldOut
	os.Stderr = oldErr

	var bufOut, bufErr bytes.Buffer
	io.Copy(&bufOut, rOut)
	io.Copy(&bufErr, rErr)

	return bufOut.String(), bufErr.String(), err
}