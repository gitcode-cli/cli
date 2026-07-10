// Package view implements the actions artifact view command.
package view

import (
	"fmt"
	"net/http"
	"strconv"
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

// ViewOptions configures the actions artifact view command.
type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	BaseRepo   func() (string, error)

	Repository string
	ArtifactID string

	JSON bool
}

// NewCmdView creates the actions artifact view command.
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   f.BaseRepo,
	}

	cmd := &cobra.Command{
		Use:   "view <artifact-id>",
		Short: "View an artifact",
		Long: heredoc.Doc(`
			View the detail of a single workflow artifact.

			The artifact id is the ` + "`id`" + ` returned by ` + "`gc actions artifact list`" + `.
			Use --json for a faithful, machine-readable copy of the API response.
		`),
		Example: heredoc.Doc(`
			# View an artifact
			$ gc actions artifact view <artifact-id> -R owner/repo

			# Faithful JSON output
			$ gc actions artifact view <artifact-id> -R owner/repo --json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ArtifactID = strings.TrimSpace(args[0])
			if opts.ArtifactID == "" {
				return cmdutil.NewUsageError("artifact id is required")
			}
			if runF != nil {
				return runF(opts)
			}
			return viewRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func viewRun(opts *ViewOptions) error {
	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}

	repository, err := cmdutil.ResolveRepo(opts.Repository, opts.BaseRepo)
	if err != nil {
		return err
	}
	owner, repo, err := cmdutil.ParseRepo(repository)
	if err != nil {
		return err
	}

	artifact, raw, err := api.GetActionsArtifact(client, owner, repo, opts.ArtifactID)
	if err != nil {
		return fmt.Errorf("failed to get artifact: %w", err)
	}

	if opts.JSON {
		if _, err := opts.IO.Out.Write(raw); err != nil {
			return fmt.Errorf("failed to write JSON output: %w", err)
		}
		_, err := fmt.Fprintln(opts.IO.Out)
		return err
	}

	return printArtifact(opts, artifact)
}

func printArtifact(opts *ViewOptions, a *api.Artifact) error {
	out := opts.IO.Out

	fmt.Fprintf(out, "%s  %s\n", orDash(a.Name), sizeLabel(a.SizeBytes))
	fmt.Fprintf(out, "  artifact id:     %s\n", orDash(a.ID))
	fmt.Fprintf(out, "  name:            %s\n", orDash(a.Name))
	fmt.Fprintf(out, "  size:            %s\n", sizeLabel(a.SizeBytes))
	fmt.Fprintf(out, "  digest:          %s\n", orDash(a.Digest))
	fmt.Fprintf(out, "  workflow id:     %s\n", orDash(a.WorkflowID))
	fmt.Fprintf(out, "  workflow run id: %s\n", orDash(a.WorkflowRunID))
	fmt.Fprintf(out, "  created:         %s\n", formatMsTimeString(a.CreatedAt))
	fmt.Fprintf(out, "  updated:         %s\n", formatMsTimeString(a.UpdatedAt))
	fmt.Fprintf(out, "  expires:         %s\n", formatMsTimeString(a.ExpiresAt))
	return nil
}

func orDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}

func sizeLabel(bytes int64) string {
	const (
		kiB = 1024
		miB = 1024 * 1024
		giB = 1024 * 1024 * 1024
	)
	switch {
	case bytes >= giB:
		return fmt.Sprintf("%.1f GiB", float64(bytes)/float64(giB))
	case bytes >= miB:
		return fmt.Sprintf("%.1f MiB", float64(bytes)/float64(miB))
	case bytes >= kiB:
		return fmt.Sprintf("%.1f KiB", float64(bytes)/float64(kiB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func formatMsTimeString(s string) string {
	if s == "" {
		return "-"
	}
	t, err := strconv.ParseInt(s, 10, 64)
	if err != nil || t <= 0 {
		return "-"
	}
	secs := t
	if t >= msTimestampThreshold {
		secs = t / 1000
	}
	return time.Unix(secs, 0).UTC().Format(time.RFC3339)
}
