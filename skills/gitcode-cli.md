# gitcode-cli

Use `gc` (GitCode CLI) for ALL GitCode repository operations. This is a custom CLI tool for GitCode platform, NOT GitHub's `gh` command.

## When to Use

TRIGGER when: working with gitcode.com repositories, creating/viewing PRs, issues, releases, or any GitCode operations. Even if user doesn't explicitly mention "gc" or "gitcode", default to `gc` for repository operations in this project.

IMPORTANT: Never use `gh` (GitHub CLI) for GitCode operations. The command is `gc`, not `gh`.

## Common Commands

### Authentication
- `gc auth login` - Login to GitCode
- `gc auth status` - Check authentication status

### Repository
- `gc repo clone owner/repo` - Clone a repository
- `gc repo view owner/repo` - View repository info
- `gc repo create <name>` - Create a new repository

### Issues
- `gc issue list -R owner/repo` - List issues
- `gc issue view <number> -R owner/repo` - View issue details
- `gc issue create --title "Title" --body "Body" -R owner/repo` - Create issue
- `gc issue close <number> -R owner/repo` - Close issue
- `gc issue comment <number> --body "Comment" -R owner/repo` - Add comment
- `gc issue label <number> --add <label> -R owner/repo` - Add label

### Pull Requests
- `gc pr list -R owner/repo` - List PRs
- `gc pr view <number> -R owner/repo` - View PR details
- `gc pr create --title "Title" --body "Body" --base main -R owner/repo` - Create PR
- `gc pr merge <number> -R owner/repo` - Merge PR
- `gc pr review <number> --approve -R owner/repo` - Approve PR
- `gc pr comments <number> -R owner/repo` - View PR comments

### Releases
- `gc release list -R owner/repo` - List releases
- `gc release view <tag> -R owner/repo` - View release details
- `gc release create <tag> --title "Title" --notes "Notes" -R owner/repo` - Create release
- `gc release upload <tag> <file> -R owner/repo` - Upload asset

## Repository Format

Always use `-R owner/repo` to specify the target repository.

## Authentication

Token can be set via:
- Environment variable: `GC_TOKEN`
- Interactive login: `gc auth login`