## 需求分析: Issue #357

### 问题定义
`api/retry.go:140` 的 `shouldRetryOnError` 用 `err.Error() == "context canceled"` 字符串比较判断 context 取消。当 context 错误被包装（如 `fmt.Errorf("...: %w", context.Canceled)`，常见于 net.OpError 内嵌），`Error()` 不再精确匹配字符串，导致 context 取消被误判为可重试。

### 影响范围
| 层级 | 受影响 |
|------|--------|
| 模块 | api/ |
| 文件 | retry.go (shouldRetryOnError) |
| 函数 | (t *retryTransport) shouldRetryOnError |
| 里程碑 | M3: Error Handling Consistency |

### 根因推断
字符串比较依赖 Go 标准库内部消息精确匹配，违背 errors.Is 的错误树遍历语义。属于错误处理一致性问题。

### 成功标准
- [x] shouldRetryOnError 改用 errors.Is(err, context.Canceled) / errors.Is(err, context.DeadlineExceeded)
- [x] 补 UT 覆盖 wrapped context error（复现测试）
- [x] go build/test/vet/fmt/race 全绿
- [x] regression-core 通过
- [ ] CI 通过
- [ ] PR 合入 main

### 约束条件
- 不得破坏: 现有重试行为（bare context + 网络错误仍正确）
- 兼容性: 完全向后兼容（errors.Is 对 bare error 行为与字符串比较一致）
