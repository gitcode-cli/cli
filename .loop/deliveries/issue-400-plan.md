## 开发计划: Issue #400

| # | 文件 | 操作 | 说明 |
|---|------|------|------|
| 1 | `pkg/config/config.go` | 修改 (+14/-1) | 新增 secureWriteFile 函数；writeConfigState 调用替换 |
| 2 | `pkg/config/auth_config.go` | 修改 (+1/-1) | writeAuthState 调用替换为 secureWriteFile |
| 3 | `pkg/config/auth_config_test.go` | 修改 (+60/-1) | 新增 3 个 UT：拒 symlink / 权限硬化 / 新建文件权限；清理遗留 strings import |

### 测试矩阵
| 类型 | 覆盖 | 状态 |
|------|------|------|
| UT | TestSecureWriteFileRejectsSymlink | ✅ |
| UT | TestSecureWriteFileHardensPermissions | ✅ |
| UT | TestSecureWriteFileCreatesNewFileWithRestrictedPermissions | ✅ |
| 系统测试 | 自动化 | ⚠️ 架构限制无法自动化（writeConfigState 无 CLI 入口；writeAuthState 需 API 验证真实 token） |
| 真实命令验证 | 人工 TTY 执行 | ⏳ 待人工执行（清单见 issue-400-design.md） |

### 同步主干
- 2026-07-06：`bugfix/issue-400` fast-forward 14 提交至 `origin/main` e65eb38，零冲突
- 远端新增 `pkg/cmdutil/token_confirm.go`（!343 token 披露保护）经核查与本改动正交
