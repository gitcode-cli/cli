# Team Agent 代码审视任务书

本目录定义一次只做“审视并提交 issue”的 Team agent 模式。

本轮不做：

- 不改代码
- 不直接修复问题
- 不在 review 过程中顺手提交重构

本轮只做：

- 分专题审视当前主干代码
- 把确认过的问题整理为结构化结论
- 查重远端已有 issue
- 将新的、明确的问题提交到远端

## 适用目标

适用于以下类型的专项审视：

- 架构和分层是否失效
- CLI 输出和 agent-friendly 契约是否不一致
- 认证、安全、危险写操作是否存在缺口
- 测试和真实回归是否存在高风险空洞

## 协作拓扑

固定使用 5 个 agent：

1. Architecture Agent
2. CLI Contract Agent
3. Security/Auth Agent
4. Testing/Regression Agent
5. Integration Reviewer Agent

前 4 个 agent 并行审视，最后由 Reviewer Agent 统一收口。

## 总体规则

- 仓库规则以 `spec/` 为准
- 事实边界以 `spec/governance/source-of-truth-matrix.md` 为准
- issue 流程以 `spec/workflows/issue-workflow.md` 为准
- 不得把 `issues-plan/PROGRESS.md` 当成远端 issue 实时状态源
- issue 只提交“确认过、可复述、可定位、可单独跟踪”的问题
- 如果远端已有等价 issue，应复用已有 issue，不重复创建
- 如果问题只是风格建议、泛泛优化、或缺少明确影响，不得创建 issue

## 统一审视范围

默认重点看：

- `pkg/cmd/`
- `pkg/cmdutil/`
- `api/`
- `internal/config/`
- `scripts/`
- `docs/COMMANDS.md`
- `spec/foundations/agent-friendly-cli.md`
- `spec/foundations/security.md`
- `spec/foundations/testing-guide.md`

## 统一 issue 准入标准

只有同时满足以下条件才可以提交 issue：

- 问题已经在当前主干代码中存在
- 有至少一个明确影响面
- 能指出根因或至少指出稳定复现路径
- 能定位到具体目录、文件或命令
- 不与远端 open issue 重复

## 统一输出格式

前 4 个专项 agent 的输出必须逐条使用以下格式：

```markdown
## Finding N

- 标题建议:
- 严重度: high / medium / low
- 类型: architecture / cli / security / testing
- 现象:
- 根因:
- 影响:
- 涉及目录或文件:
- 复现命令或核验方式:
- 是否疑似已有远端 issue:
```

Reviewer Agent 的收口输出必须逐条使用以下格式：

```markdown
## Issue Candidate N

- 标题:
- 严重度:
- 是否提交:
- 是否与已有 issue 重复:
- 对应专题 agent:
- 核心证据:
- 建议标签:
- issue 正文草案:
```

## 角色边界

Architecture Agent：

- 只审分层、初始化路径、共享抽象是否失效
- 不主导 CLI 文案和安全策略

CLI Contract Agent：

- 只审命令契约、输出格式、非交互能力、JSON/schema 一致性
- 不主导底层认证模型

Security/Auth Agent：

- 只审认证、凭证、确认语义、危险写路径
- 不把一般性重构问题包装成安全问题

Testing/Regression Agent：

- 只审测试覆盖、真实回归矩阵、空洞和脆弱路径
- 不替代其他 agent 做架构审查

Integration Reviewer Agent：

- 负责去重、排序、决定是否提交 issue
- 不新发明前 4 个 agent 未给出证据的结论

## 提交 issue 的统一要求

- 标题应直接描述问题，不写模糊优化口号
- 正文至少包含：背景、复现/核验、实际结果、预期结果、影响、建议方向
- 能归类成上游限制的，要明确写清是“API/平台限制”还是“CLI 自身缺陷”
- 若当前代码已经修复，应关闭或不创建，不得制造陈旧 issue

## 参考任务书

- [01-architecture-agent.md](./01-architecture-agent.md)
- [02-cli-contract-agent.md](./02-cli-contract-agent.md)
- [03-security-auth-agent.md](./03-security-auth-agent.md)
- [04-testing-regression-agent.md](./04-testing-regression-agent.md)
- [05-integration-reviewer-agent.md](./05-integration-reviewer-agent.md)

---

**最后更新**: 2026-04-03
