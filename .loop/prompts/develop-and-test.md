# /goal: 开发 + UT + 构建

## Prompt

```
/goal until（在 git worktree 中执行）go test ./... passes AND go build -o ./gc ./cmd/gc succeeds
AND at least one real command test passes against infra-test/*
AND pre-commit hooks all pass
完成后更新 .loop/deliveries/issue-<N>.md 和 .loop/deliveries/README.md 汇总表。
构建和测试结果同步发布到 PR 评论区。
```

## 替换参数

无需替换。如改动限定在特定包，缩小 `./...` 范围。

## 评估器检查点

- 对话中是否出现 `PASS` 标记
- 是否出现 `go build: Success`
- 是否出现 pre-commit 全部通过

## 注意事项

- docs-only 改动可跳过此阶段（见 spec 5.3 门禁表）
- 代码改动必须在此阶段通过全部 4 项检查

## .loop/ 更新

```markdown
# deliveries/issue-<N>.md
| status/verified | status/in-progress | <ts> | branch created, initial commit |
```
