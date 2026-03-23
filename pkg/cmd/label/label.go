// Package label implements the label command
package label

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/label/create"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/label/list"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/label/delete"
)

// NewCmdLabel creates the label command
func NewCmdLabel(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label <command>",
		Short: "Manage labels",
		Long: heredoc.Doc(`
			Manage labels in a GitCode repository.

			Labels are used to categorize issues and pull requests.
		`),
		Example: heredoc.Doc(`
			# List labels
			$ gc label list -R owner/repo

			# Create a label
			$ gc label create "bug" --color "#ff0000"

			# Delete a label
			$ gc label delete "old-label"
		`),
	}

	cmd.AddCommand(create.NewCmdCreate(f, nil))
	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(delete.NewCmdDelete(f, nil))

	return cmd
}