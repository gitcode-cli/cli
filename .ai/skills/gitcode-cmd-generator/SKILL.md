# 共享 skill: gitcode-cmd-generator

## 目标

定义为 gitcode-cli 生成或扩展命令实现时的共享场景。

## 统一约束

- 命令结构和编码规范以 `spec/` 为准
- 命令行为说明必须同步到 `docs/COMMANDS.md`
- 生成代码后仍需通过本地测试、真实命令验证和质量门禁

## 适配层说明

- Claude 适配：`.claude/skills/gitcode-cmd-generator/`
- Codex 适配：`.codex/skills/gitcode-cmd-generator/`
