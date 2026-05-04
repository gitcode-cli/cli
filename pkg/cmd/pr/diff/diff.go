// Package diff implements the pr diff command
package diff

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

type DiffOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	JSON bool
}

// DiffFile represents a file in the diff
type DiffFile struct {
	Path    string `json:"path"`
	OldPath string `json:"old_path,omitempty"`
	Added   int    `json:"added"`
	Removed int    `json:"removed"`
	NewFile bool   `json:"new_file,omitempty"`
	Deleted bool   `json:"deleted,omitempty"`
	Renamed bool   `json:"renamed,omitempty"`
}

// DiffResult represents the JSON output for pr diff
type DiffResult struct {
	Number       int        `json:"number"`
	Title        string     `json:"title"`
	HeadBranch   string     `json:"head_branch"`
	BaseBranch   string     `json:"base_branch"`
	AddedLines   int        `json:"added_lines"`
	RemovedLines int        `json:"removed_lines"`
	FileCount    int        `json:"file_count"`
	Files        []DiffFile `json:"files"`
}

// NewCmdDiff creates the diff command
func NewCmdDiff(f *cmdutil.Factory, runF func(*DiffOptions) error) *cobra.Command {
	opts := &DiffOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "diff <number>",
		Short: "View changes in a pull request",
		Long: heredoc.Doc(`
			View the diff of a pull request.
		`),
		Example: heredoc.Doc(`
			# View PR diff
			$ gc pr diff 123 -R owner/repo

			# JSON output
			$ gc pr diff 123 --json -R owner/repo
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return cmdutil.NewUsageError(fmt.Sprintf("invalid PR number: %s", args[0]))
			}
			opts.Number = number

			if runF != nil {
				return runF(opts)
			}
			return diffRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func diffRun(opts *DiffOptions) error {
	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := cmdutil.EnvToken()
	if token == "" {
		return cmdutil.NewAuthError("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "environment")

	// Get repository
	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}

	owner, repo, err := parseRepo(repository)
	if err != nil {
		return err
	}

	// Get PR info
	pr, err := api.GetPullRequest(client, owner, repo, opts.Number)
	if err != nil {
		return cmdutil.WrapNotFound(err, "PR #%d not found in %s/%s", opts.Number, owner, repo)
	}

	// Get PR files and diffs
	files, err := api.GetPRFiles(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to get PR diff: %w", err)
	}

	if opts.JSON {
		var diffFiles []DiffFile
		for _, diff := range files.Diffs {
			df := DiffFile{
				Path:    diff.NewPath,
				OldPath: diff.OldPath,
			}
			if diff.Statistic != nil {
				df.Added = diff.Statistic.Additions
				df.Removed = diff.Statistic.Deletions
			}
			if diff.OldPath == "" && diff.NewPath != "" {
				df.NewFile = true
			} else if diff.OldPath != "" && diff.NewPath == "" {
				df.Deleted = true
			} else if diff.OldPath != diff.NewPath {
				df.Renamed = true
			}
			diffFiles = append(diffFiles, df)
		}

		result := DiffResult{
			Number:       pr.Number,
			Title:        pr.Title,
			HeadBranch:   pr.Head.Ref,
			BaseBranch:   pr.Base.Ref,
			AddedLines:   files.AddedLines,
			RemovedLines: files.RemoveLines,
			FileCount:    files.Count,
			Files:        diffFiles,
		}
		return cmdutil.WriteJSON(opts.IO.Out, result)
	}

	// Output PR info
	fmt.Fprintf(opts.IO.Out, "PR #%d: %s\n", pr.Number, pr.Title)
	fmt.Fprintf(opts.IO.Out, "Branch: %s -> %s\n", pr.Head.Ref, pr.Base.Ref)
	fmt.Fprintf(opts.IO.Out, "Changes: +%d -%d in %d file(s)\n\n", files.AddedLines, files.RemoveLines, files.Count)

	// Output diff for each file
	for _, diff := range files.Diffs {
		if diff.NewPath != "" {
			if diff.OldPath != "" && diff.OldPath != diff.NewPath {
				fmt.Fprintf(opts.IO.Out, "diff --git a/%s b/%s\n", diff.OldPath, diff.NewPath)
			} else {
				fmt.Fprintf(opts.IO.Out, "diff --git a/%s b/%s\n", diff.NewPath, diff.NewPath)
			}
		}

		// Output diff content
		if diff.Content != nil && len(diff.Content.Text) > 0 {
			for _, line := range diff.Content.Text {
				switch line.Type {
				case "new":
					fmt.Fprintf(opts.IO.Out, "+%s\n", line.LineContent)
				case "old":
					fmt.Fprintf(opts.IO.Out, "-%s\n", line.LineContent)
				default:
					fmt.Fprintf(opts.IO.Out, " %s\n", line.LineContent)
				}
			}
		}
		fmt.Fprintf(opts.IO.Out, "\n")
	}

	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}
