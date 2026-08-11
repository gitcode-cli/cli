// Package create implements the branch-protection create command.
package create

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

type CreateOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	Wildcard   string
	Pusher     string
	Merger     string
	JSON       bool
}

func NewCmdCreate(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a branch protection rule",
		Long: heredoc.Doc(`
			Create a new branch protection rule for a branch or branch pattern.

			Pusher/merger roles: admin, develop, maintainer, or usernames
			(semicolon-separated). Empty string means no one can push/merge.
		`),
		Example: heredoc.Doc(`
			$ gc repo branch-protection create --wildcard "main" --pusher admin --merger admin -R owner/repo
			$ gc repo branch-protection create --wildcard "release/*" --pusher "admin;user1" -R owner/repo
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return createRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVar(&opts.Wildcard, "wildcard", "", "Branch name or wildcard pattern (e.g. main, release/*)")
	cmd.Flags().StringVar(&opts.Pusher, "pusher", "", "Who can push (roles/usernames, semicolon-separated)")
	cmd.Flags().StringVar(&opts.Merger, "merger", "", "Who can merge (roles/usernames, semicolon-separated)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func createRun(opts *CreateOptions) error {
	if strings.TrimSpace(opts.Wildcard) == "" {
		return cmdutil.NewUsageError("--wildcard is required")
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
		Wildcard: opts.Wildcard,
		Pusher:   opts.Pusher,
		Merger:   opts.Merger,
	}
	if err := api.CreateBranchProtection(client, owner, repo, req); err != nil {
		return fmt.Errorf("failed to create branch protection rule: %w", err)
	}

	fmt.Fprintf(opts.IO.ErrOut, "Created branch protection rule for %s.\n", opts.Wildcard)
	return nil
}
