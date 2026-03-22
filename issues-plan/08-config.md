# 配置管理需求

本文档详细描述 gitcode-cli 配置管理的设计需求和验收标准。

## 模块概述

配置管理模块负责管理 gitcode-cli 的配置信息，包括配置格式、存储位置、安全存储和环境变量支持。

---

## CFG-001: YAML 配置格式

### 功能描述

使用 YAML 格式存储配置，与 GitHub CLI 保持一致。

### 配置文件结构

```yaml
# ~/.config/gc/hosts.yml
hosts:
  gitcode.com:
    user: username
    git_protocol: https
    users:
      username:
        name: Username
        # Token 存储在 Keyring，不在配置文件中

# ~/.config/gc/config.yml
git_protocol: https
editor: vim
prompt: enabled
pager: less
```

### 配置项说明

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `git_protocol` | string | https | Git 协议 (https/ssh) |
| `editor` | string | $EDITOR | 文本编辑器 |
| `prompt` | string | enabled | 交互提示 |
| `pager` | string | less | 分页程序 |
| `browser` | string | | 浏览器 |

### 主机配置项

| 配置项 | 类型 | 说明 |
|--------|------|------|
| `user` | string | 当前活跃用户 |
| `git_protocol` | string | 该主机的 Git 协议 |
| `users` | map | 用户列表 |
| `users.{name}` | map | 用户信息 |

### 验收标准

- [ ] 支持 YAML 格式读写
- [ ] 支持默认值
- [ ] 支持主机级配置覆盖
- [ ] 配置格式与 gh 兼容

---

## CFG-002: 配置存储位置

### 功能描述

定义配置文件的存储位置，支持跨平台。

### 存储位置

| 平台 | 配置目录 | 环境变量覆盖 |
|------|----------|--------------|
| Linux/macOS | `~/.config/gc/` | `GC_CONFIG_DIR` |
| Windows | `%APPDATA%\gc\` | `GC_CONFIG_DIR` |

### 文件列表

```
~/.config/gc/
├── hosts.yml        # 主机和认证配置
├── config.yml       # 全局配置
├── aliases.yml      # 别名配置
└── state.yml        # 状态文件（上次检查更新等）
```

### 配置优先级

```
1. 命令行参数 (--flag value)
2. 环境变量 (GC_*)
3. 主机配置 (hosts.yml 中的主机配置)
4. 全局配置 (config.yml)
5. 默认值
```

### 验收标准

- [ ] 支持 Linux/macOS/Windows
- [ ] 支持环境变量覆盖
- [ ] 支持配置优先级
- [ ] 自动创建配置目录

---

## CFG-003: Keyring 安全存储

### 功能描述

使用系统密钥环安全存储敏感信息（如 Token）。

### 存储内容

- OAuth Token
- 用户凭据

### 实现方式

```go
// internal/keyring/keyring.go

// 存储服务名
func serviceName(hostname string) string {
    return fmt.Sprintf("gc:%s", hostname)
}

// 存储 Token
func Set(service, user, token string) error

// 获取 Token
func Get(service, user string) (string, error)

// 删除 Token
func Delete(service, user string) error
```

### 平台支持

| 平台 | 密钥环 |
|------|--------|
| macOS | Keychain |
| Linux | Secret Service (GNOME Keyring/KWallet) |
| Windows | Windows Credential Manager |

### 降级策略

当 Keyring 不可用时：
1. 尝试 Keyring 存储
2. 失败则降级到配置文件存储
3. 警告用户 Token 以明文存储

### 验收标准

- [ ] 支持三大平台密钥环
- [ ] 支持降级存储
- [ ] 降级时显示警告
- [ ] 支持迁移到 Keyring

---

## CFG-004: 环境变量支持

### 功能描述

支持通过环境变量配置 gitcode-cli。

### 环境变量列表

| 环境变量 | 说明 | 示例 |
|----------|------|------|
| `GC_TOKEN` | 认证 Token | `gc_token_xxxx` |
| `GITCODE_TOKEN` | 备选 Token | `gc_token_xxxx` |
| `GC_HOST` | 默认主机 | `gitcode.com` |
| `GC_CONFIG_DIR` | 配置目录 | `/home/user/.config/gc` |
| `GC_API_URL` | API URL | `https://api.gitcode.com/api/v5` |
| `GC_EDITOR` | 编辑器 | `vim` |
| `GC_BROWSER` | 浏览器 | `chrome` |
| `GC_PAGER` | 分页程序 | `less` |
| `GC_GIT_PROTOCOL` | Git 协议 | `https` 或 `ssh` |
| `NO_COLOR` | 禁用颜色 | `1` |

### 使用示例

```bash
# 使用环境变量 Token
export GC_TOKEN="your-token"
gc repo list

# 使用环境变量主机
export GC_HOST="gitcode.example.com"
gc auth status

# 禁用颜色
export NO_COLOR=1
gc issue list
```

### 验收标准

- [ ] 支持所有配置环境变量
- [ ] 环境变量优先级高于配置文件
- [ ] 支持 `NO_COLOR` 环境变量
- [ ] 文档中说明所有环境变量

---

## 配置接口设计

```go
// internal/gc/config.go

type Config interface {
    // 认证配置
    Authentication() AuthConfig

    // 别名配置
    Aliases() AliasConfig

    // 通用配置
    Get(host, key string) (string, error)
    Set(host, key, value string) error

    // Git 协议
    GitProtocol(host string) ConfigEntry

    // 编辑器
    Editor(host string) ConfigEntry

    // 浏览器
    Browser(host string) ConfigEntry

    // 写入配置
    Write() error
}

type AuthConfig interface {
    // Token 管理
    ActiveToken(hostname string) (string, string)
    HasActiveToken(hostname string) bool

    // 用户管理
    ActiveUser(hostname string) (string, error)
    UsersForHost(hostname string) []string
    Hosts() []string

    // 登录/登出
    Login(hostname, username, token, gitProtocol string, secureStorage bool) (bool, error)
    Logout(hostname, username string) error

    // 切换用户
    SwitchUser(hostname, user string) error

    // 默认主机
    DefaultHost() (string, string)
}

type ConfigEntry struct {
    Value  string
    Source string // "environment", "config", "default"
}
```

---

## 相关文档

- [gc-design/docs/config/format.md](../../gc-design/docs/config/format.md)
- [gc-design/docs/config/storage.md](../../gc-design/docs/config/storage.md)

---

**最后更新**: 2026-03-22