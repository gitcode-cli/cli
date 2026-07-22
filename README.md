# GitCode CLI

[![AI 操作指南](https://img.shields.io/badge/📖_使用_AI_操作_GitCode_指南-点击查看-FF6B6B?style=for-the-badge)](./docs/AI-GUIDE.md)

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/badge/Release-latest-blue)](https://gitcode.com/gitcode-cli/cli/releases)

GitCode CLI 把仓库、Issue、PR、Release 和 Actions 带回终端，让开发者减少页面切换，也让脚本与 AI 获得结构化、可审计、带安全边界的 GitCode 执行入口。

[快速了解核心价值与应用场景，并在五分钟内开始使用](./docs/INTRODUCTION.md)。

## 文档导航

按角色建议从以下入口开始：

| 角色 | 入口 |
|------|------|
| 使用者 | [docs/README.md](./docs/README.md) |
| 开发者 | [spec/README.md](./spec/README.md) |
| Codex / 代理 | [AGENTS.md](./AGENTS.md) |
| Claude | [CLAUDE.md](./CLAUDE.md) |

主要文档：

- [产品介绍与快速上手](./docs/INTRODUCTION.md)
- [命令手册](./docs/COMMANDS.md)
- [认证说明](./docs/AUTH.md)
- [回归说明](./docs/REGRESSION.md)
- [打包说明](./docs/PACKAGING.md)
- [AI 操作指南（外部项目）](./docs/AI-GUIDE.md)
- [应用案例库](./Example/index.md)
- [开发规范](./spec/README.md)
- [真相源矩阵](./spec/governance/source-of-truth-matrix.md)
- [AI 本地开发流程](./spec/workflows/ai-local-development-workflow.md)
- [阶段说明](./issues-plan/PROGRESS.md)

## 安装

### 从源码构建

**前置要求:**
- Go 1.22+

```bash
# 克隆仓库（需要 git clone 才能获取版本信息）
git clone https://gitcode.com/gitcode-cli/cli.git
cd cli

# 方式一：使用 go build（推荐）
go build -o gc ./cmd/gc
# 安装
mkdir -p ~/.local/bin
mv gc ~/.local/bin/

# 方式二：使用 make build（带完整版本标签）
make build
# 安装
mkdir -p ~/.local/bin
mv bin/gc ~/.local/bin/

# 添加到 PATH
export PATH="$HOME/.local/bin:$PATH"
```

> **说明**:
> - `go build` 从 `debug.ReadBuildInfo()` 自动获取 git commit 和构建时间（需要 `git clone` 源码）。
> - `make build` 使用 `-ldflags` 注入完整版本标签（如 `v0.3.11-38-g1128f2b`）。

### Linux 包管理器

**DEB (Debian/Ubuntu):**

```bash
# 从 Releases 下载 .deb 包
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.8.0/gc_0.8.0_amd64.deb

# 安装
sudo dpkg -i gc_0.8.0_amd64.deb
```

DEB/RPM packages install both `gc` and `gitcode`; on Linux they are equivalent.

**RPM (RHEL/CentOS/Fedora):**

```bash
# 从 Releases 下载 .rpm 包
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.8.0/gc-0.8.0-1.x86_64.rpm

# 安装
sudo rpm -i gc-0.8.0-1.x86_64.rpm
```

DEB/RPM packages install both `gc` and `gitcode`; on Linux they are equivalent.

### Wheel 包（跨平台，推荐）

从 Release 归档下载 wheel 包安装，**内置全平台二进制**（Linux x64/ARM、macOS Intel/Apple Silicon、Windows x64）：

```bash
# 创建虚拟环境
python3 -m venv .venv
source .venv/bin/activate  # Linux/macOS
# .venv\Scripts\activate   # Windows

# 安装（一行命令）
pip install https://gitcode.com/gitcode-cli/cli/releases/download/v0.8.0/gitcode_cli-0.8.0-py3-none-any.whl

# Windows PowerShell 中推荐使用 gitcode，避免 gc 被内置 Get-Content 别名覆盖
gitcode version
```

说明：
- wheel 会同时安装 `gc` 和 `gitcode` 两个命令入口，功能相同。
- DEB/RPM 包也会同时安装 `gc` 和 `gitcode`；Linux 上二者功能相同。
- Windows PowerShell 预置 `gc` 作为 `Get-Content` 别名；如果 `gc version` 被解析为读取文件，请改用 `gitcode version`、`gc.exe version` 或 `python -m gc_cli version`。
- Windows PowerShell 中通过 `--body-file -` / `--comment-file -` 管道传入中文或其他非 ASCII 正文时，推荐使用 UTF-8 文件；如果必须直接管道，先设置 `$OutputEncoding = [System.Text.UTF8Encoding]::new($false)`。CLI 会拦截疑似已被 PowerShell 损坏成 `???` 的输入并提示正确用法。

```powershell
Set-Content -Path body.md -Value "中文正文" -Encoding UTF8
gitcode issue create -R owner/repo --title "标题" --body-file body.md
```

### PyPI（备选）

```bash
# 创建虚拟环境
python3 -m venv .venv
source .venv/bin/activate  # Linux/macOS
# .venv\Scripts\activate   # Windows

# 安装
pip install gitcode-cli

# Windows PowerShell 中推荐使用 gitcode
gitcode version
```

### Linux 二进制文件

从 Release Assets 直接下载 Linux 二进制文件：

| 平台 | 文件 |
|------|------|
| Linux x64 | `gc_linux_amd64` |
| Linux ARM64 | `gc_linux_arm64` |

下载地址: https://gitcode.com/gitcode-cli/cli/releases

下载后赋予可执行权限，并放到 PATH 目录：

```bash
chmod +x gc_linux_amd64
mv gc_linux_amd64 ~/.local/bin/gc
gc version
```

Windows 和 macOS 用户建议使用上方 wheel 包；wheel 内置 Linux、macOS 和 Windows 二进制，并同时提供 `gc` 与 `gitcode` 两个命令入口。

### Docker 镜像

仓库已提供 `Dockerfile`、`docker-compose.yml` 和 Makefile 目标：

```bash
# 构建并运行
make docker-build
make docker-run

# 或使用 docker compose
docker compose up gc
```

认证 Token 通过环境变量传入。请在交互终端中静默读取，避免 Token 值进入 shell history：

```bash
read -rsp "GitCode token: " GC_TOKEN
export GC_TOKEN
make docker-run
unset GC_TOKEN
```

更多用法参见 Makefile 和 `docker-compose.yml`。

### 规划中的安装方式

以下安装方式正在开发中：

- [ ] Homebrew (macOS/Linux)
- [ ] Scoop (Windows)

## 快速开始

### 认证

以下示例使用安装包提供的 `gitcode` 入口；从源码构建或使用独立二进制时，请将命令名改为 `gc`。

**方式一：打开令牌页面并登录**

```bash
gitcode auth login --web
```

当前 `--web` 会打开个人令牌页面，生成令牌后仍需回到终端粘贴。当前版本不会隐藏输入，因此必须由用户本人在私有、未录制且不由 AI 控制的本地终端中执行。

**方式二：交互式 Token 登录**

```bash
gitcode auth login
```

浏览器不可用时，可在同样受控的本地交互终端中输入 Token。不要把 Token 值直接写进命令、shell history 或配置脚本。CI 场景应通过平台 Secret 注入 `GC_TOKEN` 或 `GITCODE_TOKEN`。

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
gitcode auth status
```

> 详细命令行为和完整示例请查看 [docs/COMMANDS.md](./docs/COMMANDS.md)。

## 输出格式

`gitcode` 的只读命令继续以文本输出为默认体验，同时为脚本和代理保留稳定的结构化入口。

```bash
# 结构化输出
gitcode issue list -R owner/repo --json
gitcode issue list -R owner/repo --format json
gitcode repo log -R owner/repo --file README.md --branch main --json
gitcode pr list -R owner/repo --paginate --per-page 100 --json

# 常规文本与表格
gitcode issue list -R owner/repo --format simple
gitcode issue list -R owner/repo --format table

# 时间格式切换
gitcode issue list -R owner/repo --time-format absolute
gitcode issue list -R owner/repo --time-format relative

# 自定义模板输出
gitcode issue list -R owner/repo --template '{{range .}}#{{.Number}} {{.Title}}{{"\n"}}{{end}}'

# typed command 尚未覆盖的 API，可用 gitcode api 读取原始响应
gitcode api repos/owner/repo
```

`issue view` 和 `pr view` 的文本详情展示也会保持稳定布局，而 `--json` 仍然是面向机器调用的首选入口。

## 常见任务入口

最常用的起步命令：

```bash
# 查看仓库
gitcode repo view

# 查看文件提交历史
gitcode repo log -R owner/repo --file README.md --branch main

# 创建 Issue
gitcode issue create --title "Bug report" --body "Description"

# 列出 Issues
gitcode issue list --state open

# 创建 PR
gitcode pr create --title "New feature" --base main

# 按提交信息反查 PR
gitcode pr list -R owner/repo --commit-message "fix login"

# 提交前检查 pre-commit 配置与本地环境
gitcode precommit check

# 查看流水线运行记录
gitcode actions run list -R owner/repo --status FAILED

# 查看流水线运行详情
gitcode actions run view <run-id> -R owner/repo

# 列出流水线运行的 jobs
gitcode actions job list <run-id> -R owner/repo

# 查看工作流 job 详情
gitcode actions job view <run-id> <job-id> -R owner/repo

# 下载 job 日志归档
gitcode actions job log <run-id> <job-id> -R owner/repo --output job-log.zip

# 列出仓库 artifacts
gitcode actions artifact list -R owner/repo

# 查看 artifact 详情
gitcode actions artifact view <artifact-id> -R owner/repo

# 下载 artifact
gitcode actions artifact download <artifact-id> -R owner/repo --output artifact.zip

# 删除 artifact
gitcode actions artifact delete <artifact-id> -R owner/repo --yes

# 调用 GitCode API 原始响应
gitcode api repos/owner/repo

# 查看认证状态
gitcode auth status
```

完整命令说明、参数细节、平台限制和更多示例，请直接查看：

- [docs/COMMANDS.md](./docs/COMMANDS.md)
- [docs/AUTH.md](./docs/AUTH.md)
- [docs/PACKAGING.md](./docs/PACKAGING.md)
- [docs/REGRESSION.md](./docs/REGRESSION.md)

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

补充说明：

- `docs/AI-GUIDE.md` 只服务外部项目通过 AI 使用 `gitcode`（或源码构建的 `gc`）
- gitcode-cli 仓库内部 AI 开发请看 `AGENTS.md`、`CLAUDE.md` 和 `spec/workflows/ai-local-development-workflow.md`
- `issues-plan/PROGRESS.md` 只作为阶段说明，不作为单个 issue / PR 的实时事实依据

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
