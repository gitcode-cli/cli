---
name: gitcode-repo
description: Use GitCode CLI repository commands for view, list, clone, create, fork, delete, stats, and repo sync. Trigger when users need to inspect, create, fork, clone, delete, or synchronize GitCode repositories.
---

# gitcode-repo

Operate GitCode repositories with current GitCode CLI behavior.

## Command Entry

Use `gitcode` for cross-platform instructions. Linux/macOS may use `gc`. Windows PowerShell should use `gitcode`, `gc.exe`, or `python -m gc_cli`.

## Repository Formats

Accepted repository inputs:

```text
owner/repo
https://gitcode.com/owner/repo
git@gitcode.com:owner/repo.git
```

Code download and sync should use SSH by default. Verify SSH before clone or sync:

```bash
ssh -T git@gitcode.com
```

## Read Operations

```bash
gitcode repo view owner/repo
gitcode repo view owner/repo --json
gitcode repo view

gitcode repo list
gitcode repo list --owner owner --limit 30
gitcode repo list --visibility public --format table
gitcode repo list --json

gitcode repo stats --branch main -R owner/repo
gitcode repo stats --branch main -R owner/repo --json
gitcode repo stats --branch main --since 2026-01-01 --until 2026-05-26 -R owner/repo --json
```

## Clone and Fork

```bash
# SSH is preferred
gitcode repo clone owner/repo --git-protocol ssh
gitcode repo clone owner/repo target-dir --git-protocol ssh
gitcode repo clone owner/repo --branch main --depth 1 --git-protocol ssh

# Fork
gitcode repo fork owner/repo --json
gitcode repo fork owner/repo --clone
```

If a command needs to run on both Windows and Linux, avoid shell-specific syntax in the command itself and keep file paths simple.

## Create and Delete

```bash
gitcode repo create my-repo --private --json
gitcode repo create my-repo --public --description "My project" --json

# Destructive command: preview first
gitcode repo delete owner/repo --dry-run
gitcode repo delete owner/repo --yes --json
```

## Sync Directory to Another Repo

`repo sync` clones and pushes over SSH, copies a local source directory into a target repository path, commits, pushes a sync branch, and creates a PR.

```bash
gitcode repo sync \
  --target-repo owner/target-repo \
  --source-dir docs/api \
  --target-dir mirror/api \
  --base main \
  --title "sync: update api docs" \
  --yes \
  --json
```

## Rules

- Prefer `--json` for data consumed by scripts or AI agents.
- Use `--dry-run` before destructive commands.
- Use explicit `-R owner/repo` when the current directory is not definitely the target repository.
- Do not assume the current Git remote is GitCode; verify with `git remote -v` or `gitcode repo view`.
- For sync and clone workflows, SSH access is part of the precondition.
