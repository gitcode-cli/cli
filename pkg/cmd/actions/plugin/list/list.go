// Package list implements the actions plugin list command.
package list

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/output"
)

// ListOptions configures the actions plugin list command.
type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string

	Limit      int
	Page       int
	Paginate   bool
	PerPage    int
	LimitSet   bool
	PerPageSet bool

	JSON   bool
	Format string
}

// NewCmdList creates the actions plugin list command.
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List official Actions plugins",
		Long: heredoc.Doc(`
			List official GitCode Actions plugins available for a repository.

			By default all pages are fetched automatically. Use --limit,
			--page, or --per-page to control pagination.
		`),
		Example: heredoc.Doc(`
			# List plugins for a repository
			$ gc actions plugin list -R owner/repo

			# Fetch a single page
			$ gc actions plugin list -R owner/repo --page 1 --per-page 20

			# Output as JSON
			$ gc actions plugin list -R owner/repo --json
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

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 0, "Maximum number of plugins to list (0 = no limit, fetch all)")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number to fetch")
	cmd.Flags().BoolVar(&opts.Paginate, "paginate", false, "Fetch all pages")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "API page size (default: --limit, or 100 with --paginate)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)
	cmdutil.AddFormatFlag(cmd, &opts.Format)

	return cmd
}

func listRun(opts *ListOptions) error {
	if err := validateListFlags(opts); err != nil {
		return err
	}

	format, err := resolveOutputFormat(opts.JSON, opts.Format)
	if err != nil {
		return err
	}

	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}

	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}
	owner, repo, err := cmdutil.ParseRepo(repository)
	if err != nil {
		return err
	}

	project := owner + "/" + repo
	rawEntries, err := fetchPlugins(client, project, opts)
	if err != nil {
		return fmt.Errorf("failed to list actions plugins: %w", err)
	}

	if len(rawEntries) == 0 {
		if format == output.FormatJSON {
			return cmdutil.WriteJSON(opts.IO.Out, rawEntries)
		}
		fmt.Fprintf(opts.IO.Out, "No actions plugins found\n")
		return nil
	}

	if format == output.FormatJSON {
		return cmdutil.WriteJSON(opts.IO.Out, rawEntries)
	}

	return printPluginsTable(opts.IO, rawEntries)
}

func validateListFlags(opts *ListOptions) error {
	if opts.Limit < 0 {
		return cmdutil.NewUsageError("--limit must be greater than or equal to 0")
	}
	if opts.Page < 0 {
		return cmdutil.NewUsageError("--page must be greater than or equal to 0")
	}
	if opts.PerPage < 0 {
		return cmdutil.NewUsageError("--per-page must be greater than or equal to 0")
	}
	if opts.Paginate && opts.Page > 0 {
		return cmdutil.NewUsageError("--paginate cannot be combined with --page")
	}
	return nil
}

func fetchPlugins(client *api.Client, project string, opts *ListOptions) ([]json.RawMessage, error) {
	if opts.Page > 0 && !opts.Paginate {
		return fetchPluginsPage(client, project, opts)
	}
	return fetchAllPlugins(client, project, opts)
}

func fetchPluginsPage(client *api.Client, project string, opts *ListOptions) ([]json.RawMessage, error) {
	perPage := resolvePerPage(opts)
	raw, err := api.ListActionsPlugins(client, project, &api.ActionsListPluginsOptions{
		PerPage: perPage,
		Page:    opts.Page,
	})
	if err != nil {
		return nil, err
	}
	entries, err := api.ParseActionsPluginsListRaw(raw)
	if err != nil {
		return nil, err
	}
	return trimEntries(entries, opts), nil
}

func fetchAllPlugins(client *api.Client, project string, opts *ListOptions) ([]json.RawMessage, error) {
	perPage := resolvePerPage(opts)
	if perPage == 0 {
		perPage = 100
	}
	var all []json.RawMessage
	for page := 1; ; page++ {
		raw, err := api.ListActionsPlugins(client, project, &api.ActionsListPluginsOptions{
			PerPage: perPage,
			Page:    page,
		})
		if err != nil {
			return nil, err
		}
		entries, err := api.ParseActionsPluginsListRaw(raw)
		if err != nil {
			return nil, err
		}
		all = append(all, entries...)
		if opts.Limit > 0 && len(all) >= opts.Limit {
			return trimEntries(all, opts), nil
		}
		if len(entries) < perPage {
			break
		}
	}
	return trimEntries(all, opts), nil
}

func resolvePerPage(opts *ListOptions) int {
	if opts.PerPageSet && opts.PerPage > 0 {
		return opts.PerPage
	}
	if opts.Paginate || opts.Page == 0 {
		return 100
	}
	if opts.Limit > 0 {
		return opts.Limit
	}
	return 30
}

func trimEntries(entries []json.RawMessage, opts *ListOptions) []json.RawMessage {
	if entries == nil {
		return []json.RawMessage{}
	}
	if opts.Limit > 0 && len(entries) > opts.Limit {
		return entries[:opts.Limit]
	}
	return entries
}

func printPluginsTable(io *iostreams.IOStreams, rawEntries []json.RawMessage) error {
	cs := io.ColorScheme()
	for _, raw := range rawEntries {
		var p api.ActionsPlugin
		if err := json.Unmarshal(raw, &p); err != nil {
			continue
		}
		fmt.Fprintf(io.Out, "%s\t%s\t%s\t%s\n",
			cs.Bold(orDash(p.Name)), orDash(p.DisplayName), orDash(p.Version),
			truncateDescription(p.Description))
	}
	return nil
}

func orDash(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "-"
	}
	return s
}

func truncateDescription(desc string) string {
	desc = strings.TrimSpace(desc)
	if desc == "" {
		return "-"
	}
	if utf8.RuneCountInString(desc) > 60 {
		runes := []rune(desc)
		return string(runes[:57]) + "..."
	}
	return desc
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
