// Package delete implements the milestone delete command
package delete

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

type DeleteOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	Number     int

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
		Use:   "delete <number>",
		Short: "Delete a milestone",
		Long: heredoc.Doc(`
			Delete a milestone from a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Delete a milestone
			$ gc milestone delete 1 -R owner/repo

			# Skip confirmation
			$ gc milestone delete 1 --yes
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid milestone number: %s", args[0])
			}
			opts.Number = number

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

	// Get milestone for confirmation
	ms, err := api.GetMilestone(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to get milestone: %w", err)
	}

	// Confirm deletion
	if !opts.Yes {
		fmt.Fprintf(opts.IO.ErrOut, "! This will delete milestone #%d %s\n", ms.Number, cs.Bold(ms.Title))
		fmt.Fprintf(opts.IO.ErrOut, "Type the milestone number to confirm: ")
		var input string
		fmt.Scanln(&input)
		if input != strconv.Itoa(opts.Number) {
			return fmt.Errorf("confirmation did not match milestone number")
		}
	}

	// Delete milestone
	err = api.DeleteMilestone(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to delete milestone: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Deleted milestone #%d from %s/%s\n", cs.Red("✗"), opts.Number, owner, repo)
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