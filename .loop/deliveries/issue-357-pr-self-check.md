## 作者自检

- 作者主体标识: AI 实现子代理 (glm-5.2 via opencode, 分支 bugfix/issue-357)
- 根因或实现理由: shouldRetryOnError 用 err.Error() 字符串比较判断 context 取消，对 wrapped context error（如 net.OpError 内嵌）漏判，导致 context 取消后仍发起重试。改用 errors.Is 遍历错误树。
- 主要修改: 2 文件 — retry.go import context+errors + shouldRetryOnError 改 errors.Is；retry_test.go 新增表格驱动 UT
  - api/retry.go: import 加 "context"+"errors"；shouldRetryOnError:140 字符串比较 → errors.Is(err, context.Canceled)/errors.Is(err, context.DeadlineExceeded)
  - api/retry_test.go: 新增 TestShouldRetryOnErrorWrappedContext（nil/bare/wrapped context/other error 表格驱动）
- 影响范围: api/retry.go shouldRetryOnError（HTTP 重试中间件的错误判断）；不改变用户可见命令行为（重试是内部 HTTP 行为）。
- 单元测试: ✅ TestShouldRetryOnErrorWrappedContext 全 PASS，表格驱动覆盖 nil/bare context/wrapped context/other error；-race 通过
- 构建: ✅ go build ./... 全包通过；go vet ./api/... 无问题
- 实际命令验证: ⏩ 豁免 — 内部重试逻辑，UT 充分覆盖；regression-core.sh 只读路径全过；wrapped context 在真实场景难以确定性触发
- 安全审查: ✅ 错误处理改进，不涉及 token/凭证/网络调用/删除；无硬编码 secret；改动方向是提升错误判断准确性
- 文档同步: ✅ 内部行为，不改命令行为，docs/COMMANDS.md 无需改；设计文档已补
- 风险: medium（classify-change-risk → risk=medium，runtime 路径）。按规范 medium 风险独立 AI 评审，仅在 blocker 或不确定时升级人工。
- 未覆盖项: CI 验证待 PR 推送触发
- 自检结论: 可进入 ready-for-review（medium 风险，需独立 AI 评审）

---

**关联 Issue**: #357
**Closes #357**
