# 项目总览

## 项目简介

**gitcode-cli**（命令名：`gc`）是一款面向 GitCode 平台的命令行工具，参考 GitHub CLI (gh) 的优秀架构设计，为 GitCode 用户提供便捷的命令行操作体验。

### 设计目标

1. **用户体验优先** - 提供与 GitHub CLI 类似的用户体验和命令结构
2. **平台适配** - 完全适配 GitCode 平台的 API 特性
3. **跨平台支持** - 支持 Linux、macOS、Windows
4. **易于扩展** - 模块化设计，便于添加新功能
5. **安全可靠** - 安全的认证存储，完善的错误处理

### 第一版本功能范围

| 优先级 | 功能模块 | 描述 |
|--------|----------|------|
| P0 | 认证 (auth) | Token认证、OAuth Device Flow、多账户支持 |
| P0 | 仓库 (repo) | 克隆、创建、查看仓库 |
| P0 | Issue 管理 | 创建、列表、查看、关闭 Issue |
| P0 | Pull Request | 创建、列表、检出、合并 PR |
| P0 | 代码检视 | 审阅 PR，添加评论，批准或请求修改 |

## 技术栈

| 组件 | 技术选型 | 版本要求 | 选型理由 |
|------|----------|----------|----------|
| 编程语言 | Go | 1.21+ | 跨平台编译、单二进制部署、丰富的标准库 |
| 命令框架 | Cobra | v1.8+ | 成熟的CLI框架，支持子命令、参数解析、自动补全 |
| 配置格式 | YAML | - | 可读性强，支持注释，与 gh 保持一致 |
| HTTP 客户端 | net/http | - | 标准库稳定可靠 |
| 表格输出 | go-gh/tableprinter | v2 | 成熟的表格渲染库 |
| 交互提示 | survey | v2 | 丰富的交互组件 |
| 密钥存储 | go-keyring | v0.2+ | 跨平台安全存储 |

## 项目结构

```
gitcode-cli/
├── cmd/gc/                    # 主程序入口
│   └── main.go
├── internal/                  # 内部模块（不对外暴露）
│   ├── config/               # 配置管理
│   ├── gtcmd/                # CLI 主逻辑
│   ├── authflow/             # 认证流程
│   ├── gtrepo/               # 仓库模型
│   ├── gtinstance/           # GitCode 实例管理
│   ├── browser/              # 浏览器操作
│   ├── prompter/             # 交互式提示
│   ├── tableprinter/         # 表格输出
│   ├── text/                 # 文本处理
│   ├── keyring/              # 密钥环集成
│   └── update/               # 版本更新检查
├── pkg/
│   ├── cmd/                  # 命令实现
│   │   ├── root/            # 根命令
│   │   ├── auth/            # 认证命令
│   │   ├── repo/            # 仓库命令
│   │   ├── issue/           # Issue 命令
│   │   ├── pr/              # PR 命令
│   │   ├── user/            # 用户命令
│   │   ├── config/          # 配置命令
│   │   ├── api/             # API 命令
│   │   ├── extension/       # 扩展管理
│   │   └── factory/         # 工厂模式
│   ├── cmdutil/             # 命令工具函数
│   ├── iostreams/           # IO 流管理
│   └── http/                # HTTP 工具
├── api/                      # API 客户端
│   ├── client.go            # API 客户端
│   ├── http_client.go       # HTTP 客户端
│   ├── queries_repo.go      # 仓库查询
│   ├── queries_issue.go     # Issue 查询
│   ├── queries_pr.go        # PR 查询
│   └── query_builder.go     # GraphQL 查询构建
├── git/                      # Git 操作封装
├── context/                  # 上下文管理
├── docs/                     # 文档
├── scripts/                  # 构建脚本
├── issues-plan/              # 需求管理
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── CLAUDE.md                 # AI 辅助开发指南
```

## 命令概览

### 认证命令

```bash
gc auth login     # 登录 GitCode 账户
gc auth logout    # 登出账户
gc auth status    # 查看认证状态
gc auth token     # 打印认证 Token
gc auth switch    # 切换账户
```

### 仓库命令

```bash
gc repo clone     # 克隆仓库
gc repo create    # 创建仓库
gc repo fork      # Fork 仓库
gc repo view      # 查看仓库
gc repo list      # 列出仓库
gc repo delete    # 删除仓库
```

### Issue 命令

```bash
gc issue create   # 创建 Issue
gc issue list     # 列出 Issues
gc issue view     # 查看 Issue
gc issue close    # 关闭 Issue
gc issue reopen   # 重开 Issue
gc issue comment  # 添加评论
```

### PR 命令

```bash
gc pr create      # 创建 PR
gc pr list        # 列出 PRs
gc pr view        # 查看 PR
gc pr checkout    # 检出 PR 分支
gc pr merge       # 合并 PR
gc pr close       # 关闭 PR
gc pr review      # 代码检视（重点功能）
gc pr diff        # 查看 PR 差异
gc pr ready       # 标记为就绪/WIP
```

## 环境变量

| 环境变量 | 描述 | 示例 |
|----------|------|------|
| `GC_TOKEN` | 认证 Token | `gc_token_xxxx` |
| `GITCODE_TOKEN` | 备选 Token | `gc_token_xxxx` |
| `GC_HOST` | 默认主机 | `gitcode.com` |
| `GC_CONFIG_DIR` | 配置目录 | `/home/user/.config/gc` |

## 参考项目

本项目基于 GitHub CLI (gh) 的源码分析编写，主要参考：

- [cli/cli](https://github.com/cli/cli) - GitHub CLI 官方仓库
- 源码位置: `cli/` 目录

---

**最后更新**: 2026-03-22