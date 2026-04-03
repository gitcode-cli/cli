package edit

import (
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdEdit(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "edit title",
			args:    []string{"123", "--title", "New title", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "edit body",
			args:    []string{"123", "--body", "New body", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "edit state close",
			args:    []string{"123", "--state", "close", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "edit state reopen",
			args:    []string{"123", "--state", "reopen", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "edit with labels",
			args:    []string{"123", "--label", "bug,enhancement", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "edit with assignees",
			args:    []string{"123", "--assignee", "user1", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "edit with milestone",
			args:    []string{"123", "--milestone", "5", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "edit with security-hole",
			args:    []string{"123", "--security-hole", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "edit multiple fields",
			args:    []string{"123", "--title", "Title", "--body", "Body", "--label", "bug", "-R", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "missing issue number",
			args:    []string{"-R", "owner/repo"},
			wantErr: true,
		},
		{
			name:    "invalid issue number",
			args:    []string{"abc", "-R", "owner/repo"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdEdit(f, func(opts *EditOptions) error {
				// Mock run function - just validate
				return nil
			})
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEditRun_NoEditOptions(t *testing.T) {
	f := cmdutil.TestFactory()
	opts := &EditOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Number:     123,
		Repository: "owner/repo",
	}

	err := editRun(opts)
	if err == nil {
		t.Error("Expected error when no edit options provided")
	}
	if err.Error() != "at least one edit option is required (e.g., --title, --body, --state, --assignee, --label, --milestone, --security-hole)" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestParseRepo(t *testing.T) {
	tests := []struct {
		repo      string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		{"owner/repo", "owner", "repo", false},
		{"gitcode-cli/cli", "gitcode-cli", "cli", false},
		{"", "gitcode-cli", "cli", false},
		{"invalid", "", "", true},
		{"too/many/parts", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.repo, func(t *testing.T) {
			owner, repo, err := parseRepo(tt.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if owner != tt.wantOwner {
				t.Errorf("parseRepo() owner = %v, want %v", owner, tt.wantOwner)
			}
			if repo != tt.wantRepo {
				t.Errorf("parseRepo() repo = %v, want %v", repo, tt.wantRepo)
			}
		})
	}
}
