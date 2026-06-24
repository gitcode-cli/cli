# /loop: 定期 Issue Triage

## Prompt

```
/loop 30m 检查 gitcode-cli/cli 的 open issue，对缺少 type/* scope/* status/* 标签的 issue 补全分类标签。
每次批量操作后更新 .loop/memory/INDEX.md。如有 issue 交付，同步更新 .loop/deliveries/README.md 汇总表。
```

## 替换参数

间隔可按需调整（推荐 30m 或 1h）

## .loop/ 更新

```markdown
# memory/INDEX.md
## 当前活跃 Issue
- [#N] (新发现) - <title>
- [#M] - 标签已补全
```
