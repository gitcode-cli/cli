## 开发计划: Issue #357

| # | 文件 | 操作 | 说明 |
|---|------|------|------|
| 1 | api/retry.go | 修改 (+2/-2) | import 加 context+errors；shouldRetryOnError 改 errors.Is |
| 2 | api/retry_test.go | 修改 (+28/-0) | 新增 TestShouldRetryOnErrorWrappedContext 表格驱动 UT |

### 测试矩阵
| 类型 | 覆盖 | 状态 |
|------|------|------|
| UT | wrapped context.Canceled 不重试 | ✅ |
| UT | wrapped context.DeadlineExceeded 不重试 | ✅ |
| UT | bare context.Canceled/DeadlineExceeded 不重试 | ✅ |
| UT | nil / other error | ✅ |
| regression-core | 只读路径 | ✅ |
