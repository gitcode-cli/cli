# Delivery Record: Issue #272

- **Title**: repo stats decodes commit_statistics with wrong field names
- **Type**: bug
- **Status**: merged
- **Loop**: fullflow-main (4f73b1f3)
- **PR**: [#264](https://gitcode.com/gitcode-cli/cli/pulls/264)
- **Branch**: bugfix/issue-272
- **Started**: 2026-06-24 | **Merged**: 2026-06-24

## Gate Compliance

| # | Gate | Result | Evidence |
|---|------|--------|----------|
| 1 | 验证 | ✅ | `gc repo stats` 数据全为 0，JSON tag 与 API 字段名不匹配 |
| 2 | 开发 | ✅ | branch bugfix/issue-272, 2 files +7/-7 |
| 3 | 构建 | ✅ | `go build -o ./gc ./cmd/gc` Success |
| 4 | UT | ✅ | `go test ./...` 1199 passed in 95 packages |
| 5 | Pre-commit | ✅ | 全部通过 |
| 6 | 实际命令 | ✅ | `./gc repo stats -R gitcode-cli/cli --branch main` → aflyingto 79351/13644 ✅ |
| 7 | CI | ✅ | [run 28077718073](https://github.com/gitcode-cli/cli/actions/runs/28077718073) — All 8 ✅ |
| 8 | 风险分级 | ⚠️ | risk/high (误判，仅 4 JSON tag 修正) |
| + | 多角色评审 | ✅ | inline approved |
| + | 合并 | ✅ | 人工确认后合并 |

## Key Artifacts
- Fix: api/queries_repo.go — author→user_name, commits→commit_count, additions→add_lines, deletions→delete_lines
- Test: stats_test.go updated with correct field names

## Notes
- Classic Go JSON silent failure: tags mismatch → all zero values, no error
- Data went from all-zero to correct values: aflyingto 79351/13644
