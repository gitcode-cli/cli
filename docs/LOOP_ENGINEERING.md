# Loop Engineering

gitcode-cli 是 Loop Engineering Demo v1 的参考项目。

本项目用三仓协同展示工程闭环：

- `gitcode-cli/cli`：规则、命令实现、演示文档
- `gitcode-cli/skills`：AI 执行方法
- `gitcode-cli/loop-kits`：schema、policy、hooks、templates、adapters 标准包

## 目标

Loop Engineering 让人、AI、GitCode issue / PR、GitHub mirror CI、评审和知识归档形成可重复、可审计、可停止的流程。

## Demo v1 范围

Phase 1-3 不实现 `gc loop` 命令。演示依赖：

- `spec/loop/` 定义规则
- `.loop/` 定义项目配置
- `gitcode-cli/skills` 定义 AI 执行方式
- `gitcode-cli/loop-kits` 定义标准件
- GitCode issue / PR 保存长期状态和证据
- GitHub mirror Actions 保存 CI 事实

## 事实源

- 协作状态：GitCode issue / PR
- 主干完成：GitCode merged PR + `origin/main`
- CI 结果：GitHub mirror Actions
- CI 绑定：commit SHA
- AI 方法：`gitcode-cli/skills`
- 标准件：`gitcode-cli/loop-kits`

## 继续阅读

- [Loop Engineering Architecture](./LOOP_ENGINEERING_ARCHITECTURE.md)
- [Loop Engineering Demo](./LOOP_ENGINEERING_DEMO.md)
- [Mirror CI](./MIRROR_CI.md)
- [Hooks](./HOOKS.md)
- [spec/loop](../spec/loop/README.md)
