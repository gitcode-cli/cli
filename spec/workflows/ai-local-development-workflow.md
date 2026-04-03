# AI 本地开发流程

本文档编排 gitcode-cli 仓库内 AI 协作者执行本地开发任务时的标准闭环。

本文件不新增项目规则；它把现有 `spec/`、命令文档、模板和 AI 协作资产串成一条可执行流程。

## 职责

- 定义 AI 接到开发任务后的标准读取顺序
- 明确哪些事实看本地文档，哪些事实看远端平台
- 统一从 issue 验证到本地验证到 PR 自检的执行链路

## 适用场景

- 修复 bug
- 开发功能
- 文档或流程修改
- 仓库内部 AI 协作开发

## 不适用场景

- 外部项目通过 AI 使用 `gc` 操作 GitCode
- 外部项目安装或分发 `gc-core`
- 不涉及本仓库开发的纯使用咨询

## 先看什么

AI 接到任务后，按以下顺序建立上下文：

1. `AGENTS.md` 或 `CLAUDE.md`
2. [spec/README.md](../README.md)
3. 与任务直接相关的 `spec/foundations/*`、`spec/workflows/*`、`spec/delivery/*`
4. 如涉及命令行为，继续看 [../../docs/COMMANDS.md](../../docs/COMMANDS.md)
5. 如涉及共享 AI 场景，继续看 [../../.ai/skills/README.md](../../.ai/skills/README.md)

## 事实判定边界

- 项目规则：看 `spec/`
- 命令行为：看 `docs/COMMANDS.md`
- 单个 issue / PR 当前状态：看远端 GitCode 平台
- 是否已主干合入：看 merged PR 和 `origin/main`
- 阶段背景或收口说明：可参考 `issues-plan/PROGRESS.md`
- 外部项目 AI 使用说明：`docs/AI-GUIDE.md`，不适用于本仓库内部开发规则

## 标准执行流程

```text
接收任务
→ 读取正式规范与事实源
→ 核验 issue / PR 当前状态
→ 验证问题或确认需求
→ 创建非 main 分支
→ 开发实现
→ 本地测试与构建
→ 真实命令验证
→ 安全审查
→ 文档同步
→ 补 issue 进度与 PR 自检
→ ready-for-review
→ 独立评审
→ merge
```

## 详细步骤

### 1. 核验事实

开始修改前，必须先确认：

- 当前 issue / PR 的远端状态
- 是否已有 merged PR
- `origin/main` 是否已包含对应改动
- 如果 issue 已关闭，是否仍缺少 merged PR 或主干代码

不得把 `issues-plan/PROGRESS.md` 当成上述事实的唯一依据。

若 issue 已关闭，但不存在 merged PR 或 `origin/main` 不包含对应改动，必须明确判定为“未完成主干合入”。

### 2. 验证问题

未完成验证，不得开始写代码。

最小要求：

- 复现问题或确认需求有效
- 留下结构化验证记录
- 进入 `status/verified`

### 3. 创建开发分支

只有在完成验证后，才允许创建非 `main` 开发分支并开始修改实现文件。

### 4. 本地开发

开发过程中必须遵守：

- 命令行为以 `docs/COMMANDS.md` 为准
- 质量门禁以 `spec/foundations/code-quality-gates.md` 为准
- 共享 AI 约束不得覆盖 `spec/`

### 5. 本地验证

最小验证命令集：

```bash
go test ./...
go build -o ./gc ./cmd/gc
./gc version
./scripts/regression-core.sh
```

如果改动影响具体命令行为，还必须补至少一个真实命令验证。

真实命令测试只能使用 `infra-test/*`。

### 6. AI 友好 CLI 约束

AI 或脚本消费 `gc` 时，应优先使用：

- 读取类命令：`--json`
- 命令发现：`gc schema`
- 高风险写操作：先 `--dry-run`
- 删除或确认类命令：非交互场景中显式传 `--yes`

### 7. 安全审查

进入 ready-for-review 前，必须完成最小安全审查。

至少检查：

- 无硬编码 token、password、secret
- 文档、测试和示例中未误写真实凭证
- 涉及认证、配置、权限、网络调用、删除或覆盖行为时，已对照 `spec/foundations/security.md` 检查
- 如存在安全影响，已在作者自检和评审中明确记录

### 8. 文档同步

按改动类型同步相关文档：

- 命令行为变化：`docs/COMMANDS.md`，必要时 `README.md`
- 流程或门禁变化：`spec/*`、`AGENTS.md`、`CLAUDE.md`
- agent-friendly CLI 变化：`spec/foundations/agent-friendly-cli.md`、`docs/REGRESSION.md`
- 构建打包变化：`spec/delivery/*`、`docs/PACKAGING.md`
- AI 协作边界变化：`.ai/skills/*` 与适配层

### 9. 证据留存

提交评审前，至少应留存以下证据：

- issue 验证记录
- issue 开发进度记录
- 本地测试结果
- 构建结果
- 真实命令验证结果
- 安全审查结果
- 文档同步结果
- PR 作者自检

模板见：

- [../../docs/AI-TEMPLATES.md](../../docs/AI-TEMPLATES.md)
- `docs/ai-templates/*.md`

如使用本地文件准备 issue / PR 评论，建议在提交前运行最小机器校验：

```bash
python3 scripts/validate-ai-record.py --mode record --kind pr-self-check /path/to/pr-self-check.md
```

## 禁止事项

- 不得在 `main` 直接开发
- 不得跳过 `status/verified`
- 不得缺少证据就宣称完成
- 不得把作者自检当独立评审
- 不得把 `issues-plan/PROGRESS.md` 当成实时事实源
- 不得把 `docs/AI-GUIDE.md` 当成仓库内部开发规则

## 下一步去看哪里

- Issue 级动作：看 [issue-workflow.md](./issue-workflow.md)
- PR 级动作：看 [pr-workflow.md](./pr-workflow.md)
- 独立评审：看 [review-workflow.md](./review-workflow.md)
- 本地验证：看 [test-workflow.md](./test-workflow.md)
