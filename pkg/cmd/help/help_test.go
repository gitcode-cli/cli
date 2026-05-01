package help

import (
	"testing"

	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestBuildIndex(t *testing.T) {
	// Create a mock command tree
	root := &cobra.Command{
		Use:   "root",
		Short: "Root command",
	}
	child1 := &cobra.Command{
		Use:   "child1",
		Short: "First child",
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "topic1",
		},
	}
	child2 := &cobra.Command{
		Use:     "child2",
		Short:   "Second child",
		Aliases: []string{"c2", "ch2"},
	}
	hiddenCmd := &cobra.Command{
		Use:    "hidden",
		Short:  "Hidden command",
		Hidden: true,
	}

	root.AddCommand(child1, child2, hiddenCmd)

	index := BuildIndex(root)

	// Should include root, child1, child2 but not hidden
	expectedCount := 3
	if len(index) != expectedCount {
		t.Errorf("BuildIndex returned %d commands, expected %d", len(index), expectedCount)
	}

	// Verify hidden command is excluded
	for _, cmd := range index {
		if cmd.Name == "hidden" {
			t.Error("BuildIndex included hidden command")
		}
	}

	// Verify topic annotation is captured
	for _, cmd := range index {
		if cmd.Name == "child1" && cmd.Topic != "topic1" {
			t.Errorf("child1 topic = %s, expected topic1", cmd.Topic)
		}
	}

	// Verify aliases are captured
	for _, cmd := range index {
		if cmd.Name == "child2" {
			if len(cmd.Aliases) != 2 {
				t.Errorf("child2 aliases count = %d, expected 2", len(cmd.Aliases))
			}
		}
	}
}

func TestSearch(t *testing.T) {
	index := []CommandInfo{
		{Path: "gc pr", Name: "pr", Short: "Manage pull requests", Topic: "pull-requests"},
		{Path: "gc pr view", Name: "view", Short: "View a pull request", Topic: "pull-requests"},
		{Path: "gc issue", Name: "issue", Short: "Manage issues", Topic: "issues"},
		{Path: "gc issue create", Name: "create", Short: "Create an issue", Topic: "issues", Aliases: []string{"new"}},
	}

	tests := []struct {
		name     string
		keyword  string
		expected int
	}{
		{"search by name", "pr", 2},
		{"search by short", "Manage", 2},
		{"search by alias", "new", 1},
		{"search by path", "issue", 2},
		{"no match", "nonexistent", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := Search(index, tt.keyword)
			if len(results) != tt.expected {
				t.Errorf("Search(%s) returned %d results, expected %d", tt.keyword, len(results), tt.expected)
			}
		})
	}
}

func TestFilterByTopic(t *testing.T) {
	index := []CommandInfo{
		{Path: "gc pr", Name: "pr", Short: "Manage pull requests", Topic: "pull-requests"},
		{Path: "gc pr view", Name: "view", Short: "View a pull request", Topic: "pull-requests"},
		{Path: "gc issue", Name: "issue", Short: "Manage issues", Topic: "issues"},
		{Path: "gc release", Name: "release", Short: "Manage releases", Topic: "releases"},
	}

	tests := []struct {
		name     string
		topic    string
		expected int
	}{
		{"filter by pull-requests", "pull-requests", 2},
		{"filter by issues", "issues", 1},
		{"filter by releases", "releases", 1},
		{"no matching topic", "nonexistent", 0},
		{"case insensitive", "PULL-REQUESTS", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := FilterByTopic(index, tt.topic)
			if len(results) != tt.expected {
				t.Errorf("FilterByTopic(%s) returned %d results, expected %d", tt.topic, len(results), tt.expected)
			}
		})
	}
}

func TestCollectTopics(t *testing.T) {
	root := &cobra.Command{
		Use:   "root",
		Short: "Root command",
	}
	prCmd := &cobra.Command{
		Use:   "pr",
		Short: "Manage pull requests",
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "pull-requests",
		},
	}
	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues",
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "issues",
		},
	}
	releaseCmd := &cobra.Command{
		Use:   "release",
		Short: "Manage releases",
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "releases",
		},
	}

	root.AddCommand(prCmd, issueCmd, releaseCmd)

	topics := CollectTopics(root)

	expectedTopics := []string{"issues", "pull-requests", "releases"}
	if len(topics) != len(expectedTopics) {
		t.Errorf("CollectTopics returned %d topics, expected %d", len(topics), len(expectedTopics))
	}

	// Verify topics are sorted
	for i, topic := range topics {
		if topic != expectedTopics[i] {
			t.Errorf("topics[%d] = %s, expected %s", i, topic, expectedTopics[i])
		}
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	index := []CommandInfo{
		{Path: "gc PullRequest", Name: "PullRequest", Short: "Manage Pull Requests"},
	}

	results := Search(index, "pullrequest")
	if len(results) != 1 {
		t.Errorf("Search should be case insensitive, got %d results", len(results))
	}
}

func TestBuildIndexSkipsHelp(t *testing.T) {
	root := &cobra.Command{
		Use:   "root",
		Short: "Root command",
	}
	helpCmd := &cobra.Command{
		Use:   "help",
		Short: "Help command",
	}
	childCmd := &cobra.Command{
		Use:   "child",
		Short: "Child command",
	}

	root.AddCommand(helpCmd, childCmd)

	index := BuildIndex(root)

	for _, cmd := range index {
		if cmd.Name == "help" {
			t.Error("BuildIndex should skip 'help' command")
		}
	}
}
