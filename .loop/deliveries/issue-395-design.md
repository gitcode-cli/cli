## 方案设计: Issue #395

### 方案
与同文件 SafeFetch/SafeCheckout 一致：ValidateRef + `--` 双保险。

| 位置 | 当前 | 修复后 |
|------|------|--------|
| git.go:62 | `exec.Command("git", "remote", "get-url", name)` | `ValidateRef(name)` + `exec.Command("git", "remote", "get-url", "--", name)` |

### 防护层
1. **ValidateRef(name)** — 校验非空/不以 - 开头/无控制字符或 shell metacharacter（与 SafeFetch 一致）
2. **`--` 分隔符** — 即使 ValidateRef 漏过，`--` 后 name 被视为参数而非 option

### 风险
- 无：ValidateRef 对正常 remote 名（origin 等）通过，向后兼容
- ValidateRef 的 refPattern 允许 remote 名字符（字母数字/下划线/连字符/点）
