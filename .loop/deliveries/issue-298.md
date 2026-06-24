# Delivery Record: Issue #298

- **Title**: fetchURLHost IPv6 literal address parsing
- **Type**: bug
- **Status**: merged
- **Loop**: /goal relay (4 stages)
- **PR**: [#252](https://gitcode.com/gitcode-cli/cli/pulls/252)
- **Branch**: bugfix/issue-298
- **Started**: 2026-06-24 | **Merged**: 2026-06-24

## Gate Compliance

| # | Gate | Result | Evidence |
|---|------|--------|----------|
| 1 | 验证 | ✅ | [comment #176840841](https://gitcode.com/gitcode-cli/cli/issues/298#176840841): `[::1]:22` → `[` bug 复现 |
| 2 | 开发 | ✅ | branch bugfix/issue-298, pkg/cmd/pr/checkout/ (+88/-1) |
| 3 | 构建 | ✅ | `go build -o ./gc ./cmd/gc` Success |
| 4 | UT | ✅ | `go test ./...` 16 passed (checkout) / 1187 passed (full) |
| 5 | Pre-commit | ✅ | end-of-file-fixer 修复后通过 |
| 6 | 实际命令 | not executed | fetchURLHost 是内部函数，无独立 CLI 入口 |
| 7 | CI | ✅ | [run 28071396966](https://github.com/gitcode-cli/cli/actions/runs/28071396966) — All 8 ✅ (2nd run, after fixing #302) |
| 8 | 风险分级 | ✅ | risk/low (manually assessed) |
| + | 多角色评审 | ✅ | 4/4 approved (code+security+test+docs, security agent silent → manual) |
| + | 合并 | ✅ | PR #252 merged |

## Key Artifacts
- Fix: pkg/cmd/pr/checkout/checkout.go — fetchURLHost IPv6 bracket detection
- Fix: pkg/cmd/pr/checkout/checkout_test.go — TestFetchURLHost (9 cases, 4 IPv6)
- Unblock: pkg/cmd/pr/comment/resolve/resolve_test.go — t.Setenv("GC_TOKEN") for #302

## Notes
- /goal relay exercise: issues #298 + #302 together
- Stage 2 deadlock: global CI condition blocked by pre-existing #302
- Evaluator: correctly rejected "failure ≠ success", blind to root cause
- First CI run failed (#302 auth), second run passed after fix
