# 使用 AI 操作 GitCode 指南

本指南帮助你通过 AI 助手（如 Claude Code）操作 GitCode 平台。

## 1. 安装 GitCode CLI

**Linux (DEB/RPM):**

```bash
# DEB (Debian/Ubuntu)
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.5/gc_0.3.5_amd64.deb
sudo dpkg -i gc_0.3.5_amd64.deb

# RPM (RHEL/CentOS/Fedora)
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.5/gc-0.3.5-1.x86_64.rpm
sudo rpm -i gc-0.3.5-1.x86_64.rpm
```

**Wheel 包（跨平台，推荐）:**

从 Release 归档下载 wheel 包安装：

```bash
# 创建虚拟环境
python3 -m venv .venv
source .venv/bin/activate  # Linux/macOS
# .venv\Scripts\activate   # Windows

# 安装（一行命令）
pip install https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.5/gitcode_cli-0.3.5-py3-none-any.whl
```

**PyPI（备选）:**

> ⚠️ **注意**: PyPI 官方源可能有同步延迟，推荐使用上方 wheel 包下载

```bash
# 创建虚拟环境
python3 -m venv .venv
source .venv/bin/activate  # Linux/macOS
# .venv\Scripts\activate   # Windows

# 使用官方 PyPI 源安装
pip install -i https://pypi.org/simple/ gitcode-cli
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

## 4. 安装 gitcode-cli Skill

将本项目的 skill 文件安装到你的 Claude Code 配置目录：

```bash
# 复制 skill 文件到 Claude Code 配置目录
mkdir -p ~/.claude/skills/gitcode-cli
cp .claude/skills/gitcode-cli/SKILL.md ~/.claude/skills/gitcode-cli/
```

安装后，AI 会自动识别并使用 `gc` 命令操作 GitCode。

> **Skill 文件位置**: [.claude/skills/gitcode-cli/SKILL.md](../.claude/skills/gitcode-cli/SKILL.md)

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

**最后更新**: 2026-03-25