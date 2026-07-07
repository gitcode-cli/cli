## 需求分析: Issue #400

### 问题定义
`pkg/config/auth_config.go:277` 与 `pkg/config/config.go:269` 使用 `os.WriteFile(path, data, 0o600)` 写入认证/配置状态。`os.WriteFile` 有两个缺陷：
1. 权限受进程 umask 影响 — 宽松 umask（如 000）下，新建文件可能以 0666 权限落盘，导致 GitCode 认证 token 被同主机其他用户读取
2. 对已存在文件不收紧权限 — 若文件已存在且权限宽松（如 0644），`os.WriteFile` 不会 chmod，旧权限保留

此外存在同源威胁：若攻击者将 `~/.config/gc/auth.json` 替换为指向他处的符号链接，CLI 写入时会覆盖目标文件，造成凭证内容泄露或目标文件损坏。

### 影响范围
| 层级 | 受影响 |
|------|--------|
| 模块 | pkg/config |
| 文件 | auth_config.go (writeAuthState), config.go (writeConfigState) |
| 函数/类型 | secureWriteFile (新增), (c *config).writeAuthState, (c *config).writeConfigState |
| 安全域 | scope/credential-storage, M1: Security Hardening |

### 根因推断
`os.WriteFile` 的权限语义依赖 umask 且对既有文件无 chmod，与同包 `config.go:211-212` 对 config 目录显式 `os.Chmod(c.configDir, 0o700)` 的做法不一致。属于"目录层已硬化、文件层未硬化"的遗漏。

### 成功标准
- [x] 新增 `secureWriteFile(path, data, perm)` 函数，封装三步：Lstat 检测 symlink → WriteFile → Chmod
- [x] `writeAuthState` 替换为 `secureWriteFile` 调用
- [x] `writeConfigState` 替换为 `secureWriteFile` 调用
- [x] 拒绝写入符号链接（返回 "refusing to write to symlink" 错误）
- [x] 写入后显式 `os.Chmod(path, 0o600)`，覆盖既有文件宽松权限
- [x] `go build -o ./gc ./cmd/gc` 成功
- [x] `go test ./pkg/config/...` 全部通过（含 3 个新增 secureWriteFile UT）
- [x] `gofmt` / `go vet` / `-race` 全绿
- [ ] 人工真实命令验证（TTY 输入真实 token）：symlink 拒绝 + 权限硬化
- [ ] CI 通过（PR 提交后）
- [ ] PR 合入 main

### 约束条件
- 不得破坏: `gc auth login` 现有认证流程
- 兼容性: 完全向后兼容（权限更严格不影响合法使用）
- 安全: AI 代理不得读取/打印真实 token；真实命令验证由人工在 TTY 完成
- 平台: symlink 检测与 Unix 权限在 Windows 上跳过（测试用 `runtime.GOOS == "windows"` t.Skip）
