# 里程碑 2: 认证功能

## 概述

实现完整的 GitCode 认证功能，支持 OAuth Device Flow 和 Token 认证。

**预计工期**: 1 周

**依赖**: 里程碑 1 (基础架构)

**目标**: 用户能够通过多种方式登录 GitCode，安全存储凭据

---

## 任务清单

### AUTH-001: OAuth Device Flow 登录

**优先级**: P0

**任务描述**:

实现 OAuth Device Flow 认证流程，这是主要的认证方式。

**文件**:

```
pkg/cmd/auth/login/login.go
internal/authflow/oauth.go
internal/authflow/device_flow.go
```

**认证流程**:

```
1. 调用 /oauth/device/code 获取设备码
2. 显示验证 URL 和用户码
3. 轮询 /oauth/device/token 等待用户授权
4. 获取 Access Token
5. 存储到 Keyring
```

**验收标准**:

- [ ] `gc auth login` 启动 OAuth 流程
- [ ] 显示验证 URL 和用户码
- [ ] 等待用户完成授权
- [ ] 获取并存储 Token
- [ ] 显示登录成功信息

**示例输出**:

```bash
$ gc auth login
? What account do you want to log into? GitCode
? What is your preferred protocol for Git operations? HTTPS
! First copy your one-time code: XXXX-XXXX
- Press Enter to open gitcode.com in your browser...
✓ Authentication complete. Press Enter to continue...

✓ Logged in as username
✓ Configured git protocol
✓ Logged in to gitcode.com account username
```

---

### AUTH-002: Token 认证

**优先级**: P0

**任务描述**:

支持通过 Token 直接登录，适用于 CI/CD 环境。

**文件**:

```
pkg/cmd/auth/login/login.go
internal/authflow/token.go
```

**验收标准**:

- [ ] `gc auth login --with-token` 从 stdin 读取 Token
- [ ] `gc auth login --token "xxx"` 直接指定 Token
- [ ] 验证 Token 有效性
- [ ] 获取用户信息

**示例**:

```bash
# 从 stdin 读取
echo "your-token" | gc auth login --with-token

# 从文件读取
cat token.txt | gc auth login --with-token

# 直接指定
gc auth login --token "your-token"
```

---

### AUTH-003: Keyring 集成

**优先级**: P0

**任务描述**:

集成系统密钥环，安全存储 Token。

**文件**:

```
internal/keyring/keyring.go
internal/keyring/keyring_darwin.go
internal/keyring/keyring_linux.go
internal/keyring/keyring_windows.go
```

**平台支持**:

| 平台 | 后端 |
|------|------|
| macOS | Keychain |
| Linux | Secret Service |
| Windows | Credential Manager |

**降级策略**:

```go
// Keyring 不可用时降级到文件存储
func Set(service, user, token string) error {
    err := keyring.Set(service, user, token)
    if err != nil {
        // 降级到配置文件
        return saveToConfig(service, user, token)
    }
    return nil
}
```

**验收标准**:

- [ ] 支持 macOS/Linux/Windows
- [ ] Token 不明文存储
- [ ] 支持降级策略
- [ ] 降级时显示警告

---

### AUTH-004: 认证状态查看

**优先级**: P0

**任务描述**:

实现 `gc auth status` 命令，查看登录状态。

**文件**:

```
pkg/cmd/auth/status/status.go
```

**验收标准**:

- [ ] 显示登录状态
- [ ] 显示用户名
- [ ] 显示 Token 来源
- [ ] 显示 Git 协议

**示例输出**:

```bash
$ gc auth status
gitcode.com
  ✓ Logged in to gitcode.com as username (keyring)
  ✓ Git operations protocol: https
  ✓ Token: ghp_xxxx...xxxx (from keyring)

$ gc auth status
gitcode.com
  ✗ Not logged in
  ✓ To authenticate, run: gc auth login
```

---

### AUTH-005: 登出功能

**优先级**: P1

**任务描述**:

实现 `gc auth logout` 命令，移除认证信息。

**文件**:

```
pkg/cmd/auth/logout/logout.go
```

**验收标准**:

- [ ] 移除本地 Token
- [ ] 从 Keyring 删除
- [ ] 显示登出确认

