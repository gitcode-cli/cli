// Package list implements the repo list command
package list

import (
	"fmt"
	"net/http"
	"strings"

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

	// Flags
	Limit      int
	PerPage    int
	LimitSet   bool
	PerPageSet bool
	Visibility string
	Owner      string
	JSON       bool
	Format     string
}

// NewCmdList creates the list command
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List repositories",
		Long: heredoc.Doc(`
			List repositories for the authenticated user or an organization.
		`),
		Example: heredoc.Doc(`
			# List your repositories
			$ gc repo list

			# List with limit
			$ gc repo list --limit 50

			# List only public repos
			$ gc repo list --visibility public

			# Render as a table
			$ gc repo list --format table

			# Output as JSON
			$ gc repo list --json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			opts.LimitSet = cmd.Flags().Changed("limit")
			opts.PerPageSet = cmd.Flags().Changed("per-page")
			return listRun(opts)
		},
	}

	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of repos to list")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "API page size (default: --limit)")
	cmd.Flags().StringVarP(&opts.Visibility, "visibility", "v", "", "Filter by visibility (public/private)")
	cmdutil.SetFlagEnumOrWarn(cmd, "visibility", "public", "private")
	cmd.Flags().StringVarP(&opts.Owner, "owner", "o", "", "List repos for an organization")
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
	client, err := cmdutil.AuthenticatedClient(httpClient)
	if err != nil {
		return err
	}

	if opts.Limit <= 0 {
		return cmdutil.NewUsageError("--limit must be greater than 0")
	}
	if opts.PerPage < 0 {
		return cmdutil.NewUsageError("--per-page must be greater than or equal to 0")
	}

	// List repos
	repoOpts := &api.RepoListOptions{
		PerPage:    resolvePerPage(opts),
		Visibility: opts.Visibility,
	}
	var repos []api.Repository
	if owner := strings.TrimSpace(opts.Owner); owner != "" {
		repos, err = api.ListOrgRepos(client, owner, repoOpts)
	} else {
		repos, err = api.ListUserRepos(client, repoOpts)
	}
	if err != nil {
		return fmt.Errorf("failed to list repositories: %w", err)
	}
	repos = trimRepos(repos, opts)

	// Output
	if len(repos) == 0 {
		if format == output.FormatJSON {
			return cmdutil.WriteJSON(opts.IO.Out, repos)
		}
		fmt.Fprintf(opts.IO.Out, "No repositories found\n")
		return nil
	}

	if format == output.FormatJSON {
		return cmdutil.WriteJSON(opts.IO.Out, repos)
	}
	printer, err := output.NewRepoListPrinter(output.RepoListOptions{Format: format})
	if err != nil {
		return err
	}
	return printer.Print(opts.IO.Out, repos)
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

// resolvePerPage returns the API page size: the explicit --per-page when set,
// otherwise --limit. Mirrors the pr list / issue list pattern.
func resolvePerPage(opts *ListOptions) int {
	if opts.PerPageSet && opts.PerPage > 0 {
		return opts.PerPage
	}
	return opts.Limit
}

// trimRepos truncates to --limit when both --limit and --per-page were set.
func trimRepos(repos []api.Repository, opts *ListOptions) []api.Repository {
	if opts.PerPageSet && opts.LimitSet && len(repos) > opts.Limit {
		return repos[:opts.Limit]
	}
	return repos
}
