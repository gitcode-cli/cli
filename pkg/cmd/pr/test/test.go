// Package test implements the pr test command
package test

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type TestOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	Force bool // Force test pass (admin only)
}

// NewCmdTest creates the test command
func NewCmdTest(f *cmdutil.Factory, runF func(*TestOptions) error) *cobra.Command {
	opts := &TestOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "test <number>",
		Short: "Trigger or manage PR tests",
		Long: heredoc.Doc(`
			Trigger or manage tests for a pull request.

			With --force flag, admins can force a test pass.
		`),
		Example: heredoc.Doc(`
			# Trigger tests for a PR
			$ gc pr test 123

			# Force test pass (admin only)
			$ gc pr test 123 --force
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid PR number: %s", args[0])
			}
			opts.Number = number

			if runF != nil {
				return runF(opts)
			}
			return testRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Force test pass (admin only)")

	return cmd
}

func testRun(opts *TestOptions) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := getEnvToken()
	if token == "" {
		return fmt.Errorf("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "environment")

	// Get repository
	owner, repo, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	// Trigger or force test
	err = api.TestPR(client, owner, repo, opts.Number, &api.TestPROptions{
		Force: opts.Force,
	})
	if err != nil {
		return fmt.Errorf("failed to trigger PR test: %w", err)
	}

	if opts.Force {
		fmt.Fprintf(opts.IO.Out, "%s Forced test pass for PR #%d\n", cs.Green("✓"), opts.Number)
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Triggered tests for PR #%d\n", cs.Green("✓"), opts.Number)
	}
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	if token := os.Getenv("GITCODE_TOKEN"); token != "" {
		return token
	}
	return cmdutil.EnvToken()
}
