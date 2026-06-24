---
name: gitcode-cli
description: |
  Use `gc` (GitCode CLI) for ALL GitCode repository operations. This is a custom CLI tool for GitCode platform, NOT GitHub's `gh` command.

  TRIGGER when: working with gitcode.com repositories, creating/viewing PRs, issues, releases, or any GitCode operations. Even if user doesn't explicitly mention "gc" or "gitcode", default to `gc` for repository operations in this project.

  IMPORTANT: Never use `gh` (GitHub CLI) for GitCode operations. The command is `gc`, not `gh`.
---

# 共享 skill: gitcode-cli

## 目标

定义在 gitcode-cli 项目中使用 `gc` 命令进行 GitCode 操作的共享约束。

## 统一约束

- GitCode 仓库操作优先使用 `gc`
- 不使用 `gh` 代替 `gc` 处理 GitCode 仓库
- 命令行为以 `docs/COMMANDS.md` 为准
- 项目规则以 `spec/` 为准
- 不同信息类型的事实边界以 `spec/governance/source-of-truth-matrix.md` 为准
- 流程推进以 `spec/workflows/*` 定义的状态机为准
- 质量门禁以 `spec/foundations/code-quality-gates.md` 为准
- 读取类命令优先使用 `--json`
- 需要探索命令结构时优先使用 `gc schema`
- typed command 尚未覆盖的平台能力可使用 `gc api <endpoint>`，但必须把远端原始响应当作事实，不自行补写字段
- 需要按文件或分支追踪提交历史时使用 `gc repo log --file ... --branch ... --json`
- 需要跨页扫描 PR 或从提交信息反查 PR 时使用 `gc pr list --paginate` / `gc pr list --commit-message`
- 删除类命令在自动化场景中先用 `--dry-run`，真实执行时显式传 `--yes`
- 提交代码前可用 `gc precommit check` 校验仓库 pre-commit 配置与本地环境（`--run` 实际拉起检查；非交互修改环境需 `--yes`，`--no-install` 仅诊断）
- Windows PowerShell 中优先使用 `gitcode`，避免 `gc` 被内置 `Get-Content` 别名覆盖；需要从 stdin 传中文/非 ASCII 正文时，优先使用 UTF-8 文件配合 `--body-file` / `--comment-file`，直接管道前先设置 `$OutputEncoding = [System.Text.UTF8Encoding]::new($false)`

## AI 执行硬约束

- 未进入 `status/verified` 的 issue，不得开始写代码
- 未创建非 `main` 分支前，不得修改实现文件
- 未补验证记录、自检证据和必要测试结果前，不得宣称“已完成”
- 作者自检不得充当独立执行主体评审
- PR 未合入主干前，不得把待修复 issue 视为已完成主干合入
- 不得把 `issues-plan/PROGRESS.md` 当成单个 issue / PR 实时状态真相源
- 不得把 `docs/AI-GUIDE.md` 当成 gitcode-cli 仓库内部开发流程规范

## 最小执行顺序

1. 读取 issue / PR 当前状态
2. 补齐类型、状态、范围标签
3. 补验证记录，再进入开发
4. 完成本地测试、构建和真实命令验证
5. 在 PR 中补作者自检
6. 进入 `status/ready-for-review`
7. 完成风险分级并等待独立执行主体评审后再合并

## 关键边界

- 外部项目如何通过 AI 使用 `gc`：看 `docs/AI-GUIDE.md`
- gitcode-cli 仓库内部 AI 本地开发闭环：看 `spec/workflows/ai-local-development-workflow.md`
- 阶段背景说明可参考 `issues-plan/PROGRESS.md`，但实时事实仍看远端平台
- `pr create --json` 若 warning 提示远端 body 未返回，不得把本地提交的正文当作远端事实；应使用 `gitcode pr view <number> -R owner/repo --json` 再核验
- `pr view --json` 应优先作为 PR 详情事实来源；新版本会包含 `body`、`description`、`merged_at` 并尽量补齐统计字段

## 适配层说明

- Claude 适配：`.claude/skills/gitcode-cli/`
- Codex 适配：`.codex/skills/gitcode-cli/`
