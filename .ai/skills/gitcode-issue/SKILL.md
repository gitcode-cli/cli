---
name: gitcode-issue
description: Use GitCode CLI issue commands for creating, listing, viewing, editing, closing, reopening, commenting, labeling, and querying issue relations. Trigger for GitCode issue workflows.
---

# gitcode-issue

Manage GitCode issues with stable CLI patterns.

## Command Entry

Use `gitcode` for cross-platform instructions. Linux/macOS may use `gc`; Windows PowerShell should use `gitcode`.

## Read First

Before changing an issue, inspect the current state:

```bash
gitcode issue view 123 -R owner/repo --json
gitcode issue comments 123 -R owner/repo --json
```

Use `-R owner/repo` unless you are sure the current Git remote is the target repository.

## Create

```bash
gitcode issue create -R owner/repo --title "Bug: concise title" --body "Description" --json
gitcode issue create -R owner/repo --title "Feature request" --body-file issue.md --json
echo "Description from stdin" | gitcode issue create -R owner/repo --title "Task" --body-file - --json

# Preview without creating
gitcode issue create -R owner/repo --title "Task" --body "Description" --dry-run --json
```

Advanced fields:

```bash
gitcode issue create -R owner/repo --title "Security report" --body-file report.md --security-hole
gitcode issue create -R owner/repo --title "Feature" --body-file body.md --issue-type "需求" --issue-severity "高"
gitcode issue create -R owner/repo --title "Task" --body-file body.md --custom-fields-file custom-fields.json
```

## List and View

```bash
gitcode issue list -R owner/repo --state open --json
gitcode issue list -R owner/repo --state all --search "keyword" --json
gitcode issue list -R owner/repo --label bug,enhancement --limit 20 --format table
gitcode issue view 123 -R owner/repo --comments --json
gitcode issue prs 123 -R owner/repo --json
gitcode issue relations -R owner/repo --state open --json
```

## Edit, Close, Reopen

```bash
gitcode issue edit 123 -R owner/repo --title "New title" --json
gitcode issue edit 123 -R owner/repo --body-file new-body.md --json
gitcode issue edit 123 -R owner/repo --label bug,priority-high --json
gitcode issue edit 123 -R owner/repo --milestone 5 --json

gitcode issue close 123 -R owner/repo --yes --json
gitcode issue reopen 123 -R owner/repo --yes --json
```

## Comments and Labels

```bash
gitcode issue comment 123 -R owner/repo --body "Short comment" --json
gitcode issue comment 123 -R owner/repo --body-file comment.md --json
echo "Long comment" | gitcode issue comment 123 -R owner/repo --body-file - --json

gitcode issue comment edit 166061383 -R owner/repo --body-file updated.md
gitcode issue label 123 -R owner/repo --add bug,priority-high
gitcode issue label 123 -R owner/repo --remove stale
gitcode issue label 123 -R owner/repo --list
```

## Rules

- Search before creating a duplicate issue.
- Use `--body-file` for multi-line content.
- Use `--dry-run --json` when available before creating complex issues.
- Use `--yes` for close/reopen in non-interactive automation.
- Do not include secrets in issue body, comments, titles, or labels.
