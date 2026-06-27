#378 security: README Docker 示例使用内联 GC_TOKEN=value 模式 | 2026-06-27
PR: #306 | risk: low | docs-only
Gates: G5(✅) G8(✅) | CI: skipped (docs-only)
Commit: worktree-issue-378

## 处置说明

README.md:173 将内联 `GC_TOKEN=your_token make docker-run` 改为 `export GC_TOKEN=your_token && make docker-run`，避免 token 被记录到 shell history。

同时添加注释 `# 使用 export 而非内联赋值，避免 token 被记录到 shell history`，明确说明安全原因。

## Token 消耗

| 指标 | 值 |
|------|-----|
| 输入 tokens (cache miss) | — |
| 输出 tokens | — |
| 缓存命中 | — |
| 缓存写入 | — |
| 总计 tokens | — |
| 成本 | — |
| 耗时 | — |
