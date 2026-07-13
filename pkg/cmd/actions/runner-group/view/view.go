// Package view implements the actions runner-group view command.
package view

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// msTimestampThreshold mirrors pkg/output/artifacts.go; tracked by #459 for
// factoring into a shared helper.
const msTimestampThreshold = 100_000_000_000 // 1e11

// ViewOptions configures the actions runner-group view command.
type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	Org           string
	RunnerGroupID string

	JSON bool
}

// NewCmdView creates the actions runner-group view command.
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "view <runner-group-id>",
		Short: "View a runner group",
		Long: heredoc.Doc(`
			View the detail of a single organization runner group.

			The runner group id is the ` + "`id`" + ` returned by ` + "`gc actions runner-group list`" + `.
			Use --json for a faithful, machine-readable copy of the API response.
		`),
		Example: heredoc.Doc(`
			# View a runner group
			$ gc actions runner-group view <runner-group-id> --org my-org

			# Faithful JSON output
			$ gc actions runner-group view <runner-group-id> --org my-org --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.RunnerGroupID = strings.TrimSpace(args[0])
			if opts.RunnerGroupID == "" {
				return cmdutil.NewUsageError("runner group id is required")
			}
			if runF != nil {
				return runF(opts)
			}
			return viewRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Org, "org", "", "Organization path (required)")
	cmd.MarkFlagRequired("org")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func viewRun(opts *ViewOptions) error {
	if opts.Org == "" {
		return cmdutil.NewUsageError("--org is required")
	}
	if opts.RunnerGroupID == "" {
		return cmdutil.NewUsageError("runner group id is required")
	}

	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}

	detail, rawBody, err := api.GetOrgRunnerGroup(client, opts.Org, opts.RunnerGroupID)
	if err != nil {
		return fmt.Errorf("failed to get runner group: %w", err)
	}

	if opts.JSON {
		if _, err := opts.IO.Out.Write(rawBody); err != nil {
			return fmt.Errorf("failed to write JSON output: %w", err)
		}
		_, err := fmt.Fprintln(opts.IO.Out)
		return err
	}

	printRunnerGroupDetail(opts.IO, detail)
	return nil
}

func printRunnerGroupDetail(io *iostreams.IOStreams, d *api.RunnerGroupDetail) {
	cs := io.ColorScheme()
	fmt.Fprintf(io.Out, "%s  %s\n", cs.Bold("ID:"), d.RunnerGroupID)
	fmt.Fprintf(io.Out, "%s    %s\n", cs.Bold("Name:"), d.RunnerGroupName)
	share := "no"
	if d.ShareAll {
		share = "yes"
	}
	fmt.Fprintf(io.Out, "%s   %s\n", cs.Bold("Share All:"), share)
	publicShare := "no"
	if d.ShareAllPublicRepos {
		publicShare = "yes"
	}
	fmt.Fprintf(io.Out, "%s  %s\n", cs.Bold("Share Public:"), publicShare)
	fmt.Fprintf(io.Out, "%s  %d\n", cs.Bold("Shared Repos:"), d.ExplicitSharedRepoCount)
	fmt.Fprintf(io.Out, "%s  %s\n", cs.Bold("Created:"), formatTimestamp(d.CreatedAt))
	fmt.Fprintf(io.Out, "%s  %s\n", cs.Bold("Updated:"), formatTimestamp(d.UpdatedAt))
}

func formatTimestamp(ts int64) string {
	if ts <= 0 {
		return "unknown"
	}
	secs := ts
	if ts >= msTimestampThreshold {
		secs = ts / 1000
	}
	return time.Unix(secs, 0).UTC().Format(time.RFC3339)
}
