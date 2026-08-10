package api

import "regexp"

// tokenRedactPatterns matches sensitive values (tokens, Authorization headers)
// in log messages so they can be redacted before writing to stderr. This
// prevents credential leakage if debug logging is extended to include URLs
// or headers (#334).
var tokenRedactPatterns = []struct {
	re      *regexp.Regexp
	replace string
}{
	{regexp.MustCompile(`(?i)authorization:\s*bearer\s+\S+`), "Authorization: Bearer [REDACTED]"},
	{regexp.MustCompile(`(?i)authorization:\s*basic\s+\S+`), "Authorization: Basic [REDACTED]"},
	{regexp.MustCompile(`(?i)access_token=[^&\s]+`), "access_token=[REDACTED]"},
	{regexp.MustCompile(`(?i)token=[^&\s]+`), "token=[REDACTED]"},
	{regexp.MustCompile(`(?i)token:\s*\S+`), "token: [REDACTED]"},
}

// SanitizeLogMessage redacts sensitive values (tokens, Authorization headers)
// from log messages to prevent credential leakage when debug logging is
// extended to include URLs or headers (#334).
func SanitizeLogMessage(msg string) string {
	for _, p := range tokenRedactPatterns {
		msg = p.re.ReplaceAllString(msg, p.replace)
	}
	return msg
}
