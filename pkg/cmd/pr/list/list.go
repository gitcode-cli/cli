// Package list implements the pr list command
package list

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/output"
)

type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string

	// Flags
	State     string
	Limit     int
	Head      string
	Base      string
	Sort      string
	Direction string
	Page      int
	JSON      bool
	Format    string
}

// NewCmdList creates the list command
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List pull requests",
		Long: heredoc.Doc(`
			List pull requests in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# List open PRs
			$ gc pr list -R owner/repo

			# List closed PRs
			$ gc pr list -R owner/repo --state closed

			# Filter by head and base branches
			$ gc pr list -R owner/repo --head feature/login --base main

			# Sort results
			$ gc pr list -R owner/repo --sort updated --direction desc

			# Render as a table
			$ gc pr list -R owner/repo --format table
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.State, "state", "s", "open", "Filter by state (open/closed/all)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of PRs to list")
	cmd.Flags().StringVarP(&opts.Head, "head", "H", "", "Filter by head branch")
	cmd.Flags().StringVarP(&opts.Base, "base", "B", "", "Filter by base branch")
	cmd.Flags().StringVar(&opts.Sort, "sort", "", "Sort by created/updated/popularity/long-running")
	cmd.Flags().StringVar(&opts.Direction, "direction", "", "Sort direction (asc/desc)")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number to fetch")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)
	cmdutil.AddFormatFlag(cmd, &opts.Format)

	return cmd
}

func listRun(opts *ListOptions) error {
	format, err := resolveOutputFormat(opts.JSON, opts.Format)
	if err != nil {
		return err
	}

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := cmdutil.EnvToken()
	if token == "" {
		return cmdutil.NewAuthError("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "environment")

	// Get repository
	owner, repo, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	// List PRs
	prs, err := api.ListPullRequests(client, owner, repo, &api.PRListOptions{
		State:     opts.State,
		Head:      opts.Head,
		Base:      opts.Base,
		Sort:      opts.Sort,
		Direction: opts.Direction,
		PerPage:   opts.Limit,
		Page:      opts.Page,
	})
	if err != nil {
		return fmt.Errorf("failed to list PRs: %w", err)
	}

	// Output
	if len(prs) == 0 {
		if format == output.FormatJSON {
			return cmdutil.WriteJSON(opts.IO.Out, prs)
		}
		fmt.Fprintf(opts.IO.Out, "No pull requests found\n")
		return nil
	}

	if format == output.FormatJSON {
		return cmdutil.WriteJSON(opts.IO.Out, prs)
	}
	printer, err := output.NewPRListPrinter(output.PRListOptions{
		Format: format,
		Color:  opts.IO.ColorScheme(),
	})
	if err != nil {
		return err
	}
	return printer.Print(opts.IO.Out, prs)
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

func resolveOutputFormat(jsonFlag bool, raw string) (output.Format, error) {
	format, err := output.ParseFormat(raw)
	if err != nil {
		return "", cmdutil.NewUsageError(err.Error())
	}
	if jsonFlag {
		if raw != "" && format != output.FormatJSON {
			return "", cmdutil.NewUsageError("--json cannot be combined with --format unless --format json")
		}
		return output.FormatJSON, nil
	}
	return format, nil
}
