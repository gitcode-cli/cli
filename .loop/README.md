# .loop — Loop 运行时管理

本目录管理 Claude Code `/loop` 和 `/goal` 命令的全生命周期：场景模板、活跃追踪、历史记录、交付追踪和长期记忆。

## 职责

- 提供可复用的 `/loop` 和 `/goal` prompt 模板
- 追踪当前活跃的 loop（跨会话记忆）
- 记录已完成 loop 的执行历史
- 追踪每个 issue 的交付状态
- 积累跨会话的长期记忆和经验教训

## 不是职责

- 定义项目规则（规则在 `spec/`）
- 定义 loop 工程策略（未来在 `.loop/project.yaml`、`.loop/policy.yaml`）

## 目录

| 目录 | 说明 | 提交？ |
|------|------|:--:|
| `prompts/` | 可复用 prompt 模板 | ✅ |
| `registry/` | 活跃 loop 注册表 | ❌ runtime |
| `history/` | 已完成 loop 记录 | ❌ runtime |
| `deliveries/` | 按 issue 的交付追踪 | ❌ runtime |
| `memory/` | 长期记忆和会话摘要 | ❌ runtime |

## AI 工作流

### 会话启动
1. 读 `registry/active.yaml` — 检查活跃 loop，标记已丢失的
2. 读 `memory/INDEX.md` — 获取上下文（活跃 Issue、待办、经验教训）

### 启动新 loop
1. 从 `prompts/README.md` 选模板
2. 自定义 `<placeholder>` 参数
3. 运行 `/loop` 或 `/goal`
4. 写 `registry/active.yaml`

### Loop 完成
1. 写 `history/YYYY-MM-DD-<type>.md`
2. 更新 `deliveries/issue-N.md`
3. 更新 `memory/INDEX.md`

### 会话关闭
1. 写 `memory/YYYY-MM-DD-session.md`
2. 标记活跃 loop 为 lost（CronCreate 随会话消失）

## 相关文档

- [docs/LOOP-GOAL-GUIDE.md](../docs/LOOP-GOAL-GUIDE.md) — /loop /goal 使用指南
- [spec/workflows/development-workflow.md](../spec/workflows/development-workflow.md) — 开发流程状态机
- [docs/superpowers/specs/2026-06-24-loop-fullflow-validation-report.md](../docs/superpowers/specs/2026-06-24-loop-fullflow-validation-report.md) — 验证报告
