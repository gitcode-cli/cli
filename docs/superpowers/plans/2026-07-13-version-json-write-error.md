# Version JSON Write Error Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `gc version --json` return stdout write failures through Cobra without changing successful output.

**Architecture:** Keep the existing command constructor and JSON serialization helper. Change only the Cobra callback contract from `Run` to `RunE`, then return the existing `cmdutil.WriteJSON` error; verify behavior with a deterministic failing writer.

**Tech Stack:** Go 1.22, Cobra, standard-library `errors` and `testing` packages.

## Global Constraints

- Work only on `bugfix/issue-416` in the isolated worktree.
- Follow test-first RED-GREEN-REFACTOR; no production edit before the focused test fails for the expected reason.
- Preserve successful JSON fields and human-readable output exactly.
- Do not expand scope to human-readable `fmt.Fprint*` error propagation.
- Do not add or expose credentials; real remote operations source the local untracked credential environment.
- Deliver through the contributor fork `zhaogev5_87/cli`; upstream labels and merge remain maintainer-owned.

---

### Task 1: Propagate version JSON writer errors

**Files:**
- Modify: `pkg/cmd/version/version_test.go`
- Modify: `pkg/cmd/version/version.go`

**Interfaces:**
- Consumes: `NewCmdVersion(ver, commit, date string, commandName ...string) *cobra.Command`
- Produces: `cmd.Execute()` returns the error produced by `cmdutil.WriteJSON` when `--json` output fails.

- [ ] **Step 1: Write the failing test**

Add the `errors` import, a package-local failing writer, and this focused test:

```go
type errorWriter struct {
	err error
}

func (w errorWriter) Write([]byte) (int, error) {
	return 0, w.err
}

func TestVersionJSONReturnsWriteError(t *testing.T) {
	wantErr := errors.New("write failed")
	cmd := NewCmdVersion("v1.0.0", "abc123", "2026-07-13")
	cmd.SetArgs([]string{"--json"})
	cmd.SetOut(errorWriter{err: wantErr})

	err := cmd.Execute()
	if !errors.Is(err, wantErr) {
		t.Fatalf("cmd.Execute() error = %v, want %v", err, wantErr)
	}
}
```

- [ ] **Step 2: Run the focused test and verify RED**

Run:

```bash
PATH=/tmp/codex-go1.22.12/bin:$PATH go test ./pkg/cmd/version -run TestVersionJSONReturnsWriteError -count=1 -v
```

Expected: `FAIL`; `cmd.Execute() error = <nil>, want write failed`, proving the current JSON branch discards the writer error.

- [ ] **Step 3: Implement the minimal production change**

Change the Cobra callback in `NewCmdVersion` to:

```go
RunE: func(cmd *cobra.Command, args []string) error {
	out := cmd.OutOrStdout()

	if jsonOutput {
		info := VersionInfo{
			Version: ver,
			Commit:  commit,
			Built:   date,
			URL:     "https://gitcode.com/gitcode-cli/cli",
		}
		return cmdutil.WriteJSON(out, info)
	}

	fmt.Fprintf(out, "%s version %s\n", displayName, ver)
	fmt.Fprintf(out, "  commit: %s\n", commit)
	fmt.Fprintf(out, "  built:  %s\n", date)
	fmt.Fprintln(out, "https://gitcode.com/gitcode-cli/cli")
	return nil
},
```

- [ ] **Step 4: Format and verify GREEN**

Run:

```bash
PATH=/tmp/codex-go1.22.12/bin:$PATH gofmt -w pkg/cmd/version/version.go pkg/cmd/version/version_test.go
PATH=/tmp/codex-go1.22.12/bin:$PATH go test ./pkg/cmd/version -run TestVersionJSONReturnsWriteError -count=1 -v
PATH=/tmp/codex-go1.22.12/bin:$PATH go test ./pkg/cmd/version -count=1
```

Expected: focused test and package test both `PASS`.

- [ ] **Step 5: Run repository gates**

Run:

