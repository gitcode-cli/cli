package git

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Repo represents a parsed repository reference
type Repo struct {
	Owner string
	Name  string
	Host  string
}

// String returns the full repository path (owner/name)
func (r *Repo) String() string {
	return r.Owner + "/" + r.Name
}

// URL returns the full URL for the repository
func (r *Repo) URL() string {
	host := r.Host
	if host == "" {
		host = "gitcode.com"
	}
	return fmt.Sprintf("https://%s/%s", host, r.String())
}

// GitURL returns the git URL for the repository
func (r *Repo) GitURL(protocol string) string {
	host := r.Host
	if host == "" {
		host = "gitcode.com"
	}
	if protocol == "ssh" {
		return fmt.Sprintf("git@%s:%s.git", host, r.String())
	}
	return fmt.Sprintf("https://%s/%s.git", host, r.String())
}

// ParseRepo parses a repository reference
// Supports formats:
// - owner/repo
// - https://gitcode.com/owner/repo
// - git@gitcode.com:owner/repo.git
func ParseRepo(ref string) (*Repo, error) {
	// Try URL format
	if strings.HasPrefix(ref, "http://") || strings.HasPrefix(ref, "https://") {
		return parseHTTPURL(ref)
	}

	// Try SSH format
	if strings.HasPrefix(ref, "git@") {
		return parseSSHURL(ref)
	}

	// Try owner/repo format
	parts := strings.Split(ref, "/")
	if len(parts) == 2 {
		return &Repo{
			Owner: parts[0],
			Name:  strings.TrimSuffix(parts[1], ".git"),
			Host:  "",
		}, nil
	}

	return nil, fmt.Errorf("invalid repository format: %s", ref)
}

func parseHTTPURL(ref string) (*Repo, error) {
	u, err := url.Parse(ref)
	if err != nil {
		return nil, err
	}

	path := strings.TrimPrefix(u.Path, "/")
	path = strings.TrimSuffix(path, ".git")
	parts := strings.Split(path, "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository path: %s", ref)
	}

	return &Repo{
		Owner: parts[0],
		Name:  parts[1],
		Host:  u.Host,
	}, nil
}

func parseSSHURL(ref string) (*Repo, error) {
	// Format: git@host:owner/repo.git
	re := regexp.MustCompile(`^git@([^:]+):([^/]+)/(.+?)(?:\.git)?$`)
	matches := re.FindStringSubmatch(ref)
	if len(matches) != 4 {
		return nil, fmt.Errorf("invalid SSH URL format: %s", ref)
	}

	return &Repo{
		Host:  matches[1],
		Owner: matches[2],
		Name:  matches[3],
	}, nil
}

// CurrentRepo returns the current repository from git remote
func CurrentRepo() (*Repo, error) {
	remote, err := DefaultRemote()
	if err != nil {
		return nil, err
	}

	remoteURL, err := RemoteURL(remote)
	if err != nil {
		return nil, err
	}

	return ParseRepo(remoteURL)
}

// Clone clones a repository to the specified directory
func Clone(repo *Repo, dir string, depth int) error {
	args := []string{"clone"}
	if depth > 0 {
		args = append(args, "--depth", fmt.Sprintf("%d", depth))
	}
	args = append(args, repo.GitURL("https"))
	if dir != "" {
		args = append(args, dir)
	}

	_, err := Run(args...)
	return err
}

// CloneWithProtocol clones a repository using the specified protocol
func CloneWithProtocol(repo *Repo, dir string, protocol string, depth int) error {
	args := []string{"clone"}
	if depth > 0 {
		args = append(args, "--depth", fmt.Sprintf("%d", depth))
	}
	args = append(args, repo.GitURL(protocol))
	if dir != "" {
		args = append(args, dir)
	}

	_, err := Run(args...)
	return err
}