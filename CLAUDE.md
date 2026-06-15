# CLAUDE.md

本文件是 Claude 在 gitcode-cli 仓库中的项目级入口文档。

如果任务涉及代码、文档、流程、评审或发布，请先阅读：

1. [spec/README.md](./spec/README.md)
2. [docs/README.md](./docs/README.md)
3. [README.md](./README.md)

## 1. 入口职责

`CLAUDE.md` 的职责是：

- 为 Claude 提供仓库级入口
- 指向正式规范、用户文档和后续 skill / loop 分层
- 约束代理不得绕过项目正式规则

`CLAUDE.md` 不是项目规则源。

项目正式规则以 [spec/README.md](./spec/README.md) 和 `spec/` 目录中的规范文档为准。

## 2. 必读文档

先读：

1. [spec/README.md](./spec/README.md)

再根据任务进入对应规范，不要机械顺序通读全部文档。

常用任务入口：

- 改命令行为：`spec/workflows/development-workflow.md`、`spec/governance/docs-governance.md`、`spec/foundations/code-quality-gates.md`
- 改 agent / script 可消费性：`spec/foundations/agent-friendly-cli.md`、`spec/foundations/code-quality-gates.md`
- 改 API / auth / config：`spec/foundations/coding-standards.md`、`spec/foundations/security.md`、`spec/foundations/testing-guide.md`
- 补测试或做真实命令验证：`spec/foundations/testing-guide.md`、`spec/workflows/test-workflow.md`
- 提交 PR / 做 review：`spec/workflows/pr-workflow.md`、`spec/workflows/review-workflow.md`
- 改构建 / 打包 / 发布：`spec/delivery/build-and-package.md`、`spec/delivery/release-process.md`
- 改 Loop Engineering 规则：`spec/loop/README.md`、`spec/governance/source-of-truth-matrix.md`

具体流程任务再进入：

- [Issue 流程](./spec/workflows/issue-workflow.md)
- [PR 流程](./spec/workflows/pr-workflow.md)
- [评审流程](./spec/workflows/review-workflow.md)
- [测试流程](./spec/workflows/test-workflow.md)

## 3. 核心执行规则

Claude 在本仓库中必须遵守：

- 项目命令固定为 `gc`；Windows PowerShell 中实际调用优先使用 `gitcode`
- 项目正式规范以 `spec/` 为准
- 命令行为以 [docs/COMMANDS.md](./docs/COMMANDS.md) 为准
- AI skill 真相源以独立仓库 `gitcode-cli/skills` 为准
- Loop Engineering 标准包以独立仓库 `gitcode-cli/loop-kits` 为准
- 本仓库内 `.ai/skills/`、`.claude/skills/`、`.codex/skills/` 只作为历史兼容/迁移参考
- 项目阶段说明可参考 [issues-plan/PROGRESS.md](./issues-plan/PROGRESS.md)，但该文档可能滞后，不作为单个 issue / PR 实时状态真相源
- 流程推进以 `spec/workflows/*` 和 `spec/loop/*` 定义的状态机为准，不能只把 checklist 当完成标准
- 状态标签必须随流程推进实时更新，参考 [status-label-checklist.md](./spec/workflows/status-label-checklist.md)
- 判断“某个 issue / 功能是否已合入主干”时，必须以 merged PR 和 `origin/main` 为准，不能只依据 issue 状态、issue comment、release 文案或功能分支存在与否
- GitHub mirror Actions 只作为 CI 执行事实源，不能替代 GitCode 主仓事实
- 如果 issue 已关闭但没有 merged PR 或 `origin/main` 不包含对应代码，必须明确判定为“未完成主干合入”
- 外部项目使用 AI 操作 GitCode 的说明以 `docs/AI-GUIDE.md` 为准，但该文档不定义本仓库内部开发流程
- 代码或流程变化后必须同步检查相关文档
- 实际命令测试只能使用 `infra-test/*`
- 不得在 `main` 直接开发
- 不得提交构建产物、评估输出、本地 loop runtime 或真实凭证
- 不得在缺少验证记录、自检证据或独立执行主体评审的情况下宣称“已完成”
- 遇到规范未覆盖的情况，必须先向用户确认，不得自由发挥
- 发现规范之间有冲突，必须先向用户报告，以用户确认为准
- MEMORY.md 是会话级记忆摘要，如与 `spec/` 冲突，以 `spec/` 为准

## 4. Claude 入口边界

当前仓库内的 Claude 项目级入口是：

- `CLAUDE.md`

当前仓库的 AI / Loop 分层为：

- `spec/`：项目正式规则
- `docs/`：用户文档
- `gitcode-cli/skills`：AI skill 真相源
- `gitcode-cli/loop-kits`：Loop Engineering 标准包
- `.ai/skills/`、`.claude/skills/`：历史兼容/迁移参考

Claude 应先以 `spec/` 和本文件为主要入口，再根据任务进入 `docs/`、独立 skills 仓或 loop-kits 仓。

## 5. 常用入口

- 用户文档入口：[docs/README.md](./docs/README.md)
- 命令手册：[docs/COMMANDS.md](./docs/COMMANDS.md)
- Loop Engineering：[docs/LOOP_ENGINEERING.md](./docs/LOOP_ENGINEERING.md)
- 认证说明：[docs/AUTH.md](./docs/AUTH.md)
- 回归说明：[docs/REGRESSION.md](./docs/REGRESSION.md)
- 打包说明：[docs/PACKAGING.md](./docs/PACKAGING.md)
- 发布说明：[RELEASE.md](./RELEASE.md)
- 贡献说明：[CONTRIBUTING.md](./CONTRIBUTING.md)
- Codex 入口：[AGENTS.md](./AGENTS.md)

## 6. 当前阶段说明

当前阶段目标是把 gitcode-cli 打造成 Loop Engineering Demo v1：

- Phase 1：`gitcode-cli/cli` 定义规则、事实源、状态机和演示入口
- Phase 2：`gitcode-cli/skills` 定义 AI 执行 skill
- Phase 3：`gitcode-cli/loop-kits` 定义 schema、policy、hooks、templates、adapters

旧的仓内 skill 分层只保留兼容和迁移说明，不再作为正式资产增长方向。
