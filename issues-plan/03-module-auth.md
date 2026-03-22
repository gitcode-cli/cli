# 认证模块需求 (auth)

本文档详细描述 gitcode-cli 认证模块的功能需求、验收标准和 API 映射。

## 模块概述

认证模块提供 GitCode 平台的身份认证功能，支持 OAuth Device Flow 和 Token 认证两种方式，并支持多账户管理。

### 命令结构

```
gc auth <command>

Commands:
  login    Log in to a GitCode account
  logout   Log out of a GitCode account
  status   View authentication status
  token    Print an authentication token
  switch   Switch between GitCode accounts
  refresh  Refresh authentication token
```

---

## AUTH-001: auth login - OAuth Device Flow 认证

### 功能描述

通过 OAuth Device Flow 进行浏览器认证登录。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --hostname | -h | string | gitcode.com | GitCode 实例主机名 |
| --web | -w | bool | false | 打开浏览器认证 |
| --scopes | | []string | | 额外的权限范围 |
| --git-protocol | | string | https | Git 协议 (https/ssh) |

### 认证流程

1. 获取设备码（Device Code）
2. 显示验证码给用户
3. 打开浏览器进行授权
4. 轮询等待授权完成
5. 获取 Token 并存储

### 使用示例

```bash
# 交互式登录
gc auth login

# 指定主机名登录
gc auth login --hostname gitcode.example.com

# 使用 SSH 协议
gc auth login --git-protocol ssh
```

### 验收标准

- [ ] `gc auth login` 能启动 OAuth Device Flow
- [ ] 正确显示设备码和验证 URL
- [ ] 能打开浏览器进行授权
- [ ] 授权成功后正确保存 Token
- [ ] Token 存储在 Keyring 中（降级到配置文件）
- [ ] 显示登录成功的用户名

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/oauth/device/code` | POST | 获取设备码 |
| `/oauth/device/token` | POST | 轮询获取 Token |
| `/api/v5/user` | GET | 获取用户信息验证 Token |

### 测试用例映射

- 参考 `gc-api-doc/test/test_users.py`

---

## AUTH-002: auth login --with-token - Token 认证

### 功能描述

通过 Token 直接登录，适用于 CI/CD 环境或已有 Token 的场景。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --with-token | | bool | false | 从标准输入读取 Token |
| --hostname | -h | string | gitcode.com | GitCode 实例主机名 |

### 使用示例

```bash
# 从标准输入读取 Token
echo "your-token" | gc auth login --with-token

# 从文件读取 Token
cat token.txt | gc auth login --with-token

# 环境变量方式
export GC_TOKEN="your-token"
gc auth login --with-token
```

### 验收标准

- [ ] 能从标准输入正确读取 Token
- [ ] 验证 Token 有效性
- [ ] Token 无效时显示清晰错误信息
- [ ] 正确保存 Token 和用户信息

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/user` | GET | 验证 Token 并获取用户信息 |

---

## AUTH-003: auth logout - 登出账户

### 功能描述

登出指定的 GitCode 账户，删除本地存储的认证信息。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --hostname | -h | string | | 指定登出的主机 |
| --user | -u | string | | 指定登出的用户 |
| --confirm | -y | bool | false | 跳过确认提示 |

### 使用示例

```bash
# 登出当前账户
gc auth logout

# 登出指定主机
gc auth logout --hostname gitcode.example.com

# 登出指定用户
gc auth logout --user username

# 跳过确认
gc auth logout --confirm
```

### 验收标准

- [ ] 正确删除 Keyring 中的 Token
- [ ] 正确删除配置文件中的用户信息
- [ ] 多账户时正确切换到其他账户
- [ ] 显示登出成功的确认信息

---

## AUTH-004: auth status - 查看认证状态

### 功能描述

查看当前登录状态和账户信息。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --hostname | -h | string | | 查看指定主机的状态 |
| --show-token | | bool | false | 显示完整 Token（谨慎使用） |

### 使用示例

```bash
# 查看认证状态
gc auth status

# 查看指定主机
gc auth status --hostname gitcode.example.com

# 显示 Token
gc auth status --show-token
```

### 输出示例

```
gitcode.com:
  ✓ Logged in as username
  Token source: keyring
```

### 验收标准

- [ ] 正确显示已登录的主机列表
- [ ] 显示每个主机的登录用户名
- [ ] 显示 Token 来源（keyring/config/environment）
- [ ] 未登录时显示提示信息
- [ ] 验证 Token 是否仍然有效

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/user` | GET | 验证当前 Token |

---

## AUTH-005: auth token - 打印认证 Token

### 功能描述

打印当前账户的认证 Token，用于脚本或 CI/CD 场景。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --hostname | -h | string | | 指定主机 |
| --security | | string | filter | 安全级别 (none/filter/all) |

### 使用示例

```bash
# 打印遮蔽的 Token
gc auth token

# 打印完整 Token
gc auth token --security all
```

### 输出示例

```
gt_1...abcd
Token source: keyring
```

### 验收标准

- [ ] 默认输出遮蔽的 Token（前4位...后4位）
- [ ] `--security all` 输出完整 Token
- [ ] 显示 Token 来源
- [ ] 未登录时返回错误

---

## AUTH-006: auth switch - 切换账户

### 功能描述

在多个账户之间切换。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --hostname | -h | string | | 指定主机 |
| --user | -u | string | | 切换到的用户名 |

### 使用示例

```bash
# 交互式选择账户
gc auth switch

# 切换到指定用户
gc auth switch --user another-user
```

### 验收标准

- [ ] 显示可切换的账户列表
- [ ] 正确切换到选定账户
- [ ] 更新配置文件中的活跃用户
- [ ] 显示切换成功的确认信息

---

## Token 存储策略

### 存储位置优先级

1. **Keyring**（首选）- 使用系统密钥环安全存储
2. **配置文件**（降级）- 明文存储在 `~/.config/gc/hosts.yml`

### 环境变量支持

| 环境变量 | 优先级 | 说明 |
|----------|--------|------|
| `GC_TOKEN` | 高 | GitCode CLI Token |
| `GITCODE_TOKEN` | 中 | 备选 Token |
| `GC_HOST` | - | 默认主机 |

### Token 优先级

```
1. GC_TOKEN 环境变量
2. GITCODE_TOKEN 环境变量
3. Keyring 存储（加密）
4. 配置文件存储（明文）
```

---

## 权限范围 (Scopes)

| Scope | 说明 |
|-------|------|
| `api` | 完整 API 访问 |
| `read_user` | 读取用户信息 |
| `write_repository` | 写入仓库 |
| `read_ssh_key` | 读取 SSH 公钥 |

---

## 相关文档

- [gc-api-doc/doc/01-authentication.md](../../gc-api-doc/doc/01-authentication.md)
- [gc-api-doc/test/test_users.py](../../gc-api-doc/test/test_users.py)

---

**最后更新**: 2026-03-22