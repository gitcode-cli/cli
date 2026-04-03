# Integration Reviewer Agent 任务书

## 目标

汇总前 4 个专项 agent 的结论，做去重、排序、查重和 issue 提交决策。

## 本 agent 只做什么

- 读取前 4 个 agent 输出
- 对照远端已有 issue 做查重
- 合并重复 finding
- 判断哪些应提交、哪些应丢弃、哪些应挂到已有 issue
- 产出最终 issue 列表和正文草案

## 本 agent 不做什么

- 不自行发明新的技术结论
- 不扩大没有证据的风险描述
- 不把模糊建议包装成 issue

## 必看输入

- `issues-plan/team-agent-review/01-architecture-agent.md`
- `issues-plan/team-agent-review/02-cli-contract-agent.md`
- `issues-plan/team-agent-review/03-security-auth-agent.md`
- `issues-plan/team-agent-review/04-testing-regression-agent.md`
- 当前轮次各 agent 实际输出
- 远端相关 issue 列表
  - 至少覆盖当前 open issue
  - 如某条 finding 疑似历史回归或曾被关闭，必须补查相关 closed issue

## 必答问题

1. 哪些 finding 本质相同，应该合并成一个 issue？
2. 哪些 finding 已被远端 open issue 覆盖？
3. 哪些 finding 已在主干修复，不应继续创建 issue？
4. 哪些 finding 只是重构建议，证据不够，不应提交？
5. 最终哪些 issue 具备立即提交条件？

## 去重规则

- 同一根因、多个命令触发，优先合并成一个系统性 issue
- 同一用户影响、同一修复方向，优先合并
- 已有 open issue 能覆盖的，不重复创建
- 已关闭 issue 若主干仍未修复，可重开或新建，但必须写清原因
- 若创建新 issue，必须提醒执行主体在创建后立即补 `type/*`、`status/triage`、`scope/*` 标签

## 最终输出格式

```markdown
## Issue Candidate N

- 标题:
- 严重度:
- 是否提交: 是 / 否
- 是否与已有 issue 重复:
- 对应专题 agent:
- 核心证据:
- 建议标签:
- issue 正文草案:
```

## 远端提交要求

- 只提交 `是否提交: 是` 的候选项
- issue 标题必须和正文结论一致
- issue 正文必须能被未参与本轮审视的人独立理解
- 若结论是“上游 API 限制”，正文必须明确写清不是 CLI 自身已实现能力

## 合格标准

- 最终 issue 列表没有重复项
- 每条 issue 都能在代码、命令或文档里找到直接证据
- 最终 issue 列表里不包含“纯建议”

---

**最后更新**: 2026-04-03
