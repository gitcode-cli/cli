// Package browser provides utilities for opening URLs in a web browser
package browser

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Open opens a URL in the default web browser
func Open(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		// Use rundll32 to avoid cmd.exe metacharacter injection (CWE-78).
		// cmd /c start parses & | ^ % as command operators.
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			return fmt.Errorf("unsupported URL scheme: %s", url)
		}
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}
