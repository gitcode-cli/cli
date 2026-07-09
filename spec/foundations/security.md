# 安全规范

本文档定义 gitcode-cli 项目的安全要求。

## 职责

定义 token、敏感信息、测试安全和提交前安全检查要求。

## 适用场景

- 涉及认证、token、config 持久化
- 提交代码前做安全检查
- 调整发布或自动化相关凭证处理

## 必须

- 不硬编码 token 或密钥
- 通过环境变量或本地配置管理认证信息
- 在提交前检查敏感信息和敏感文件

## 禁止

- 在源码、测试或文档中提交真实凭证
- 使用非 `infra-test/*` 仓库做真实命令测试

## 同步要求

- 认证模型变化时同步 `docs/AUTH.md`、相关命令文档和 AI 入口文档
- 自动化或发布涉及凭证时同步相关交付规范

## 不负责什么

- 一般代码风格
- PR 提交流程
- 命令行为说明

## Token 管理

### 基本原则

1. **禁止硬编码** - Token 绝对不能写入源代码
2. **允许本地配置存储** - Token 可以保存到用户本地配置目录，但不能提交到版本控制系统
3. **环境变量优先** - 若同时存在环境变量和本地配置，环境变量优先
4. **禁止提交** - Token 绝对不能提交到版本控制系统
5. **显示 Token 必须人工确认** - `auth token` 和 `auth status --show-token` 会输出完整 token，必须在交互式 TTY 中输入 hostname 确认；非交互环境一律拒绝，AI 代理默认不得调用这些命令
6. **禁止提取 token 走 curl/脚本** - AI 代理不得提取 `gh auth token` 走 curl/脚本调 API（token 会泄露到进程列表/命令历史/日志）；gh 已封装认证，网络问题用 gh 多次重试，多次重试失败报告人工解决，不自行绕过

### Token 获取方式

```bash
# 方式一：环境变量（推荐）
export GC_TOKEN="your_gitcode_token"
# 或
export GITCODE_TOKEN="your_gitcode_token"

# 方式二：stdin 登录
echo "YOUR_TOKEN" | gc auth login --with-token
```

本地配置默认写入：

```bash
~/.config/gc/auth.json
```

### 凭证文件落盘约束

写入凭证文件（`auth.json`）与配置状态文件（`config.json`）时必须满足：

1. **文件权限 0600** - 凭证文件落盘权限必须为 `0600`，不得依赖 `umask`；对既有文件也必须显式 `Chmod` 收紧（`os.WriteFile` 受 umask 影响且不收紧既有文件权限，必须额外 `Chmod`）
2. **目录权限 0700** - 配置目录 `~/.config/gc/` 权限必须为 `0700`，创建后显式 `Chmod` 收紧
3. **拒绝符号链接** - 写入路径若为符号链接必须拒绝，防止凭证重定向攻击（攻击者将 `auth.json` 替换为指向他处的软链，CLI 写入时覆盖目标或泄露内容）
4. **原子检测** - Unix 实现必须用 `O_NOFOLLOW` 原子打开 + fd `fchmod`，消除 `Lstat`→`WriteFile` 的 TOCTOU 竞态窗口；Windows 无 `O_NOFOLLOW`，保留 `Lstat` 检测（依赖 ACL 模型缓解）
5. **威胁模型假设** - 依赖目录 `0700` 防止跨用户父目录/硬链接攻击；同用户攻击者本可读凭证文件，不在防护范围；父目录组件为符号链接的彻底防护需 `openat2(RESOLVE_NO_SYMLINKS)`，非跨平台，不纳入

