// Package repo implements repository commands
package repo

import (
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/repo/clone"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/repo/create"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/repo/delete"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/repo/fork"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/repo/list"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/repo/stats"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/repo/sync"
	"gitcode.com/gitcode-cli/cli/pkg/cmd/repo/view"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
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
  sync     Sync a local directory into another repository and create a PR
  delete   Delete a repository`,
	}

	cmd.AddCommand(clone.NewCmdClone(f, nil))
	cmd.AddCommand(create.NewCmdCreate(f, nil))
	cmd.AddCommand(fork.NewCmdFork(f, nil))
	cmd.AddCommand(view.NewCmdView(f, nil))
	cmd.AddCommand(list.NewCmdList(f, nil))
	cmd.AddCommand(sync.NewCmdSync(f, nil))
	cmd.AddCommand(delete.NewCmdDelete(f, nil))
	cmd.AddCommand(stats.NewCmdStats(f, nil))

	return cmd
}
