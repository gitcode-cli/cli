// Package view implements the actions plugin view command.
package view

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// ViewOptions configures the actions plugin view command.
type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	PluginName string

	JSON bool
}

// NewCmdView creates the actions plugin view command.
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "view <plugin-name>",
		Short: "View an Actions plugin detail and README",
		Long: heredoc.Doc(`
			View the metadata, versions and Markdown README of a specific
			official GitCode Actions plugin.

			By default the README of each version in vision_content is printed.
			Use --json to output the full API response with all fields preserved.
		`),
		Example: heredoc.Doc(`
			# View a plugin's README
			$ gc actions plugin view checkout -R owner/repo

			# Output as JSON
			$ gc actions plugin view checkout -R owner/repo --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.PluginName = args[0]
			if runF != nil {
				return runF(opts)
			}
			return viewRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func viewRun(opts *ViewOptions) error {
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

	raw, err := api.ViewActionsPlugin(client, owner+"/"+repo, opts.PluginName)
	if err != nil {
		return fmt.Errorf("failed to view actions plugin: %w", err)
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, json.RawMessage(raw))
	}

	var detail api.ActionsPluginDetail
	if err := json.Unmarshal(raw, &detail); err != nil {
		return fmt.Errorf("failed to parse plugin detail: %w", err)
	}

	return printPluginDetail(opts.IO, &detail, opts.PluginName)
}

func printPluginDetail(io *iostreams.IOStreams, d *api.ActionsPluginDetail, requestedName string) error {
	cs := io.ColorScheme()
	name := d.Name
	if name == "" {
		name = requestedName
	}
	fmt.Fprintf(io.Out, "%s\t%s\n", cs.Bold("Plugin:"), name)
	fmt.Fprintf(io.Out, "%s\t%s\n", "Display Name:", orValue(d.DisplayName))
	fmt.Fprintf(io.Out, "%s\t%s\n", "Description:", orValue(d.Description))
	fmt.Fprintf(io.Out, "%s\t%d\n", "Versions:", len(d.VisionContent))
	fmt.Fprintln(io.Out)

	for _, v := range d.VisionContent {
		label := v.Version
		if label == "" {
			label = "unknown"
		}
		fmt.Fprintf(io.Out, "%s\n", cs.Bold(fmt.Sprintf("== Version %s ==", label)))
		readme := v.Readme
		if readme == "" {
			fmt.Fprintf(io.Out, "(no README available)\n\n")
			continue
		}
		fmt.Fprintf(io.Out, "%s\n\n", readme)
	}
	return nil
}

func orValue(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "-"
	}
	return s
}
