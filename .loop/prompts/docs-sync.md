# /goal: 文档同步检查

## Prompt

```
/goal until documentation sync is confirmed:
  - 命令行为变化 → docs/COMMANDS.md 已更新
  - 流程变化 → spec/* 已更新
  - AI 协作变化 → AGENTS.md / CLAUDE.md 已更新
  - 无变化 → 在自检中说明无需更新的依据
```

## 评估器检查点

- 对话中是否记录了文档同步决策（已更新 or 无需更新）
- 是否说明了依据
