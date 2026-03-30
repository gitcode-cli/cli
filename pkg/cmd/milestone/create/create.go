// Package create implements the milestone create command
package create

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type CreateOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	Title      string

	// Flags
	Description string
	DueDate     string
}

// NewCmdCreate creates the create command
func NewCmdCreate(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "create <title>",
		Short: "Create a milestone",
		Long: heredoc.Doc(`
			Create a new milestone in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Create a milestone (due date defaults to 30 days from now)
			$ gc milestone create "v1.0" -R owner/repo

			# Create with specific due date
			$ gc milestone create "v1.0" --due-date "2024-12-31" -R owner/repo

			# Create with description
			$ gc milestone create "v2.0" --description "Next release" --due-date "2025-01-31" -R owner/repo
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Title = args[0]

			if runF != nil {
				return runF(opts)
			}
			return createRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Description, "description", "d", "", "Description")
	cmd.Flags().StringVarP(&opts.DueDate, "due-date", "D", "", "Due date (YYYY-MM-DD)")

	return cmd
}

func createRun(opts *CreateOptions) error {
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

	// Parse due date (required by GitCode API)
	var dueOn string
	if opts.DueDate != "" {
		_, err := time.Parse("2006-01-02", opts.DueDate)
		if err != nil {
			return fmt.Errorf("invalid due date format, use YYYY-MM-DD: %w", err)
		}
		dueOn = opts.DueDate
	} else {
		// Default to 30 days from now
		dueOn = time.Now().AddDate(0, 0, 30).Format("2006-01-02")
	}

	// Create milestone
	ms, err := api.CreateMilestone(client, owner, repo, &api.CreateMilestoneOptions{
		Title:       opts.Title,
		Description: opts.Description,
		DueOn:       dueOn,
	})
	if err != nil {
		return fmt.Errorf("failed to create milestone: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Created milestone #%d %s in %s/%s\n", cs.Green("✓"), ms.Number, cs.Bold(ms.Title), owner, repo)
	return nil
}

func parseRepo(repo string) (string, string, error) {
	return cmdutil.ParseRepo(repo)
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITCODE_TOKEN")
}
