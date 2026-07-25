package edit

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/testutil"
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

func TestNewCmdEditTracksLabelFlags(t *testing.T) {
	ioStreams, _, _, _ := iostreams.Test()
	f := &cmdutil.Factory{IOStreams: ioStreams}

	var got *EditOptions
	cmd := NewCmdEdit(f, func(opts *EditOptions) error {
		got = opts
		return nil
	})
	cmd.SetArgs([]string{
		"123", "-R", "owner/repo",
		"--labels", "legacy",
		"--add-label", "added",
		"--remove-label", "removed",
	})
	if _, err := cmd.ExecuteC(); err != nil {
		t.Fatalf("ExecuteC() error = %v", err)
	}
	if got == nil || !got.LabelsSet || !got.AddLabelsSet || !got.RemoveLabelsSet {
		t.Fatalf("label flag presence was not recorded: %#v", got)
	}
}

func TestEditRunUpdatesLabels(t *testing.T) {
	tests := []struct {
		name          string
		configure     func(*EditOptions)
		wantLabels    string
		wantRequests  int
		wantFirstVerb string
	}{
		{
			name: "legacy labels add and deduplicate",
			configure: func(opts *EditOptions) {
				opts.Labels = []string{"risk/high", "status/approved"}
			},
			wantLabels:    "type/bug,risk/high,status/approved",
			wantRequests:  2,
			wantFirstVerb: http.MethodGet,
		},
		{
			name: "add and remove preserve other labels",
			configure: func(opts *EditOptions) {
				opts.AddLabels = []string{"status/approved"}
				opts.RemoveLabels = []string{"risk/high"}
			},
			wantLabels:    "type/bug,status/approved",
			wantRequests:  2,
			wantFirstVerb: http.MethodGet,
		},
		{
			name: "replace labels skips current labels lookup",
			configure: func(opts *EditOptions) {
				opts.ReplaceLabels = []string{"status/approved", "status/approved"}
				opts.ReplaceLabelsSet = true
				opts.Yes = true
			},
			wantLabels:    "status/approved",
			wantRequests:  1,
			wantFirstVerb: http.MethodPatch,
		},
		{
			name: "replace with empty value clears labels",
			configure: func(opts *EditOptions) {
				opts.ReplaceLabelsSet = true
				opts.Yes = true
			},
			wantLabels:    ",",
			wantRequests:  1,
			wantFirstVerb: http.MethodPatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("GC_TOKEN", "test-token")
			var methods []string
			var gotLabels string

			ioStreams, _, _, _ := iostreams.Test()
			opts := &EditOptions{
				IO:         ioStreams,
				Repository: "owner/repo",
				Number:     123,
			}
			tt.configure(opts)
			opts.HttpClient = func() (*http.Client, error) {
				return &http.Client{Transport: testutil.NewRoundTripFunc(func(req *http.Request) (*http.Response, error) {
					methods = append(methods, req.Method)
					switch req.Method {
					case http.MethodGet:
						return editTestResponse(`{"number":123,"labels":[{"name":"type/bug"},{"name":"risk/high"}]}`), nil
					case http.MethodPatch:
						body, err := io.ReadAll(req.Body)
						if err != nil {
							t.Fatalf("read request body: %v", err)
						}
						values, err := url.ParseQuery(string(body))
						if err != nil {
							t.Fatalf("parse request body: %v", err)
						}
						gotLabels = values.Get("labels")
						if _, ok := values["labels"]; !ok {
							t.Fatal("PATCH request omitted labels field")
						}
						return editTestResponse(`{"number":123,"title":"updated"}`), nil
					default:
						t.Fatalf("unexpected method %s", req.Method)
						return nil, nil
					}
				})}, nil
			}

			if err := editRun(opts); err != nil {
				t.Fatalf("editRun() error = %v", err)
			}
			if len(methods) != tt.wantRequests {
				t.Fatalf("request count = %d, want %d (%v)", len(methods), tt.wantRequests, methods)
			}
			if methods[0] != tt.wantFirstVerb {
				t.Fatalf("first method = %q, want %q", methods[0], tt.wantFirstVerb)
			}
			if gotLabels != tt.wantLabels {
				t.Fatalf("labels = %q, want %q", gotLabels, tt.wantLabels)
			}
		})
	}
}

func TestEditRunRejectsUnsafeLabelChangesBeforeHTTP(t *testing.T) {
	tests := []struct {
		name      string
		configure func(*EditOptions)
		want      string
	}{
		{
			name: "replace requires confirmation",
			configure: func(opts *EditOptions) {
				opts.ReplaceLabels = []string{"status/approved"}
				opts.ReplaceLabelsSet = true
			},
			want: "confirmation required",
		},
		{
			name: "replace conflicts with add",
			configure: func(opts *EditOptions) {
				opts.ReplaceLabelsSet = true
				opts.AddLabels = []string{"status/approved"}
			},
			want: "cannot be combined",
		},
		{
			name: "same label cannot be added and removed",
			configure: func(opts *EditOptions) {
				opts.AddLabels = []string{"status/approved"}
				opts.RemoveLabels = []string{"status/approved"}
			},
			want: "both added and removed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ioStreams, _, _, _ := iostreams.Test()
			httpCalled := false
			opts := &EditOptions{
				IO:         ioStreams,
				Repository: "owner/repo",
				Number:     123,
				HttpClient: func() (*http.Client, error) {
					httpCalled = true
					return &http.Client{}, nil
				},
			}
			tt.configure(opts)

			err := editRun(opts)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("editRun() error = %v, want containing %q", err, tt.want)
			}
			if httpCalled {
				t.Fatal("HTTP client was created before label validation or confirmation")
			}
		})
	}
}

func editTestResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
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

func TestEditRunScansInlineBodyForSecrets(t *testing.T) {
	t.Setenv("GC_TOKEN", "secret-token-abc123")
	f := cmdutil.TestFactory()
	opts := &EditOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		BaseRepo:   func() (string, error) { return "owner/repo", nil },
		Number:     1,
		Body:       "leaked: secret-token-abc123",
	}
	err := editRun(opts)
	if err == nil || !strings.Contains(err.Error(), "secret") {
		t.Fatalf("editRun() error = %v, want secret detection error", err)
	}
}
