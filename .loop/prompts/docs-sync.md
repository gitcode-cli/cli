# /goal: 文档同步检查

## Prompt

```
/goal until documentation sync is confirmed:
  - 命令行为变化 → docs/COMMANDS.md 已更新
  - 流程变化 → spec/* 已更新
  - 无变化 → 在自检中说明无需更新的依据
完成后更新 .loop/deliveries/issue-<N>.md 记录文档同步决策
```

## .loop/ 更新

```markdown
| docs_sync | <ts> | no change / updated COMMANDS.md / updated spec |
```
