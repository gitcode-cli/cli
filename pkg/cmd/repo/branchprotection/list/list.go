// Package list implements the branch-protection list command.
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

type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	JSON       bool
}

func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List branch protection rules",
		Long:    heredoc.Doc(`List all branch protection rules in a repository.`),
		Example: heredoc.Doc(`$ gc repo branch-protection list -R owner/repo`),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func listRun(opts *ListOptions) error {
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

	rules, err := api.ListBranchProtections(client, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to list branch protection rules: %w", err)
	}

	if len(rules) == 0 {
		if opts.JSON {
			return cmdutil.WriteJSON(opts.IO.Out, rules)
		}
		fmt.Fprintln(opts.IO.Out, "No branch protection rules found.")
		return nil
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, rules)
	}

	for _, r := range rules {
		fmt.Fprintf(opts.IO.Out, "%s\tpush=%s\tmerge=%s\n", r.Name, orDash(r.PushUsers), orDash(r.MergeUsers))
	}
	return nil
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
