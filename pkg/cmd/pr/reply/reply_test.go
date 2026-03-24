// Package reply_test tests the pr reply command
package reply

import (
	"net/http"
	"testing"

	"gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

func TestNewCmdReply(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid command",
			args:    []string{"123", "--discussion", "abc123", "--body", "test reply", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "missing discussion flag",
			args:    []string{"123", "--body", "test reply", "-R", "owner/repo"},
			wantErr: true,
		},
		{
			name:    "missing body flag",
			args:    []string{"123", "--discussion", "abc123", "-R", "owner/repo"},
			wantErr: true,
		},
		{
			name:    "missing PR number",
			args:    []string{"--discussion", "abc123", "--body", "test reply", "-R", "owner/repo"},
			wantErr: true,
		},
		{
			name:    "invalid PR number",
			args:    []string{"abc", "--discussion", "abc123", "--body", "test reply", "-R", "owner/repo"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, _, _ := iostreams.Test()
			f := &cmdutil.Factory{
				IOStreams: io,
				HttpClient: func() (*http.Client, error) {
					return &http.Client{}, nil
				},
			}

			cmd := NewCmdReply(f, nil)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				// Note: This will fail because we don't have a real API
				// but we're testing that the command parses args correctly
				if err != nil {
					// Check that error is not "unknown flag" or parsing error
					if contains(err.Error(), "unknown flag") {
						t.Errorf("unexpected error: %v", err)
					}
					if contains(err.Error(), "invalid PR number") {
						t.Errorf("unexpected error: %v", err)
					}
				}
			}
		})
	}
}

func TestReplyFlags(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	f := &cmdutil.Factory{
		IOStreams: io,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{}, nil
		},
	}

	cmd := NewCmdReply(f, nil)

	// Check that flags are registered
	discussionFlag := cmd.Flags().Lookup("discussion")
	if discussionFlag == nil {
		t.Error("--discussion flag should be registered")
	} else if discussionFlag.Shorthand != "d" {
		t.Errorf("--discussion shorthand should be 'd', got '%s'", discussionFlag.Shorthand)
	}

	bodyFlag := cmd.Flags().Lookup("body")
	if bodyFlag == nil {
		t.Error("--body flag should be registered")
	} else if bodyFlag.Shorthand != "b" {
		t.Errorf("--body shorthand should be 'b', got '%s'", bodyFlag.Shorthand)
	}

	repoFlag := cmd.Flags().Lookup("repo")
	if repoFlag == nil {
		t.Error("--repo flag should be registered")
	} else if repoFlag.Shorthand != "R" {
		t.Errorf("--repo shorthand should be 'R', got '%s'", repoFlag.Shorthand)
	}
}

func TestFormatID(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string ID",
			input:    "abc123",
			expected: "abc123",
		},
		{
			name:     "int ID",
			input:    12345,
			expected: "12345",
		},
		{
			name:     "int64 ID",
			input:    int64(123456789),
			expected: "123456789",
		},
		{
			name:     "float64 ID",
			input:    float64(123.0),
			expected: "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatID(tt.input)
			if result != tt.expected {
				t.Errorf("formatID(%v) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseRepo(t *testing.T) {
	tests := []struct {
		name      string
		repo      string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		{
			name:      "valid repo",
			repo:      "owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:    "empty repo",
			repo:    "",
			wantErr: true,
		},
		{
			name:    "invalid format - no slash",
			repo:    "owner",
			wantErr: true,
		},
		{
			name:    "invalid format - too many slashes",
			repo:    "owner/repo/extra",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := parseRepo(tt.repo)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if owner != tt.wantOwner {
					t.Errorf("owner = %s, want %s", owner, tt.wantOwner)
				}
				if repo != tt.wantRepo {
					t.Errorf("repo = %s, want %s", repo, tt.wantRepo)
				}
			}
		})
	}
}

func TestReplyOptions(t *testing.T) {
	opts := &ReplyOptions{
		PRNumber:     123,
		DiscussionID: "test-discussion",
		Body:         "test body",
	}

	if opts.PRNumber != 123 {
		t.Errorf("PRNumber = %d, want 123", opts.PRNumber)
	}
	if opts.DiscussionID != "test-discussion" {
		t.Errorf("DiscussionID = %s, want test-discussion", opts.DiscussionID)
	}
	if opts.Body != "test body" {
		t.Errorf("Body = %s, want 'test body'", opts.Body)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}