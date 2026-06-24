# /loop: 定期 Issue Triage

## Prompt

```
/loop 30m 检查 gitcode-cli/cli 的 open issue，对缺少 type/* scope/* status/* 标签的 issue 补全分类标签
```

## 替换参数

间隔可按需调整（推荐 30m 或 1h）

## 门禁检查

- 每个 open issue 是否至少有 type/* + scope/* + status/* 标签
- 缺标签的在 issue comment 中说明补充了什么

## 预期输出

- 每个 issue 的分类标签完整
- Loop 自动收敛（无事可做时停止）
