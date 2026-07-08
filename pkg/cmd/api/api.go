// Package api implements the generic api command.
package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type Options struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	Endpoint  string
	Method    string
	MethodSet bool
	Headers   []string
	Input     string
}

// NewCmdAPI creates the generic api command.
func NewCmdAPI(f *cmdutil.Factory, runF func(*Options) error) *cobra.Command {
	opts := &Options{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Method:     http.MethodGet,
	}

	cmd := &cobra.Command{
		Use:   "api <endpoint>",
		Short: "Make an authenticated GitCode REST API request",
		Long: heredoc.Doc(`
			Make an authenticated GitCode REST API request.

			The endpoint may be written with or without a leading slash. If it does
			not start with /api/, the CLI prefixes /api/v5 automatically.
		`),
		Example: heredoc.Doc(`
			# Get PR changed files
			$ gc api repos/owner/repo/pulls/123/files

			# Query commits for a file on a branch
			$ gc api 'repos/owner/repo/commits?path=src/main.go&sha=main'

			# Send a PATCH request with a JSON body
			$ gc api repos/owner/repo/pulls/123 --method PATCH --input body.json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Endpoint = args[0]
			opts.MethodSet = cmd.Flags().Changed("method")
			if runF != nil {
				return runF(opts)
			}
			return run(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Method, "method", "X", http.MethodGet, "HTTP method")
	cmd.Flags().StringArrayVarP(&opts.Headers, "header", "H", nil, "Additional HTTP header (Name: value)")
	cmd.Flags().StringVar(&opts.Input, "input", "", "Read request body from file (use - for stdin)")

	return cmd
}

func run(opts *Options) error {
	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}

	body, err := readInput(opts)
	if err != nil {
		return err
	}
	headers, err := parseHeaders(opts.Headers)
	if err != nil {
		return err
	}

	method := strings.ToUpper(strings.TrimSpace(opts.Method))
	if method == "" {
		method = http.MethodGet
	}
	if !opts.MethodSet && body != nil && method == http.MethodGet {
		method = http.MethodPost
	}

	resp, err := client.RawREST(method, opts.Endpoint, body, headers)
	if err != nil {
		return err
	}
	if len(resp.Body) == 0 {
		return nil
	}
	if _, err := opts.IO.Out.Write(resp.Body); err != nil {
		return err
	}
	if !bytes.HasSuffix(resp.Body, []byte("\n")) {
		_, err = fmt.Fprintln(opts.IO.Out)
	}
	return err
}

func readInput(opts *Options) (io.Reader, error) {
	if strings.TrimSpace(opts.Input) == "" {
		return nil, nil
	}
	if opts.Input == "-" {
		data, err := io.ReadAll(opts.IO.In)
		if err != nil {
			return nil, fmt.Errorf("failed to read stdin: %w", err)
		}
		return bytes.NewReader(data), nil
	}
	data, err := cmdutil.ReadInputFile(opts.Input)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %w", err)
	}
	return bytes.NewReader(data), nil
}

func parseHeaders(values []string) (map[string]string, error) {
	headers := map[string]string{}
	for _, raw := range values {
		name, value, ok := strings.Cut(raw, ":")
		if !ok || strings.TrimSpace(name) == "" {
			return nil, cmdutil.NewUsageError("--header must be in 'Name: value' format")
		}
		headers[strings.TrimSpace(name)] = strings.TrimSpace(value)
	}
	return headers, nil
}
