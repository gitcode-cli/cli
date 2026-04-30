// Package merge implements the pr merge command
package merge

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type MergeOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	MergeMethod  string
	DeleteBranch bool
	Yes          bool
	JSON         bool
}

type mergeResult struct {
	Number        int              `json:"number"`
	Merged        bool             `json:"merged"`
	PullRequest   *api.PullRequest `json:"pull_request"`
	DeletedBranch string           `json:"deleted_branch,omitempty"`
}

// NewCmdMerge creates the merge command
func NewCmdMerge(f *cmdutil.Factory, runF func(*MergeOptions) error) *cobra.Command {
	opts := &MergeOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "merge <number>",
		Short: "Merge a pull request",
		Long: heredoc.Doc(`
			Merge a pull request in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Merge a PR
			$ gc pr merge 123 -R owner/repo

			# Merge with squash
			$ gc pr merge 123 -R owner/repo --method squash

			# Merge and delete branch
			$ gc pr merge 123 -R owner/repo --delete-branch
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid PR number: %s", args[0])
			}
			opts.Number = number

			if runF != nil {
				return runF(opts)
			}
			return mergeRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.MergeMethod, "method", "m", "merge", "Merge method (merge/squash/rebase)")
	cmdutil.SetFlagEnum(cmd, "method", "merge", "squash", "rebase")
	cmd.Flags().BoolVarP(&opts.DeleteBranch, "delete-branch", "d", false, "Delete branch after merge")
	cmd.Flags().BoolVar(&opts.Yes, "yes", false, "Skip confirmation prompt")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func mergeRun(opts *MergeOptions) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := cmdutil.EnvToken()
	if token == "" {
		return cmdutil.NewAuthError("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "active")

	// Get repository
	owner, repo, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	// Get PR first
	pr, err := api.GetPullRequest(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to get PR: %w", err)
	}
	if err := cmdutil.ConfirmOrAbort(cmdutil.ConfirmOptions{
		IO:       opts.IO,
		Yes:      opts.Yes,
		Expected: strconv.Itoa(pr.Number),
		Prompt:   fmt.Sprintf("! This will merge PR #%d %s\nType the PR number to confirm: ", pr.Number, cs.Bold(pr.Title)),
	}); err != nil {
		return err
	}

	deleteOwner, deleteRepo, deleteRef, err := deleteBranchTarget(pr, opts.DeleteBranch)
	if err != nil {
		return err
	}

	// Merge PR
	mergedPR, err := api.MergePullRequest(client, owner, repo, opts.Number, &api.MergePROptions{
		MergeMethod: opts.MergeMethod,
	})
	if err != nil {
		return fmt.Errorf("failed to merge PR: %w", err)
	}

	// Delete branch if requested
	if opts.DeleteBranch {
		if err := api.DeleteBranch(client, deleteOwner, deleteRepo, deleteRef); err != nil {
			return fmt.Errorf("failed to delete branch %s: %w", deleteRef, err)
		}
	}

	if opts.JSON {
		result := mergeResult{
			Number:      opts.Number,
			Merged:      true,
			PullRequest: mergedPR,
		}
		if mergedPR != nil && mergedPR.Merged {
			result.Merged = mergedPR.Merged
		}
		if opts.DeleteBranch {
			result.DeletedBranch = deleteRef
		}
		return cmdutil.WriteJSON(opts.IO.Out, result)
	}

	fmt.Fprintf(opts.IO.Out, "%s Merged PR #%d\n", cs.Green("✓"), opts.Number)
	if opts.DeleteBranch {
		fmt.Fprintf(opts.IO.Out, "  Deleted branch %s\n", deleteRef)
	}

	return nil
}

func deleteBranchTarget(pr *api.PullRequest, enabled bool) (string, string, string, error) {
	if !enabled {
		return "", "", "", nil
	}
	if pr == nil || pr.Head == nil {
		return "", "", "", fmt.Errorf("failed to delete branch: PR head metadata is missing")
	}
	if strings.TrimSpace(pr.Head.Ref) == "" {
		return "", "", "", fmt.Errorf("failed to delete branch: PR head branch is empty")
	}
	if pr.Head.Repo == nil || strings.TrimSpace(pr.Head.Repo.FullName) == "" {
		return "", "", "", fmt.Errorf("failed to delete branch %s: PR head repository is missing", pr.Head.Ref)
	}
	owner, repo, err := cmdutil.ParseRepo(pr.Head.Repo.FullName)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to delete branch %s: invalid PR head repository: %w", pr.Head.Ref, err)
	}
	return owner, repo, pr.Head.Ref, nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
