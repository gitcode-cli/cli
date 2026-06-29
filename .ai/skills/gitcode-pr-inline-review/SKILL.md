---
name: gitcode-pr-inline-review
description: Submit inline code review comments on GitCode PRs using gc pr comment --path --position. Use when the user asks for inline review, line-level code review, posting review findings on specific code lines, or 行内检视/代码行级评审. DO NOT use for general PR comments without line references.
---

# gitcode-pr-inline-review

Submit inline code review comments directly on PR diff lines using GitCode CLI's `--path` and `--position` flags. Inline comments anchor to specific code lines in the PR diff view — reviewers see exactly which line each finding targets.

## Core command

```bash
gitcode pr comment <PR> --body "<finding>" --path <file> --position <N> -R owner/repo
```

| Flag | Purpose |
|------|---------|
| `--path` | File path exactly as in the PR diff (`pkg/precommit/detect.go`) |
| `--position` | 1-indexed line number counting from the first `@@` hunk header in the unified diff |

## What `--position` means

Position is the **line number in the unified diff output**, not the source code file. Count every line starting from the first `@@`:

```
@@ -10,6 +10,8 @@ func Foo() {    ← position 1 (hunk header)
     x := 1                                       ← position 2 (context)
     y := 2                                       ← position 3 (context)
+    z := 3                                       ← position 4 (addition)
+    w := 4                                       ← position 5 (addition)
     return x + y                                 ← position 6 (context)
@@ -20,4 +24,6 @@ func Bar() {   ← position 7 (next hunk)
```

Position ≠ source code line. Position counts hunk headers, context lines, and +/- lines alike.

## How to compute position

### For a new file

If the file is new (no deletion lines), position ≈ source line number + hunk header offset. The `@@` header itself takes 1 position. Each preceding `@@` line also counts.

**Simplest approach**: Run `git diff <base>..<head> -- <file>`, find the hunk containing the target line, and count lines from the first `@@` to the target.

### For a modified file

Find the hunk whose `+new_start` range contains the target source line. Count diff lines from the first `@@` header down to the matching line:

1. Run `git diff <base>..<head> -- <file>`
2. Locate the `@@` hunk containing the target line (the `+<new_start>` value tells you which source lines this hunk covers)
3. Count lines from position 1 (first `@@`) to the target line — every diff line adds 1 to the position counter
4. For context lines (no `+`/`-` prefix): they correspond to both old and new, advancing both counters
5. For `+` lines: they correspond to new source lines
6. For `-` lines: they exist only in the diff, do not advance the new source counter

### Quick rule for new files

If the file is entirely new in the PR, the diff consists only of `+` lines and `@@` headers. Position ≈ source line number + (number of `@@` headers before the target line). Usually within ±3 of the source line.

### Verification

After computing, sanity-check: submit the first comment and look at the PR diff view. If it lands on the wrong line, adjust by the offset.

## Workflow

### 1. Get the PR diff

GitCode PRs expose their head commit via merge request refs even when the author's branch is not pushed:

```bash
git fetch origin refs/merge-requests/<PR>/head:refs/remotes/pr/<PR>
git diff origin/main...refs/remotes/pr/<PR> -- <file>
```

### 2. Compute position from the diff

Read the diff output. For each finding with a source code line, find the corresponding diff line number by counting from the first `@@` header.

### 3. Submit the comment

```bash
gitcode pr comment 220 \
  --body "[需修改] hookPath 返回路径未调 filepath.Clean()" \
  --path pkg/precommit/detect.go \
  --position 68 \
  -R gitcode-cli/cli
```

### 4. Submit a summary review

After all inline comments, post a general review comment:

```bash
gitcode pr review <PR> -R owner/repo --comment-file review-summary.md
```

## Finding format

Prefix each body with a severity tag:

| Tag | Meaning |
|-----|---------|
| `[需修改]` | Must fix before merge |
| `[中]` | Should fix |
| `[建议]` | Nice to have |
| `[低]` | Minor, can defer |

## Rules

- Compute position by reading the diff and counting. You're an AI — you can read and count lines more adaptively than any script.
- Use the exact file path from the diff output.
- Each comment must state what's wrong AND suggest the fix.
- Don't post the same finding as both inline AND general comment.
- If unknown about a position, err on the side of a slightly higher number — it will still land close to the right code.
- For new files, the first `@@` line is position 1; the first source line (`+new_start`) is roughly position 2 or 3.
