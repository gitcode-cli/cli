package browse

import (
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/config"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

func TestNewCmdBrowse(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{name: "no args", args: []string{}, wantErr: false},
		{name: "number arg", args: []string{"42"}, wantErr: false},
		{name: "path arg", args: []string{"releases"}, wantErr: false},
		{name: "too many args", args: []string{"a", "b"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdBrowse(f, func(opts *BrowseOptions) error {
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

func TestNewCmdBrowseFlagsExist(t *testing.T) {
	cmd := NewCmdBrowse(cmdutil.TestFactory(), func(opts *BrowseOptions) error {
		return nil
	})
	for _, flag := range []string{"repo", "branch", "commit", "pr", "no-browser"} {
		f := cmd.Flags().Lookup(flag)
		if f == nil {
			t.Fatalf("flag %q missing", flag)
		}
	}
}

func TestBuildURLRepoHome(t *testing.T) {
	opts := &BrowseOptions{}
	got := buildURL("gitcode.com", "owner", "repo", opts)
	want := "https://gitcode.com/owner/repo"
	if got != want {
		t.Fatalf("buildURL() = %q, want %q", got, want)
	}
}

func TestBuildURLIssue(t *testing.T) {
	opts := &BrowseOptions{Arg: "42"}
	got := buildURL("gitcode.com", "owner", "repo", opts)
	want := "https://gitcode.com/owner/repo/issues/42"
	if got != want {
		t.Fatalf("buildURL() issue = %q, want %q", got, want)
	}
}

func TestBuildURLPullRequest(t *testing.T) {
	opts := &BrowseOptions{Arg: "42", PR: true}
	got := buildURL("gitcode.com", "owner", "repo", opts)
	want := "https://gitcode.com/owner/repo/pulls/42"
	if got != want {
		t.Fatalf("buildURL() pr = %q, want %q", got, want)
	}
}

func TestBuildURLPath(t *testing.T) {
	opts := &BrowseOptions{Arg: "releases"}
	got := buildURL("gitcode.com", "owner", "repo", opts)
	want := "https://gitcode.com/owner/repo/releases"
	if got != want {
		t.Fatalf("buildURL() path = %q, want %q", got, want)
	}
}

func TestBuildURLBranch(t *testing.T) {
	opts := &BrowseOptions{Branch: "feature/test"}
	got := buildURL("gitcode.com", "owner", "repo", opts)
	want := "https://gitcode.com/owner/repo/tree/feature/test"
	if got != want {
		t.Fatalf("buildURL() branch = %q, want %q", got, want)
	}
}

func TestBuildURLCommit(t *testing.T) {
	opts := &BrowseOptions{Commit: "abc123"}
	got := buildURL("gitcode.com", "owner", "repo", opts)
	// PathEscape escapes "abc123" as "abc123" (no special chars)
	want := "https://gitcode.com/owner/repo/commit/abc123"
	if got != want {
		t.Fatalf("buildURL() commit = %q, want %q", got, want)
	}
}

func TestBuildURLCommitPriority(t *testing.T) {
	opts := &BrowseOptions{Commit: "abc123", Branch: "main", Arg: "42"}
	got := buildURL("gitcode.com", "owner", "repo", opts)
	want := "https://gitcode.com/owner/repo/commit/abc123"
	if got != want {
		t.Fatalf("buildURL() commit should take priority = %q, want %q", got, want)
	}
}

func TestBuildURLBranchPriority(t *testing.T) {
	opts := &BrowseOptions{Branch: "main", Arg: "42"}
	got := buildURL("gitcode.com", "owner", "repo", opts)
	want := "https://gitcode.com/owner/repo/tree/main"
	if got != want {
		t.Fatalf("buildURL() branch should take priority over arg = %q, want %q", got, want)
	}
}

func TestBrowseRunNoBrowserPrintsURL(t *testing.T) {
	io, _, out, _ := iostreams.Test()
	opts := &BrowseOptions{
		IO: io,
		Config: func() (config.Config, error) {
			return config.New(), nil
		},
		Repository: "owner/repo",
		NoBrowser:  true,
		Browser: func(url string) error {
			t.Fatalf("Browser should not be called with --no-browser")
			return nil
		},
	}

	if err := browseRun(opts); err != nil {
		t.Fatalf("browseRun() error = %v", err)
	}
	got := out.String()
	if !strings.Contains(got, "https://gitcode.com/owner/repo") {
		t.Fatalf("output should contain URL; got: %s", got)
	}
}

func TestBrowseRunNonTTYPrintsURL(t *testing.T) {
	io, _, out, _ := iostreams.Test()
	opts := &BrowseOptions{
		IO: io,
		Config: func() (config.Config, error) {
			return config.New(), nil
		},
		Repository: "owner/repo",
		Browser: func(url string) error {
			t.Fatalf("Browser should not be called in non-TTY")
			return nil
		},
	}

	if err := browseRun(opts); err != nil {
		t.Fatalf("browseRun() error = %v", err)
	}
	if !strings.Contains(out.String(), "https://gitcode.com/owner/repo") {
		t.Fatalf("output should contain URL; got: %s", out.String())
	}
}

func TestBrowseRunOpensBrowser(t *testing.T) {
	io, _, _, _ := iostreams.TestTTY()
	var openedURL string
	opts := &BrowseOptions{
		IO: io,
		Config: func() (config.Config, error) {
			return config.New(), nil
		},
		Repository: "owner/repo",
		Browser: func(url string) error {
			openedURL = url
			return nil
		},
	}

	if err := browseRun(opts); err != nil {
		t.Fatalf("browseRun() error = %v", err)
	}
	if openedURL != "https://gitcode.com/owner/repo" {
		t.Fatalf("browser opened %q, want %q", openedURL, "https://gitcode.com/owner/repo")
	}
}

func TestBrowseRunIssueNumber(t *testing.T) {
	io, _, out, _ := iostreams.Test()
	opts := &BrowseOptions{
		IO: io,
		Config: func() (config.Config, error) {
			return config.New(), nil
		},
		Repository: "owner/repo",
		Arg:        "42",
		NoBrowser:  true,
	}

	if err := browseRun(opts); err != nil {
		t.Fatalf("browseRun() error = %v", err)
	}
	want := "https://gitcode.com/owner/repo/issues/42"
	if !strings.Contains(out.String(), want) {
		t.Fatalf("output should contain %q; got: %s", want, out.String())
	}
}

func TestBrowseRunPRNumber(t *testing.T) {
	io, _, out, _ := iostreams.Test()
	opts := &BrowseOptions{
		IO: io,
		Config: func() (config.Config, error) {
			return config.New(), nil
		},
		Repository: "owner/repo",
		Arg:        "42",
		PR:         true,
		NoBrowser:  true,
	}

	if err := browseRun(opts); err != nil {
		t.Fatalf("browseRun() error = %v", err)
	}
	want := "https://gitcode.com/owner/repo/pulls/42"
	if !strings.Contains(out.String(), want) {
		t.Fatalf("output should contain %q; got: %s", want, out.String())
	}
}

func TestBrowseRunMissingRepo(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := &BrowseOptions{
		IO:         io,
		Config:     func() (config.Config, error) { return config.New(), nil },
		Repository: "",
		BaseRepo:   func() (string, error) { return "", cmdutil.NewUsageError("no repo") },
		NoBrowser:  true,
	}

	err := browseRun(opts)
	if err == nil {
		t.Fatal("browseRun() error = nil, want error for missing repo")
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"42", true},
		{"12345", true},
		{"releases", false},
		{"", false},
		{"abc", false},
		{"42abc", false},
	}
	for _, tt := range tests {
		if got := isNumeric(tt.input); got != tt.want {
			t.Errorf("isNumeric(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestResolveHostDefault(t *testing.T) {
	host, err := resolveHost(func() (config.Config, error) {
		return config.New(), nil
	})
	if err != nil {
		t.Fatalf("resolveHost() error = %v", err)
	}
	if host != "gitcode.com" {
		t.Fatalf("resolveHost() = %q, want %q", host, "gitcode.com")
	}
}

func TestResolveHostNil(t *testing.T) {
	host, err := resolveHost(nil)
	if err != nil {
		t.Fatalf("resolveHost(nil) error = %v", err)
	}
	if host != "gitcode.com" {
		t.Fatalf("resolveHost(nil) = %q, want %q", host, "gitcode.com")
	}
}
