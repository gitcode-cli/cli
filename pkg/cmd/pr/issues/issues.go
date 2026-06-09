// Package issues implements the pr issues command
package issues

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

type IssuesOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	JSON bool
}

// NewCmdIssues creates the pr issues command
func NewCmdIssues(f *cmdutil.Factory, runF func(*IssuesOptions) error) *cobra.Command {
	opts := &IssuesOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "issues <number>",
		Short: "List issues linked to a pull request",
		Long: heredoc.Doc(`
			List issues linked to a pull request in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# List issues linked to a PR
			$ gc pr issues 123 -R owner/repo

			# Output as JSON
			$ gc pr issues 123 -R owner/repo --json
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
			return issuesRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func issuesRun(opts *IssuesOptions) error {
	cs := opts.IO.ColorScheme()

	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}

	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}
	owner, repo, err := parseRepo(repository)
	if err != nil {
		return err
	}

	issues, err := api.ListPRIssues(client, owner, repo, opts.Number)
	if err != nil {
		return cmdutil.WrapNotFound(err, "PR #%d not found in %s/%s", opts.Number, owner, repo)
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, issues)
	}

	if len(issues) == 0 {
		fmt.Fprintf(opts.IO.Out, "\nNo linked issues found for PR #%d\n\n", opts.Number)
		return nil
	}

	fmt.Fprintf(opts.IO.Out, "\nLinked issues for PR #%d (%d):\n\n", opts.Number, len(issues))
	for _, issue := range issues {
		state := issue.State
		if state == "" {
			state = "unknown"
		}
		num := issue.Number
		title := issue.Title
		fmt.Fprintf(opts.IO.Out, "  %s #%s  %s\n", stateLabel(cs, state), num, title)
	}
	fmt.Fprintf(opts.IO.Out, "\n")

	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

func stateLabel(cs *iostreams.ColorScheme, state string) string {
	switch strings.ToLower(state) {
	case "open":
		return cs.Green("open")
	case "closed":
		return cs.Red("closed")
	case "merged":
		return cs.Magenta("merged")
	default:
		return state
	}
}
