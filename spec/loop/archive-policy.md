# Archive Policy

## 目标

归档用于判断 loop 产物应该留在哪里，避免把一次性证据提交为长期规则，也避免可复用经验丢失在聊天中。

## 归档目标

| 内容 | 位置 |
| --- | --- |
| 一次性执行证据 | GitCode issue / PR comment |
| 长期工程规则 | `spec/` |
| 用户可读说明 | `docs/` |
| AI 可复用执行方法 | `gitcode-cli/skills` |
| schema / policy / hook / template / adapter | `gitcode-cli/loop-kits` |
| CI 原始事实 | GitHub Actions run URL |
| 代码能力 | 项目源码 |

## 判断问题

- 这会影响未来工程行为吗？
- 用户是否需要阅读理解？
- AI 是否会复用执行方法？
- 机器是否需要校验格式？
- 它是否只属于本次执行？

## 禁止

- 把原始 CI 日志全文提交进主仓
- 把聊天总结当作事实源
- 把单次执行过程写成通用规则
- 在未 merged 前宣称完成归档
