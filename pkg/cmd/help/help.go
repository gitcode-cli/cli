// Package help implements the help command with search and discovery features.
package help

import (
	"fmt"
	"io"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

type HelpOptions struct {
	Root   *cobra.Command
	Search string
	Topics bool
	Topic  string
}

// NewCmdHelp creates the help command with search and discovery features.
func NewCmdHelp(root *cobra.Command) *cobra.Command {
	opts := &HelpOptions{Root: root}

	cmd := &cobra.Command{
		Use:   "help [command]",
		Short: "Help about any command",
		Long: heredoc.Doc(`
			Help provides help for any command in the application.
			Just type gc help [path to command] for full details.

			Additional discovery features:
			--search: Search commands by keyword
			--topics: List all available topics
			--topic:  Filter commands by topic
		`),
		Example: heredoc.Doc(`
			# Show help for a command
			$ gc help issue
			$ gc help pr view

			# Search for commands containing "issue"
			$ gc help --search issue

			# List all available topics
			$ gc help --topics

			# Show commands related to pull-requests topic
			$ gc help --topic pull-requests
		`),
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				// Standard help for specific command
				target, _, err := root.Find(args)
				if err != nil {
					return cmdutil.NewUsageError(err.Error())
				}
				return standardHelp(target)
			}

			// Discovery features
			if opts.Search != "" {
				return searchCommands(root, opts.Search, cmd.OutOrStdout())
			}
			if opts.Topics {
				return listTopics(root, cmd.OutOrStdout())
			}
			if opts.Topic != "" {
				return filterByTopic(root, opts.Topic, cmd.OutOrStdout())
			}

			// Default: show root help
			return standardHelp(root)
		},
	}

	cmd.Flags().StringVar(&opts.Search, "search", "", "Search commands by keyword")
	cmd.Flags().BoolVar(&opts.Topics, "topics", false, "List all available topics")
	cmd.Flags().StringVar(&opts.Topic, "topic", "", "Filter commands by topic")

	return cmd
}

func standardHelp(cmd *cobra.Command) error {
	cmd.Help()
	return nil
}

func searchCommands(root *cobra.Command, keyword string, out io.Writer) error {
	index := BuildIndex(root)
	results := Search(index, keyword)

	if len(results) == 0 {
		fmt.Fprintf(out, "No commands found matching '%s'\n", keyword)
		return nil
	}

	fmt.Fprintf(out, "Commands matching '%s':\n\n", keyword)
	for _, cmd := range results {
		fmt.Fprintf(out, "  %s\t%s\n", cmd.Path, cmd.Short)
		if len(cmd.Aliases) > 0 {
			fmt.Fprintf(out, "    Aliases: %s\n", strings.Join(cmd.Aliases, ", "))
		}
	}
	fmt.Fprintf(out, "\n")
	return nil
}

func listTopics(root *cobra.Command, out io.Writer) error {
	topics := CollectTopics(root)

	if len(topics) == 0 {
		fmt.Fprintf(out, "No topics defined\n")
		return nil
	}

	fmt.Fprintf(out, "Available topics:\n\n")
	for _, topic := range topics {
		fmt.Fprintf(out, "  %s\n", topic)
	}
	fmt.Fprintf(out, "\nUse 'gc help --topic <topic>' to see commands in that topic.\n")
	return nil
}

func filterByTopic(root *cobra.Command, topic string, out io.Writer) error {
	index := BuildIndex(root)
	results := FilterByTopic(index, topic)

	if len(results) == 0 {
		fmt.Fprintf(out, "No commands found for topic '%s'\n", topic)
		return nil
	}

	fmt.Fprintf(out, "Commands in topic '%s':\n\n", topic)
	for _, cmd := range results {
		fmt.Fprintf(out, "  %s\t%s\n", cmd.Path, cmd.Short)
		if len(cmd.Aliases) > 0 {
			fmt.Fprintf(out, "    Aliases: %s\n", strings.Join(cmd.Aliases, ", "))
		}
	}
	fmt.Fprintf(out, "\n")
	return nil
}
