# Loop Engineering 规范入口

`spec/loop/` 定义 gitcode-cli 作为 Loop Engineering reference implementation 的项目级规则。

Loop Engineering 的目标是让人、AI、GitCode issue / PR、GitHub mirror CI、skills、loop-kits 和知识归档形成可重复、可审计、可停止的工程闭环。

## 阅读顺序

1. [loop-engineering.md](./loop-engineering.md)
2. [state-machine.md](./state-machine.md)
3. [evidence-model.md](./evidence-model.md)
4. [mirror-ci-contract.md](./mirror-ci-contract.md)
5. [archive-policy.md](./archive-policy.md)

按任务补充阅读：

- 设计 hooks：看 [hook-contract.md](./hook-contract.md)
- 设计 skills：看 [skill-contract.md](./skill-contract.md)
- 判断人工确认点：看 [human-approval-policy.md](./human-approval-policy.md)

## 分层

- `gitcode-cli/cli`：项目规则、命令实现和演示入口
- `gitcode-cli/skills`：AI 执行方法
- `gitcode-cli/loop-kits`：schema、policy、hooks、templates、adapters 标准包
- GitCode issue / PR：长期协作状态和证据
- GitHub mirror Actions：CI 执行事实

## 边界

Phase 1-3 不实现 `gc loop` 命令。演示由 `spec/loop`、`gitcode-cli/skills`、`gitcode-cli/loop-kits`、GitCode 和 GitHub mirror CI 共同完成。
