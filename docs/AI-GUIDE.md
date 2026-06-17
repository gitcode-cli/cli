# 使用 AI 操作 GitCode 指南

本指南帮助你通过 AI 助手（如 Claude Code）操作 GitCode 平台。

边界说明：

- 本文档只适用于“外部项目通过 AI 使用 `gc` 操作 GitCode”
- 本文档不定义 gitcode-cli 仓库内部开发流程
- 在 gitcode-cli 仓库内部参与开发时，应以 `AGENTS.md`、`CLAUDE.md`、`spec/README.md` 和 `spec/workflows/*` 为准

## 1. 安装 GitCode CLI

**Linux (DEB/RPM):**

```bash
# DEB (Debian/Ubuntu)
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.5.9/gc_0.5.9_amd64.deb
sudo dpkg -i gc_0.5.9_amd64.deb

# RPM (RHEL/CentOS/Fedora)
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.5.9/gc-0.5.9-1.x86_64.rpm
sudo rpm -i gc-0.5.9-1.x86_64.rpm
```

DEB/RPM packages install both `gc` and `gitcode`; on Linux they are equivalent.

**Wheel 包（跨平台，推荐）:**

从 Release 归档下载 wheel 包安装，**内置全平台二进制**（Linux x64/ARM、macOS Intel/Apple Silicon、Windows x64）：

```bash
# 创建虚拟环境
python3 -m venv .venv
source .venv/bin/activate  # Linux/macOS
# .venv\Scripts\activate   # Windows

# 安装（一行命令）
pip install https://gitcode.com/gitcode-cli/cli/releases/download/v0.5.9/gitcode_cli-0.5.9-py3-none-any.whl

# Windows PowerShell 中推荐使用 gitcode，避免 gc 被内置 Get-Content 别名覆盖
gitcode version
```

说明：
- wheel 会同时安装 `gc` 和 `gitcode` 两个命令入口，功能相同。
- DEB/RPM 包也会同时安装 `gc` 和 `gitcode`；Linux 上二者功能相同。
- Windows PowerShell 预置 `gc` 作为 `Get-Content` 别名；如果 `gc version` 被解析为读取文件，请改用 `gitcode version`、`gc.exe version` 或 `python -m gc_cli version`。
- Windows PowerShell 中让 AI 直接执行命令时，优先使用 `gitcode`。中文或其他非 ASCII 正文需要传给 `--body-file -` / `--comment-file -` 时，优先写入 UTF-8 临时文件再传文件路径；若必须直接管道，先设置 `$OutputEncoding = [System.Text.UTF8Encoding]::new($false)`。

**PyPI（备选）:**

> ⚠️ **注意**: PyPI 官方源可能有同步延迟，推荐使用上方 wheel 包下载

```bash
# 创建虚拟环境
python3 -m venv .venv
source .venv/bin/activate  # Linux/macOS
# .venv\Scripts\activate   # Windows

# 使用官方 PyPI 源安装
pip install -i https://pypi.org/simple/ gitcode-cli

# Windows PowerShell 中推荐使用 gitcode
gitcode version
```

**从源码构建:**

```bash
git clone https://gitcode.com/gitcode-cli/cli.git
cd cli
go build -o gc ./cmd/gc
```

## 2. 认证配置

```bash
# 设置 Token 环境变量
export GC_TOKEN=your_gitcode_token

# 添加到 shell 配置文件（永久生效）
echo 'export GC_TOKEN="your_gitcode_token"' >> ~/.bashrc
source ~/.bashrc
```

