// Package milestone implements the milestone command
package milestone

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/milestone/create"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/milestone/list"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/milestone/view"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/milestone/delete"
)

// NewCmdMilestone creates the milestone command
func NewCmdMilestone(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "milestone <command>",
		Short:   "Manage milestones",
		Long:    "Manage milestones in a GitCode repository.",
		Aliases: []string{"ms"},
		Example: heredoc.Doc(`
			# List milestones
			$ gc milestone list -R owner/repo

			# Create a milestone
			$ gc milestone create "v1.0"

			# View a milestone
			$ gc milestone view 1
		`),
	}

	cmd.AddCommand(create.NewCmdCreate(f, nil))
	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(view.NewCmdView(f, nil))
	cmd.AddCommand(delete.NewCmdDelete(f, nil))

	return cmd
}