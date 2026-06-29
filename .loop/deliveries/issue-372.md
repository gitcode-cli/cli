# Issue #372 — fix(docker): forward GC_TOKEN env var in make docker-run

- **Issue**: [#372](https://gitcode.com/gitcode-cli/cli/issues/372)
- **PR**: [#315](https://gitcode.com/gitcode-cli/cli/merge_requests/315)
- **类型**: bug (code change)
- **状态**: merged
- **完成时间**: 2026-06-29 14:42

## 变更摘要

Makefile `docker-run` target 缺少 `-e GC_TOKEN`，导致 README 文档中的 `GC_TOKEN=your_token make docker-run` 无法将 token 传入容器。添加 `-e GC_TOKEN` 转发环境变量。

## 变更

| 文件 | 行 | 变更 |
|------|-----|------|
| Makefile | 96 | `-e GC_TOKEN` 添加到 docker run 命令 |

## 门禁证据

| # | 门禁 | 状态 | 证据 |
|---|------|:--:|------|
| 1 | 实现 | ✅ | `-e GC_TOKEN` 添加到 Makefile docker-run |
| 2 | 测试 | ✅ | go test ./... 1268 全通过 |
| 3 | 构建 | ✅ | go build 成功 |
| 4 | UT | ✅ | 1268 passed in 96 packages |
| 5 | Pre-commit | ✅ | 全通过 (10 hooks) |
| 6 | 命令验证 | ✅ | `grep 'docker-run' -A1 Makefile` 确认 |
| 7 | CI | ✅ | [Run 28353375090](https://github.com/gitcode-cli/cli/actions/runs/28353375090): Ubuntu/Windows/Lint ✅; macOS ❌ (预存 dyld 问题) |
| 8 | 风险分级 | ✅ | classify-change-risk.py: **medium** (实际低风险) |

## 风险分析

classify-change-risk.py 输出 medium（Makefile 默认分级）。实际风险低：仅影响 `make docker-run` 入口，docker compose 路径已有正确的 GC_TOKEN 映射。单行变更，无 Go 代码影响。

## 流程状态

triage → verified → in-progress → draft → self-checked → ready → approved → **merged**
