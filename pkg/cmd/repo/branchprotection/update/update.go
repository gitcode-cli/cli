// Package update implements the branch-protection update command.
package update

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

type UpdateOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	Wildcard   string
	Pusher     string
	Merger     string
}

func NewCmdUpdate(f *cmdutil.Factory, runF func(*UpdateOptions) error) *cobra.Command {
	opts := &UpdateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:     "update <wildcard>",
		Short:   "Update a branch protection rule",
		Long:    heredoc.Doc(`Update an existing branch protection rule.`),
		Example: heredoc.Doc(`$ gc repo branch-protection update "main" --pusher develop -R owner/repo`),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Wildcard = strings.TrimSpace(args[0])
			if opts.Wildcard == "" {
				return cmdutil.NewUsageError("wildcard is required")
			}
			if runF != nil {
				return runF(opts)
			}
			return updateRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVar(&opts.Pusher, "pusher", "", "Who can push (roles/usernames)")
	cmd.Flags().StringVar(&opts.Merger, "merger", "", "Who can merge (roles/usernames)")

	return cmd
}

func updateRun(opts *UpdateOptions) error {
	if opts.Pusher == "" && opts.Merger == "" {
		return cmdutil.NewUsageError("at least one of --pusher or --merger must be provided")
	}

	if err := cmdutil.ScanContentForSecrets(opts.Pusher); err != nil {
		return err
	}
	if err := cmdutil.ScanContentForSecrets(opts.Merger); err != nil {
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

	req := &api.BranchProtectionRequest{
		Pusher: opts.Pusher,
		Merger: opts.Merger,
	}
	if err := api.UpdateBranchProtection(client, owner, repo, opts.Wildcard, req); err != nil {
		return fmt.Errorf("failed to update branch protection rule: %w", err)
	}

	fmt.Fprintf(opts.IO.ErrOut, "Updated branch protection rule for %s.\n", opts.Wildcard)
	return nil
}
