// Package label implements the issue label command
package label

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

type LabelOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	Add    []string
	Remove string
	List   bool
	JSON   bool
}

// LabelResult represents the JSON output for issue label operations
type LabelResult struct {
	Number int      `json:"number"`
	Owner  string   `json:"owner"`
	Repo   string   `json:"repo"`
	Action string   `json:"action"`
	Labels []string `json:"labels,omitempty"`
	Label  string   `json:"label,omitempty"` // For remove action
}

// NewCmdLabel creates the label command
func NewCmdLabel(f *cmdutil.Factory, runF func(*LabelOptions) error) *cobra.Command {
	opts := &LabelOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "label <number>",
		Short: "Manage labels on an issue",
		Long: heredoc.Doc(`
			Add, remove, or list labels on an issue.
		`),
		Example: heredoc.Doc(`
			# Add labels to an issue
			$ gc issue label 123 --add bug,enhancement -R owner/repo

			# Remove a label from an issue
			$ gc issue label 123 --remove bug -R owner/repo

			# List labels on an issue
			$ gc issue label 123 --list -R owner/repo
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return cmdutil.NewUsageError(fmt.Sprintf("invalid issue number: %s", args[0]))
			}
			opts.Number = number

			if runF != nil {
				return runF(opts)
			}
			return labelRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringSliceVarP(&opts.Add, "add", "a", nil, "Add labels (comma-separated)")
	cmd.Flags().StringVarP(&opts.Remove, "remove", "r", "", "Remove a label")
	cmd.Flags().BoolVarP(&opts.List, "list", "l", false, "List labels")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func labelRun(opts *LabelOptions) error {
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

	// List labels
	if opts.List {
		issue, err := api.GetIssue(client, owner, repo, opts.Number)
		if err != nil {
			return cmdutil.WrapNotFound(err, "issue #%d not found in %s/%s", opts.Number, owner, repo)
		}

		var labelNames []string
		for _, label := range issue.Labels {
			labelNames = append(labelNames, label.Name)
		}

		if opts.JSON {
			result := LabelResult{
				Number: opts.Number,
				Owner:  owner,
				Repo:   repo,
				Action: "list",
				Labels: labelNames,
			}
			return cmdutil.WriteJSON(opts.IO.Out, result)
		}

		if len(issue.Labels) == 0 {
			fmt.Fprintf(opts.IO.Out, "No labels on issue #%s\n", issue.Number)
			return nil
		}

		fmt.Fprintf(opts.IO.Out, "Labels on issue #%s:\n", issue.Number)
		for _, label := range issue.Labels {
			fmt.Fprintf(opts.IO.Out, "  %s\n", cs.Cyan(label.Name))
		}
		return nil
	}

	// Add labels
	if len(opts.Add) > 0 {
		// Parse comma-separated labels
		var labels []string
		for _, l := range opts.Add {
			labels = append(labels, strings.Split(l, ",")...)
		}

		added, err := api.AddIssueLabels(client, owner, repo, opts.Number, labels)
		if err != nil {
			return fmt.Errorf("failed to add labels: %w", err)
		}

		var addedNames []string
		for _, l := range added {
			addedNames = append(addedNames, l.Name)
		}

		if opts.JSON {
			result := LabelResult{
				Number: opts.Number,
				Owner:  owner,
				Repo:   repo,
				Action: "add",
				Labels: addedNames,
			}
			return cmdutil.WriteJSON(opts.IO.Out, result)
		}

		fmt.Fprintf(opts.IO.Out, "%s Added labels to issue #%d: %s\n", cs.Green("✓"), opts.Number, formatLabels(added))
		return nil
	}

	// Remove label
	if opts.Remove != "" {
		err := api.RemoveIssueLabel(client, owner, repo, opts.Number, opts.Remove)
		if err != nil {
			return fmt.Errorf("failed to remove label: %w", err)
		}

		if opts.JSON {
			result := LabelResult{
				Number: opts.Number,
				Owner:  owner,
				Repo:   repo,
				Action: "remove",
				Label:  opts.Remove,
			}
			return cmdutil.WriteJSON(opts.IO.Out, result)
		}

		fmt.Fprintf(opts.IO.Out, "%s Removed label '%s' from issue #%d\n", cs.Red("✗"), opts.Remove, opts.Number)
		return nil
	}

	return cmdutil.NewUsageError("specify --add, --remove, or --list")
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

func formatLabels(labels []*api.Label) string {
	names := make([]string, len(labels))
	for i, l := range labels {
		names[i] = l.Name
	}
	return strings.Join(names, ", ")
}