**获取 Token：**
1. 登录 [GitCode](https://gitcode.com)
2. 进入 设置 -> 私人令牌
3. 生成新令牌并复制

## 3. 验证安装

```bash
gc version
gc auth status
```

## 4. 安装 gc-core Skill

外部项目推荐使用 `gc-core` 通用 skill 包，而不是仓库内部协作 skill。

详细安装与分发说明见：

- [gc-core 安装与分发说明](../.ai/distribution/gc-core/INSTALL.md)

常见安装方式：

```bash
# Claude
mkdir -p ~/.claude/skills/gc-pr
cp .ai/distribution/gc-core/pr/SKILL.md ~/.claude/skills/gc-pr/SKILL.md

# Codex
mkdir -p ~/.codex/skills/gc-pr
cp .ai/distribution/gc-core/pr/SKILL.md ~/.codex/skills/gc-pr/SKILL.md
```

你也可以按同样方式安装 `gc-auth`、`gc-issue`、`gc-review` 等其他通用 skill。

安装后，AI 就可以通过 `gc` 命令操作 GitCode。

## 5. 面向 AI 的使用建议

为了让 AI 和脚本更稳定地消费 `gc`，优先使用以下模式：

```bash
# 读取类命令优先使用 JSON
gc repo view owner/repo --json
gc repo log -R owner/repo --file README.md --branch main --json
gc issue list -R owner/repo --json
gc pr view 123 -R owner/repo --json
gc pr list -R owner/repo --paginate --per-page 100 --json
gc pr list -R owner/repo --commit-message "fix login" --json
gc pr comments 123 -R owner/repo --json
gc issue prs 123 -R owner/repo --json
gc repo stats -R owner/repo --branch main --json
gc milestone view 1 -R owner/repo --json
gc commit comments list -R owner/repo --json
gc commit comments list-by-sha <sha> -R owner/repo --json
gc commit comments view <id> -R owner/repo --json

# 探索命令结构优先使用 schema
gc schema
gc schema "issue view"

# typed command 尚未覆盖时，可使用 gc api 读取 GitCode API 原始响应
gc api repos/owner/repo
gc api 'repos/owner/repo/commits?path=README.md&sha=main'

# 高风险删除命令先 dry-run，再决定是否执行
gc repo delete owner/repo --dry-run
gc release delete v1.0.0 -R owner/repo --dry-run

# 高频写路径需要解析结果时使用 JSON
gc issue create -R owner/repo --title "Bug" --body "..." --json
gc pr create -R owner/repo --head feature-branch --title "Feature" --body "..." --json
gc issue edit 123 -R owner/repo --title "Updated" --json
gc pr merge 123 -R owner/repo --yes --json
gc repo fork owner/repo --json
gc release create v1.0.0 -R owner/repo --title "v1.0.0" --notes "..." --json
gc release upload v1.0.0 app.zip -R owner/repo --json
```

说明：

- 删除、关闭、重开、状态切换、合并、同步推送/建 PR 等高风险写操作在非交互环境中不会再隐式等待输入；如果未显式传 `--yes`，会直接失败。
- 当前默认文本输出仍保留；代理和脚本应优先使用 `--json`。
- `repo log --json` 适合按文件和分支追踪提交历史；`pr list --paginate` 适合跨页扫描，`--commit-message` 适合从提交信息反查 PR。
- `gc api` 输出远端原始响应，适合 typed command 尚未覆盖的 API；使用含 `&` 的查询参数时建议整体加引号。
- 写路径 `--json` 只在操作成功后输出结构化结果；执行失败时不要从 stdout 解析半成品结果。
- `pr create --json` 会尽量回读新建 PR 以补齐创建响应缺失的正文；如果远端仍未返回 body，会在 stderr 给 warning，并保持 JSON 中的远端事实为空，脚本可再运行 `gitcode pr view <number> -R owner/repo --json` 核验。
- 当前基础退出码语义：`0` 成功，`1` 通用错误，`2` 参数/用法错误，`3` 资源不存在，`4` 认证/权限错误，`5` 资源冲突。

## 6. 使用 RTK 优化 Token 消耗（可选）

[RTK（Rust Token Killer）](https://github.com/rtk-ai/rtk) 是一个轻量级 CLI 代理工具，可在 CLI 命令的输出到达 LLM 上下文之前进行智能过滤与压缩，减少 60-90% 的 Token 消耗。

### 安装 RTK

```bash
# 从 GitHub 安装 RTK
cargo install rtk
# 或下载预编译二进制: https://github.com/rtk-ai/rtk/releases
```

### 配置 gc 过滤器

```bash
# 复制参考配置
mkdir -p ~/.config/rtk
cp contrib/rtk/config.toml ~/.config/rtk/config.toml

# 编辑配置以自定义过滤规则
$EDITOR ~/.config/rtk/config.toml
```

### 在 AI 工具中启用

```bash
# Claude Code
rtk init -g

# Hook 会自动拦截 gc 命令输出:
# gc pr list -R owner/repo → rtk gc pr list -R owner/repo
```

### 配置示例

参考配置文件位于 `contrib/rtk/config.toml`，预置了以下 gc 命令的过滤规则：

| 命令 | 策略 | 效果 |
|------|------|------|
| `gc pr list` / `gc issue list` | table-compact | 仅保留 number/state/title 列，限制 20 行 |
| `gc pr view` / `gc issue view` | strip-ansi-and-condense | 剥离颜色，仅保留标题/状态/正文 |
| `gc repo list` | table-compact | 仅保留 name/visibility/description |
| `gc release list` | table-compact | 仅保留 tag/title/date |
| `gc auth status` / `gc version` | one-line | 单行输出 |

### 效果对比

```
# 标准输出（~120 tokens）
$ gc pr list -R owner/repo
Showing 15 of 15 pull requests in owner/repo (filtered)
#1  open    Fix login bug                    bugfix/login ...
#2  merged  Add dark mode support           feature/dark-mode ...
...

# RTK 压缩后（~40 tokens，节省 ~67%）
#1 open Fix login bug
#2 merged Add dark mode support
...
```

### 注意事项

- RTK 是可选的外部工具，不影响 `gc` 的默认行为
- RTK 未安装时，`gc` 命令完全不受影响
- 错误输出默认透传，不会被过滤
- [参考配置文件](../contrib/rtk/config.toml) 可根据需要自定义

## 7. 在规范化仓库中的协作提醒

如果目标仓库本身已经定义了开发规范，AI 应直接遵守目标仓库自己的正式规则、状态机和证据门禁，而不是套用本文档。

## 8. 可参考的固定模板

外部项目如需固定模板，应由目标项目自己定义。

如果你只是参考 gitcode-cli 仓库的内部模板结构，请查看：

- [AI-TEMPLATES.md](./AI-TEMPLATES.md)

但这些模板不应被默认视为外部项目的正式流程模板。

## 完成后的使用方式

安装完成后，直接告诉 AI 你想做什么：

```
查看 owner/repo 仓库的所有 Issue
创建一个 PR，标题是"新增功能"
发布 v1.0.0 版本
```

AI 会自动使用 `gc` 命令执行操作。

## 更多信息

- [命令详细文档](./COMMANDS.md)
- [GitCode CLI 仓库](https://gitcode.com/gitcode-cli/cli)

---

**最后更新**: 2026-06-17
