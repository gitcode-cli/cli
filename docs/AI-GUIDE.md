# 使用 AI 操作 GitCode 指南

本指南帮助你通过 AI 助手（如 Claude Code）操作 GitCode 平台，实现仓库管理、Issue 处理、PR 操作等常见任务。

## 目录

- [快速开始](#快速开始)
- [安装 GitCode CLI](#安装-gitcode-cli)
- [配置 AI 助手](#配置-ai-助手)
- [常用操作示例](#常用操作示例)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)

## 快速开始

### 1. 安装 GitCode CLI

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

### 2. 认证配置

```bash
# 登录 GitCode
gc auth login

# 或设置环境变量
export GC_TOKEN=your_gitcode_token
```

### 3. 验证安装

```bash
gc version
gc auth status
```

## 配置 AI 助手

### Claude Code 配置

在项目根目录创建 `CLAUDE.md` 文件，添加以下内容：

```markdown
# GitCode CLI 配置

## 命令工具
使用 `gc` 命令操作 GitCode 仓库，而不是 `gh`。

## 常用命令参考
- Issue 管理: `gc issue create/list/view/close/comment`
- PR 管理: `gc pr create/list/view/merge/checkout/review`
- Release 管理: `gc release create/list/view/upload`
- 仓库操作: `gc repo clone/create/fork/view`

## 认证
Token 通过环境变量 `GC_TOKEN` 设置，或使用 `gc auth login` 登录。

## 测试仓库
只使用指定的测试仓库进行测试，避免影响生产环境。
```

### 其他 AI 助手

对于其他 AI 助手（如 Cursor、GitHub Copilot 等），可以在系统提示或自定义指令中添加类似说明：

```
使用 gc 命令操作 GitCode 仓库：
- 创建 Issue: gc issue create --title "标题" --body "内容" -R owner/repo
- 查看 PR: gc pr view <number> -R owner/repo
- 创建 Release: gc release create <tag> --title "标题" -R owner/repo
```

## 常用操作示例

### Issue 管理

#### 创建 Issue

告诉 AI：

```
创建一个 Issue，标题是"修复登录页面样式问题"，内容描述问题现象和复现步骤
```

AI 会执行：

```bash
gc issue create --title "修复登录页面样式问题" --body "## 问题描述\n..." -R owner/repo
```

#### 查看 Issue 列表

```
查看项目中所有未关闭的 Issue
```

```bash
gc issue list -R owner/repo --state open
```

#### 为 Issue 添加标签

```
给 Issue #5 添加 bug 标签
```

```bash
gc issue label 5 --add bug -R owner/repo
```

#### 关闭 Issue

```
关闭 Issue #5
```

```bash
gc issue close 5 -R owner/repo
```

#### 添加 Issue 评论

```
在 Issue #5 中评论说明已修复
```

```bash
gc issue comment 5 --body "问题已在 commit abc123 中修复" -R owner/repo
```

### Pull Request 管理

#### 创建 PR

```
创建一个 PR，标题是"feat: 添加用户认证功能"，描述改动内容
```

```bash
gc pr create --title "feat: 添加用户认证功能" --body "## 改动内容\n- 新增登录接口\n- 添加权限验证" --base main -R owner/repo
```

#### 查看 PR 列表

```
查看所有待审查的 PR
```

```bash
gc pr list -R owner/repo --state open
```

#### 查看 PR 详情

```
查看 PR #10 的详细信息
```

```bash
gc pr view 10 -R owner/repo
```

#### 查看 PR 代码变更

```
查看 PR #10 的代码改动
```

```bash
gc pr diff 10 -R owner/repo
```

#### 审查 PR

```
批准 PR #10
```

```bash
gc pr review 10 --approve -R owner/repo
```

```
在 PR #10 中提交修改建议
```

```bash
gc pr review 10 --comment "建议修改变量命名" -R owner/repo
```

#### 查看 PR 评论

```
查看 PR #10 的所有评论
```

```bash
gc pr comments 10 -R owner/repo
```

#### 合并 PR

```
合并 PR #10
```

```bash
gc pr merge 10 -R owner/repo
```

### Release 管理

#### 创建 Release

```
创建 v1.0.0 版本发布
```

```bash
gc release create v1.0.0 --title "v1.0.0" --notes "## 更新内容\n- 首次发布" -R owner/repo
```

#### 查看 Release 列表

```
查看所有已发布的版本
```

```bash
gc release list -R owner/repo
```

#### 上传发布资源

```
上传构建产物到 v1.0.0 版本
```

```bash
gc release upload v1.0.0 ./dist/app_1.0.0_amd64.deb -R owner/repo
```

### 仓库操作

#### 克隆仓库

```
克隆 gitcode-cli/cli 仓库
```

```bash
gc repo clone gitcode-cli/cli
```

#### 查看仓库信息

```
查看仓库 owner/repo 的详细信息
```

```bash
gc repo view owner/repo
```

#### 创建仓库

```
创建一个名为 my-project 的新仓库
```

```bash
gc repo create my-project --public
```

## 最佳实践

### 1. 使用项目配置文件

在项目根目录创建 `CLAUDE.md`，让 AI 自动了解项目上下文：

```markdown
# 项目说明

## GitCode 配置
- 命令工具: gc
- 仓库: owner/repo
- 测试仓库: test-org/test-repo

## 开发规范
- 分支命名: feature/xxx, bugfix/xxx
- 提交信息: feat/fix/docs: 描述
- PR 必须关联 Issue

## 工作流程
1. 创建 Issue 并打标签
2. 创建开发分支
3. 开发完成后创建 PR
4. PR 审查通过后合并
```

### 2. 明确指定仓库

始终使用 `-R owner/repo` 参数明确指定操作目标仓库，避免操作错误的仓库。

### 3. 使用测试仓库验证

开发新功能或测试命令时，使用指定的测试仓库，不要在生产仓库测试。

### 4. 敏感信息保护

- Token 不要写入配置文件
- 使用环境变量传递 Token
- 不要在公开仓库中暴露 Token

### 5. 完整工作流程示例

告诉 AI 执行完整的开发流程：

```
实现一个新功能并创建 PR：
1. 创建 Issue 描述功能需求，打上 enhancement 标签
2. 创建开发分支 feature/issue-xxx
3. 编写代码和测试
4. 提交代码
5. 创建 PR 并关联 Issue
6. 在 Issue 中添加进度评论
```

## 常见问题

### Q: AI 使用了 `gh` 命令而不是 `gc`

在项目配置中明确说明：

```markdown
重要: 使用 `gc` 命令操作 GitCode，不要使用 `gh`（GitHub CLI）
```

### Q: 命令执行失败

检查：
1. `gc` 是否正确安装 (`gc version`)
2. 是否已认证 (`gc auth status`)
3. 仓库名称是否正确 (`owner/repo` 格式)
4. Token 是否有足够权限

### Q: 如何查看命令帮助

```bash
gc help
gc issue --help
gc pr create --help
```

### Q: API 返回 404 错误

可能原因：
1. 仓库不存在或无权限访问
2. PR/Issue 编号不存在
3. API 端点路径错误

检查仓库和资源是否存在：

```bash
gc repo view owner/repo
gc issue view <number> -R owner/repo
```

## 命令速查表

| 操作 | 命令 |
|------|------|
| 登录 | `gc auth login` |
| 查看状态 | `gc auth status` |
| 克隆仓库 | `gc repo clone owner/repo` |
| 查看仓库 | `gc repo view owner/repo` |
| 创建 Issue | `gc issue create -R owner/repo` |
| 列出 Issue | `gc issue list -R owner/repo` |
| 查看 Issue | `gc issue view <number> -R owner/repo` |
| 关闭 Issue | `gc issue close <number> -R owner/repo` |
| Issue 添加标签 | `gc issue label <number> --add <label> -R owner/repo` |
| 创建 PR | `gc pr create -R owner/repo` |
| 列出 PR | `gc pr list -R owner/repo` |
| 查看 PR | `gc pr view <number> -R owner/repo` |
| 查看 PR 评论 | `gc pr comments <number> -R owner/repo` |
| 审查 PR | `gc pr review <number> --approve -R owner/repo` |
| 合并 PR | `gc pr merge <number> -R owner/repo` |
| 创建 Release | `gc release create <tag> -R owner/repo` |
| 上传资源 | `gc release upload <tag> <file> -R owner/repo` |

## 更多资源

- [GitCode CLI 仓库](https://gitcode.com/gitcode-cli/cli)
- [命令详细文档](./COMMANDS.md)
- [API 文档](https://gitcode.com/gitcode-cli/cli/tree/main/docs)

---

**最后更新**: 2026-03-24