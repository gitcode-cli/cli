// Package issue implements the issue command
package issue

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/issue/create"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/issue/list"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/issue/view"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/issue/close"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/issue/reopen"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/issue/comment"
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
	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(view.NewCmdView(f, nil))
	cmd.AddCommand(close.NewCmdClose(f, nil))
	cmd.AddCommand(reopen.NewCmdReopen(f, nil))
	cmd.AddCommand(comment.NewCmdComment(f, nil))

	return cmd
}