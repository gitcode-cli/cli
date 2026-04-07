package edit

import (
	"fmt"
	"io"
	"net/http"
	"strings"
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

func TestEditRunUsesFormEncodedLabels(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	var gotPath string
	var gotContentType string
	var gotBody string

	ioStreams, _, _, _ := iostreams.Test()
	opts := &EditOptions{
		IO:         ioStreams,
		HttpClient: func() (*http.Client, error) { return &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			gotPath = req.URL.Path
			gotContentType = req.Header.Get("Content-Type")
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("failed to read request body: %v", err)
			}
			gotBody = string(body)
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     http.StatusText(http.StatusOK),
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"number":123,"title":"updated"}`)),
			}, nil
		})}, nil },
		Repository: "owner/repo",
		Number:     123,
		Labels:     []string{"type/feature", "risk/medium"},
	}

	if err := editRun(opts); err != nil {
		t.Fatalf("editRun() error = %v", err)
	}

	if gotPath != "/api/v5/repos/owner/repo/pulls/123" {
		t.Fatalf("request path = %q, want %q", gotPath, "/api/v5/repos/owner/repo/pulls/123")
	}
	if gotContentType != "application/x-www-form-urlencoded" {
		t.Fatalf("Content-Type = %q, want %q", gotContentType, "application/x-www-form-urlencoded")
	}
	for _, pair := range []string{
		"labels%5B%5D=type%2Ffeature",
		"labels%5B%5D=risk%2Fmedium",
	} {
		if !strings.Contains(gotBody, pair) {
			t.Fatalf("request body %q does not contain %q", gotBody, pair)
		}
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

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
