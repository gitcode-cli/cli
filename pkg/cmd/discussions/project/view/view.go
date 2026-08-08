// Package view implements the discussions project view command.
package view

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/discussions/render"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// ViewOptions configures the discussions project view command.
type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	Number     int

	JSON bool
}

// NewCmdView creates the discussions project view command.
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "view <number>",
		Short: "View a repository discussion",
		Long: heredoc.Doc(`
			View a single GitCode repository (project-level) discussion (discuss)
			via the v5 API (GET /api/v5/repos/{owner}/{repo}/discuss/{number}).
			Use --json for the raw API response.
		`),
		Example: heredoc.Doc(`
			# View a discussion
			$ gc discussions project view 42 -R owner/repo

			# Output as JSON
			$ gc discussions project view 42 -R owner/repo --json
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

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func viewRun(opts *ViewOptions) error {
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

	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}
	owner, repo, err := cmdutil.ParseRepo(repository)
	if err != nil {
		return err
	}

	d, err := api.GetRepoDiscussion(client, owner, repo, opts.Number)
	if err != nil {
		return cmdutil.WrapNotFound(err, "discussion #%d not found in %s", opts.Number, repository)
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, d)
	}

	render.PrintDiscussion(opts.IO, d)
	return nil
}
