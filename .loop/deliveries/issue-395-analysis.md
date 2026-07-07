## 需求分析: Issue #395

### 问题定义
`git/git.go:62` 的 `RemoteURL` 直接传 name 给 `exec.Command("git", "remote", "get-url", name)`，无 `--` 分隔符也无 ValidateRef 校验。若 name 以 `-` 开头（如 `--upload-pack=/tmp/evil`，从恶意仓库的 .git/config 读取的远程名），git 会将其解释为选项。

### 影响范围
| 层级 | 受影响 |
|------|--------|
| 模块 | git/ |
| 文件 | git.go (RemoteURL) |
| 函数 | RemoteURL(name string) |
| 里程碑 | M1: Security Hardening |

### 根因推断
RemoteURL 实现时未与同文件 SafeFetch(248)/SafeCheckout(208) 保持一致（都用了 ValidateRef + `--`）。属于"同文件部分命令硬化、部分遗漏"。

### 成功标准
- [x] RemoteURL 加 ValidateRef(name) 校验
- [x] RemoteURL 加 `--` 分隔符
- [x] 补 UT 覆盖 option 注入/dash 前缀/空/shell metacharacter
- [x] go build/test/vet/fmt/race/lint 全绿
- [ ] CI 通过
- [ ] PR 合入 main

### 约束条件
- 不得破坏: 正常 remote 名（如 origin）的查询
- 兼容性: ValidateRef 对正常 remote 名通过，向后兼容
