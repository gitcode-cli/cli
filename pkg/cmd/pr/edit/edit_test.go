package edit

import (
	"fmt"
	"net/http"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

func TestNewCmdEdit(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no arguments",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "valid number with title",
			args:    []string{"123", "--repo", "owner/repo", "--title", "New title"},
			wantErr: false,
		},
		{
			name:    "invalid number",
			args:    []string{"abc"},
			wantErr: true,
		},
		{
			name:    "no flags specified",
			args:    []string{"123", "--repo", "owner/repo"},
			wantErr: true, // No changes specified
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

			var runCalled bool
			cmd := NewCmdEdit(f, func(opts *EditOptions) error {
				runCalled = true
				// Check if no changes specified
				if opts.Title == "" && opts.Body == "" && opts.BodyFile == "" &&
					opts.Base == "" && opts.Draft == "" &&
					len(opts.Labels) == 0 && opts.Milestone == 0 &&
					opts.CloseRelatedIssue == "" {
					return fmt.Errorf("no changes specified. Use flags to specify what to edit")
				}
				return nil
			})
			cmd.SetArgs(tt.args)

			_, err := cmd.ExecuteC()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteC() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !runCalled {
				t.Error("run function was not called")
			}
		})
	}
}

func TestEditRun(t *testing.T) {
	tests := []struct {
		name    string
		opts    *EditOptions
		wantErr bool
	}{
		{
			name: "no repository",
			opts: &EditOptions{
				Repository: "",
				Number:     123,
				Title:      "New title",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, _, _ := iostreams.Test()
			tt.opts.IO = io
			tt.opts.HttpClient = func() (*http.Client, error) {
				return &http.Client{}, nil
			}

			err := editRun(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("editRun() error = %v, wantErr %v", err, tt.wantErr)
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
			name:    "empty repo requires explicit repo",
			repo:    "",
			wantErr: true,
		},
		{
			name:    "invalid format",
			repo:    "invalid",
			wantErr: true,
		},
		{
			name:    "too many parts",
			repo:    "owner/repo/extra",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := parseRepo(tt.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if owner != tt.wantOwner {
					t.Errorf("parseRepo() owner = %v, want %v", owner, tt.wantOwner)
				}
				if repo != tt.wantRepo {
					t.Errorf("parseRepo() repo = %v, want %v", repo, tt.wantRepo)
				}
			}
		})
	}
}
