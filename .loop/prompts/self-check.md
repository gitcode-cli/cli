# /goal: 作者自检模板

## Prompt

```
/goal until PR self-check is complete for issue #<ISSUE_NUMBER>:
  - 根因或实现理由
  - 修改范围
  - 测试结果
  - CI 证据（run ID + Job 状态）
  - 安全审查结果
  - 文档同步说明
  - 风险分级
  - 未覆盖项
  AND gc issue comment <ISSUE_NUMBER> with self-check record -R gitcode-cli/cli
  AND gc issue label <ISSUE_NUMBER> --add status/self-checked -R gitcode-cli/cli
```

## 替换参数

- `<ISSUE_NUMBER>`: 目标 issue 编号

## 评估器检查点

- Issue comment 是否包含全部 9 项字段
- 每项字段是否有非空内容
- 标签是否包含 status/self-checked
