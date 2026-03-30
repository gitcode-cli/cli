# GitCode CLI

[![AI 操作指南](https://img.shields.io/badge/📖_使用_AI_操作_GitCode_指南-点击查看-FF6B6B?style=for-the-badge)](./docs/AI-GUIDE.md)

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/badge/Release-v0.3.6-blue)](https://gitcode.com/gitcode-cli/cli/releases)

GitCode 命令行工具，为 GitCode 用户提供便捷的命令行操作体验。

## 文档导航

按角色建议从以下入口开始：

| 角色 | 入口 |
|------|------|
| 使用者 | [docs/README.md](./docs/README.md) |
| 开发者 | [spec/README.md](./spec/README.md) |
| Codex / 代理 | [AGENTS.md](./AGENTS.md) |
| Claude | [CLAUDE.md](./CLAUDE.md) |

主要文档：

- [命令手册](./docs/COMMANDS.md)
- [认证说明](./docs/AUTH.md)
- [回归说明](./docs/REGRESSION.md)
- [打包说明](./docs/PACKAGING.md)
- [开发规范](./spec/README.md)
- [当前进度](./issues-plan/PROGRESS.md)

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
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gc_0.3.6_amd64.deb

# 安装
sudo dpkg -i gc_0.3.6_amd64.deb
```

**RPM (RHEL/CentOS/Fedora):**

```bash
# 从 Releases 下载 .rpm 包
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gc-0.3.6-1.x86_64.rpm

# 安装
sudo rpm -i gc-0.3.6-1.x86_64.rpm
```

### Wheel 包（跨平台，推荐）

从 Release 归档下载 wheel 包安装：

```bash
# 创建虚拟环境
python3 -m venv .venv
source .venv/bin/activate  # Linux/macOS
# .venv\Scripts\activate   # Windows

# 安装（一行命令）
pip install https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gitcode_cli-0.3.6-py3-none-any.whl
```

### PyPI（备选）

```bash
# 创建虚拟环境
python3 -m venv .venv
source .venv/bin/activate  # Linux/macOS
# .venv\Scripts\activate   # Windows

# 安装
pip install gitcode-cli
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

当前版本认证优先级：

1. `GC_TOKEN`
2. `GITCODE_TOKEN`
3. 本地登录配置

说明：
- `gc auth login` 会将认证信息持久化到本地配置目录
- 如果设置了环境变量，环境变量始终覆盖本地配置
- `gc auth logout` 只清理本地配置，不会自动取消环境变量
- 详细规则见 [docs/AUTH.md](./docs/AUTH.md)

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

> 详细命令行为和完整示例请查看 [docs/COMMANDS.md](./docs/COMMANDS.md)。

## 常用示例

### 仓库操作

```bash
# 克隆仓库
gc repo clone owner/repo

# 创建新仓库
gc repo create my-repo --public

# 查看仓库信息
gc repo view owner/repo
gc repo view

# Fork 仓库
gc repo fork owner/repo

# 删除仓库
gc repo delete owner/repo
```

### Issue 管理

```bash
# 创建 Issue
gc issue create --title "Bug report" --body "Description"

# 列出 Issues
gc issue list --state open

# 查看 Issue 详情
gc issue view 123

# 编辑 Issue
gc issue edit 123 --title "New title"
gc issue edit 123 --body "New description"
gc issue edit 123 --state close
gc issue edit 123 --assignee username
gc issue edit 123 --label bug,enhancement
gc issue edit 123 --milestone 5

# 关闭 Issue
gc issue close 123

# 重开 Issue
gc issue reopen 123

# 添加评论
gc issue comment 123 --body "Comment text"

# 列出 Issue 评论
gc issue comments 123
gc issue comments 123 --limit 10
gc issue comments 123 --order desc

# 编辑评论
gc issue comment edit 166061383 --body "Updated comment"
gc issue comment edit --id 166061383 --body "Updated comment"

# 从文件添加评论
gc issue comment 123 --body-file comment.txt

# 管理 Issue 标签
gc issue label 123 --add bug,enhancement

# 查看 Issue 关联的 PRs
gc issue prs 123
```

### Pull Request 管理

```bash
# 创建 PR（自动检测当前分支）
gc pr create --title "New feature" --base main

# 创建跨仓库 PR（从 fork 到 upstream）
gc pr create -R upstream/repo --fork myfork/repo --title "Feature"

# 列出 PR
gc pr list --state open

# 查看 PR 详情
gc pr view 456

# 查看 PR 评论
gc pr comments 456

# 回复 PR 评论
gc pr reply 456 --discussion <discussion_id> --body "Reply text"

# 当前不支持通过 CLI 解决/取消解决 PR 评论
# resolved 状态需要在 GitCode Web UI 中手动处理

# 检出 PR 分支
gc pr checkout 456

# 合并 PR
gc pr merge 456 --squash

# 关闭/重开 PR
gc pr close 456
gc pr reopen 456

# 编辑 PR
gc pr edit 456 --title "New title"
gc pr edit 456 --draft true

# 代码检视
gc pr review 456 --approve
gc pr review 456 --comment "Review comment"
gc pr review 456 --approve --comment "LGTM"
```

更多命令细节、仓库参数格式、自动推断行为、Release 命令和平台限制，请查看 [docs/COMMANDS.md](./docs/COMMANDS.md)。

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

## 项目定位

当前仓库已经建立：

- 用户文档入口：[`docs/`](./docs/README.md)
- 正式规范入口：[`spec/`](./spec/README.md)
- Codex 入口：[`AGENTS.md`](./AGENTS.md)
- Claude 入口：[`CLAUDE.md`](./CLAUDE.md)

如果你要看完整规范、构建与发布规则、质量门禁和 AI 协作边界，请直接进入对应入口，不要仅依赖本 README。

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

欢迎贡献代码。开始前请查看 [贡献指南](./CONTRIBUTING.md) 和 [spec/README.md](./spec/README.md)。

## 许可证

[MIT License](./LICENSE)

## 致谢

本项目参考了 [GitHub CLI](https://github.com/cli/cli) 的设计与实现，感谢 GitHub 团队的开源贡献。

## 相关链接

- [GitCode](https://gitcode.com) - GitCode 平台
- [API 文档](https://gitcode.com/docs/api) - GitCode API 参考
- [问题反馈](https://gitcode.com/gitcode-cli/cli/issues) - 提交 Bug 或建议
