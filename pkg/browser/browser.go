// Package browser provides utilities for opening URLs in a web browser
package browser

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
)

// Open opens a URL in the default web browser.
// Only http and https schemes are allowed to prevent command injection
// via custom URL scheme handlers (CWE-78).
func Open(rawURL string) error {
	if err := validateURLScheme(rawURL); err != nil {
		return err
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", rawURL)
	case "windows":
		// Use rundll32 to avoid cmd.exe metacharacter injection (CWE-78).
		// cmd /c start parses & | ^ % as command operators.
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", rawURL)
	case "darwin":
		cmd = exec.Command("open", rawURL)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

// validateURLScheme ensures the URL uses http or https scheme (case-insensitive).
func validateURLScheme(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %s", parsed.Scheme)
	}
	return nil
}
