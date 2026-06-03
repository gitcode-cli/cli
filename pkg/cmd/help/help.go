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
	JSON   bool
}

// NewCmdHelp creates the help command with search and discovery features.
func NewCmdHelp(root *cobra.Command) *cobra.Command {
	opts := &HelpOptions{Root: root}
	commandName := root.Name()

	cmd := &cobra.Command{
		Use:   "help [command]",
		Short: "Help about any command",
		Long: heredoc.Docf(`
			Help provides help for any command in the application.
			Just type %[1]s help [path to command] for full details.

			Additional discovery features:
			--search: Search commands by keyword
			--topics: List all available topics
			--topic:  Filter commands by topic
			--json:   Output in JSON format (for discovery features only)
		`, commandName),
		Example: heredoc.Docf(`
			# Show help for a command
			$ %[1]s help issue
			$ %[1]s help pr view

			# Search for commands containing "issue"
			$ %[1]s help --search issue

			# Search with JSON output
			$ %[1]s help --search issue --json

			# List all available topics
			$ %[1]s help --topics

			# List topics with JSON output
			$ %[1]s help --topics --json

			# Show commands related to pull-requests topic
			$ %[1]s help --topic pull-requests

			# List all commands in JSON format
			$ %[1]s help --json
		`, commandName),
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				// Standard help for specific command - JSON not supported for specific command help
				if opts.JSON {
					return cmdutil.NewUsageError("--json is only supported for discovery features (--search, --topics, --topic, or no args)")
				}
				target, _, err := root.Find(args)
				if err != nil {
					return cmdutil.NewUsageError(err.Error())
				}
				return standardHelp(target)
			}

			// Discovery features with optional JSON output
			if opts.Search != "" {
				if opts.JSON {
					return searchCommandsJSON(root, opts.Search, cmd.OutOrStdout())
				}
				return searchCommands(root, opts.Search, cmd.OutOrStdout())
			}
			if opts.Topics {
				if opts.JSON {
					return listTopicsJSON(root, cmd.OutOrStdout())
				}
				return listTopics(root, cmd.OutOrStdout())
			}
			if opts.Topic != "" {
				if opts.JSON {
					return filterByTopicJSON(root, opts.Topic, cmd.OutOrStdout())
				}
				return filterByTopic(root, opts.Topic, cmd.OutOrStdout())
			}

			// Default: show root help (JSON output for command list)
			if opts.JSON {
				return listCommandsJSON(root, cmd.OutOrStdout())
			}
			return standardHelp(root)
		},
	}

	cmd.Flags().StringVar(&opts.Search, "search", "", "Search commands by keyword")
	cmd.Flags().BoolVar(&opts.Topics, "topics", false, "List all available topics")
	cmd.Flags().StringVar(&opts.Topic, "topic", "", "Filter commands by topic")
	RegisterTopicEnum(cmd, "topic")
	cmd.Flags().BoolVar(&opts.JSON, "json", false, "Output in JSON format")

	return cmd
}

func standardHelp(cmd *cobra.Command) error {
	_ = cmd.Help()

	// Only append discovery hints for root command
	if cmd.Parent() == nil {
		fmt.Fprintf(cmd.OutOrStdout(), "\n%s", discoveryHints(cmd.Name()))
	}

	return nil
}

func discoveryHints(commandName ...string) string {
	name := "gc"
	if len(commandName) > 0 && commandName[0] != "" {
		name = commandName[0]
	}
	return fmt.Sprintf(`Additional discovery features:
  %[1]s help --search <keyword>  Search commands by keyword
  %[1]s help --topics            List all available topics
  %[1]s help --topic <topic>     Filter commands by topic
  %[1]s help --json              Output in JSON format
  %[1]s schema                   Print machine-readable command metadata

For AI agents: Use "%[1]s schema" for structured command discovery.
Use "--json" flag on commands for machine-readable output.
`, name)
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
	fmt.Fprintf(out, "\nUse '%s help --topic <topic>' to see commands in that topic.\n", root.Name())
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

// JSON output types
type commandJSON struct {
	Path    string   `json:"path"`
	Name    string   `json:"name"`
	Short   string   `json:"short"`
	Topic   string   `json:"topic,omitempty"`
	Aliases []string `json:"aliases,omitempty"`
}

type commandsListJSON struct {
	Commands []commandJSON `json:"commands"`
}

type searchResultsJSON struct {
	Query   string        `json:"query"`
	Results []commandJSON `json:"results"`
}

type topicsListJSON struct {
	Topics []string `json:"topics"`
}

type topicCommandsJSON struct {
	Topic    string        `json:"topic"`
	Commands []commandJSON `json:"commands"`
}

func listCommandsJSON(root *cobra.Command, out io.Writer) error {
	index := BuildIndex(root)
	commands := make([]commandJSON, 0, len(index))
	for _, cmd := range index {
		commands = append(commands, commandJSON{
			Path:    cmd.Path,
			Name:    cmd.Name,
			Short:   cmd.Short,
			Topic:   cmd.Topic,
			Aliases: cmd.Aliases,
		})
	}
	return cmdutil.WriteJSON(out, commandsListJSON{Commands: commands})
}

func searchCommandsJSON(root *cobra.Command, keyword string, out io.Writer) error {
	index := BuildIndex(root)
	results := Search(index, keyword)
	commands := make([]commandJSON, 0, len(results))
	for _, cmd := range results {
		commands = append(commands, commandJSON{
			Path:    cmd.Path,
			Name:    cmd.Name,
			Short:   cmd.Short,
			Topic:   cmd.Topic,
			Aliases: cmd.Aliases,
		})
	}
	return cmdutil.WriteJSON(out, searchResultsJSON{
		Query:   keyword,
		Results: commands,
	})
}

func listTopicsJSON(root *cobra.Command, out io.Writer) error {
	topics := CollectTopics(root)
	return cmdutil.WriteJSON(out, topicsListJSON{Topics: topics})
}

func filterByTopicJSON(root *cobra.Command, topic string, out io.Writer) error {
	index := BuildIndex(root)
	results := FilterByTopic(index, topic)
	commands := make([]commandJSON, 0, len(results))
	for _, cmd := range results {
		commands = append(commands, commandJSON{
			Path:    cmd.Path,
			Name:    cmd.Name,
			Short:   cmd.Short,
			Topic:   cmd.Topic,
			Aliases: cmd.Aliases,
		})
	}
	return cmdutil.WriteJSON(out, topicCommandsJSON{
		Topic:    topic,
		Commands: commands,
	})
}
