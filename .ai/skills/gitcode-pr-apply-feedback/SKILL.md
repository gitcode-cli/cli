---
name: gitcode-pr-apply-feedback
description: Pull all review feedback (inline comments, PR discussion, issue comments) from a GitCode PR, organize into an actionable fix list, apply fixes locally, verify, push, and reply. Use when the user receives PR review feedback, needs to apply code review suggestions, fix inline review findings, or says 处理评审意见/按评审修改代码/apply review.
---

# gitcode-pr-apply-feedback

Pull all review feedback from a GitCode PR, organize it into a fix list, apply each fix locally, verify, push, and reply to reviewers.

## Workflow overview

```
Pull all feedback → Triage & prioritize → Fix each item → Verify → Push → Reply
```

## Step 1: Pull all feedback

### Inline comments (includes file + position info)

```bash
gitcode pr comments <PR> -R owner/repo --json
```

Each comment with `path` and `position` fields is an inline finding. The `path` gives the file, `position` gives the diff line.

### PR discussion comments

```bash
gitcode pr view <PR> -R owner/repo --comments --json
```

General comments without file/line — these are broader feedback, design suggestions, or overall concerns.

### Associated issue comments (if the PR has a linked issue)

```bash
# Find linked issue from PR body
gitcode pr view <PR> -R owner/repo --json | jq '.body'

# Pull issue comments
gitcode issue comments <issue> -R owner/repo --json
```

## Step 2: Triage

Organize all findings into a fix list ordered by severity:

| Severity | Action |
|----------|--------|
| `需修改` / P0 / blocking | Fix now, cannot merge without |
| `中` / P1 | Fix now, strongly recommended |
| `建议` / P2 | Fix if reasonable, can defer |
| `低` / P3 | Note and move on |

Group by file so you can fix all issues in one file before moving to the next.

Create a working list — just a mental or written checklist:

```
[需修改] pkg/precommit/detect.go:62 → add filepath.Clean
[需修改] pkg/cmd/precommit/check/check.go:79 → use AddJSONFlag
[中]     pkg/precommit/check.go:34 → clarify OK semantics in comment
[建议]   pkg/precommit/detect.go:26 → switch to cmd.Output()
```

## Step 3: Apply fixes

### Checkout the PR branch

```bash
git fetch origin refs/merge-requests/<PR>/head:refs/remotes/pr/<PR>
git checkout -b fix/pr-<PR>-review refs/remotes/pr/<PR>
```

Or if the author's branch is already pushed:

```bash
git fetch origin <pr-branch>
git checkout <pr-branch>
```

### Map inline comments to local code

An inline comment's `path` and `position` point to the PR diff. To find the corresponding line in the local file:

- **New files**: position ≈ source line (within ±3). The line in the comment maps directly to the local file line.
- **Modified files**: Look at the diff hunk (`@@ -old,count +new,count @@`). The `+new_start` tells you which source line the hunk starts at. The comment's position minus the `@@` line count gives the offset within the hunk.
- When in doubt: read the file, find the code being referenced by content, not by line number.

### Fix each item

For each finding:
1. Read the relevant code
2. Understand the reviewer's concern
3. Apply the fix — don't blindly follow; verify the fix makes sense
4. If the fix would break something the reviewer didn't consider, note it

### Reply as you go

After fixing a batch of related items, reply to the reviewer:

```bash
gitcode pr comment <PR> --body "Fixed: added filepath.Clean to hookPath return. [done]" -R owner/repo
```

For inline comments, reply to the specific comment thread:

```bash
gitcode pr reply <PR> --comment-id <id> --body "Done: switched to AddJSONFlag. [fixed]" -R owner/repo
```

## Step 4: Verify

Before pushing, run the minimal verification:

```bash
go build ./...
go test ./...
go vet ./...
```

If the PR has associated CI (GitHub Actions), trigger it:

```bash
gh workflow run ci.yml --ref <branch>
```

## Step 5: Push

```bash
git add <changed-files>
git commit -m "fix: address review feedback for PR #<N>"
git push origin <branch> --force
```

If this is the author's branch you're pushing to, use `--force-with-lease` to avoid overwriting concurrent changes.

## Step 6: Summary reply

After all fixes are pushed, post a summary comment listing what was addressed, what was deferred, and why:

```bash
gitcode pr review <PR> -R owner/repo --comment-file summary.md
```

Summary format:

```markdown
## Review feedback applied

### Fixed
- [需修改] pkg/precommit/detect.go:62 — added filepath.Clean
- [需修改] pkg/cmd/precommit/check/check.go:79 — switched to AddJSONFlag
- [中] pkg/precommit/check.go:34 — clarified comment

### Deferred
- [建议] pkg/precommit/detect.go:26 — cmd.Output change needs broader discussion

### Verification
- go build: passed
- go test: passed
- CI: triggered (link)
```

## Rules

- Triage before coding. Don't fix P3 items before P0 items.
- Fix by understanding, not by rote. If a reviewer's suggestion would break something, explain why and propose an alternative.
- Reply to each finding. Never leave a reviewer wondering whether their feedback was seen.
- Push after each batch of related fixes, not after every single line change.
- If a finding is unclear, ask for clarification before coding — don't guess.
- Verify before pushing. A broken push is worse than no push.
- The PR author's branch may be read-protected on GitCode. In that case, create a new branch and comment the branch name so the author can pull from it.
