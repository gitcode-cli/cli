# Issue #352 交付记录

## 基本信息

| 字段 | 值 |
|------|-----|
| Issue | [#352](https://gitcode.com/gitcode-cli/cli/issues/352) |
| 标题 | test: ExitCode 映射测试覆盖不完整 — 缺少 409/403/400/网络错误等关键路径 |
| 类型 | type/test |
| 范围 | scope/errors |
| 风险 | risk/medium (classify-change-risk.py) |
| 原风险标签 | risk/low (已由脚本修正) |

## 状态流转

| 步骤 | 状态 | 时间 | 证据 |
|------|------|------|------|
| Triage | status/triage → status/verified | 2026-06-27 10:01 | [Comment](https://gitcode.com/gitcode-cli/cli/issues/352#note_177409566) + [Verify](https://gitcode.com/gitcode-cli/cli/issues/352#note_177409645) |
| Verified | status/verified | 2026-06-27 10:02 | 覆盖率 75% ExitCode, 0% WrapNotFound, 0% Unwrap |
| In Progress | status/in-progress | 2026-06-27 10:02 | Branch `fix/issue-352-exitcode-tests` (worktree) |
| Self Check | status/self-checked | 2026-06-27 10:08 | [PR #8 comment](https://github.com/gitcode-cli/cli/pull/8#issuecomment-1) |
| Multi-role Review | approved | 2026-06-27 10:10 | [Review Summary](https://github.com/gitcode-cli/cli/pull/8#issuecomment-2) |
| Merged | status/merged | 2026-06-27 10:11 | GitHub PR [#8](https://github.com/gitcode-cli/cli/pull/8) + GitCode PR [#295](https://gitcode.com/gitcode-cli/cli/pulls/295) |

## 8 项门禁证据

| # | 门禁 | 结果 | 证据 |
|---|------|:----:|------|
| 1 | 开发实现 | ✅ | +224/-1 lines: nil guard fix (errors.go) + 27 test cases (errors_test.go) |
| 2 | 测试 | ✅ | 27 新测试用例 (14 ExitCode + 5 WrapNotFound + 5 CLIError.Error + 3 CLIError.Unwrap) |
| 3 | 本地构建 | ✅ | `go build -o ./gc ./cmd/gc` 成功 |
| 4 | 单元测试 | ✅ | `go test ./...` — 1244 passed, 96 packages |
| 5 | Pre-commit | ✅ | 10/10 hooks passed |
| 6 | 实际命令验证 | ✅ | `./gc issue list -R infra-test/gctest1` exit 0 |
| 7 | 远端 CI | ⚠️ | [Run 28275334484](https://github.com/gitcode-cli/cli/actions/runs/28275334484) — ubuntu ✅, windows ✅, macos ❌ (pre-existing dyld env issue) |
| 8 | 风险分级 | ✅ | `python3 scripts/classify-change-risk.py --base origin/main` → risk=medium |

## CI 详细

- **Run URL**: https://github.com/gitcode-cli/cli/actions/runs/28275334484
- **Lint**: ✅
- **Test (ubuntu)**: ✅
- **Test (windows)**: ✅
- **Test (macos)**: ❌ — pre-existing dyld environment issue (`missing LC_UUID load command`), NOT caused by this change
- **Build**: ⏭️ (skipped due to macOS test dependency)
- **Docker**: ⏭️ (skipped due to macOS test dependency)

## 关联 PR

| 平台 | PR | 状态 |
|------|-----|:----:|
| GitHub Mirror | [#8](https://github.com/gitcode-cli/cli/pull/8) | MERGED |
| GitCode | [#295](https://gitcode.com/gitcode-cli/cli/pulls/295) | MERGED |

## 修改文件

| 文件 | 变更 | 说明 |
|------|------|------|
| `pkg/cmdutil/errors.go` | +1/-1 | nil guard: `&& cliErr != nil` 防止 typed nil dereference |
| `pkg/cmdutil/errors_test.go` | +223 | 27 新测试用例 |

## 评审结论

| 角色 | 结论 |
|------|:----:|
| 代码审查 | ✅ approved |
| 安全审查 | ✅ approved |
| 测试审查 | ✅ approved (3 non-blocking suggestions, 1 fixed) |
| 文档审查 | ✅ approved |

全部 4 角色独立评审通过，无需第二轮评审。

---

**完成时间**: 2026-06-27 10:11 UTC
**执行主体**: Claude (aflyingto)
