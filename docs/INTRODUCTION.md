# GitCode CLI：把 GitCode 变成开发者和 AI 都能直接调用的工程能力

当团队每天都在 GitCode 上查看 Issue、评审 PR、追踪流水线、整理发布时，真正消耗时间的往往不是某一次点击，而是重复查找、复制信息、切换页面，以及把同一套操作重新做一遍。

GitCode CLI 把这些工作带回终端：仓库、Issue、Pull Request、Commit、标签、里程碑、Release 和 Actions 都可以通过统一命令完成。对开发者，它减少上下文切换；对团队，它让操作可以复用和审计；对 AI 代理，它提供结构化、可发现、带安全边界的 GitCode 执行入口。本文统一使用跨平台入口 `gitcode`；通过 wheel、PyPI、DEB 或 RPM 安装时也会提供等价的 `gc` 入口。

- 项目仓库：[gitcode-cli/cli](https://gitcode.com/gitcode-cli/cli)
- 下载最新版本：[Releases](https://gitcode.com/gitcode-cli/cli/releases)
- 完整命令手册：[docs/COMMANDS.md](https://gitcode.com/gitcode-cli/cli/blob/main/docs/COMMANDS.md)

## 为什么值得使用

### 少离开终端，多完成一次完整交付

从发现问题到合并发布，常见动作都可以留在当前工作区中完成：

```bash
# 看仓库和待处理事项
gitcode repo view owner/repo
gitcode issue list -R owner/repo --state open
gitcode pr list -R owner/repo --state open

# 创建 Issue 和 PR
gitcode issue create -R owner/repo --title "修复登录超时" --body-file issue.md --json
gitcode pr create -R owner/repo --base main --title "fix: 修复登录超时" --body-file pr.md --json

# 评审与发布
gitcode pr diff 42 -R owner/repo
gitcode pr review 42 -R owner/repo --comment-file review.md
gitcode release create v1.0.0 -R owner/repo --title "v1.0.0" --notes-file CHANGELOG.md --json
```

命令参数、输出字段和平台限制以[完整命令手册](https://gitcode.com/gitcode-cli/cli/blob/main/docs/COMMANDS.md)为准。

### 自动化不必依赖脆弱的页面脚本

高频只读命令和主要写操作支持 `--json`。脚本可以解析稳定的标准输出，错误则通过标准错误和明确退出码返回。需要探索能力时，`gitcode schema` 可以直接给出命令树、参数和元数据；尚未封装成专用命令的接口，还可以通过 `gitcode api` 调用。

```bash
# 机器可读结果
gitcode issue list -R owner/repo --paginate --per-page 100 --json
gitcode pr view 42 -R owner/repo --json
gitcode actions run list -R owner/repo --status FAILED --json

# 让脚本或 AI 发现命令，而不是猜参数
gitcode schema
gitcode schema "pr create"

# 专用命令尚未覆盖时读取原始 API 响应
gitcode api repos/owner/repo
```

### AI 不只是“告诉你怎么做”，而是可以在边界内完成操作

网页适合人浏览，CLI 更适合 AI 执行。GitCode CLI 的结构化输出、命令元数据、非交互行为和确认机制，让 AI 可以完成“读取事实、分析、执行、核验”的闭环，同时避免因为等待交互输入而卡住。

例如，你可以直接对 AI 说：

> 查看 `owner/repo` 当前所有开放 PR，按风险排序，逐个总结改动和 CI 状态，只执行只读命令。

> 根据本地改动起草一个 Issue 和 PR 描述，先展示给我确认，再使用 GitCode CLI 提交。

> 找出 `main` 分支最近失败的 Actions 运行，定位失败 job，下载日志并给出根因判断。

面向 Codex、Claude 等 AI 客户端的可安装 skills 已独立维护在 [gitcode-cli/skills](https://gitcode.com/gitcode-cli/skills)。其中包括 Issue 创建与评审、PR 创建与评审、反馈修复、Release 发布、安全检查、流水线分析等端到端工作流。更完整的 AI 使用约定见[使用 AI 操作 GitCode 指南](https://gitcode.com/gitcode-cli/cli/blob/main/docs/AI-GUIDE.md)。

### 自动化有边界，危险动作不会悄悄发生

GitCode CLI 对删除等高风险操作提供 `--dry-run` 和确认保护。在非交互环境中，未明确确认的破坏性操作会直接失败，不会无限等待输入。对 AI 来说，这意味着“能执行”不等于“可自行授权”：只有用户明确批准后，AI 才应使用 `--yes` 跳过确认。

```bash
# 先预演，再由人确认是否执行
gitcode repo delete owner/repo --dry-run
gitcode release delete v1.0.0 -R owner/repo --dry-run
```

认证信息不应出现在聊天、Prompt、脚本参数、Issue 或 PR 正文中。登录必须由用户本人在私有、未录制且不由 AI 控制的本地终端完成，AI 只运行 `gitcode auth status` 确认认证是否可用。详细规则见[认证说明](https://gitcode.com/gitcode-cli/cli/blob/main/docs/AUTH.md)。

## 适合哪些场景

| 使用者 | 典型任务 | GitCode CLI 带来的价值 |
| --- | --- | --- |
| 日常开发者 | 查 Issue、创建 PR、查看 diff、回复评审 | 减少页面切换，让工作流留在代码旁边 |
| 项目维护者 | Issue 分诊、标签与里程碑、批量审查、版本发布 | 形成一致、可复用的项目治理动作 |
| 测试与发布人员 | 追踪变更、核对 Release、下载发布资产 | 用命令和 JSON 构建可重复的发布检查 |
| CI/CD 运维人员 | 查看 Actions run/job、下载日志和 Artifact、检查 Runner | 更快定位流水线和运行环境问题 |
| AI 编码代理 | 获取远端事实、提交 Issue/PR、评审、核验结果 | 获得可发现、结构化、受约束的执行接口 |
| 企业自动化平台 | 跨仓库统计、流水线巡检、标准化交付 | 以统一 CLI 代替零散 API 脚本和页面自动化 |

仓库中已经整理了可直接复用的真实场景，包括 Issue 到 PR 的完整链路、发布评审、CI 流水线定位、安全检查和 AI 全流程交付，见 [GitCode CLI 应用案例库](https://gitcode.com/gitcode-cli/cli/tree/main/Example)。

## 五分钟开始使用

### 1. 安装

GitCode CLI 提供跨平台 wheel、DEB、RPM 和独立二进制等安装方式。为避免版本号过期，请直接从[项目 README 的安装章节](https://gitcode.com/gitcode-cli/cli#安装)或 [Releases](https://gitcode.com/gitcode-cli/cli/releases)选择最新版。

安装后先确认命令可用：

```bash
gitcode version
```

Windows PowerShell 已将 `gc` 用作 `Get-Content` 的别名，因此推荐使用 `gitcode`。通过 wheel、PyPI、DEB 或 RPM 安装时会同时提供 `gitcode` 和 `gc`；从源码构建或使用独立二进制时，默认产物通常只有 `gc`。

### 2. 在私有终端登录并确认认证

```bash
gitcode auth login --web
gitcode auth status
```

`--web` 会打开 GitCode 的新建访问令牌页面，生成后仍需回到终端粘贴；当前版本不会隐藏输入。请只在用户本人控制的私有、未录制本地终端中执行。浏览器不可用时可运行 `gitcode auth login`，但同样需要遵守这一限制。不要把 Token 交给 AI，也不要把 Token 直接写进命令、脚本或聊天内容。认证来源、优先级和安全注意事项见[认证说明](https://gitcode.com/gitcode-cli/cli/blob/main/docs/AUTH.md)。

### 3. 先从只读命令开始

```bash
gitcode repo view gitcode-cli/cli
gitcode issue list -R gitcode-cli/cli --state open
gitcode pr list -R gitcode-cli/cli --state open --json
```

`-R owner/repo` 可以让命令在任意目录中操作指定仓库；进入本地 Git 仓库后，多数命令也可以自动识别当前仓库。

### 4. 完成一次真实协作

先把正文写入 Markdown 文件，便于审阅、复用，也能减少 shell 转义和中文编码问题：

```bash
gitcode issue create -R owner/repo --title "问题标题" --body-file issue.md --json
gitcode pr create -R owner/repo --base main --title "feat: 功能标题" --body-file pr.md --json
```

创建前应遵守目标仓库自己的 Issue、分支、测试和 PR 规范。CLI 负责可靠执行平台操作，不替代项目本身的工程规则。

### 5. 按任务找到下一条命令

```bash
gitcode help --json
gitcode schema
gitcode help issue create
gitcode help pr review
gitcode help actions run list
```

常用入口：

- 所有命令和参数：[命令手册](https://gitcode.com/gitcode-cli/cli/blob/main/docs/COMMANDS.md)
- 登录与 Token 安全：[认证说明](https://gitcode.com/gitcode-cli/cli/blob/main/docs/AUTH.md)
- AI 操作建议：[AI 使用指南](https://gitcode.com/gitcode-cli/cli/blob/main/docs/AI-GUIDE.md)
- 安装、构建与平台说明：[项目 README](https://gitcode.com/gitcode-cli/cli)
- 可复制的业务场景：[应用案例库](https://gitcode.com/gitcode-cli/cli/tree/main/Example)
- AI 工作流 skills：[gitcode-cli/skills](https://gitcode.com/gitcode-cli/skills)
- 问题反馈与功能建议：[Issues](https://gitcode.com/gitcode-cli/cli/issues)

## 从今天的一件小事开始

不必先改造整套研发流程。可以先用 `gitcode pr list --json` 做一次待评审 PR 汇总，用 `gitcode actions run list` 找一次失败流水线，或者让 AI 在只读模式下完成一次 Issue/PR 分析。

当一个操作可以被命令准确表达，它就可以被保存、复用、审计，也可以安全地交给自动化和 AI。GitCode CLI 的价值，正是把 GitCode 上分散的协作动作，变成团队可以持续积累的工程能力。
