// Package edit implements the repo edit command.
package edit

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// EditOptions configures the repo edit command.
type EditOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository    string
	Description   string
	Homepage      string
	DefaultBranch string
	Name          string
	Private       bool
	Public        bool
	JSON          bool
	privateSet    bool
}

// NewCmdEdit creates the repo edit command.
func NewCmdEdit(f *cmdutil.Factory, runF func(*EditOptions) error) *cobra.Command {
	opts := &EditOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit repository settings",
		Long: heredoc.Doc(`
			Update a GitCode repository's settings. Only provided fields
			are updated.

			Use --private or --public to toggle visibility.
		`),
		Example: heredoc.Doc(`
			# Update description
			$ gc repo edit --description "New description" -R owner/repo

			# Make repository private
			$ gc repo edit --private -R owner/repo

			# Update default branch
			$ gc repo edit --default-branch main -R owner/repo

			# Rename repository
			$ gc repo edit --name new-name -R owner/repo
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return editRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Repository description")
	cmd.Flags().StringVar(&opts.Homepage, "homepage", "", "Homepage URL")
	cmd.Flags().StringVar(&opts.DefaultBranch, "default-branch", "", "Default branch")
	cmd.Flags().StringVar(&opts.Name, "name", "", "Repository name (rename)")
	cmd.Flags().BoolVar(&opts.Private, "private", false, "Make repository private")
	cmd.Flags().BoolVar(&opts.Public, "public", false, "Make repository public")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func editRun(opts *EditOptions) error {
	if opts.Private && opts.Public {
		return cmdutil.NewUsageError("--private and --public cannot be used together")
	}

	req := &api.UpdateRepoRequest{
		Name:          opts.Name,
		Description:   opts.Description,
		Homepage:      opts.Homepage,
		DefaultBranch: opts.DefaultBranch,
	}
	if opts.Private {
		p := true
		req.Private = &p
	}
	if opts.Public {
		p := false
		req.Private = &p
	}

	if req.Name == "" && req.Description == "" && req.Homepage == "" &&
		req.DefaultBranch == "" && req.Private == nil {
		return cmdutil.NewUsageError("at least one field must be provided to update")
	}

	if req.Description != "" {
		if err := cmdutil.ScanContentForSecrets(req.Description); err != nil {
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

	result, err := api.UpdateRepo(client, owner, repo, req)
	if err != nil {
		return fmt.Errorf("failed to update repository: %w", err)
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, result)
	}

	fmt.Fprintf(opts.IO.ErrOut, "Updated repository %s/%s.\n", owner, repo)
	return nil
}
