// Package create implements the repo create command
package create

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type CreateOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Flags
	Name        string
	Description string
	Public      bool
	Private     bool
}

// NewCmdCreate creates the create command
func NewCmdCreate(f *cmdutil.Factory, runF func(*CreateOptions) error) *cobra.Command {
	opts := &CreateOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new repository",
		Long: heredoc.Doc(`
			Create a new GitCode repository.
		`),
		Example: heredoc.Doc(`
			# Create a public repository
			$ gc repo create my-repo --public

			# Create a private repository
			$ gc repo create my-repo --private

			# Create with description
			$ gc repo create my-repo --description "My project"
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

	cmd.Flags().StringVarP(&opts.Description, "description", "d", "", "Description of the repository")
	cmd.Flags().BoolVarP(&opts.Public, "public", "p", false, "Make the repository public")
	cmd.Flags().BoolVarP(&opts.Private, "private", "P", false, "Make the repository private")

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

	// Create repo
	private := opts.Private
	if opts.Public {
		private = false
	}

	// Check if name contains org prefix (org/repo format)
	var repo *api.Repository
	name := opts.Name

	if strings.Contains(name, "/") {
		parts := strings.SplitN(name, "/", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid repository name format: %s (expected 'repo' or 'org/repo')", name)
		}
		org, repoName := parts[0], parts[1]
		repo, err = api.CreateOrgRepo(client, org, &api.CreateRepoOptions{
			Name:        repoName,
			Description: opts.Description,
			Private:     private,
			AutoInit:    true,
		})
	} else {
		repo, err = api.CreateRepo(client, &api.CreateRepoOptions{
			Name:        name,
			Description: opts.Description,
			Private:     private,
			AutoInit:    true,
		})
	}

	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	fmt.Fprintf(opts.IO.Out, "%s Created repository %s\n", cs.Green("✓"), repo.FullName)
	fmt.Fprintf(opts.IO.Out, "  %s\n", repo.HTMLURL)

	return nil
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	if token := os.Getenv("GITCODE_TOKEN"); token != "" {
		return token
	}
	return cmdutil.EnvToken()
}
