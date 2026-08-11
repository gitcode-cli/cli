// Package view implements the repo pr-settings view command.
package view

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// ViewOptions configures the pr-settings view command.
type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	JSON       bool
}

// NewCmdView creates the pr-settings view command.
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "view",
		Short: "View pull request settings",
		Long:  heredoc.Doc(`View the pull request merge gate settings for a repository.`),
		Example: heredoc.Doc(`
			$ gc repo pr-settings view -R owner/repo
			$ gc repo pr-settings view -R owner/repo --json
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
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

	settings, err := api.GetPRSettings(client, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to get PR settings: %w", err)
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, settings)
	}

	printSettings(opts, settings)
	return nil
}

func printSettings(opts *ViewOptions, s *api.PRSettings) {
	out := opts.IO.Out
	fmt.Fprintf(out, "Pull Request Settings:\n")
	fmt.Fprintf(out, "  Merge gate:\n")
	fmt.Fprintf(out, "    Pipeline required:        %v\n", s.OnlyAllowMergeIfPipelineSucceeds)
	fmt.Fprintf(out, "    Discussions resolved:     %v\n", s.OnlyAllowMergeIfAllDiscussionsResolved)
	fmt.Fprintf(out, "    Disable self-merge:       %v\n", s.DisableMergeBySelf)
	fmt.Fprintf(out, "    Force merge (admin):      %v\n", s.CanForceMerge)
	fmt.Fprintf(out, "    Merge method:             %s\n", orDash(s.MergeMethod))
	fmt.Fprintf(out, "    Delete source on merge:   %v\n", s.DeleteSourceBranchWhenMerged)
	fmt.Fprintf(out, "  Reviewers:\n")
	fmt.Fprintf(out, "    Reviewers enabled:        %v\n", s.ApprovalRequiredReviewersEnable)
	fmt.Fprintf(out, "    Min reviewers:            %d\n", s.ApprovalRequiredReviewers)
	fmt.Fprintf(out, "    Min approvers:            %d\n", s.ApprovalRequiredApprovers)
	fmt.Fprintf(out, "    Min testers:              %d\n", s.ApprovalRequiredTesters)
	fmt.Fprintf(out, "  Squash:\n")
	fmt.Fprintf(out, "    Disable squash merge:     %v\n", s.DisableSquashMerge)
	fmt.Fprintf(out, "    Auto squash merge:        %v\n", s.AutoSquashMerge)
	fmt.Fprintf(out, "    Squash no merge commit:   %v\n", s.SquashMergeWithNoMergeCommit)
	fmt.Fprintf(out, "  Other:\n")
	fmt.Fprintf(out, "    Check CLA:                %v\n", s.IsCheckCla)
	fmt.Fprintf(out, "    Lite PR enabled:          %v\n", s.IsAllowLiteMergeRequest)
	fmt.Fprintf(out, "    Close issue on merge:     %v\n", s.CloseIssueWhenMrMerged)
	fmt.Fprintf(out, "    Can reopen:               %v\n", s.CanReopen)
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
