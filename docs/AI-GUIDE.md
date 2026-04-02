# 使用 AI 操作 GitCode 指南

本指南帮助你通过 AI 助手（如 Claude Code）操作 GitCode 平台。

## 1. 安装 GitCode CLI

**Linux (DEB/RPM):**

```bash
# DEB (Debian/Ubuntu)
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.8/gc_0.3.8_amd64.deb
sudo dpkg -i gc_0.3.8_amd64.deb

# RPM (RHEL/CentOS/Fedora)
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.8/gc-0.3.8-1.x86_64.rpm
sudo rpm -i gc-0.3.8-1.x86_64.rpm
```

**Wheel 包（跨平台，推荐）:**

从 Release 归档下载 wheel 包安装：

```bash
# 创建虚拟环境
python3 -m venv .venv
source .venv/bin/activate  # Linux/macOS
# .venv\Scripts\activate   # Windows

# 安装（一行命令）
pip install https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.8/gitcode_cli-0.3.8-py3-none-any.whl
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
gc issue list -R owner/repo --json
gc pr view 123 -R owner/repo --json

# 探索命令结构优先使用 schema
gc schema
gc schema "issue view"

# 高风险删除命令先 dry-run，再决定是否执行
gc repo delete owner/repo --dry-run
gc release delete v1.0.0 -R owner/repo --dry-run
```

说明：

- 删除类命令在非交互环境中不会再隐式等待输入；如果未显式传 `--yes`，会直接失败。
- 当前默认文本输出仍保留；代理和脚本应优先使用 `--json`。
- 当前基础退出码语义：`0` 成功，`2` 参数/用法错误，`3` 资源不存在，`4` 认证/权限错误。

## 6. 面向 AI 的开发流程约束

如果 AI 在本仓库或类似规范化仓库内参与开发，不应只把流程理解成一份 checklist，而应遵守“状态机 + 证据门禁”。

最低要求：

- 未进入 `status/verified` 的 issue，不得开始写代码
- 未创建非 `main` 分支前，不得修改实现文件
- 未完成测试、构建和必要的真实命令验证前，不得宣称“已完成”
- 作者自检不是独立评审
- PR 未合入主干前，不得把 issue 视为已完成主干合入

推荐状态流：

```text
Issue: status/triage -> status/verified -> status/in-progress -> status/ready-for-review -> status/merged
PR:    status/draft -> status/self-checked -> status/ready-for-review -> status/approved -> status/merged
```

## 7. 推荐给 AI 的固定模板

### Issue 验证记录

```markdown
## 验证记录

- 当前版本或分支:
- 复现命令:
- 实际结果:
- 结论:
```

### Issue 开发进度

```markdown
## 开发进度

- 根因:
- 主要修改:
- 测试:
- 实际命令验证:
- 风险或未覆盖项:
- 关联 PR:
```

### PR 作者自检

```markdown
## 作者自检

- 根因或实现理由:
- 主要修改:
- 单元测试:
- 构建:
- 实际命令验证:
- 文档同步:
- 风险:
- 未覆盖项:
```

### PR 评审结论

```markdown
## 评审结论

- 发现:
- blocker:
- 结论:
```

如果你希望直接复用完整模板而不是从本文复制，请查看：

- [AI-TEMPLATES.md](./AI-TEMPLATES.md)

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

**最后更新**: 2026-04-02
