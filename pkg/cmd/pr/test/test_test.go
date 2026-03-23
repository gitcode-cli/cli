package test

import (
	"net/http"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

func TestNewCmdTest(t *testing.T) {
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
			name:    "valid number",
			args:    []string{"123", "--repo", "owner/repo"},
			wantErr: false,
		},
		{
			name:    "invalid number",
			args:    []string{"abc"},
			wantErr: true,
		},
		{
			name:    "with force flag",
			args:    []string{"123", "--repo", "owner/repo", "--force"},
			wantErr: false,
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
			cmd := NewCmdTest(f, func(opts *TestOptions) error {
				runCalled = true
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

func TestTestRun(t *testing.T) {
	tests := []struct {
		name    string
		opts    *TestOptions
		wantErr bool
	}{
		{
			name: "no repository",
			opts: &TestOptions{
				Repository: "",
				Number:     123,
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

			err := testRun(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("testRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseRepo(t *testing.T) {
	tests := []struct {
		name       string
		repo       string
		wantOwner  string
		wantRepo   string
		wantErr    bool
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
			name:    "invalid format",
			repo:    "invalid",
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