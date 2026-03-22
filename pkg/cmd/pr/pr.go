// Package pr implements the pr command
package pr

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/pr/create"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/pr/list"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/pr/view"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/pr/checkout"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/pr/merge"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/pr/close"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/pr/reopen"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/pr/review"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/pr/diff"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/pr/ready"
)

// NewCmdPR creates the pr command
func NewCmdPR(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pr <command>",
		Short: "Manage pull requests",
		Long: heredoc.Doc(`
			Work with GitCode pull requests.

			A pull request is a proposal to merge changes from one branch into another.
		`),
		Example: heredoc.Doc(`
			# Create a new pull request
			$ gc pr create --title "Feature" --body "Description"

			# List pull requests
			$ gc pr list

			# View a pull request
			$ gc pr view 123

			# Review a pull request
			$ gc pr review 123 --approve
		`),
		Annotations: map[string]string{
			"IsCore": "true",
		},
	}

	cmd.AddCommand(create.NewCmdCreate(f, nil))
	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(view.NewCmdView(f, nil))
	cmd.AddCommand(checkout.NewCmdCheckout(f, nil))
	cmd.AddCommand(merge.NewCmdMerge(f, nil))
	cmd.AddCommand(close.NewCmdClose(f, nil))
	cmd.AddCommand(reopen.NewCmdReopen(f, nil))
	cmd.AddCommand(review.NewCmdReview(f, nil))
	cmd.AddCommand(diff.NewCmdDiff(f, nil))
	cmd.AddCommand(ready.NewCmdReady(f, nil))

	return cmd
}