package fork

import (
	"os"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/git"
	"gitcode.com/gitcode-cli/cli/internal/config"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdFork(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "fork with repo",
			args:    []string{"owner/repo"},
			wantErr: false,
		},
		{
			name:    "fork with clone flag",
			args:    []string{"owner/repo", "--clone"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdFork(f, func(opts *ForkOptions) error {
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

func TestForkRunUsesRequestedRepository(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	var gotOwner string
	var gotRepo string

	opts := &ForkOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Config: func() (config.Config, error) {
			return fakeConfig{}, nil
		},
		ParseRepo: git.ParseRepo,
		ForkRepo: func(client *api.Client, owner, name string) (*api.Repository, error) {
			gotOwner = owner
			gotRepo = name
			if client.Token() != "test-token" {
				t.Fatalf("Token() = %q, want %q", client.Token(), "test-token")
			}
			return &api.Repository{
				FullName: "fork-owner/fork-repo",
				HTMLURL:  "https://gitcode.com/fork-owner/fork-repo",
			}, nil
		},
		CloneRepo: func(repo *git.Repo, dir string, protocol string, depth int) error {
			t.Fatalf("CloneRepo() should not be called when --clone is false")
			return nil
		},
		Repository: "infra-test/gctest1",
	}

	err := forkRun(opts)
	if err != nil {
		t.Fatalf("forkRun() error = %v", err)
	}

	if gotOwner != "infra-test" || gotRepo != "gctest1" {
		t.Fatalf("ForkRepo() called with %s/%s, want infra-test/gctest1", gotOwner, gotRepo)
	}
}

func TestForkRunCloneUsesForkedRepository(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")
	t.Setenv("GC_GIT_PROTOCOL", "ssh")

	f := cmdutil.TestFactory()
	var clonedRepo *git.Repo
	var clonedProtocol string

	opts := &ForkOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Config: func() (config.Config, error) {
			return fakeConfig{}, nil
		},
		ParseRepo: git.ParseRepo,
		ForkRepo: func(client *api.Client, owner, name string) (*api.Repository, error) {
			return &api.Repository{
				FullName: "fork-owner/forked-gctest1",
				HTMLURL:  "https://gitcode.com/fork-owner/forked-gctest1",
			}, nil
		},
		CloneRepo: func(repo *git.Repo, dir string, protocol string, depth int) error {
			clonedRepo = repo
			clonedProtocol = protocol
			return nil
		},
		Repository: "infra-test/gctest1",
		Clone:      true,
	}

	err := forkRun(opts)
	if err != nil {
		t.Fatalf("forkRun() error = %v", err)
	}

	if clonedRepo == nil {
		t.Fatalf("CloneRepo() was not called")
	}
	if clonedRepo.Owner != "fork-owner" || clonedRepo.Name != "forked-gctest1" {
		t.Fatalf("CloneRepo() repo = %s/%s, want fork-owner/forked-gctest1", clonedRepo.Owner, clonedRepo.Name)
	}
	if clonedProtocol != "ssh" {
		t.Fatalf("CloneRepo() protocol = %q, want ssh", clonedProtocol)
	}
}

func TestForkRunRejectsInvalidRepository(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	opts := &ForkOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Config: func() (config.Config, error) {
			return fakeConfig{}, nil
		},
		ParseRepo: git.ParseRepo,
		ForkRepo: func(client *api.Client, owner, name string) (*api.Repository, error) {
			t.Fatalf("ForkRepo() should not be called for invalid repository")
			return nil, nil
		},
		CloneRepo: func(repo *git.Repo, dir string, protocol string, depth int) error {
			t.Fatalf("CloneRepo() should not be called for invalid repository")
			return nil
		},
		Repository: "invalid",
	}

	err := forkRun(opts)
	if err == nil {
		t.Fatalf("forkRun() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "invalid repository") {
		t.Fatalf("forkRun() error = %q, want substring %q", err.Error(), "invalid repository")
	}
}

type fakeConfig struct{}

func (fakeConfig) Get(host, key string) (string, error) { return "", nil }

func (fakeConfig) Set(host, key, value string) error { return nil }

func (fakeConfig) GitProtocol(host string) config.ConfigEntry {
	if value := strings.TrimSpace(strings.ToLower(os.Getenv("GC_GIT_PROTOCOL"))); value != "" {
		return config.ConfigEntry{Value: value, Source: "environment"}
	}
	return config.ConfigEntry{Value: "https", Source: "default"}
}

func (fakeConfig) Editor(host string) config.ConfigEntry { return config.ConfigEntry{} }

func (fakeConfig) Browser(host string) config.ConfigEntry { return config.ConfigEntry{} }

func (fakeConfig) Pager(host string) config.ConfigEntry { return config.ConfigEntry{} }

func (fakeConfig) Authentication() config.AuthConfig { return nil }

func (fakeConfig) Write() error { return nil }
