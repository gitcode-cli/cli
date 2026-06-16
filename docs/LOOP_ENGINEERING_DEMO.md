# Loop Engineering Demo

本文记录 gitcode-cli Loop Engineering Demo v1 的一次真实分支验证。

本次演示不合并 `main`，也不依赖 `gc loop` 命令，而是在 `loop-engineering` 分支上验证 Phase 1-3 是否已经形成可执行、可审计、可恢复的工程闭环。

## 演示目标

验证以下能力是否已经可用：

- `cli/spec + .loop` 能定义规则、状态机和门禁
- `skills workflow` 能定义 AI 如何执行 loop
- `loop-kits` 能提供 schema、policy、hooks、templates、adapters 标准件
- GitCode issue / PR 能保存长期协作事实和证据
- GitHub mirror CI 能以 commit SHA 作为后续 CI 证据绑定边界

## 真实输入

| 类型 | 事实 |
| --- | --- |
| 总 Issue | `gitcode-cli/cli#299` |
| 规则层 PR | `gitcode-cli/cli!242` |
| 执行层 PR | `gitcode-cli/skills!11` |
| 标准层 PR | `gitcode-cli/loop-kits!1` |
| 验证分支 | `loop-engineering` |

## 三仓分工

| 仓库 | 职责 | 当前验证点 |
| --- | --- | --- |
| `gitcode-cli/cli` | 项目规则、`.loop` 配置、用户文档和演示入口 | `spec/loop`、`.loop/project.yaml`、`.loop/policy.yaml` 可读 |
| `gitcode-cli/skills` | AI 执行方法 | `gitcode-loop-engineering`、`gitcode-loop-ci`、`gitcode-loop-archive` 存在 |
| `gitcode-cli/loop-kits` | 标准资源包 | schemas、policies、templates、hooks、adapters 存在并可解析 |

## 执行过程

本次按 Loop Engineering v1 的人工编排路径执行：

```text
read issue / PR facts
-> verify loop-engineering branches
-> read spec/loop and .loop policy
-> check skills and loop-kits assets
-> generate local verification evidence
-> write verification result back to GitCode issue
```

对应状态推进：

```text
discovered
-> triaged
-> verified
-> planned
-> executing
-> self_checked
-> review_requested
```

本次没有推进到 `ci_waiting`、`ci_passed`、`approved`、`merged` 或 `archived`。

## 验证结果

已验证：

- GitCode Issue / PR 事实可通过 `gitcode` 读取
- 三仓 `loop-engineering` 分支存在，并且 PR 处于 open / mergeable 状态
- `spec/loop` 状态机覆盖 13 个标准状态
- `.loop/project.yaml` 指向 `gitcode-cli/skills` 和 `gitcode-cli/loop-kits`
- `.loop/policy.yaml` 定义了分支、证据、review、CI SHA、Issue/PR 和 `origin/main` 门禁
- 三个 loop skills 已存在
- loop-kits 的 JSON / YAML 文件可解析
- 样例 loop event、verification evidence、archive decision 满足当前契约的 required 字段和枚举约束
- 分支验证结果已写回 GitCode Issue #299 评论

未宣称：

- 未实现 `gc loop`
- 未合并 `main`
- 未抓取 GitHub mirror Actions run
- 未完成完整 JSON Schema validator 校验
- 未进入 `merged` 或 `archived` 状态

## 证据位置

| 证据 | 位置 |
| --- | --- |
| 长期协作事实 | GitCode Issue #299 |
| 分支验证评论 | Issue #299 comment `175852472` |
| 本地临时证据 | `.loop-output/loop-engineering-branch-verification.json` |
| 用户汇报材料 | `outputs/loop-engineering-demo-v1.pptx`，不提交 |

`.loop-output/` 是本地临时输出，已被 `.gitignore` 忽略，不作为长期事实源。

## 结论

Phase 1-3 已经可以支撑真实需求的 v1 手动编排闭环：

```text
需求
-> Issue
-> plan
-> branch
-> change
-> self-check
-> PR
-> evidence comment
```

当前能力适合演示和评审 Loop Engineering 的工程约束、事实源和证据流。

后续进入 Phase 4 时，再把高频检查产品化为 `gc loop doctor`、`gc loop scan`、`gc loop ci` 和 `gc loop evidence`，但仍不应过早实现无人工确认的 `gc loop run`。