**示例**:

```bash
$ gc auth logout
? Confirm logout of gitcode.com account username? Yes
✓ Logged out of gitcode.com account username
```

---

### AUTH-006: Token 打印

**优先级**: P1

**任务描述**:

实现 `gc auth token` 命令，打印认证 Token。

**文件**:

```
pkg/cmd/auth/token/token.go
```

**验收标准**:

- [ ] 打印有效 Token
- [ ] 支持安全输出（stderr）
- [ ] 支持环境变量 Token

**示例**:

```bash
$ gc auth token
ghp_xxxxxxxxxxxxxxxxxxxx

$ gc auth token --hostname gitcode.com
ghp_xxxxxxxxxxxxxxxxxxxx
```

---

### AUTH-007: 多账户支持

**优先级**: P1

**任务描述**:

支持多账户登录和切换。

**文件**:

```
pkg/cmd/auth/switch/switch.go
internal/config/hosts_config.go
```

**验收标准**:

- [ ] 支持多个账户登录
- [ ] `gc auth switch` 切换账户
- [ ] 显示当前活跃账户

**示例**:

```bash
$ gc auth login --hostname gitcode.com
# 登录第一个账户

$ gc auth login --hostname gitcode.com
# 登录第二个账户

$ gc auth switch
? Switch to account: [user1, user2]

$ gc auth status
gitcode.com
  ✓ Logged in as user2 (active)
  Other accounts: user1
```

---

### AUTH-008: 环境变量支持

**优先级**: P0

**任务描述**:

支持通过环境变量提供 Token。

**环境变量**:

| 变量 | 优先级 | 说明 |
|------|--------|------|
| `GC_TOKEN` | 最高 | GitCode Token |
| `GITCODE_TOKEN` | 次高 | 备选 Token |

**验收标准**:

- [ ] `GC_TOKEN` 环境变量有效
- [ ] 环境变量优先级高于配置文件
- [ ] `gc auth status` 显示 Token 来源

**示例**:

```bash
export GC_TOKEN="your-token"
gc auth status
# gitcode.com
#   ✓ Logged in as username (env)
#   ✓ Token: from GC_TOKEN environment variable
```

---

## 依赖关系

```
AUTH-001 (OAuth) ─┬─→ AUTH-003 (Keyring)
                  │
AUTH-002 (Token) ─┘

AUTH-003 (Keyring) → AUTH-004 (Status)
                     AUTH-005 (Logout)
                     AUTH-006 (Token)

AUTH-003 (Keyring) → AUTH-007 (多账户)

AUTH-008 (环境变量) → AUTH-004 (Status)
```

---

## 完成标准

里程碑 M2 完成需满足：

1. ✅ `gc auth login` 完成 OAuth 认证
2. ✅ `gc auth login --with-token` Token 登录成功
3. ✅ `gc auth status` 正确显示状态
4. ✅ Token 安全存储在 Keyring
5. ✅ 环境变量 Token 支持
6. ✅ 单元测试覆盖率 ≥ 80%

---

## 测试用例

### 单元测试

```bash
go test ./pkg/cmd/auth/... -v
go test ./internal/authflow/... -v
go test ./internal/keyring/... -v
```

### 集成测试

```bash
# 需要 GC_TEST_TOKEN
go test -tags=integration ./pkg/cmd/auth/... -v
```

### 手动测试清单

- [ ] OAuth Device Flow 完整流程
- [ ] Token 登录
- [ ] 环境变量 Token
- [ ] Keyring 存储
- [ ] 配置文件降级
- [ ] 登出功能
- [ ] 多账户切换

---

## API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/oauth/device/code` | POST | 获取设备码 |
| `/oauth/device/token` | POST | 获取 Token |
| `/api/v5/user` | GET | 获取用户信息 |

---

## 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| Keyring 不可用 | 中 | 实现降级策略 |
| OAuth 端点变化 | 高 | 端点可配置 |
| Token 过期 | 中 | 实现 refresh 机制 |

---

## 相关文档

- [issues-plan/03-module-auth.md](../03-module-auth.md)
- [gc-api-doc/doc/01-authentication.md](../../../gc-api-doc/doc/01-authentication.md)

---

**最后更新**: 2026-03-22