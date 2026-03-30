# Issue 拆分建议

## 1. 拆分原则

不建议只建一个“大而全”的总 issue。

更合理的做法是：

- 建一个总控 issue 作为 epic
- 再按“规范 / 基座 / 输出 / 错误 / 交互 / 发现层”拆成可执行 issue
- 每个 issue 都有清晰边界、影响面和验收标准

## 2. 建议里程碑

- `v0.6.0 Agent-Friendly CLI`
- 已创建：`#307769`

## 3. 建议 Issue 列表

### Issue 1

标题：

- `规划：建立 gitcode-cli 的 Agent-Friendly CLI 专项改进基线`
- 已创建：`#88`

类型：

- `enhancement`

目标：

- 作为本专项总控 issue
- 链接其余所有子 issue
- 明确参考基线、范围和非目标

建议内容：

- 背景：参考 `agent-cli-guide`
- 当前差距摘要
- 本专项范围
- 里程碑目标
- 子 issue 清单

验收：

- 总控 issue 可作为 roadmap 入口
- 后续 issue 全部挂接到该 issue

### Issue 2

标题：

- `规范：新增 agent-friendly CLI 设计规范并纳入质量门禁`
- 已创建：`#89`

类型：

- `documentation`
- `enhancement`

目标：

- 在 `spec/` 中新增正式规范
- 将输出契约、错误契约、非 TTY 规则、帮助文本要求纳入门禁

验收：

- 新规范文档落地
- `spec/README.md`、质量门禁、入口文档同步更新

### Issue 3

标题：

- `基础设施：抽象统一的输出模式、非交互策略与确认组件`
- 已创建：`#90`

类型：

- `enhancement`

目标：

- 在共享层建立统一运行时能力
- 结束各命令各自处理 JSON / prompt / TTY 的状态

验收：

- 新共享组件可复用
- 至少一组命令接入

### Issue 4

标题：

- `能力：为高频只读命令补齐统一 --json 输出契约`
- 已创建：`#91`

类型：

- `enhancement`

目标：

- 优先覆盖 `repo view/list`、`issue list/view`、`pr list/view`、`release list/view`

验收：

- 第一批高频读命令支持统一 `--json`
- 帮助和文档同步更新

### Issue 5

标题：

- `能力：建立 CLI 统一错误模型与语义化退出码`
- 已创建：`#92`

类型：

- `enhancement`

目标：

- 统一 CLI 错误类型、错误建议、retryable 语义
- 建立正式退出码说明

验收：

- 核心错误路径已有稳定退出码
- 文档可引用

### Issue 6

标题：

- `能力：统一删除类命令的非交互确认与安全执行语义`
- 已创建：`#93`

类型：

- `enhancement`

目标：

- 收口 `repo delete`、`label delete`、`milestone delete`、`release delete`
- 解决当前直接读 stdin 的分散实现

验收：

- 删除类命令统一确认流程
- 非 TTY 缺输入时明确失败

### Issue 7

标题：

- `能力：为第一批写命令引入 dry-run / 预检查执行模式`
- 已创建：`#94`

类型：

- `enhancement`

目标：

- 在 delete / create / review 等高风险路径先试点

验收：

- 至少一批写命令支持 dry-run 或预检查
- 行为与文档一致

### Issue 8

标题：

- `能力：新增 gc schema 命令提供命令元数据与参数结构发现能力`
- 已创建：`#95`

类型：

- `enhancement`

目标：

- 支持代理按需查询命令结构，而不是被迫全量读文档

验收：

- `gc schema` 至少支持命令树和单命令查询
- 输出为 JSON

### Issue 9

标题：

- `文档：重构 COMMANDS / AI-GUIDE / skills 以适配 agent-friendly CLI`
- 已创建：`#96`

类型：

- `documentation`

目标：

- 将新契约写入用户文档和 AI 协作入口

验收：

- `docs/COMMANDS.md`
- `docs/AI-GUIDE.md`
- `AGENTS.md`
- `CLAUDE.md`
- `.ai/skills/`
- `.codex/skills/`
- `.claude/skills/`
  完成必要同步

### Issue 10

标题：

- `测试：扩展回归矩阵以覆盖 JSON、非 TTY、错误码与 dry-run 契约`
- 已创建：`#97`

类型：

- `enhancement`
- `documentation`

目标：

- 把“代理友好性”从设计文档变成可回归验证的能力

验收：

- 回归脚本或测试矩阵扩展
- `docs/REGRESSION.md` 更新

## 4. 建议执行顺序

建议按以下顺序推进：

1. Issue 1 总控
2. Issue 2 规范
3. Issue 3 基础设施
4. Issue 5 错误与退出码
5. Issue 4 高频读命令 JSON
6. Issue 6 删除类命令确认
7. Issue 7 dry-run / 预检查
8. Issue 8 schema
9. Issue 9 文档与 skills
10. Issue 10 回归矩阵

## 5. 标签建议

建议至少使用以下标签组合：

- `enhancement`
- `documentation`
- `question` 仅用于需要设计讨论的 issue

如果仓库未来有更细标签，也可以补：

- `cli`
- `agent`
- `runtime`
- `testing`

## 6. 推荐做法

执行创建时，建议先：

1. 创建 milestone
2. 创建总控 issue
3. 再创建子 issue 并全部挂入 milestone
4. 最后把 issue 编号回填到总控 issue 和本目录文档
