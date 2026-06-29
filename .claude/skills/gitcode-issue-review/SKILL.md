---
name: gitcode-issue-review
description: Analyze a GitCode issue before implementation using GitCode CLI, repository inspection, comments, linked PRs, and acceptance criteria. Trigger when users ask to review, refine, validate, or prepare an issue for development.
---

# gitcode-issue-review

Review an issue before development.

## Command Entry

Use `gitcode` for cross-platform instructions. Linux/macOS may use `gc`; Windows PowerShell should use `gitcode`.

## Workflow

1. Read the issue, comments, linked PRs, labels, and milestone.
2. Verify the current repository behavior if code is available.
3. Identify missing requirements, risks, and acceptance criteria.
4. Produce an implementation-ready issue review.
5. Optionally post a comment with the review.

## Commands

```bash
gitcode issue view 123 -R owner/repo --comments --json
gitcode issue comments 123 -R owner/repo --json
gitcode issue prs 123 -R owner/repo --json
gitcode issue relations -R owner/repo --json
gitcode label list -R owner/repo --json
gitcode milestone list -R owner/repo --json
```

If code is needed:

```bash
ssh -T git@gitcode.com
gitcode repo clone owner/repo --git-protocol ssh
```

Post a review:

```bash
gitcode issue comment 123 -R owner/repo --body-file issue-review.md --json
```

## Output

```markdown
## Issue Review

### Current Understanding
- ...

### Missing Information
- ...

### Acceptance Criteria
- [ ] ...

### Implementation Notes
- ...

### Risks
- ...

### Suggested Next Step
- ...
```

## Rules

- Do not start implementation from an ambiguous issue without calling out ambiguity.
- Use issue comments and linked PRs as evidence, not just the issue title.
- Separate facts from assumptions.
- Do not overwrite labels or milestone unless asked.
