// Package list implements the actions workflow list command.
package list

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// ListOptions configures the actions workflow list command.
type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	Limit      int
	Page       int
	JSON       bool
}

// NewCmdList creates the actions workflow list command.
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List workflows",
		Long:  heredoc.Doc(`List all workflow files in a GitCode repository.`),
		Example: heredoc.Doc(`
			$ gc actions workflow list -R owner/repo
			$ gc actions workflow list -R owner/repo --json
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum results per page")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "Page number")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func listRun(opts *ListOptions) error {
	if opts.Limit <= 0 {
		return cmdutil.NewUsageError("--limit must be greater than 0")
	}
	if opts.Page < 0 {
		return cmdutil.NewUsageError("--page must be greater than or equal to 0")
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

	resp, err := api.ListActionsWorkflows(client, owner, repo, opts.Page, opts.Limit)
	if err != nil {
		return fmt.Errorf("failed to list workflows: %w", err)
	}

	if len(resp.Workflows) == 0 {
		if opts.JSON {
			return cmdutil.WriteJSON(opts.IO.Out, resp.Workflows)
		}
		fmt.Fprintln(opts.IO.Out, "No workflows found.")
		return nil
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, resp.Workflows)
	}

	cs := opts.IO.ColorScheme()
	for _, w := range resp.Workflows {
		state := w.State
		if state == "active" || state == "Active" {
			state = cs.Green(state)
		}
		fmt.Fprintf(opts.IO.Out, "%s\t%s\t%s\n", w.Name, w.FilePath, state)
	}
	return nil
}
