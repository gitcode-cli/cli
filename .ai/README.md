# .ai 目录说明

`.ai/` 是 gitcode-cli 仓库内保留的历史 AI 协作目录。

## 当前定位

当前正式 skill 真相源已经迁移到独立仓库：

- https://gitcode.com/gitcode-cli/skills

因此，本仓库内的 `.ai/skills/`、`.claude/skills/`、`.codex/skills/` 只作为历史兼容层和迁移参考，不再作为新增正式 skill 的落点。

## 权威边界

- 项目正式规则：`spec/`
- 命令行为：`docs/COMMANDS.md`
- AI skill 真相源：`gitcode-cli/skills`
- Loop Engineering 标准包：`gitcode-cli/loop-kits`
- 旧仓内 skill 目录：历史兼容/迁移参考

`.ai/` 不定义项目规则，也不覆盖 `spec/`。

不同信息类型的事实边界见 [spec/governance/source-of-truth-matrix.md](../spec/governance/source-of-truth-matrix.md)。

## 迁移规则

- 新增或更新通用 GitCode CLI skill 时，优先修改 `gitcode-cli/skills`。
- 新增 Loop Engineering schema、policy、hook、template、adapter 时，优先修改 `gitcode-cli/loop-kits`。
- 本仓库内旧 skill 文档只在需要保持历史兼容或迁移说明时调整。
- 不再把 `.ai/skills/` 当作跨 AI 的唯一 skill 来源。

## 历史内容

`.ai/distribution/` 和 `.ai/skills/` 中的内容可作为历史参考，但后续可复用资产应迁移到独立仓库维护。
