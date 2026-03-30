# Agent CLI 改进方案

本目录用于沉淀 `gitcode-cli` 参考 `agent-cli-guide` 后形成的专项改进方案。

目标不是立即改代码，而是先把“为什么改、改什么、先后顺序、如何拆 issue”说清楚，并且保持与当前仓库的 `spec/`、`docs/`、`issues-plan/` 分层一致。

## 文档索引

| 文档 | 说明 |
|------|------|
| [gap-analysis.md](./gap-analysis.md) | 现状扫描、差距判断、问题分层 |
| [target-state.md](./target-state.md) | 目标能力模型、架构落点、兼容策略 |
| [roadmap.md](./roadmap.md) | 分阶段改进路线、优先级、验收标准 |
| [issue-breakdown.md](./issue-breakdown.md) | 建议创建的 milestone / issues 拆分方案 |

## 参考基线

本方案主要参考：

- `https://github.com/Johnixr/agent-cli-guide`
- `https://github.com/Johnixr/agent-cli-guide/blob/main/GUIDE.md`
- `https://github.com/Johnixr/agent-cli-guide/blob/main/CHECKLIST.md`

同时以本仓库当前正式规则为约束：

- 项目规则以 `spec/` 为准
- 命令行为真相源以 `docs/COMMANDS.md` 为准
- 当前项目状态以 `issues-plan/PROGRESS.md` 为准

## 结论摘要

`gitcode-cli` 当前已经具备较好的资源型命令树基础，例如 `repo`、`issue`、`pr`、`release`、`milestone`、`commit`。这意味着它天然适合继续向 agent-friendly CLI 演进。

但从“AI 代理可稳定消费”的角度看，当前能力还主要停留在“命令可用”，尚未形成统一的“契约层”：

- 结构化输出覆盖很低，只有少数命令显式支持 `--json`
- 非 TTY 默认行为尚未统一成规则
- 错误输出仍以字符串为主，没有稳定错误类型 / 退出码契约
- 破坏性操作确认存在，但没有统一的 `--dry-run` / 幂等策略
- 帮助文本质量参差不齐，尚未把 `--help` 视作 agent 的主发现入口
- 文档和 AI skill 还没有围绕“代理友好 CLI”建立专项规范

因此，本专项建议作为 `v0.6.x` 阶段的治理和产品能力收口主题，而不是零散地逐命令修补。
