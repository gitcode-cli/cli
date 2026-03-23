# GitCode CLI

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/gitcode-cli/cli?include_prereleases)](https://gitcode.com/gitcode-cli/cli/releases)

GitCode 命令行工具，为 GitCode 用户提供便捷的命令行操作体验。

## 安装

### 从源码构建

**前置要求:**
- Go 1.22+

```bash
# 克隆仓库
git clone https://gitcode.com/gitcode-cli/cli.git
cd cli

# 构建
go build -o gc ./cmd/gc

# 安装到用户目录
mkdir -p ~/.local/bin
mv gc ~/.local/bin/
export PATH="$HOME/.local/bin:$PATH"
```

### Linux 包管理器

**DEB (Debian/Ubuntu):**

```bash
# 从 Releases 下载 .deb 包
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.2.1/gc_0.2.1_amd64.deb

# 安装
sudo dpkg -i gc_0.2.1_amd64.deb
```

**RPM (RHEL/CentOS/Fedora):**

```bash
# 从 Releases 下载 .rpm 包
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.2.1/gc-0.2.1-1.x86_64.rpm

# 安装
sudo rpm -i gc-0.2.1-1.x86_64.rpm
```

### 规划中的安装方式

以下安装方式正在开发中：

- [ ] 预编译二进制文件（Linux/macOS/Windows）
- [ ] Homebrew (macOS/Linux)
- [ ] Scoop (Windows)
- [ ] Docker 镜像

## 快速开始

### 认证

**方式一：设置环境变量（推荐）**

```bash
# 设置 Token 环境变量
export GC_TOKEN="your_gitcode_token"

# 或使用备用变量名
export GITCODE_TOKEN="your_gitcode_token"

# 添加到 shell 配置文件（永久生效）
echo 'export GC_TOKEN="your_gitcode_token"' >> ~/.bashrc
source ~/.bashrc
```

**方式二：交互式登录**

```bash
# 交互式登录（需输入 Token）
gc auth login

# 使用 Token 参数登录
gc auth login --token YOUR_TOKEN
```

**获取 Token：**
1. 登录 [GitCode](https://gitcode.com)
2. 进入 设置 -> 私人令牌
3. 点击"生成新令牌"，选择所需权限
4. 复制生成的 Token

**验证认证：**

```bash
# 查看认证状态
gc auth status
```

### 仓库操作

```bash
# 克隆仓库
gc repo clone owner/repo

# 创建新仓库
gc repo create my-repo --public

# 查看仓库信息
gc repo view owner/repo

# Fork 仓库
gc repo fork owner/repo
```

### Issue 管理

```bash
# 创建 Issue
gc issue create --title "Bug report" --body "Description"

# 列出 Issues
gc issue list --state open

# 查看 Issue 详情
gc issue view 123

# 关闭 Issue
gc issue close 123
```

### Pull Request 管理

```bash
# 创建 PR
gc pr create --title "New feature" --base main

# 列出 PR
gc pr list --state open

# 查看 PR 详情
gc pr view 456

# 检出 PR 分支
gc pr checkout 456

# 合并 PR
gc pr merge 456 --squash

# 代码检视
gc pr review 456 --approve
gc pr review 456 --request-changes
```

### Release 管理

```bash
# 创建 Release（建议包含 --notes 参数）
gc release create v1.0.0 --title "Version 1.0" --notes "Release notes"

# 列出 Releases
gc release list -R owner/repo

# 查看 Release 详情
gc release view v1.0.0 -R owner/repo

# 上传资产到 Release
gc release upload v1.0.0 app.zip -R owner/repo
gc release upload v1.0.0 file1.tar.gz file2.rpm -R owner/repo

# 下载 Release 资产
gc release download v1.0.0 -R owner/repo
gc release download v1.0.0 app.zip -R owner/repo -o ./downloads/

# 删除 Release
gc release delete v1.0.0 -R owner/repo
```

> **注意**: `release create` 命令建议包含 `--notes` 参数，否则可能返回错误。

## 功能特性

| 功能 | 描述 |
|------|------|
| 🔐 认证管理 | Token 认证、多账户支持、安全存储 |
| 📦 仓库操作 | 克隆、创建、Fork、查看、删除 |
| 🐛 Issue 管理 | 创建、列表、查看、关闭、重开、评论 |
| 🔀 PR 管理 | 创建、列表、查看、检出、合并、关闭 |
| 👀 代码检视 | 批准、请求修改、添加评论 |
| 🏷️ 标签管理 | 创建、列表、删除 |
| 🎯 里程碑管理 | 创建、列表、查看、删除 |
| 🚀 Release 管理 | 创建、列表、查看、上传资产、下载资产 |

## 命令概览

```
gc <command> <subcommand> [flags]

Commands:
  auth        认证管理 (login, logout, status, token)
  repo        仓库操作 (clone, create, list, view, fork, delete)
  issue       Issue 管理 (create, list, view, close, reopen, comment)
  pr          PR 管理 (create, list, view, checkout, merge, close, reopen, review, diff, ready)
  label       标签管理 (create, list, delete)
  milestone   里程碑管理 (create, list, view, delete)
  release     Release 管理 (create, list, view, upload, download, delete)
  version     显示版本信息
```

## 配置

配置文件位置: `~/.config/gc/config.yaml`

```yaml
# 默认主机
host: gitcode.com

# Git 协议
git_protocol: https

# 默认编辑器
editor: vim

# 分页器
pager: less
```

### 环境变量

| 变量 | 描述 |
|------|------|
| `GC_TOKEN` | 认证 Token |
| `GITCODE_TOKEN` | 备用 Token |
| `GC_HOST` | 默认主机 |
| `NO_COLOR` | 禁用颜色输出 |

## Shell 补全

```bash
# Bash
gc completion bash > /etc/bash_completion.d/gc
source ~/.bashrc

# Zsh
gc completion zsh > "${fpath[1]}/_gc"
source ~/.zshrc

# Fish
gc completion fish > ~/.config/fish/completions/gc.fish
source ~/.config/fish/config.fish
```

## 文档

- [贡献指南](./CONTRIBUTING.md) - 开发和发布流程
- [安全策略](./SECURITY.md) - 敏感信息保护和安全规范
- [CLAUDE.md](./CLAUDE.md) - AI 辅助开发指南
- [需求文档](./issues-plan/) - 完整需求规格和里程碑规划

## 技术栈

| 组件 | 技术 |
|------|------|
| 语言 | Go 1.22+ |
| 命令框架 | Cobra |
| 配置格式 | YAML |
| 目标平台 | Linux, macOS, Windows |

## 开发

```bash
# 克隆仓库
git clone https://gitcode.com/gitcode-cli/cli.git
cd gitcode-cli

# 安装依赖
make deps

# 构建
make build

# 运行测试
make test

# 代码检查
make lint

# 运行
make run
```

## 贡献

欢迎贡献代码！请查看 [贡献指南](./CONTRIBUTING.md)。

## 许可证

[MIT License](./LICENSE)

## 致谢

本项目参考了 [GitHub CLI](https://github.com/cli/cli) 的设计与实现，感谢 GitHub 团队的开源贡献。

## 相关链接

- [GitCode](https://gitcode.com) - GitCode 平台
- [API 文档](https://gitcode.com/docs/api) - GitCode API 参考
- [问题反馈](https://gitcode.com/gitcode-cli/cli/issues) - 提交 Bug 或建议