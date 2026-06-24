# .loop — Loop 运行时管理

本目录管理 Claude Code `/loop` 和 `/goal` 命令的全生命周期：场景模板、交付追踪、历史记录、长期记忆。

## 职责

- 提供可复用的 `/loop` 和 `/goal` prompt 模板
- 追踪每个 issue 的交付状态（团队共享）
- 积累跨会话的经验教训（团队共享）
- 记录个人 loop 运行历史和活跃状态

## 不是职责

- 定义项目规则（规则在 `spec/`）
- 定义 loop 工程策略（未来在 `.loop/project.yaml`、`.loop/policy.yaml`）

## 目录结构

```
.loop/
├── README.md                  [commit] 本文件
│
├── prompts/                   [commit] 团队共享 — prompt 模板
│   ├── README.md              模板索引表（11 个模板）
│   ├── triage.md              /goal: Issue triage → verified
│   ├── develop-and-test.md    /goal: 开发 + UT + 构建 + pre-commit
│   ├── self-check.md          /goal: 作者 9 项自检
│   ├── security-review.md     /goal: 安全扫描
│   ├── risk-classify.md       /goal: 风险分级
│   ├── docs-sync.md           /goal: 文档同步检查
│   ├── full-flow.md           /loop: 全流程 issue→merged（核心，Cron 30min）
│   ├── ci-monitor.md          /loop: CI 监控 + 修复
│   ├── batch-triage.md        /loop: 定期 issue triage
│   ├── review-response.md     /loop: 评审意见响应
│   └── pr-review-patrol.md    /loop: 定期评审他人 PR
│
├── deliverables/               [commit] 团队共享 — 按 issue 的交付记录
│   └── README.md              使用说明
│   └── issue-<N>.md           状态流转 + PR/CI 证据链接
│
├── memory/                     [混合]
│   ├── lessons.md             [commit] 团队共享 — 经验教训池
│   ├── INDEX.md               [ignore] 个人 — 跨会话上下文索引
│   └── YYYY-MM-DD-session.md  [ignore] 个人 — 会话摘要
│
├── history/                    [ignore] 个人 — 已完成 loop 记录
│   └── YYYY-MM-DD-<type>-<subject>.md
│
└── registry/                   [ignore] 个人 — 活跃 loop 注册表
    └── active.yaml            当前运行的 loop 状态
```

## 团队共享 vs 个人数据

| 类型 | 内容 | 提交？ | 说明 |
|------|------|:--:|------|
| 团队共享 | prompts/ | ✅ | 全员可用的 prompt 知识库 |
| 团队共享 | deliveries/ | ✅ | 按 issue 的交付链路，可审计 |
| 团队共享 | memory/lessons.md | ✅ | 跨会话积累的操作经验 |
| 个人 | registry/active.yaml | ❌ | 个人 loop 状态，不同步 |
| 个人 | history/*.md | ❌ | 个人 loop 运行记录 |
| 个人 | memory/INDEX.md | ❌ | 个人上下文，不同步 |
| 个人 | memory/*-session.md | ❌ | 个人会话摘要 |

## AI 工作流

### 会话启动
1. 读 `registry/active.yaml` — 检查活跃 loop，标记已丢失的
2. 读 `memory/INDEX.md` — 获取个人上下文
3. 读 `memory/lessons.md` — 加载团队经验教训

### 启动新 loop
1. 从 `prompts/README.md` 选模板
2. 自定义 `<placeholder>` 参数
3. 运行 `/loop` 或 `/goal`
4. 写 `registry/active.yaml`

### Loop 执行中
- 每完成一个 issue 的处理阶段，更新 `deliveries/issue-<N>.md`
- 发现新问题/教训，追加到 `memory/lessons.md`

### Loop 完成
1. 写 `history/YYYY-MM-DD-<type>.md`（loop 执行记录）
2. 更新 `memory/INDEX.md`（待办、上下文）

### 会话关闭
1. 写 `memory/YYYY-MM-DD-session.md`
2. 标记活跃 loop 为 lost

## 相关文档

- [docs/LOOP-GOAL-GUIDE.md](../docs/LOOP-GOAL-GUIDE.md) — /loop /goal 使用指南
- [spec/workflows/development-workflow.md](../spec/workflows/development-workflow.md) — 开发流程状态机
- [docs/superpowers/specs/2026-06-24-loop-fullflow-validation-report.md](../docs/superpowers/specs/2026-06-24-loop-fullflow-validation-report.md) — /loop 验证报告
- [docs/superpowers/specs/2026-06-24-loop-goal-validation-report.md](../docs/superpowers/specs/2026-06-24-loop-goal-validation-report.md) — /goal 验证报告

---

**最后更新**: 2026-06-24
