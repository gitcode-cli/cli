// Package repo implements repository commands
package repo

import (
	"github.com/spf13/cobra"

	"github.com/gitcode-com/gitcode-cli/pkg/cmd/repo/clone"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/repo/create"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/repo/delete"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/repo/fork"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/repo/list"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/repo/view"
	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
)

// NewCmdRepo creates the repo command
func NewCmdRepo(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo <command>",
		Short: "Manage GitCode repositories",
		Long: `Work with GitCode repositories.

Available commands:
  clone    Clone a repository locally
  create   Create a new repository
  fork     Fork a repository
  view     View a repository
  list     List repositories
  delete   Delete a repository`,
	}

	cmd.AddCommand(clone.NewCmdClone(f, nil))
	cmd.AddCommand(create.NewCmdCreate(f, nil))
	cmd.AddCommand(fork.NewCmdFork(f, nil))
	cmd.AddCommand(view.NewCmdView(f, nil))
	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(delete.NewCmdDelete(f, nil))

	return cmd
}