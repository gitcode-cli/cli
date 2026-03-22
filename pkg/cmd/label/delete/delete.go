// Package delete implements the label delete command
package delete

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

type DeleteOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	Name       string

	// Flags
	Yes bool
}

// NewCmdDelete creates the delete command
func NewCmdDelete(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a label",
		Long: heredoc.Doc(`
			Delete a label from a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Delete a label
			$ gc label delete old-label

			# Skip confirmation
			$ gc label delete old-label --yes
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]

			if runF != nil {
				return runF(opts)
			}
			return deleteRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVarP(&opts.Yes, "yes", "y", false, "Skip confirmation")

	return cmd
}

func deleteRun(opts *DeleteOptions) error {
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

	// Confirm deletion
	if !opts.Yes {
		fmt.Fprintf(opts.IO.ErrOut, "! This will delete label %s\n", cs.Bold(opts.Name))
		fmt.Fprintf(opts.IO.ErrOut, "Type the label name to confirm: ")
		var input string
		fmt.Scanln(&input)
		if input != opts.Name {
			return fmt.Errorf("confirmation did not match label name")
		}
	}

	// Delete label
	err = api.DeleteLabel(client, owner, repo, opts.Name)
	if err != nil {
		return fmt.Errorf("failed to delete label: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Deleted label %s from %s/%s\n", cs.Red("✗"), opts.Name, owner, repo)
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