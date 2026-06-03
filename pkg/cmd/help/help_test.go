package help

import (
	"bytes"
	"encoding/json"
	"strings"
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

func TestNewCmdHelp(t *testing.T) {
	root := &cobra.Command{
		Use:   "gc",
		Short: "GitCode CLI",
	}
	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues",
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "issues",
		},
	}
	prCmd := &cobra.Command{
		Use:   "pr",
		Short: "Manage pull requests",
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "pull-requests",
		},
	}
	root.AddCommand(issueCmd, prCmd)

	tests := []struct {
		name       string
		args       []string
		contains   string
		notContain string
	}{
		{
			name:     "search by keyword",
			args:     []string{"--search", "issue"},
			contains: "Commands matching 'issue'",
		},
		{
			name:     "search no match",
			args:     []string{"--search", "nonexistent"},
			contains: "No commands found matching 'nonexistent'",
		},
		{
			name:     "list topics",
			args:     []string{"--topics"},
			contains: "Available topics:",
		},
		{
			name:     "filter by topic",
			args:     []string{"--topic", "issues"},
			contains: "Commands in topic 'issues'",
		},
		{
			name:       "filter by nonexistent topic",
			args:       []string{"--topic", "nonexistent"},
			contains:   "No commands found for topic 'nonexistent'",
			notContain: "Commands in topic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new help command for each test to avoid flag state pollution
			testHelpCmd := NewCmdHelp(root)
			buf := &bytes.Buffer{}
			testHelpCmd.SetOut(buf)
			testHelpCmd.SetArgs(tt.args)
			_ = testHelpCmd.Execute()

			output := buf.String()
			if tt.contains != "" && !bytes.Contains(buf.Bytes(), []byte(tt.contains)) {
				t.Errorf("output should contain '%s', got:\n%s", tt.contains, output)
			}
			if tt.notContain != "" && bytes.Contains(buf.Bytes(), []byte(tt.notContain)) {
				t.Errorf("output should not contain '%s', got:\n%s", tt.notContain, output)
			}
		})
	}
}

func TestSearchCommands(t *testing.T) {
	root := &cobra.Command{
		Use:   "gc",
		Short: "GitCode CLI",
	}
	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues",
	}
	prCmd := &cobra.Command{
		Use:   "pr",
		Short: "Manage pull requests",
	}
	root.AddCommand(issueCmd, prCmd)

	tests := []struct {
		name     string
		keyword  string
		expected int
	}{
		{"match issue", "issue", 1},
		{"match pr", "pr", 1},
		{"match both", "manage", 2},
		{"empty keyword", "", 3},
		{"no match", "xyz", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := searchCommands(root, tt.keyword, buf)
			if err != nil {
				t.Errorf("searchCommands returned error: %v", err)
			}
			output := buf.String()
			if tt.expected == 0 && !bytes.Contains(buf.Bytes(), []byte("No commands found")) {
				t.Errorf("expected 'No commands found' for keyword '%s', got:\n%s", tt.keyword, output)
			}
		})
	}
}

func TestListTopics(t *testing.T) {
	root := &cobra.Command{
		Use:   "gc",
		Short: "GitCode CLI",
	}
	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues",
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "issues",
		},
	}
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication",
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "auth",
		},
	}
	root.AddCommand(issueCmd, authCmd)

	buf := &bytes.Buffer{}
	err := listTopics(root, buf)
	if err != nil {
		t.Errorf("listTopics returned error: %v", err)
	}

	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("Available topics:")) {
		t.Errorf("output should contain 'Available topics:', got:\n%s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("auth")) {
		t.Errorf("output should contain topic 'auth', got:\n%s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("issues")) {
		t.Errorf("output should contain topic 'issues', got:\n%s", output)
	}
}

func TestFilterByTopicOutput(t *testing.T) {
	root := &cobra.Command{
		Use:   "gc",
		Short: "GitCode CLI",
	}
	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues",
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "issues",
		},
	}
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication",
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "auth",
		},
	}
	root.AddCommand(issueCmd, authCmd)

	tests := []struct {
		name     string
		topic    string
		contains string
	}{
		{"match issues topic", "issues", "Commands in topic 'issues'"},
		{"match auth topic", "auth", "Commands in topic 'auth'"},
		{"no match topic", "nonexistent", "No commands found for topic 'nonexistent'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := filterByTopic(root, tt.topic, buf)
			if err != nil {
				t.Errorf("filterByTopic returned error: %v", err)
			}
			if !bytes.Contains(buf.Bytes(), []byte(tt.contains)) {
				t.Errorf("output should contain '%s', got:\n%s", tt.contains, buf.String())
			}
		})
	}
}

