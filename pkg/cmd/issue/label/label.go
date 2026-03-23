// Package label implements the issue label command
package label

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

type LabelOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	Add    []string
	Remove string
	List   bool
}

// NewCmdLabel creates the label command
func NewCmdLabel(f *cmdutil.Factory, runF func(*LabelOptions) error) *cobra.Command {
	opts := &LabelOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
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
				return fmt.Errorf("invalid issue number: %s", args[0])
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

	return cmd
}

func labelRun(opts *LabelOptions) error {
	cs := opts.IO.ColorScheme()

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

	// List labels
	if opts.List {
		issue, err := api.GetIssue(client, owner, repo, opts.Number)
		if err != nil {
			return fmt.Errorf("failed to get issue: %w", err)
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

		fmt.Fprintf(opts.IO.Out, "%s Added labels to issue #%d: %s\n", cs.Green("✓"), opts.Number, formatLabels(added))
		return nil
	}

	// Remove label
	if opts.Remove != "" {
		err := api.RemoveIssueLabel(client, owner, repo, opts.Number, opts.Remove)
		if err != nil {
			return fmt.Errorf("failed to remove label: %w", err)
		}

		fmt.Fprintf(opts.IO.Out, "%s Removed label '%s' from issue #%d\n", cs.Red("✗"), opts.Remove, opts.Number)
		return nil
	}

	return fmt.Errorf("specify --add, --remove, or --list")
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

func formatLabels(labels []*api.Label) string {
	names := make([]string, len(labels))
	for i, l := range labels {
		names[i] = l.Name
	}
	return strings.Join(names, ", ")
}