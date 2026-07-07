## 方案设计: Issue #357

### 方案
唯一方案：替换字符串比较为 errors.Is。

| 位置 | 当前 | 修复后 |
|------|------|--------|
| retry.go:140 | `err.Error() == "context canceled" \|\| err.Error() == "context deadline exceeded"` | `errors.Is(err, context.Canceled) \|\| errors.Is(err, context.DeadlineExceeded)` |
| retry.go import | 无 errors/context | 加 "context" + "errors" |

### 风险
- 无：errors.Is 是 Go 标准错误判断方式，对 bare error 与字符串比较行为一致，对 wrapped error 更准确
- 不影响其他重试逻辑（shouldRetryOnStatus 不变）