func TestStandardHelpWithDiscoveryHints(t *testing.T) {
	root := &cobra.Command{
		Use:   "gc",
		Short: "GitCode CLI",
	}
	childCmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues",
	}
	root.AddCommand(childCmd)

	buf := &bytes.Buffer{}
	root.SetOut(buf)

	err := standardHelp(root)
	if err != nil {
		t.Errorf("standardHelp returned error: %v", err)
	}

	output := buf.String()

	// Verify discovery hints are present for root
	if !strings.Contains(output, "Additional discovery features:") {
		t.Error("root help should contain discovery hints section")
	}
	if !strings.Contains(output, "--search") {
		t.Error("root help should mention --search flag")
	}
	if !strings.Contains(output, "gc schema") {
		t.Error("root help should mention gc schema command")
	}
	if !strings.Contains(output, "For AI agents:") {
		t.Error("root help should contain AI agent guidance")
	}
}

func TestStandardHelpNoHintsForSubcommand(t *testing.T) {
	root := &cobra.Command{
		Use:   "gc",
		Short: "GitCode CLI",
	}
	childCmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues",
	}
	root.AddCommand(childCmd)

	buf := &bytes.Buffer{}
	childCmd.SetOut(buf)

	err := standardHelp(childCmd)
	if err != nil {
		t.Errorf("standardHelp returned error: %v", err)
	}

	output := buf.String()

	// Verify no discovery hints for subcommand
	if strings.Contains(output, "Additional discovery features:") {
		t.Error("subcommand help should NOT contain discovery hints")
	}
}

func TestDiscoveryHintsFormat(t *testing.T) {
	hints := discoveryHints()

	// Verify all required elements are present
	requiredStrings := []string{
		"Additional discovery features:",
		"--search",
		"--topics",
		"--topic",
		"gc schema",
		"For AI agents:",
		"--json",
		"machine-readable",
	}

	for _, s := range requiredStrings {
		if !strings.Contains(hints, s) {
			t.Errorf("discoveryHints should contain '%s'", s)
		}
	}
}

// JSON output tests
func TestListCommandsJSON(t *testing.T) {
	root := &cobra.Command{
		Use:   "gc",
		Short: "GitCode CLI",
	}
	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues",
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "issues",
		},
	}
	prCmd := &cobra.Command{
		Use:   "pr",
		Short: "Manage pull requests",
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "pull-requests",
		},
	}
	root.AddCommand(issueCmd, prCmd)

	buf := &bytes.Buffer{}
	err := listCommandsJSON(root, buf)
	if err != nil {
		t.Errorf("listCommandsJSON returned error: %v", err)
	}

	// Verify JSON is valid and has expected structure
	var result commandsListJSON
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Errorf("JSON output is invalid: %v, output:\n%s", err, buf.String())
	}

	// Should have 3 commands (root, issue, pr)
	if len(result.Commands) < 3 {
		t.Errorf("expected at least 3 commands, got %d", len(result.Commands))
	}

	// Verify each command has required fields
	for _, cmd := range result.Commands {
		if cmd.Path == "" {
			t.Error("command should have 'path' field")
		}
		if cmd.Name == "" {
			t.Error("command should have 'name' field")
		}
		if cmd.Short == "" {
			t.Error("command should have 'short' field")
		}
	}
}

