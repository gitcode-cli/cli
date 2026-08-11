// Package delete implements the branch-protection delete command.
package delete

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type DeleteOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	Wildcard   string
	Yes        bool
}

func NewCmdDelete(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:     "delete <wildcard>",
		Short:   "Delete a branch protection rule",
		Long:    heredoc.Doc(`Delete an existing branch protection rule.`),
		Example: heredoc.Doc(`$ gc repo branch-protection delete "main" -R owner/repo`),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Wildcard = strings.TrimSpace(args[0])
			if opts.Wildcard == "" {
				return cmdutil.NewUsageError("wildcard is required")
			}
			if runF != nil {
				return runF(opts)
			}
			return deleteRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVarP(&opts.Yes, "yes", "y", false, "Skip confirmation")

	return cmd
}

func deleteRun(opts *DeleteOptions) error {
	if err := cmdutil.ConfirmOrAbort(cmdutil.ConfirmOptions{
		IO:       opts.IO,
		Yes:      opts.Yes,
		Expected: opts.Wildcard,
		Prompt:   fmt.Sprintf("Type the wildcard to confirm deletion: %s\n", opts.Wildcard),
	}); err != nil {
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

	if err := api.DeleteBranchProtection(client, owner, repo, opts.Wildcard); err != nil {
		return fmt.Errorf("failed to delete branch protection rule: %w", err)
	}

	fmt.Fprintf(opts.IO.ErrOut, "Deleted branch protection rule for %s.\n", opts.Wildcard)
	return nil
}
