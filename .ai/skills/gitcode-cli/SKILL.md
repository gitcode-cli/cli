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
- 删除类命令在自动化场景中先用 `--dry-run`，真实执行时显式传 `--yes`

## AI 执行硬约束

- 未进入 `status/verified` 的 issue，不得开始写代码
- 未创建非 `main` 分支前，不得修改实现文件
- 未补验证记录、自检证据和必要测试结果前，不得宣称“已完成”
- 作者自检不得充当独立评审
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
7. 等待独立评审后再合并

## 关键边界

- 外部项目如何通过 AI 使用 `gc`：看 `docs/AI-GUIDE.md`
- gitcode-cli 仓库内部 AI 本地开发闭环：看 `spec/workflows/ai-local-development-workflow.md`
- 阶段背景说明可参考 `issues-plan/PROGRESS.md`，但实时事实仍看远端平台

## 适配层说明

- Claude 适配：`.claude/skills/gitcode-cli/`
- Codex 适配：`.codex/skills/gitcode-cli/`
