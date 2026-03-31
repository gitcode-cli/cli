package cmdutil

import "os"

// EnvToken returns the active token from supported environment variables.
func EnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITCODE_TOKEN")
}
