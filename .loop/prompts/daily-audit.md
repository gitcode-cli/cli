# /loop: Daily Audit — 每日主干全面审查

## Prompt

```
/loop 1d 在 git worktree 中审查 gitcode-cli/cli 的 main 最新代码。
启动 6 个独立 Agent 并行审查不同维度，每个 Agent 输出结构化发现：

Agent 1 - 代码质量: 逻辑 bug、错误处理缺失、边界条件、死代码、资源泄漏
Agent 2 - 安全审查: 凭证泄漏、注入风险、权限绕过、敏感路径无保护
Agent 3 - 测试审查: 覆盖缺口、flaky 测试、mock 滥用、缺少错误路径测试
Agent 4 - Agent友好: --json 输出不一致、缺 --dry-run/--yes、非交互阻塞
Agent 5 - 文档一致性: COMMANDS.md 漂移、spec 与实现不符、skill 过时
Agent 6 - 架构审查: Options 结构一致性、函数命名规范、API 契约、import 层级

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
