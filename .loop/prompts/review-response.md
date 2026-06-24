# /loop: 评审意见响应

## Prompt

```
/loop 处理 PR #<PR_NUMBER> 的 review 意见：
  - 获取最新 review comments
  - 逐条分析并修复
  - 修复后运行 go test ./... && go build ./... 验证
  - commit + push
  - 回复 review 线程
  - 全部解决后更新 PR 标签，停止
  - 完成后更新 .loop/deliveries/issue-<N>.md 和 .loop/deliveries/README.md 汇总表
```

## .loop/ 更新

```markdown
| review | <ts> | 4/4 approved | PR comment <url> |
| review | <ts> | 2 issues found → fixed → re-reviewed |
```
