// Package view implements the discussions view command.
package view

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// ViewOptions configures the discussions view command.
type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	Org    string
	Number int

	JSON bool
}

// NewCmdView creates the discussions view command.
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "view <number>",
		Short: "View an organization discussion",
		Long: heredoc.Doc(`
			View a single GitCode organization discussion (discuss) via the v5
			API (GET /api/v5/orgs/{org}/discuss/{number}). Use --json for the
			raw API response.
		`),
		Example: heredoc.Doc(`
			# View a discussion
			$ gc discussions view 42 --org my-org

			# Output as JSON
			$ gc discussions view 42 --org my-org --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return cmdutil.NewUsageError(fmt.Sprintf("invalid discussion number: %s", args[0]))
			}
			opts.Number = number
			if runF != nil {
				return runF(opts)
			}
			return viewRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Org, "org", "", "Organization path (required)")
	cmd.MarkFlagRequired("org")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func viewRun(opts *ViewOptions) error {
	if opts.Org == "" {
		return cmdutil.NewUsageError("--org is required")
	}
	if opts.Number < 1 {
		return cmdutil.NewUsageError("discussion number must be a positive integer")
	}

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}
	client, err := cmdutil.AuthenticatedClient(httpClient)
	if err != nil {
		return err
	}

	d, err := api.GetOrgDiscussion(client, opts.Org, opts.Number)
	if err != nil {
		return cmdutil.WrapNotFound(err, "discussion #%d not found in org %s", opts.Number, opts.Org)
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, d)
	}

	printDiscussion(opts.IO, d)
	return nil
}

func printDiscussion(io *iostreams.IOStreams, d *api.Discussion) {
	cs := io.ColorScheme()
	state := "open"
	if d.IsClosed != 0 {
		state = "closed"
	}
	author := ""
	if d.Author != nil {
		author = d.Author.Login
	}
	fmt.Fprintf(io.Out, "%s #%d  %s\n", cs.Bold("Discussion"), d.Number, d.Title)
	fmt.Fprintf(io.Out, "state: %s  comments: %d  author: %s\n", cs.Cyan(state), d.CommentTotal, author)
	if d.CreatedAt != "" {
		fmt.Fprintf(io.Out, "created: %s", d.CreatedAt)
		if d.UpdatedAt != "" {
			fmt.Fprintf(io.Out, "  updated: %s", d.UpdatedAt)
		}
		fmt.Fprintln(io.Out)
	}
	if d.MdContent != "" {
		fmt.Fprintln(io.Out)
		fmt.Fprintf(io.Out, "%s\n", d.MdContent)
	}
}
