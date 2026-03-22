# GitCode CLI

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/gitcode-com/gitcode-cli?include_prereleases)](https://github.com/gitcode-com/gitcode-cli/releases)

GitCode 官方命令行工具，为 GitCode 用户提供便捷的命令行操作体验。

## 安装

### 预编译二进制文件

从 [Releases](https://github.com/gitcode-com/gitcode-cli/releases) 页面下载对应平台的二进制文件：

| 平台 | 架构 | 文件 |
|------|------|------|
| Linux | amd64 | `gc_X.X.X_linux_amd64.tar.gz` |
| Linux | arm64 | `gc_X.X.X_linux_arm64.tar.gz` |
| macOS | amd64 | `gc_X.X.X_darwin_amd64.tar.gz` |
| macOS | arm64 (M1/M2) | `gc_X.X.X_darwin_arm64.tar.gz` |
| Windows | amd64 | `gc_X.X.X_windows_amd64.zip` |

```bash
# Linux/macOS 解压并安装
tar -xzf gc_X.X.X_linux_amd64.tar.gz
sudo mv gc /usr/local/bin/

# Windows 解压后将 gc.exe 放入 PATH 目录
```

### Homebrew (macOS/Linux)

```bash
# 添加 Tap
brew tap gitcode-com/tap

# 安装
brew install gc

# 或一键安装
brew install gitcode-com/tap/gc
```

### Scoop (Windows)

```bash
# 添加 Bucket
scoop bucket add gitcode-com https://github.com/gitcode-com/scoop-bucket

# 安装
scoop install gc
```

### Linux 包管理器

**DEB (Debian/Ubuntu):**

```bash
# 下载 .deb 包
wget https://github.com/gitcode-com/gitcode-cli/releases/download/vX.X.X/gc_X.X.X_linux_amd64.deb

# 安装
sudo dpkg -i gc_X.X.X_linux_amd64.deb
```

**RPM (RHEL/CentOS/Fedora):**

```bash
# 下载 .rpm 包
wget https://github.com/gitcode-com/gitcode-cli/releases/download/vX.X.X/gc_X.X.X_linux_amd64.rpm

# 安装
sudo rpm -i gc_X.X.X_linux_amd64.rpm
```

### Docker

```bash
# 拉取镜像
docker pull gitcode/gc:latest

# 或使用 GitHub Container Registry
docker pull ghcr.io/gitcode-com/gc:latest

# 运行
docker run --rm -it gitcode/gc:latest version

# 挂载配置目录
docker run --rm -it -v ~/.config/gc:/root/.config/gc gitcode/gc:latest auth status
```

### 从源码构建

**前置要求:**
- Go 1.22+
- Make (可选)

```bash
# 克隆仓库
git clone https://github.com/gitcode-com/gitcode-cli.git
cd gitcode-cli

# 安装依赖并构建
make deps
make build

# 安装到系统
make install

# 或直接使用 go install
go install ./cmd/gc
```

## 快速开始

### 认证

```bash
# 交互式登录
gc auth login

# 使用 Token 登录
gc auth login --token YOUR_TOKEN

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
# 创建 Release
gc release create v1.0.0 --title "Version 1.0"

# 列出 Releases
gc release list

# 查看 Release 详情
gc release view v1.0.0

# 删除 Release
gc release delete v1.0.0
```

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
| 🚀 Release 管理 | 创建、列表、查看、删除 |

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
  release     Release 管理 (create, list, view, delete)
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
git clone https://github.com/gitcode-com/gitcode-cli.git
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

## 相关链接

- [GitCode](https://gitcode.com) - GitCode 平台
- [API 文档](https://gitcode.com/docs/api) - GitCode API 参考
- [问题反馈](https://github.com/gitcode-com/gitcode-cli/issues) - 提交 Bug 或建议