package help

import (
	"sort"
	"strings"

	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// CommandInfo holds metadata about a command for search and discovery.
type CommandInfo struct {
	Path    string   // Full command path (e.g., "gc pr view")
	Name    string   // Command name (e.g., "view")
	Short   string   // Short description
	Topic   string   // Topic category (from annotation)
	Aliases []string // Command aliases
	Hidden  bool     // Whether command is hidden
}

// BuildIndex traverses the command tree and builds a searchable index.
func BuildIndex(root *cobra.Command) []CommandInfo {
	var index []CommandInfo
	buildIndexRecursive(root, "", &index)
	return index
}

func buildIndexRecursive(cmd *cobra.Command, parentPath string, index *[]CommandInfo) {
	// Skip the help command itself
	if cmd.Name() == "help" {
		return
	}

	// Build current command path
	path := parentPath
	if path != "" {
		path += " "
	}
	path += cmd.Name()

	// Get topic from annotation
	topic := ""
	if cmd.Annotations != nil {
		topic = cmd.Annotations[cmdutil.TopicAnnotation]
	}

	// Get aliases from the Aliases field
	aliases := cmd.Aliases

	// Add to index if not hidden
	if !cmd.Hidden {
		*index = append(*index, CommandInfo{
			Path:    path,
			Name:    cmd.Name(),
			Short:   cmd.Short,
			Topic:   topic,
			Aliases: aliases,
			Hidden:  cmd.Hidden,
		})
	}

	// Recursively process subcommands
	for _, child := range cmd.Commands() {
		buildIndexRecursive(child, path, index)
	}
}

// Search finds commands matching the given keyword.
// Matches against path, name, short description, and aliases.
func Search(index []CommandInfo, keyword string) []CommandInfo {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	var results []CommandInfo

	for _, cmd := range index {
		// Match against path
		if strings.Contains(strings.ToLower(cmd.Path), keyword) {
			results = append(results, cmd)
			continue
		}

		// Match against name
		if strings.Contains(strings.ToLower(cmd.Name), keyword) {
			results = append(results, cmd)
			continue
		}

		// Match against short description
		if strings.Contains(strings.ToLower(cmd.Short), keyword) {
			results = append(results, cmd)
			continue
		}

		// Match against aliases
		for _, alias := range cmd.Aliases {
			if strings.Contains(strings.ToLower(alias), keyword) {
				results = append(results, cmd)
				break
			}
		}
	}

	// Sort results by path
	sort.Slice(results, func(i, j int) bool {
		return results[i].Path < results[j].Path
	})

	return results
}

// FilterByTopic returns commands belonging to a specific topic.
func FilterByTopic(index []CommandInfo, topic string) []CommandInfo {
	topic = strings.ToLower(strings.TrimSpace(topic))
	var results []CommandInfo

	for _, cmd := range index {
		if strings.ToLower(cmd.Topic) == topic {
			results = append(results, cmd)
		}
	}

	// Sort results by path
	sort.Slice(results, func(i, j int) bool {
		return results[i].Path < results[j].Path
	})

	return results
}

// CollectTopics returns all unique topics defined in the command tree.
func CollectTopics(root *cobra.Command) []string {
	index := BuildIndex(root)
	topicSet := make(map[string]bool)

	for _, cmd := range index {
		if cmd.Topic != "" {
			topicSet[cmd.Topic] = true
		}
	}

	topics := make([]string, 0, len(topicSet))
	for topic := range topicSet {
		topics = append(topics, topic)
	}

	// Sort topics alphabetically
	sort.Strings(topics)

	return topics
}

// SetTopicAnnotation adds a topic annotation to a command.
// Deprecated: Use cmdutil.SetTopicAnnotation instead.
func SetTopicAnnotation(cmd *cobra.Command, topic string) {
	cmdutil.SetTopicAnnotation(cmd, topic)
}

// RegisterTopicEnum registers the standard topics as enum values for --topic flag.
func RegisterTopicEnum(cmd *cobra.Command, flagName string) {
	values := append([]string{}, cmdutil.StandardTopics...)
	// Also include any dynamically discovered topics
	topics := CollectTopics(cmd.Root())
	for _, t := range topics {
		found := false
		for _, v := range values {
			if v == t {
				found = true
				break
			}
		}
		if !found {
			values = append(values, t)
		}
	}
	cmdutil.SetFlagEnum(cmd, flagName, values...)
}
