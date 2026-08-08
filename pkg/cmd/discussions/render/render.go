// Package render provides shared renderers for discussion commands, kept in a
// standalone package to avoid import cycles between a command group and its
// subcommands.
package render

import (
	"fmt"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// PrintComments renders a list of discussion comments/replies in human-readable
// form. Shared by the organization and repository comment commands.
func PrintComments(io *iostreams.IOStreams, cs []*api.DiscussionComment) {
	c := io.ColorScheme()
	if len(cs) == 0 {
		fmt.Fprintf(io.Out, "No comments found\n")
		return
	}
	fmt.Fprintf(io.Out, "%s\n", c.Bold("Comments"))
	for _, cm := range cs {
		author := ""
		if cm.Author != nil {
			author = cm.Author.Login
		}
		hidden := ""
		if cm.IsHide != 0 {
			hidden = c.Yellow("[hidden] ")
		}
		deleted := ""
		if cm.IsDeleted != 0 {
			deleted = c.Red("[deleted] ")
		}
		fmt.Fprintf(io.Out, "  %s%s  %s  (likes=%d replies=%d)\n", deleted, hidden, author, cm.LikeTotal, cm.ReplyTotal)
		body := cm.MdContent
		if body == "" {
			body = cm.Content
		}
		if body != "" {
			fmt.Fprintf(io.Out, "       %s\n", body)
		}
		if cm.CreatedAt != "" {
			fmt.Fprintf(io.Out, "       created %s\n", cm.CreatedAt)
		}
	}
	fmt.Fprintf(io.Out, "\nTotal: %d\n", len(cs))
}
