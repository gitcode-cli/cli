// Package resolve implements the PR comment resolve/unresolve commands
package resolve

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type resolveOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository   string
	PRNumber     int
	DiscussionID string
	Resolved     bool
}

// NewCmdResolve creates the resolve command
func NewCmdResolve(f *cmdutil.Factory, runF func(*resolveOptions) error) *cobra.Command {
	opts := &resolveOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
		Resolved:   true,
	}

	cmd := &cobra.Command{
		Use:   "resolve <pr-number> <discussion-id>",
		Short: "Mark a PR comment discussion as resolved",
		Long: heredoc.Doc(`
			Mark a pull request comment discussion as resolved.

			The discussion ID can be found in the output of "gc pr comments".
		`),
		Example: heredoc.Doc(`
			# Resolve a comment discussion
			$ gc pr comment resolve 123 d1 -R owner/repo
		`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			num, err := strconv.Atoi(args[0])
			if err != nil {
				return cmdutil.NewUsageError(fmt.Sprintf("invalid PR number: %s", args[0]))
			}
			opts.PRNumber = num
			opts.DiscussionID = args[1]

			if runF != nil {
				return runF(opts)
			}
			return resolveRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")

	return cmd
}

// NewCmdUnresolve creates the unresolve command
func NewCmdUnresolve(f *cmdutil.Factory, runF func(*resolveOptions) error) *cobra.Command {
	opts := &resolveOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
		Resolved:   false,
	}

	cmd := &cobra.Command{
		Use:   "unresolve <pr-number> <discussion-id>",
		Short: "Mark a PR comment discussion as unresolved",
		Long: heredoc.Doc(`
			Mark a pull request comment discussion as unresolved.

			The discussion ID can be found in the output of "gc pr comments".
		`),
		Example: heredoc.Doc(`
			# Unresolve a comment discussion
			$ gc pr comment unresolve 123 d1 -R owner/repo
		`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			num, err := strconv.Atoi(args[0])
			if err != nil {
				return cmdutil.NewUsageError(fmt.Sprintf("invalid PR number: %s", args[0]))
			}
			opts.PRNumber = num
			opts.DiscussionID = args[1]

			if runF != nil {
				return runF(opts)
			}
			return resolveRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")

	return cmd
}

func resolveRun(opts *resolveOptions) error {
	cs := opts.IO.ColorScheme()

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

	resolveOpts := &api.ResolvePRCommentOptions{
		Resolved: opts.Resolved,
	}
	if err := api.ResolvePRComment(client, owner, repo, opts.PRNumber, opts.DiscussionID, resolveOpts); err != nil {
		return fmt.Errorf("failed to update comment resolution: %w", err)
	}

	if opts.Resolved {
		fmt.Fprintf(opts.IO.Out, "%s Marked discussion #%s as resolved\n", cs.Green("✓"), opts.DiscussionID)
	} else {
		fmt.Fprintf(opts.IO.Out, "%s Marked discussion #%s as unresolved\n", cs.Green("✓"), opts.DiscussionID)
	}
	return nil
}
