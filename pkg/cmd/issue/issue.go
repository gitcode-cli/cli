// Package issue implements the issue command
package issue

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/issue/close"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/issue/comment"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/issue/create"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/issue/edit"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/issue/label"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/issue/list"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/issue/prs"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/issue/reopen"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/issue/view"
)

// NewCmdIssue creates the issue command
func NewCmdIssue(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue <command>",
		Short: "Manage issues",
		Long: heredoc.Doc(`
			Work with GitCode issues.

			An issue is a place to discuss ideas, enhancements, tasks, and bugs.
		`),
		Example: heredoc.Doc(`
			# Create a new issue
			$ gc issue create --title "Bug in login" --body "Description"

			# List issues
			$ gc issue list

			# View an issue
			$ gc issue view 123
		`),
		Annotations: map[string]string{
			"IsCore": "true",
		},
	}

	cmd.AddCommand(create.NewCmdCreate(f, nil))
	cmd.AddCommand(edit.NewCmdEdit(f, nil))
	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(view.NewCmdView(f, nil))
	cmd.AddCommand(close.NewCmdClose(f, nil))
	cmd.AddCommand(reopen.NewCmdReopen(f, nil))
	cmd.AddCommand(comment.NewCmdComment(f, nil))
	cmd.AddCommand(label.NewCmdLabel(f, nil))
	cmd.AddCommand(prs.NewCmdPrs(f, nil))

	return cmd
}