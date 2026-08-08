// Package render provides shared renderers for discussion commands, kept in a
// standalone package to avoid import cycles between a command group and its
// subcommands.
package render

import (
	"fmt"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// PrintDiscussions renders a list of discussions in human-readable form. Shared
// by the organization and repository list commands.
func PrintDiscussions(io *iostreams.IOStreams, discussions []*api.Discussion) {
	cs := io.ColorScheme()
	if len(discussions) == 0 {
		fmt.Fprintf(io.Out, "No discussions found\n")
		return
	}
	fmt.Fprintf(io.Out, "%s\n", cs.Bold("Discussions"))
	for _, d := range discussions {
		state := "open"
		if d.IsClosed != 0 {
			state = "closed"
		}
		author := ""
		if d.Author != nil {
			author = d.Author.Login
		}
		pin := ""
		if d.IsPin != 0 {
			pin = cs.Yellow("📌 ")
		}
		fmt.Fprintf(io.Out, "  %s#%d  %s  (%s)  comments=%d  %s\n",
			pin,
			d.Number,
			d.Title,
			author,
			d.CommentTotal,
			cs.Cyan(state),
		)
		if d.UpdatedAt != "" {
			fmt.Fprintf(io.Out, "       updated %s\n", d.UpdatedAt)
		}
	}
	fmt.Fprintf(io.Out, "\nTotal: %d\n", len(discussions))
}

// PrintDiscussion renders a single discussion in human-readable form. Shared by
// the organization and repository view commands.
func PrintDiscussion(io *iostreams.IOStreams, d *api.Discussion) {
	cs := io.ColorScheme()
	state := "open"
	if d.IsClosed != 0 {
		state = "closed"
	}
	author := ""
	if d.Author != nil {
		author = d.Author.Login
	}
	fmt.Fprintf(io.Out, "%s #%d  %s\n", cs.Bold("Discussion"), d.Number, d.Title)
	fmt.Fprintf(io.Out, "state: %s  comments: %d  author: %s\n", cs.Cyan(state), d.CommentTotal, author)
	if d.CreatedAt != "" {
		fmt.Fprintf(io.Out, "created: %s", d.CreatedAt)
		if d.UpdatedAt != "" {
			fmt.Fprintf(io.Out, "  updated: %s", d.UpdatedAt)
		}
		fmt.Fprintln(io.Out)
	}
	if d.MdContent != "" {
		fmt.Fprintln(io.Out)
		fmt.Fprintf(io.Out, "%s\n", d.MdContent)
	}
}

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
