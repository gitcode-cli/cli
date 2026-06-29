---
name: gitcode-pr
description: Use GitCode CLI pull request commands for creating, listing, viewing, editing, commenting, checking out, closing, reopening, marking ready or WIP, merging, testing, and syncing PRs. Trigger for GitCode PR workflows.
---

# gitcode-pr

Manage GitCode pull requests with current CLI behavior.

## Command Entry

Use `gitcode` for cross-platform instructions. Linux/macOS may use `gc`; Windows PowerShell should use `gitcode`.

## Create PR

Always verify branch and diff before creating:

```bash
git status --short --branch
git diff --stat origin/main...HEAD
```

Create:

```bash
gitcode pr create -R owner/repo --title "feat: add capability" --body-file pr.md --json
gitcode pr create -R owner/repo --head feature-branch --base main --title "fix: bug" --body-file pr.md --json
echo "PR body" | gitcode pr create -R owner/repo --title "docs: update guide" --body-file - --json
gitcode pr create -R owner/repo --fill --json
```

Fork PR:

```bash
gitcode pr create -R upstream/repo --fork myfork/repo --head feature-branch --title "feat: change" --body-file pr.md --json
```

Notes:

- `--body-file` is supported for `pr create`; `-` reads from stdin.
- `--body` and `--body-file` are mutually exclusive.
- `--json` cannot be combined with `--web`.

## Read PRs

```bash
gitcode pr list -R owner/repo --state open --json
gitcode pr list -R owner/repo --state merged --limit 50 --json
gitcode pr view 123 -R owner/repo --json
gitcode pr view 123 -R owner/repo --comments --json
gitcode pr diff 123 -R owner/repo
gitcode pr comments 123 -R owner/repo --json
```

## Edit and State

```bash
gitcode pr edit 123 -R owner/repo --title "New title" --json
gitcode pr edit 123 -R owner/repo --body-file pr.md --json
gitcode pr edit 123 -R owner/repo --labels bug,priority-high --json
gitcode pr ready 123 -R owner/repo --ready --yes --json
gitcode pr ready 123 -R owner/repo --wip --json
gitcode pr close 123 -R owner/repo --yes --json
gitcode pr reopen 123 -R owner/repo --yes --json
```

## Checkout, Test, Merge

```bash
gitcode pr checkout 123 -R owner/repo
gitcode pr test 123 -R owner/repo
gitcode pr test 123 -R owner/repo --force

gitcode pr merge 123 -R owner/repo --yes --json
gitcode pr merge 123 -R owner/repo --method squash --delete-branch --yes --json
```

## Sync PR to Another Repo

`pr sync` clones, fetches, and pushes repositories over SSH.

```bash
gitcode pr sync \
  --source-pr owner/source-repo#123 \
  --target-repo owner/target-repo \
  --base main \
  --title "[sync] Fix login bug" \
  --yes \
  --json
```

## Rules

- Use `--body-file` for multi-line PR descriptions.
- Use `--json` when the result will be parsed.
- Use `--yes` for merge, close, ready-state changes, and sync in automation.
- Confirm branch protection, CI, and review requirements before merge.
- Use SSH for code transfer workflows.
