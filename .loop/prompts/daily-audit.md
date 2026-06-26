# /loop: Daily Audit — 每日主干全面审查

## Prompt

```
/loop 1d 在 git worktree 中审查 gitcode-cli/cli 的 main 最新代码。每次只执行一轮审查，完成后停止，等待下次触发。
启动 6 个独立 Agent 并行审查不同维度，每个 Agent 输出结构化发现：

Agent 1 - 代码质量: 依据 spec/foundations/coding-standards.md + code-quality-gates.md
Agent 2 - 安全审查: 依据 spec/foundations/security.md
Agent 3 - 测试审查: 依据 spec/foundations/testing-guide.md + spec/workflows/test-workflow.md
Agent 4 - Agent友好: 依据 spec/foundations/agent-friendly-cli.md
Agent 5 - 文档一致性: 依据 spec/governance/docs-governance.md
Agent 6 - 架构审查: 依据以下 spec 组合
  - spec/foundations/command-template.md（Options 模式、命令结构）
  - spec/foundations/coding-standards.md（命名规范、错误处理、包组织）
  - spec/governance/docs-governance.md（文档分层，检查 import 跨层引用）
  检查项: Options 结构一致性、函数/类型命名规范、API 契约、import 层级、跨包耦合度、接口隔离

过滤规则:
- 忽略纯 style nit（空格、命名偏好、注释语气）
- 忽略已有 issue 跟踪的已知问题

每个有效发现提交一个独立 Issue，包含:
- 标题: 类型前缀(bug/docs/security/perf/refactor)
- 复现步骤或证据定位
- 影响评估
- 建议修复方向
- 标签: type/* scope/* risk/* status/triage

提交后更新 .loop/deliveries/README.md 汇总表。
```

## .loop/ 维护

每发现一个 issue，创建 `.loop/deliveries/issue-<N>.md` 记录来源（daily-audit）。

## 与其他 loop 的分工

| Loop | 职责 | 阶段 |
|------|------|------|
| `daily-audit.md` | 发现 → 提交 Issue | 上游：生产需求 |
| `full-flow.md` | Triage → 解决 → Merge | 中游：消化需求 |
| `pr-review-patrol.md` | 审查他人 PR → 合并 | 下游：质量把关 |
