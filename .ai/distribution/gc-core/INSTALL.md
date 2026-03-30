# gc-core 安装与分发说明

本文档说明如何在不依赖 gitcode-cli 仓库内部上下文的情况下，安装和分发 `gc-core` 通用 skill 包。

## 1. 适用场景

适用于以下场景：

- 在其他项目中使用 `gc` 操作 GitCode
- 将 `gc-core` 复制到本地 AI skill 目录
- 将 `gc-core` 作为团队内部可复用 skill 包分发

不适用于以下场景：

- 继续开发 gitcode-cli 仓库本身
- 使用 gitcode-cli 仓库内部的 `spec/`、`docs/`、`.claude/skills/`、`.codex/skills/`

## 2. 前置条件

安装前请确认：

- `gc` 已安装并可执行
- 已完成 GitCode 认证，或可以在使用时完成认证
- 目标 AI 客户端支持本地 skill 目录

最小验证命令：

```bash
gc version
gc auth status
```

## 3. 复制方式

`gc-core` 是一个目录级 skill 包。最简单的安装方式是复制整个目录：

```bash
cp -R .ai/distribution/gc-core /path/to/target/
```

也可以只复制单个 skill 子目录，例如：

```bash
cp -R .ai/distribution/gc-core/review /path/to/target/
cp -R .ai/distribution/gc-core/pr /path/to/target/
```

## 4. Claude 推荐安装方式

如果目标环境使用 Claude 本地 skill 目录，可按目录逐个安装：

```bash
mkdir -p ~/.claude/skills/gc-auth
cp .ai/distribution/gc-core/auth/SKILL.md ~/.claude/skills/gc-auth/SKILL.md

mkdir -p ~/.claude/skills/gc-pr
cp .ai/distribution/gc-core/pr/SKILL.md ~/.claude/skills/gc-pr/SKILL.md
```

按同样方式安装：

- `gc-repo`
- `gc-issue`
- `gc-release`
- `gc-review`
- `gc-regression`

## 5. Codex 推荐安装方式

如果目标环境使用 Codex 本地 skill 目录，可按目录逐个安装：

```bash
mkdir -p ~/.codex/skills/gc-auth
cp .ai/distribution/gc-core/auth/SKILL.md ~/.codex/skills/gc-auth/SKILL.md

mkdir -p ~/.codex/skills/gc-pr
cp .ai/distribution/gc-core/pr/SKILL.md ~/.codex/skills/gc-pr/SKILL.md
```

按同样方式安装：

- `gc-repo`
- `gc-issue`
- `gc-release`
- `gc-review`
- `gc-regression`

## 6. 团队分发建议

如果要在团队内部复用，推荐：

1. 保留 `gc-core/` 目录原样
2. 在团队自己的仓库中引入该目录
3. 由每位成员复制所需 skill 到本地客户端目录

这样做的好处是：

- skill 可以被版本控制
- skill 更新可以被 code review
- 不依赖 gitcode-cli 仓库内部私有路径

## 7. 更新方式

当 `gc-core` 更新后，建议按以下方式同步：

1. 重新复制更新后的 `SKILL.md`
2. 替换本地客户端目录中的旧版本
3. 重新验证关键命令是否符合当前 `gc` 行为

最小验证建议：

```bash
gc auth status
gc repo view owner/repo
gc issue list -R owner/repo --state open
gc pr list -R owner/repo --state open
```

## 8. 边界提醒

`gc-core` 是通用 skill 包，不负责：

- 目标项目自己的开发流程
- 目标项目自己的 PR 审查规则
- 目标项目自己的 CI 规范

这些内容应由目标项目自己定义。
