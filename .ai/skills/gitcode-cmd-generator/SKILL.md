---
name: gitcode-cmd-generator
description: Generate command code templates and test files for gitcode-cli project.

Use this skill when:
- Creating a new command for gitcode-cli (e.g., "create a new gc release command")
- Adding a subcommand to existing command (e.g., "add list subcommand to release")
- Generating command template (e.g., "generate command template for workflow list")
- Creating test files for commands (e.g., "create test for auth logout")
- Scaffolding CLI commands following project conventions

This skill generates production-ready Go code following gitcode-cli patterns with Cobra framework, Factory dependency injection, and IOStreams for output handling.
---

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
