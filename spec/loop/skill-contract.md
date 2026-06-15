# Loop Skill 契约

## 定位

Loop skill 定义 AI 如何执行某类 loop，但不定义项目硬规则。

正式 skill 真相源：

- `gitcode-cli/skills`

## 必须

- 读取目标项目 `AGENTS.md` / `CLAUDE.md`
- 读取目标项目 `spec/` 和 `.loop/policy.yaml`
- 使用 GitCode issue / PR 作为长期状态源
- 使用 GitHub mirror Actions 作为 CI 执行源
- 使用 `gitcode` 命令示例，避免 Windows PowerShell 中裸 `gc` 冲突
- 明确人工确认点和停止条件

## 禁止

- 覆盖目标项目 `spec/`
- 把 `gitcode-cli/cli#299` 的具体方案写成所有项目通用硬规则
- 编造 CI 结果、review 结论或 merge 事实
- 把一次性运行记录沉淀为长期规则

## Demo v1 skills

- `gitcode-loop-engineering`
- `gitcode-loop-ci`
- `gitcode-loop-archive`
