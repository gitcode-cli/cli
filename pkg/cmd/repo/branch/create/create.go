// Package create implements the repo branch create command.
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

// CreateOptions configures the repo branch create command.
type CreateOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository  string
	BranchName  string
	Ref         string
	Description string
	JSON        bool
}

// NewCmdCreate creates the repo branch create command.
func NewCmdCreate(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "create <branch-name>",
		Short: "Create a branch in a repository",
		Long: heredoc.Doc(`
			Create a new branch in a GitCode repository from a specified
			reference (branch or tag). The default source is the default
			branch of the repository.
		`),
		Example: heredoc.Doc(`
			# Create a branch from the default branch
			$ gc repo branch create feature -R owner/repo

			# Create from a specific ref
			$ gc repo branch create feature --ref develop -R owner/repo

			# With description
			$ gc repo branch create feature --description "feature branch" -R owner/repo
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.BranchName = strings.TrimSpace(args[0])
			if opts.BranchName == "" {
				return cmdutil.NewUsageError("branch name is required")
			}
			if runF != nil {
				return runF(opts)
			}
			return createRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVar(&opts.Ref, "ref", "", "Source branch or tag name (default: repository default branch)")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Branch description")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func createRun(opts *CreateOptions) error {
	if opts.Description != "" {
		if err := cmdutil.ScanContentForSecrets(opts.Description); err != nil {
			return err
		}
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

	branch, err := api.CreateBranch(client, owner, repo, &api.CreateBranchRequest{
		Refs:        opts.Ref,
		BranchName:  opts.BranchName,
		Description: opts.Description,
	})
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, branch)
	}

	fmt.Fprintf(opts.IO.ErrOut, "Created branch %s from %s\n", branch.Name, orDash(opts.Ref))
	return nil
}

func orDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "default"
	}
	return s
}
