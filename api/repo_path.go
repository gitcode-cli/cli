package api

import "net/url"

// escapedRepoPath returns the "/repos/{owner}/{repo}" path prefix with
// url.PathEscape applied to owner and repo to prevent path injection via
// crafted owner/repo values containing '/', '?', '#' etc. (#316).
func escapedRepoPath(owner, repo string) string {
	return "/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo)
}
