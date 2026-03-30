# 共享 skill 真相源

本目录存放 gitcode-cli 的共享 skill 真相源。

## 设计目标

每个共享 skill 负责定义：

- 适用场景
- 触发条件
- 与 `spec/`、`docs/` 的关系
- 对 Claude / Codex 适配层的统一约束

共享源不直接假设某个具体客户端的能力或语法。

## 当前映射

| 共享 skill | Claude 适配层 | Codex 适配层 |
|------------|---------------|--------------|
| `gitcode-cli` | `.claude/skills/gitcode-cli/` | `.codex/skills/gitcode-cli/` |
| `pr-reviewer` | `.claude/skills/pr-reviewer/` | `.codex/skills/pr-reviewer/` |
| `issue-reviewer` | `.claude/skills/issue-reviewer/` | `.codex/skills/issue-reviewer/` |
| `gc-dev-setup` | `.claude/skills/gc-dev-setup/` | `.codex/skills/gc-dev-setup/` |
| `gitcode-cmd-generator` | `.claude/skills/gitcode-cmd-generator/` | `.codex/skills/gitcode-cmd-generator/` |

## 使用原则

- 共享源定义场景，不定义与 `spec/` 冲突的项目规则
- 客户端适配层可以补充入口和工具差异
- 行为变化后，应先更新共享源，再同步适配层

## 同步工具

仓库提供：

```bash
./scripts/sync-ai-skills.sh
```

当前边界：

- 会基于共享源重写 `.codex/skills/*` 的基础适配文件
- 只会为缺失的 `.claude/skills/*` 生成占位文件
- 不会覆盖现有 Claude skill 正文
