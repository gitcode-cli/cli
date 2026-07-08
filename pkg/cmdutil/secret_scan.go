package cmdutil

import (
	"fmt"
	"os"
	"strings"
)

// ErrSecretDetected is returned when ScanContentForSecrets detects a secret.
var ErrSecretDetected = fmt.Errorf("content contains secret")

// ScanContentForSecrets scans content for the current GC_TOKEN/GITCODE_TOKEN
// environment variable value. It is called before submitting user-provided
// body/comment content to the GitCode API to prevent AI agents (or humans)
// from accidentally leaking the current GitCode credential into issues, PRs,
// comments, or releases.
//
// GitCode tokens have no fixed prefix (unlike GitHub's ghp_), so pattern
// matching is not viable; the only reliable check is whether the current
// token value appears in the content. This is the highest-value defense
// against the common AI mistake of pasting $GC_TOKEN into an issue body.
//
// On detection, returns an error. The content is NOT modified (no redaction);
// submission is refused so the caller can surface the error before any API
// call is made.
func ScanContentForSecrets(content string) error {
	for _, env := range []string{"GC_TOKEN", "GITCODE_TOKEN"} {
		if v := os.Getenv(env); v != "" && strings.Contains(content, v) {
			return fmt.Errorf("%w: content contains the current %s value", ErrSecretDetected, env)
		}
	}
	return nil
}