### 获取 Token
1. 登录 [GitCode](https://gitcode.com)
2. 进入 设置 → 私人令牌
3. 生成新令牌并复制

## 敏感信息保护

### 禁止提交的内容

| 文件类型 | 说明 |
|----------|------|
| `*.pem`, `*.key`, `*.p12`, `*.pfx` | 私钥/证书文件 |
| `id_rsa*`, `id_ed*` | SSH 密钥 |
| `.env`, `.env.*` | 环境变量文件 |
| `*.secret` | 密钥文件 |
| `credentials.json` | 凭证文件 |
| `token.txt`, `*.token` | Token 文件 |
| `secrets.yaml`, `secrets.yml` | 密钥配置 |

### .gitignore 配置

确保以下内容已在 `.gitignore` 中：
```gitignore
# Secrets
.env
.env.*
*.pem
*.key
credentials.json
token.txt
secrets.yaml
```

### 提交前密钥扫描（pre-commit + pre-push）

仓库通过 `.pre-commit-config.yaml` 配置自动化密钥扫描，分两个阶段执行：

| 阶段 | 触发时机 | 执行的 hook |
|------|----------|-------------|
| pre-commit | `git commit` | 全部 hook（gofmt + 语法校验 + 私钥检测 + gitleaks(workspace) + 大文件 + 空白/换行） |
| pre-push | `git push` | 安全子集：`gitleaks(workspace)` + `gitleaks-history` + `detect-private-key` + `check-added-large-files` |

**密钥扫描能力**：
- `detect-private-key` — 检测私钥**文件**（SSH/PEM/RSA 等）
- `gitleaks`（workspace）— 使用 `--no-git` 扫描工作区（含 gitignored 文件）中的 **token/密钥字符串**，覆盖 `GC_TOKEN`/`glpat-`/`gho_`/AWS key 等模式
- `gitleaks-history`— 扫描 **git 提交历史**中的 token/密钥字符串，补充 workspace 扫描无法覆盖已提交但已删除的泄漏

**安装**（首次或 clone 后）：

```bash
pre-commit install --hook-type pre-commit --hook-type pre-push
```

或用 `gc precommit check` 检测工具是否安装 + hook 是否初始化。

**注意**：
- pre-push hook 需单独安装（默认 `pre-commit install` 只装 pre-commit）
- 未安装 pre-push 时，`stages: [pre-push]` 的 hook 不生效，push 前密钥扫描缺失
- gitleaks 是 local hook（`language: system`），需手动安装并加入 PATH：`brew install gitleaks` / `go install github.com/gitleaks/gitleaks/v8/cmd/gitleaks@latest` / 从 https://github.com/gitleaks/gitleaks/releases 下载二进制。`pre-commit autoupdate` 仅对 remote repo hook 生效，不更新 local hook 的 gitleaks 二进制

## 测试安全

### 测试仓库限制

**只能使用以下测试仓库：**

| 仓库 | 用途 |
|------|------|
| `infra-test/gctest1` | 主要测试仓库（首选） |
| `infra-test` 组织下其他仓库 | 其他测试场景 |

**禁止行为：**
- ❌ 使用个人仓库测试
- ❌ 使用其他组织或用户的仓库测试
- ❌ 使用 `gitcode-cli/cli` 测试

### 测试 Token 来源

- 通过环境变量 `GC_TOKEN` 传递
- 或运行时手动输入

### 测试配置

```
测试组织: infra-test
API 基础 URL: https://api.gitcode.com/api/v5
Token 来源: 环境变量或运行时输入
```

## 代码审查检查项

提交代码前必须确认：

- [ ] 没有硬编码的 Token 或密钥
- [ ] 配置文件中不包含敏感信息
- [ ] `.gitignore` 已忽略敏感文件
- [ ] 测试代码不包含真实 Token
- [ ] 文档中不包含真实凭证

## CI/CD 安全

### GitHub Actions

使用 Secrets 存储敏感信息：

```yaml
env:
  GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  GC_TOKEN: ${{ secrets.GC_TOKEN }}
```

### PyPI 发布

使用 Trusted Publishing（OIDC），无需配置 API Token。

## 安全检查命令

```bash
# 检查即将提交的内容
git diff --cached

# 检查历史中的敏感信息
git log -p | grep -iE "token|password|secret|api_key"

# 检查是否有敏感文件被追踪
git ls-files | grep -iE "\.pem|\.key|\.env|credentials|secret"
```

## 报告安全问题

如果您发现安全问题，请：

1. 不要在公开 Issue/PR/comment 中披露漏洞细节、PoC 或攻击细节
2. 提交私密 Issue 报告（使用 `gc issue create --security-hole` 标记为私有 issue）
3. 提供详细的问题描述和复现步骤

---

**最后更新**: 2026-07-09
