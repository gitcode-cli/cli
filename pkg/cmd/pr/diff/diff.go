// Package diff implements the pr diff command
package diff

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type DiffOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	Number     int
}

// NewCmdDiff creates the diff command
func NewCmdDiff(f *cmdutil.Factory, runF func(*DiffOptions) error) *cobra.Command {
	opts := &DiffOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "diff <number>",
		Short: "View changes in a pull request",
		Long: heredoc.Doc(`
			View the diff of a pull request.
		`),
		Example: heredoc.Doc(`
			# View PR diff
			$ gc pr diff 123
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
			return diffRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")

	return cmd
}

func diffRun(opts *DiffOptions) error {
	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := getEnvToken()
	if token == "" {
		return fmt.Errorf("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "environment")

	// Get repository
	owner, repo, err := parseRepo(opts.Repository)
	if err != nil {
		return err
	}

	// Get PR info
	pr, err := api.GetPullRequest(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to get PR: %w", err)
	}

	// Get PR files and diffs
	files, err := api.GetPRFiles(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to get PR diff: %w", err)
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
	if repo == "" {
		return "", "", fmt.Errorf("no repository specified. Use -R owner/repo")
	}

	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repository format: %s", repo)
	}
	return parts[0], parts[1], nil
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITCODE_TOKEN")
}