package view

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/output"
)

func TestNewCmdView(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "view PR",
			args:    []string{"123"},
			wantErr: false,
		},
		{
			name:    "view with web flag",
			args:    []string{"123", "--web"},
			wantErr: false,
		},
		{
			name:    "view with relative time",
			args:    []string{"123", "--time-format", "relative"},
			wantErr: false,
		},
		{
			name:    "no PR number",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "invalid PR number",
			args:    []string{"abc"},
			wantErr: true,
		},
		{
			name:    "invalid time format",
			args:    []string{"123", "--time-format", "yaml"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdView(f, func(opts *ViewOptions) error {
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

func TestParseTimeFormat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    output.TimeFormat
		wantErr bool
	}{
		{name: "default", input: "", want: output.TimeFormatAbsolute},
		{name: "absolute", input: "absolute", want: output.TimeFormatAbsolute},
		{name: "relative", input: "relative", want: output.TimeFormatRelative},
		{name: "invalid", input: "yaml", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTimeFormat(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("parseTimeFormat() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("parseTimeFormat() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("parseTimeFormat() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatViewTime(t *testing.T) {
	now := time.Date(2026, 3, 26, 12, 0, 0, 0, time.UTC)
	ts := api.FlexibleTime{Time: now.Add(-90 * time.Minute)}

	if got := formatViewTime(ts.Time, output.TimeFormatAbsolute, now); got != "2026-03-26 10:30" {
		t.Fatalf("absolute format = %q", got)
	}
	if got := formatViewTime(ts.Time, output.TimeFormatRelative, now); got != "1 hour ago" {
		t.Fatalf("relative format = %q", got)
	}
	if got := formatViewTime(time.Time{}, output.TimeFormatAbsolute, now); got != "unknown" {
		t.Fatalf("zero time format = %q", got)
	}
}

func TestRenderPRView(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	now := time.Date(2026, 3, 26, 12, 0, 0, 0, time.UTC)
	pr := &api.PullRequest{
		Title:        "Improve PR details",
		Number:       7,
		State:        "open",
		HTMLURL:      "https://gitcode.com/owner/repo/pulls/7",
		Body:         "PR body",
		Additions:    10,
		Deletions:    2,
		ChangedFiles: 3,
		Commits:      4,
		Comments:     1,
		CreatedAt:    api.FlexibleTime{Time: now.Add(-4 * time.Hour)},
		UpdatedAt:    api.FlexibleTime{Time: now.Add(-45 * time.Minute)},
		User:         &api.User{Login: "alice"},
		Head:         &api.PRBranch{Ref: "feature"},
		Base:         &api.PRBranch{Ref: "main"},
		Assignees:    []*api.User{&api.User{Login: "bob"}},
		Reviewers:    []*api.User{&api.User{Login: "carol"}},
		Labels:       []*api.Label{&api.Label{Name: "bug"}},
	}
	comments := []api.PRComment{
		{
			Body:      "Looks solid",
			User:      &api.User{Login: "dave"},
			CreatedAt: api.FlexibleTime{Time: now.Add(-10 * time.Minute)},
		},
	}

	var buf bytes.Buffer
	if err := renderPRView(&buf, io.ColorScheme(), pr, comments, output.TimeFormatRelative, now); err != nil {
		t.Fatalf("renderPRView() error = %v", err)
	}

	output := buf.String()
	for _, want := range []string{
		"Improve PR details #7",
		"State: open",
		"Author: alice",
		"Branch: feature -> main",
		"Created: 4 hours ago",
		"Updated: 45 minutes ago",
		"Additions: +10  Deletions: -2  Files: 3",
		"Commits: 4  Comments: 1",
		"Assignees: bob",
		"Reviewers: carol",
		"Labels: bug",
		"PR body",
		"--- Comments (1) ---",
		"dave at 10 minutes ago",
		"Looks solid",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("renderPRView() output missing %q: %s", want, output)
		}
	}
}

func TestRenderPRCommentsEmpty(t *testing.T) {
	var buf bytes.Buffer
	now := time.Date(2026, 3, 26, 12, 0, 0, 0, time.UTC)

	if err := renderPRComments(&buf, nil, output.TimeFormatAbsolute, now); err != nil {
		t.Fatalf("renderPRComments() error = %v", err)
	}

	if got := buf.String(); got != "\n--- No comments ---\n\n" {
		t.Fatalf("renderPRComments() = %q", got)
	}
}
