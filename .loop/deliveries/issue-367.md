# Delivery: Issue #367

- **Issue**: [#367](https://gitcode.com/gitcode-cli/cli/issues/367) — test: pkg/cmd/repo/clone cloneRun 无行为测试
- **PR**: [#313](https://gitcode.com/gitcode-cli/cli/merge_requests/313) — test(clone): add behavior tests for cloneRun
- **Type**: code (test-only)
- **Risk**: risk/medium
- **Branch**: worktree-issue-367-20260629
- **Date**: 2026-06-29

## Changes

Added `GitClone` function field to `CloneOptions` following project Pattern B (function field injection), enabling behavior testing of `cloneRun`.

### Source changes
| File | Change |
|------|--------|
| `pkg/cmd/repo/clone/clone.go` | Added `GitClone` function field to `CloneOptions`, extracted `execGitClone` |
| `pkg/cmd/repo/clone/clone_test.go` | Added 6 new behavior tests |

### New Tests
| Test | Coverage |
|------|----------|
| `TestCloneRunUsesCorrectURL` | URL construction (owner/repo, full URL, protocol variants - 3 subcases) |
| `TestCloneRunWithDepthFlag` | `--depth` flag propagation |
| `TestCloneRunWithBranchFlag` | `--branch` flag propagation |
| `TestCloneRunGitFailure` | Git clone error propagation |
| `TestCloneRunInvalidDepth` | Negative depth validation |
| `TestCloneRunInvalidBranch` | Invalid branch name validation |

## Gates

| # | Gate | Status |
|---|------|--------|
| 1 | Implementation | ✅ Added GitClone function field + 6 tests |
| 2 | Tests | ✅ `go test ./...` 78 packages pass (0 failures) |
| 3 | Build | ✅ `go build ./...` success |
| 4 | UT | ✅ All 78 packages pass |
| 5 | Pre-commit | ✅ 10/10 pass |
| 6 | Command verification | N/A (test-only, no command behavior change) |
| 7 | CI | ✅ [Run #28352147011](https://github.com/gitcode-cli/cli/actions/runs/28352147011) — test(3 OS)✅ lint✅ docker✅ build(ubuntu)✅; build(macOS) failure is pre-existing dyld LC_UUID |
| 8 | Risk | risk/medium |

## CI Summary

- **URL**: https://github.com/gitcode-cli/cli/actions/runs/28352147011
- **Test (ubuntu-latest)**: ✅
- **Test (macos-latest)**: ✅
- **Test (windows-latest)**: ✅
- **Lint**: ✅
- **Docker**: ✅
- **Build (ubuntu-latest)**: ✅
- **Build (macos-latest)**: ❌ (pre-existing dyld LC_UUID issue, not related)
- **Build (windows-latest)**: Cancelled (not related)

## Verification

- `go test ./pkg/cmd/repo/clone/ -v`: All 6 new + 6 existing tests pass
- `gofmt -l`: No formatting issues
- Pre-commit: 10/10 pass
- `classify-change-risk.py`: risk=medium
