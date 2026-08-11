// Package setdefault implements the repo set-default command.
package setdefault

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/config"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

const defaultHost = "gitcode.com"

// SetDefaultOptions configures the repo set-default command.
type SetDefaultOptions struct {
	IO       *iostreams.IOStreams
	Config   func() (config.Config, error)
	BaseRepo func() (string, error)

	Repository string
	View       bool
	Unset      bool
	JSON       bool
}

// NewCmdSetDefault creates the repo set-default command.
func NewCmdSetDefault(f *cmdutil.Factory, runF func(*SetDefaultOptions) error) *cobra.Command {
	opts := &SetDefaultOptions{
		IO:       f.IOStreams,
		Config:   f.Config,
		BaseRepo: f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "set-default [<owner>/<repo>]",
		Short: "Set or view the default repository",
		Long: heredoc.Doc(`
			Configure or display the default repository for the current
			host. Other commands use this value when no -R/--repo flag
			is provided and the current directory is not a git repository.

			With no arguments, sets the default to the repository inferred
			from the git remote. Pass an explicit <owner>/<repo> to
			override. Use --view to display the current default, --unset
			to clear it.
		`),
		Example: heredoc.Doc(`
			# Set default to current git remote's repository
			$ gc repo set-default

			# Set an explicit default
			$ gc repo set-default owner/repo

			# View the current default
			$ gc repo set-default --view

			# Clear the default
			$ gc repo set-default --unset
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Repository = args[0]
			}
			if runF != nil {
				return runF(opts)
			}
			return setDefaultRun(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.View, "view", false, "View the current default repository")
	cmd.Flags().BoolVar(&opts.Unset, "unset", false, "Clear the default repository")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func setDefaultRun(opts *SetDefaultOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if opts.View {
		return viewDefault(cfg, opts)
	}
	if opts.Unset {
		return unsetDefault(cfg, opts)
	}
	return setDefault(cfg, opts)
}

func viewDefault(cfg config.Config, opts *SetDefaultOptions) error {
	value, err := cfg.Get(defaultHost, "default_repo")
	if err != nil {
		return fmt.Errorf("failed to read default repo: %w", err)
	}
	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, map[string]string{"default_repo": value})
	}
	if value == "" {
		fmt.Fprintln(opts.IO.Out, "No default repository is set.")
		return nil
	}
	fmt.Fprintf(opts.IO.Out, "%s\n", value)
	return nil
}

func unsetDefault(cfg config.Config, opts *SetDefaultOptions) error {
	if err := cfg.Set(defaultHost, "default_repo", ""); err != nil {
		return fmt.Errorf("failed to clear default repo: %w", err)
	}
	if err := cfg.Write(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	fmt.Fprintln(opts.IO.ErrOut, "Default repository cleared.")
	return nil
}

func setDefault(cfg config.Config, opts *SetDefaultOptions) error {
	repo, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}
	if _, _, perr := cmdutil.ParseRepo(repo); perr != nil {
		return perr
	}
	if err := cfg.Set(defaultHost, "default_repo", repo); err != nil {
		return fmt.Errorf("failed to set default repo: %w", err)
	}
	if err := cfg.Write(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	fmt.Fprintf(opts.IO.ErrOut, "Default repository set to %s.\n", repo)
	return nil
}
