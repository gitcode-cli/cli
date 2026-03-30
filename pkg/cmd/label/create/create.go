// Package create implements the label create command
package create

import (
	"fmt"
	"net/http"
	"os"

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
	Name       string

	// Flags
	Color       string
	Description string
}

// NewCmdCreate creates the create command
func NewCmdCreate(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a label",
		Long: heredoc.Doc(`
			Create a new label in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Create a bug label
			$ gc label create bug --color "#ff0000"

			# Create with description
			$ gc label create enhancement --color "#00ff00" --description "New features"
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]

			if runF != nil {
				return runF(opts)
			}
			return createRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.Color, "color", "c", "#ededed", "Color in hex format")
	cmd.Flags().StringVarP(&opts.Description, "description", "d", "", "Description")

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

	// Create label
	label, err := api.CreateLabel(client, owner, repo, &api.CreateLabelOptions{
		Name:        opts.Name,
		Color:       opts.Color,
		Description: opts.Description,
	})
	if err != nil {
		return fmt.Errorf("failed to create label: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Created label %s in %s/%s\n", cs.Green("✓"), cs.Bold(label.Name), owner, repo)
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
