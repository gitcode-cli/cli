package cmdutil

import (
	"fmt"

	"github.com/spf13/cobra"
)

const FlagEnumAnnotation = "gc.enum"
const TopicAnnotation = "gc.topic"

// StandardTopics is the list of standard topic categories.
var StandardTopics = []string{
	"auth",
	"issues",
	"pull-requests",
	"releases",
	"milestones",
	"repo",
	"commits",
	"labels",
}

// SetFlagEnum records a stable enum set for schema/export consumers.
func SetFlagEnum(cmd *cobra.Command, name string, values ...string) {
	if err := cmd.Flags().SetAnnotation(name, FlagEnumAnnotation, values); err != nil {
		panic(fmt.Sprintf("failed to annotate flag %q: %v", name, err))
	}
}

// SetTopicAnnotation adds a topic annotation to a command.
func SetTopicAnnotation(cmd *cobra.Command, topic string) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations[TopicAnnotation] = topic
}
