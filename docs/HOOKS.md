# Loop Hooks

Loop Hooks 是 Loop Engineering 的阶段门禁，不是 Git hooks。

## 标准阶段

- `pre-loop`
- `pre-change`
- `post-change`
- `pre-pr`
- `post-pr`
- `pre-merge`
- `post-merge`
- `archive`

## Demo v1

Phase 1-3 中，hooks 由 `gitcode-cli/loop-kits` 定义契约，由 AI skills 按需执行或引用。

未来 `gc loop` 产品化后，可以由命令直接加载 `.loop/hooks.yaml` 并执行标准 hook。

## 输出

每个 hook 只产生以下结果之一：

- `pass`
- `fail`
- `blocked`
- `needs_human`

需要长期保存的 hook 结果必须写回 GitCode issue / PR。
