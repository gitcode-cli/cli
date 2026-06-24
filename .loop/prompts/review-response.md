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
```

## 替换参数

- `<PR_NUMBER>`: 目标 PR 编号

## 预期输出

- 每个 review comment 已回复/已解决
- 修复后的代码通过 UT + 构建
- PR 标签更新