func TestSearchCommandsJSON(t *testing.T) {
	root := &cobra.Command{
		Use:   "gc",
		Short: "GitCode CLI",
	}
	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues",
	}
	prCmd := &cobra.Command{
		Use:   "pr",
		Short: "Manage pull requests",
	}
	root.AddCommand(issueCmd, prCmd)

	buf := &bytes.Buffer{}
	err := searchCommandsJSON(root, "issue", buf)
	if err != nil {
		t.Errorf("searchCommandsJSON returned error: %v", err)
	}

	var result searchResultsJSON
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Errorf("JSON output is invalid: %v, output:\n%s", err, buf.String())
	}

	// Verify query field
	if result.Query != "issue" {
		t.Errorf("expected query 'issue', got '%s'", result.Query)
	}

	// Should match at least 1 command (issue)
	if len(result.Results) < 1 {
		t.Errorf("expected at least 1 result for 'issue', got %d", len(result.Results))
	}
}

func TestListTopicsJSON(t *testing.T) {
	root := &cobra.Command{
		Use:   "gc",
		Short: "GitCode CLI",
	}
	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues",
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "issues",
		},
	}
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication",
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "auth",
		},
	}
	root.AddCommand(issueCmd, authCmd)

	buf := &bytes.Buffer{}
	err := listTopicsJSON(root, buf)
	if err != nil {
		t.Errorf("listTopicsJSON returned error: %v", err)
	}

	var result topicsListJSON
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Errorf("JSON output is invalid: %v, output:\n%s", err, buf.String())
	}

	// Should have 2 topics (auth, issues)
	if len(result.Topics) != 2 {
		t.Errorf("expected 2 topics, got %d", len(result.Topics))
	}

	// Verify topics are sorted
	expectedTopics := []string{"auth", "issues"}
	for i, topic := range result.Topics {
		if topic != expectedTopics[i] {
			t.Errorf("topics[%d] = '%s', expected '%s'", i, topic, expectedTopics[i])
		}
	}
}

func TestFilterByTopicJSON(t *testing.T) {
	root := &cobra.Command{
		Use:   "gc",
		Short: "GitCode CLI",
	}
	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues",
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "issues",
		},
	}
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication",
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "auth",
		},
	}
	root.AddCommand(issueCmd, authCmd)

	buf := &bytes.Buffer{}
	err := filterByTopicJSON(root, "issues", buf)
	if err != nil {
		t.Errorf("filterByTopicJSON returned error: %v", err)
	}

	var result topicCommandsJSON
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Errorf("JSON output is invalid: %v, output:\n%s", err, buf.String())
	}

	// Verify topic field
	if result.Topic != "issues" {
		t.Errorf("expected topic 'issues', got '%s'", result.Topic)
	}

	// Should have at least 1 command (issue)
	if len(result.Commands) < 1 {
		t.Errorf("expected at least 1 command for 'issues' topic, got %d", len(result.Commands))
	}
}

func TestJSONFlagWithSpecificCommand(t *testing.T) {
	root := &cobra.Command{
		Use:   "gc",
		Short: "GitCode CLI",
	}
	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues",
	}
	root.AddCommand(issueCmd)

	testHelpCmd := NewCmdHelp(root)
	buf := &bytes.Buffer{}
	testHelpCmd.SetOut(buf)
	testHelpCmd.SetArgs([]string{"issue", "--json"})

	err := testHelpCmd.Execute()
	if err == nil {
		t.Error("expected error when using --json with specific command")
	}
	// Should be a usage error
	if !strings.Contains(err.Error(), "--json is only supported for discovery features") {
		t.Errorf("expected usage error message, got: %v", err)
	}
}
