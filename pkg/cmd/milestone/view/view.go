// Package view implements the milestone view command
package view

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/browser"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Repository string
	Number     int

	// Flags
	Web bool
}

// NewCmdView creates the view command
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "view <number>",
		Short: "View a milestone",
		Long: heredoc.Doc(`
			View a milestone in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# View a milestone
			$ gc milestone view 1 -R owner/repo
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
			return viewRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVarP(&opts.Web, "web", "w", false, "Open in browser")

	return cmd
}

func viewRun(opts *ViewOptions) error {
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

	// Get milestone
	ms, err := api.GetMilestone(client, owner, repo, opts.Number)
	if err != nil {
		return fmt.Errorf("failed to get milestone: %w", err)
	}

	milestoneURL := fmt.Sprintf("https://gitcode.com/%s/%s/milestones/%d", owner, repo, ms.Number)

	// Open in browser if --web flag is set
	if opts.Web {
		fmt.Fprintf(opts.IO.Out, "Opening %s in your browser.\n", milestoneURL)
		return browser.Open(milestoneURL)
	}

	// Output
	fmt.Fprintf(opts.IO.Out, "\n")
	fmt.Fprintf(opts.IO.Out, "%s #%d\n", cs.Bold(ms.Title), ms.Number)
	fmt.Fprintf(opts.IO.Out, "  State: %s\n", ms.State)
	if ms.DueOn != "" {
		fmt.Fprintf(opts.IO.Out, "  Due: %s\n", ms.DueOn)
	}
	fmt.Fprintf(opts.IO.Out, "\n")
	if ms.Description != "" {
		fmt.Fprintf(opts.IO.Out, "%s\n", ms.Description)
		fmt.Fprintf(opts.IO.Out, "\n")
	}
	fmt.Fprintf(opts.IO.Out, "  https://gitcode.com/%s/%s/milestones/%d\n", owner, repo, ms.Number)
	fmt.Fprintf(opts.IO.Out, "\n")

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
