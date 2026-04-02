# 真相源矩阵

本文档定义 gitcode-cli 仓库中不同信息类型的真相源、适用边界和判定优先级。

## 职责

- 明确“哪类信息该去哪看”
- 防止把入口文档、阶段说明或外部使用说明误当成正式规则源
- 统一人工与 AI 在判定事实时的依据

## 适用场景

- 判断某份文档能否直接作为事实依据
- 识别仓库内开发与外部项目使用 `gc` 的边界
- 处理文档之间出现冲突或信息滞后时的优先级

## 真相源矩阵

| 信息类型 | 真相源 | 可否直接判事实 | 说明 |
|------|------|------|------|
| 项目正式规则 | `spec/` | 是 | 项目规则唯一正式来源 |
| 命令行为 | `docs/COMMANDS.md` | 是 | `gc` 命令行为唯一真相源 |
| 测试、门禁、评审规则 | `spec/foundations/*`、`spec/workflows/*` | 是 | 包括测试、状态机、门禁、评审边界 |
| 构建与打包规则 | `spec/delivery/*` | 是 | 当前本地构建与打包规则以此为准 |
| AI 共享场景定义 | `.ai/skills/*` | 有条件地是 | 只对共享 skill 场景定义生效，不得覆盖 `spec/` |
| 项目入口导航 | `README.md` | 否 | 入口文档，不是规则源 |
| AI 客户端入口 | `AGENTS.md`、`CLAUDE.md` | 否 | 入口文档，不是规则源 |
| 外部项目 AI 使用说明 | `docs/AI-GUIDE.md` | 否 | 只服务外部项目通过 AI 使用 `gc` |
| 外部通用 skill 分发 | `.ai/distribution/gc-core/*` | 否 | 只服务外部项目复用，不定义本仓库开发规则 |
| 项目阶段说明 | `issues-plan/PROGRESS.md` | 否 | 可能滞后，只能作阶段说明与背景参考 |
| 单个 issue / PR 的实时状态 | GitCode 远端 issue、PR、label、comment | 是 | 远端平台是实时事实源 |
| 是否已主干合入 | merged PR + `origin/main` | 是 | 不能只看 issue 状态、comment 或 release 文案 |

## 判定优先级

当不同文档或信息源出现冲突时，按以下顺序理解：

1. `spec/` 定义项目正式规则
2. `docs/COMMANDS.md` 定义命令行为
3. GitCode 远端平台、merged PR 和 `origin/main` 定义实时事实
4. `.ai/skills/*` 定义共享 AI 场景
5. `README.md`、`AGENTS.md`、`CLAUDE.md`、`issues-plan/PROGRESS.md` 仅作入口、导航或阶段说明

## 必须

- 不得把入口文档当成正式规则源
- 不得把 `issues-plan/PROGRESS.md` 当成实时状态真相源
- 不得把 `docs/AI-GUIDE.md` 当成 gitcode-cli 仓库内部开发流程规范
- 判断交付完成度时，必须检查远端平台事实和 `origin/main`

## 下一步去看哪里

- 如果你在改文档边界，继续看 [docs-governance.md](./docs-governance.md)
- 如果你在改 AI 入口或 skill 分层，继续看 [ai-collaboration.md](./ai-collaboration.md)
- 如果你在执行仓库内 AI 开发流程，继续看 [../workflows/ai-local-development-workflow.md](../workflows/ai-local-development-workflow.md)