```bash
PATH=/tmp/codex-go1.22.12/bin:$PATH go test ./...
PATH=/tmp/codex-go1.22.12/bin:$PATH go build -o /tmp/gc-issue416-fixed ./cmd/gc
/tmp/gc-issue416-fixed version --json
pre-commit run --files pkg/cmd/version/version.go pkg/cmd/version/version_test.go docs/superpowers/specs/2026-07-13-version-json-write-error-design.md docs/superpowers/plans/2026-07-13-version-json-write-error.md
python3 scripts/classify-change-risk.py --base origin/main
git diff --check origin/main...HEAD
```

Expected: tests/build/pre-commit/diff checks pass; smoke command returns valid version JSON; risk classifier output is recorded verbatim for the PR.

- [ ] **Step 6: Commit the focused fix**

```bash
git add pkg/cmd/version/version.go pkg/cmd/version/version_test.go
git commit -m "fix(version): propagate JSON write errors"
```

Expected: one focused implementation commit after the design and plan commits.

---

### Task 2: Deliver the fork PR evidence

**Files:**
- Create locally only: `/tmp/issue-416-pr-body.md`
- No additional tracked source files.

**Interfaces:**
- Consumes: branch `bugfix/issue-416`, fork `zhaogev5_87/cli`, upstream `gitcode-cli/cli:main`.
- Produces: a cross-repository GitCode PR containing author self-check and maintainer-ready evidence.

- [ ] **Step 1: Push the branch to the contributor fork**

Add or update a `fork` remote pointing at `https://gitcode.com/zhaogev5_87/cli.git`. Source the untracked credential environment and use a temporary `GIT_ASKPASS` helper that reads `GC_TOKEN` from the environment without printing it. Then run:

```bash
git push -u fork bugfix/issue-416
```

Expected: the branch exists at `zhaogev5_87/cli:bugfix/issue-416`.

- [ ] **Step 2: Prepare the PR body**

Create `/tmp/issue-416-pr-body.md` with these sections in order:

```markdown
## 作者自检

- 作者主体标识: 夜雨听蝉（Codex 辅助）
- 根因或实现理由: version JSON 分支丢弃序列化写入错误，调用方无法识别 stdout 写入失败
- 主要修改: Cobra callback 改为可返回错误；JSON 分支返回 WriteJSON 错误；新增失败 writer 回归测试
- 单元测试: focused writer-error regression, version package, and full repository suite all PASS（附上 Step 4/5 的准确命令）
- 构建: `go build -o /tmp/gc-issue416-fixed ./cmd/gc` PASS
- 实际命令验证: `/tmp/gc-issue416-fixed version --json` returns valid version/commit/built/url JSON; writer failure is deterministically covered by unit test because a real terminal cannot reliably induce this path
- 安全审查: no credentials, API, auth, filesystem write, or remote mutation path changed
- 文档同步: no command schema, flag, successful output, or documented behavior changed; COMMANDS update not required
- 风险: paste the exact `classify-change-risk.py` result recorded in Task 1 Step 5, followed by the rationale that only version JSON error propagation changes
- 未覆盖项: human-readable fmt.Fprint error propagation remains out of scope
- 外部贡献说明: contributor account cannot manage upstream labels; maintainer should apply workflow labels after reviewing evidence

**Closes #416**
```

Replace the risk sentence with the exact classifier output before creating the PR. Do not place any other bare issue reference in the body.

- [ ] **Step 3: Create the cross-repository PR**

Run with the repository-built CLI and sourced credential environment:

```bash
/tmp/gc-issue416-fixed pr create \
  -R gitcode-cli/cli \
  --fork zhaogev5_87/cli \
  --head bugfix/issue-416 \
  --base main \
  --title "fix(version): propagate JSON write errors" \
  --body-file /tmp/issue-416-pr-body.md \
  --json
```

Expected: PR source is `zhaogev5_87:bugfix/issue-416`, target is `gitcode-cli/cli:main`, initially awaiting maintainer labels/review.

- [ ] **Step 4: Report maintainer handoff**

Record the PR URL, head SHA, exact local gate results, risk result, and any unavailable remote CI. State explicitly that upstream labels, independent maintainer review, CI mirror handling, approval, and merge require maintainer action.
