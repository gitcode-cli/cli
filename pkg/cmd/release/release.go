// Package release implements the release command
package release

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/release/create"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/release/delete"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/release/download"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/release/edit"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/release/list"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/release/upload"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/release/view"
)

// NewCmdRelease creates the release command
func NewCmdRelease(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release <command>",
		Short: "Manage releases",
		Long: heredoc.Doc(`
			Manage GitCode releases.

			Releases are immutable snapshots of your code at a specific point in time,
			often associated with version tags.
		`),
		Example: heredoc.Doc(`
			# Create a new release
			$ gc release create v1.0.0

			# List releases
			$ gc release list

			# View a release
			$ gc release view v1.0.0

			# Delete a release
			$ gc release delete v1.0.0
		`),
		Annotations: map[string]string{
			"IsCore": "true",
		},
	}

	cmd.AddCommand(create.NewCmdCreate(f, nil))
	cmd.AddCommand(delete.NewCmdDelete(f, nil))
	cmd.AddCommand(download.NewCmdDownload(f, nil))
	cmd.AddCommand(edit.NewCmdEdit(f, nil))
	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(upload.NewCmdUpload(f, nil))
	cmd.AddCommand(view.NewCmdView(f, nil))

	return cmd
}