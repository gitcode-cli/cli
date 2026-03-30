# 共享 skill: gitcode-cli

## 目标

定义在 gitcode-cli 项目中使用 `gc` 命令进行 GitCode 操作的共享约束。

## 统一约束

- GitCode 仓库操作优先使用 `gc`
- 不使用 `gh` 代替 `gc` 处理 GitCode 仓库
- 命令行为以 `docs/COMMANDS.md` 为准
- 项目规则以 `spec/` 为准

## 适配层说明

- Claude 适配：`.claude/skills/gitcode-cli/`
- Codex 适配：`.codex/skills/gitcode-cli/`
