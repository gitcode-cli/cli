# CLI Contract Agent 任务书

## 目标

审视 CLI 输出、非交互能力和 agent-friendly 契约是否一致，只产出可以直接提交 issue 的命令契约问题。

## 本 agent 只看什么

- `--json` 是否覆盖读取类命令
- `gc schema` 是否能覆盖相关命令结构
- 写命令是否具备稳定正文输入能力，例如 `--body-file` / stdin
- 高风险命令是否具备 `--yes`
- 文本输出和 JSON 输出的语义是否一致
- docs / spec / 实现三者是否出现已落地与未落地的漂移

## 必看目录

- `pkg/cmd/`
- `pkg/cmdutil/`
- `docs/`
- `spec/foundations/`

## 必看文件

- `pkg/cmdutil/output.go`
- `pkg/cmd/schema/schema.go`
- `docs/COMMANDS.md`
- `docs/REGRESSION.md`
- `spec/foundations/agent-friendly-cli.md`
- `spec/foundations/code-quality-gates.md`

## 建议重点命令族

- `pkg/cmd/pr/`
- `pkg/cmd/issue/`
- `pkg/cmd/release/`
- `pkg/cmd/repo/`
- `pkg/cmd/label/`
- `pkg/cmd/milestone/`

## 必答问题

1. 读取类命令里，哪些仍缺 `--json`？
2. 哪些写命令仍只能靠 shell 拼接长正文，缺少 `--body-file` 或 stdin？
3. 哪些高风险写命令没有 `--yes` 或未使用统一确认语义？
4. 哪些命令文本输出会把“未知值”错误显示成确定值？
5. 文档里承诺了什么，但实现并未兑现？

## 明确算 issue 的问题类型

- 命令帮助、文档和真实行为不一致
- 读取类命令缺结构化输出
- 真实 API 返回未知值，CLI 却错误格式化成伪确定值
- 非交互脚本无法稳定驱动命令
- 高风险命令与统一确认语义不一致

## 不要提交 issue 的内容

- 纯文案偏好
- “JSON 字段也许可以再多一点” 但没有明确使用障碍
- 需要依赖远端平台新增能力才能判断的问题

## 输出格式

```markdown
## Finding N

- 标题建议:
- 严重度: high / medium / low
- 类型: cli
- 现象:
- 根因:
- 影响:
- 涉及目录或文件:
- 复现命令或核验方式:
- 是否疑似已有远端 issue:
```

## 通过标准

- 每条 finding 都能用至少一个实际命令或代码路径证明
- 若怀疑已有 issue，必须在输出中注明 issue 号或“需查重”
- 若问题是契约缺口，必须说明对 AI / script / 非交互调用的具体影响

---

**最后更新**: 2026-04-03
