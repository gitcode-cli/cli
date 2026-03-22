// Package list implements the label list command
package list

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/gitcode-com/gitcode-cli/api"
	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
	"github.com/gitcode-com/gitcode-cli/pkg/iostreams"
)

type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string

	// Flags
	Limit int
}

// NewCmdList creates the list command
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List labels",
		Long: heredoc.Doc(`
			List labels in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# List labels
			$ gc label list -R owner/repo
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 50, "Maximum number of labels to list")

	return cmd
}

func listRun(opts *ListOptions) error {
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
	labels, err := api.ListRepoLabels(client, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to list labels: %w", err)
	}

	// Output
	if len(labels) == 0 {
		fmt.Fprintf(opts.IO.Out, "No labels found\n")
		return nil
	}

	fmt.Fprintf(opts.IO.Out, "\n")
	for _, label := range labels {
		color := label.Color
		if !strings.HasPrefix(color, "#") {
			color = "#" + color
		}
		fmt.Fprintf(opts.IO.Out, "%s  %s\n", cs.Bold(label.Name), label.Description)
	}
	fmt.Fprintf(opts.IO.Out, "\n")

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