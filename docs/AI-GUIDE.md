# 使用 AI 操作 GitCode 指南

本指南帮助你通过 AI 助手（如 Claude Code）操作 GitCode 平台。

## 1. 安装 GitCode CLI

**Linux (DEB/RPM):**

```bash
# DEB (Debian/Ubuntu)
wget https://gitcode.com/gitcode-cli/cli/releases/latest/download/gc_0.2.8_amd64.deb
sudo dpkg -i gc_0.2.8_amd64.deb

# RPM (RHEL/CentOS/Fedora)
wget https://gitcode.com/gitcode-cli/cli/releases/latest/download/gc-0.2.8-1.x86_64.rpm
sudo rpm -i gc-0.2.8-1.x86_64.rpm
```

**从源码构建:**

```bash
git clone https://gitcode.com/gitcode-cli/cli.git
cd cli
go build -o gc ./cmd/gc
```

## 2. 认证配置

```bash
# 方式一：交互式登录
gc auth login

# 方式二：设置环境变量
export GC_TOKEN=your_gitcode_token
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

在 Claude Code 中安装 gitcode-cli skill，让 AI 自动使用 `gc` 命令操作 GitCode：

**方式一：在项目中配置**

在项目根目录创建或编辑 `.claude/settings.json`：

```json
{
  "permissions": {
    "allow": [
      "Bash(gc *)"
    ]
  }
}
```

**方式二：使用 Skill 文件**

将以下内容添加到项目的 `CLAUDE.md` 文件中：

```markdown
## GitCode CLI

使用 `gc` 命令操作 GitCode 仓库，不要使用 `gh`（GitHub CLI）。

常用命令：
- `gc issue create/list/view/close/comment`
- `gc pr create/list/view/merge/review`
- `gc release create/list/upload`
- `gc repo clone/view/create`
```

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

**最后更新**: 2026-03-24