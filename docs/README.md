# docs 文档入口

`docs/` 是 gitcode-cli 的用户文档目录。

如果你是：

- 使用者：从本目录开始
- 开发者或 AI 协作者：先看 [spec/README.md](../spec/README.md)

## 权威边界

本目录的职责是：

- 解释用户如何使用 `gc`
- 提供命令行为说明
- 提供认证、回归、打包等使用文档

本目录不是项目规则源。

项目正式规则以 [spec/README.md](../spec/README.md) 和 `spec/` 目录为准。

其中：

- [COMMANDS.md](./COMMANDS.md) 是命令行为唯一真相源

## 推荐阅读顺序

1. [COMMANDS.md](./COMMANDS.md)
2. [AUTH.md](./AUTH.md)
3. [REGRESSION.md](./REGRESSION.md)
4. [PACKAGING.md](./PACKAGING.md)
5. [AI-GUIDE.md](./AI-GUIDE.md)
6. [LOOP_ENGINEERING.md](./LOOP_ENGINEERING.md)
7. [LOOP_ENGINEERING_DEMO.md](./LOOP_ENGINEERING_DEMO.md)
8. [AI-TEMPLATES.md](./AI-TEMPLATES.md)
9. `docs/ai-templates/*.md`
10. [应用案例库](../Example/index.md)

说明：

- [AI-GUIDE.md](./AI-GUIDE.md) 只服务“外部项目用 AI 操作 GitCode”
- 本仓库内部开发流程请看 `AGENTS.md`、`CLAUDE.md` 和 `spec/workflows/*`
- [AI-TEMPLATES.md](./AI-TEMPLATES.md) 提供 gitcode-cli 仓库内部协作模板，不是项目规则源

## 当前包含内容

| 文档 | 说明 |
|------|------|
| [COMMANDS.md](./COMMANDS.md) | 命令行为说明和示例 |
| [AUTH.md](./AUTH.md) | 认证来源和行为说明 |
| [REGRESSION.md](./REGRESSION.md) | 核心回归矩阵说明 |
| [PACKAGING.md](./PACKAGING.md) | 本地打包与发布使用说明 |
| [AI-GUIDE.md](./AI-GUIDE.md) | 外部项目使用 AI 操作 GitCode 的场景指南 |
| [LOOP_ENGINEERING.md](./LOOP_ENGINEERING.md) | gitcode-cli Loop Engineering Demo v1 说明 |
| [LOOP_ENGINEERING_DEMO.md](./LOOP_ENGINEERING_DEMO.md) | Loop Engineering 分支验证演示记录 |
| [AI-TEMPLATES.md](./AI-TEMPLATES.md) | gitcode-cli 仓库内部协作的固定模板 |
| `docs/ai-templates/*.md` | 可直接复用的 AI 评论与检查模板文件 |
| [Example/](../Example/index.md) | GitCode CLI 在业务场景中的应用案例和可复制 prompt |
