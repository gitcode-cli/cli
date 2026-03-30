// Package main is the entry point for gitcode-cli
package main

import (
	"os"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/root"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Execute the root command
	if err := root.Execute(version, commit, date); err != nil {
		os.Stderr.WriteString("Error: " + err.Error() + "\n")
		os.Exit(cmdutil.ExitCode(err))
	}
}
