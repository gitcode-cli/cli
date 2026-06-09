// Package branch implements the repo branch command
package branch

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	Repository string
	Branch     string

	// Flags
	JSON bool
}

// NewCmdBranch creates the repo branch command group with view subcommand
func NewCmdBranch(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "branch <command>",
		Short: "Manage branches",
		Long:  "Work with repository branches.",
	}

	viewCmd := newCmdView(f, runF)
	cmd.AddCommand(viewCmd)

	return cmd
}

func newCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "view <branch>",
		Short: "View a branch in a repository",
		Long: heredoc.Doc(`
			View information about a branch in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# View a branch
			$ gc repo branch view main -R owner/repo

			# View a branch in current repository
			$ gc repo branch view main

			# Output as JSON
			$ gc repo branch view main -R owner/repo --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Branch = args[0]

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
	cs := opts.IO.ColorScheme()

	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}

	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}
	owner, repo, err := parseRepo(repository)
	if err != nil {
		return err
	}

	b, err := api.GetBranch(client, owner, repo, opts.Branch)
	if err != nil {
		return cmdutil.WrapNotFound(err, "branch %q not found in %s/%s", opts.Branch, owner, repo)
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, b)
	}

	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "%s\n", cs.Bold(b.Name))
	if b.Protected {
		fmt.Fprintf(opts.IO.Out, "  Protected: yes\n")
	}
	if b.Commit != nil {
		fmt.Fprintf(opts.IO.Out, "  Commit: %s\n", b.Commit.ID)
		if b.Commit.ShortID != "" {
			fmt.Fprintf(opts.IO.Out, "  Short ID: %s\n", b.Commit.ShortID)
		}
		if b.Commit.Title != "" {
			fmt.Fprintf(opts.IO.Out, "  Title: %s\n", b.Commit.Title)
		}
		if b.Commit.Author != nil && b.Commit.Author.Login != "" {
			fmt.Fprintf(opts.IO.Out, "  Author: %s\n", b.Commit.Author.Login)
		}
	}
	fmt.Fprintf(opts.IO.Out, "\n")

	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
