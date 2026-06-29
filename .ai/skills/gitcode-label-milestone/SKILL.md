---
name: gitcode-label-milestone
description: Use GitCode CLI label and milestone commands to list, create, edit, delete, and apply repository planning metadata. Trigger for GitCode labels, milestones, issue labels, release planning, or queue organization.
---

# gitcode-label-milestone

Manage GitCode labels and milestones.

## Command Entry

Use `gitcode` for cross-platform instructions. Linux/macOS may use `gc`; Windows PowerShell should use `gitcode`.

## Labels

```bash
gitcode label list -R owner/repo
gitcode label list -R owner/repo --json

gitcode label create "bug" -R owner/repo --color "#ff0000" --description "Bug report"

# Destructive command: preview first
gitcode label delete bug -R owner/repo --dry-run
gitcode label delete bug -R owner/repo --yes
```

Apply labels to issues:

```bash
gitcode issue label 123 -R owner/repo --list
gitcode issue label 123 -R owner/repo --add bug,priority-high
gitcode issue label 123 -R owner/repo --remove stale
```

Apply labels to PRs:

```bash
gitcode pr edit 123 -R owner/repo --labels bug,priority-high --json
```

## Milestones

```bash
gitcode milestone list -R owner/repo --json
gitcode milestone create "v1.0" -R owner/repo --description "First release"
gitcode milestone view 1 -R owner/repo --json
gitcode milestone view 1 -R owner/repo --issues=false

gitcode milestone edit 1 -R owner/repo --title "v2.0" --json
gitcode milestone edit 1 -R owner/repo --description-file milestone.md --json
gitcode milestone edit 1 -R owner/repo --due-date "2026-06-30" --json
gitcode milestone edit 1 -R owner/repo --state closed --json

# Destructive command: preview first
gitcode milestone delete 1 -R owner/repo --dry-run
gitcode milestone delete 1 -R owner/repo --yes
```

Attach milestones:

```bash
gitcode issue edit 123 -R owner/repo --milestone 1 --json
gitcode pr edit 123 -R owner/repo --milestone 1 --json
```

## Rules

- List existing labels before creating new ones to avoid duplicates.
- Use repository naming conventions for type, priority, scope, and status labels.
- Use `--dry-run` before deleting labels or milestones.
- Avoid changing milestones for active work without checking open issues and PRs.
