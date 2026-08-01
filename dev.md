# GitCode CLI 开发模式指南

> 本文档总结项目的开发模式、规范体系和 AI 自动化流程。

---

## 1. 项目架构概览

| 项目属性 | 值 |
|---------|-----|
| 项目名 | gitcode-cli |
| 命令名 | `gc` |
| 语言 | Go 1.22+ |
| 框架 | Cobra |
| 配置目录 | `~/.config/gc/` |
| 环境变量前缀 | `GC_*` |
| 完成状态 | 100% (56/56 任务) |

---

## 2. Spec 规范体系

项目建立了完整的规范文档体系，位于 `spec/` 目录。

### 目录结构

```
spec/
├── README.md                # 文档索引
├── workflows/               # 操作流程
│   └── development-workflow.md
├── foundations/             # 基础规范
│   ├── coding-standards.md
│   ├── testing-guide.md
│   ├── command-template.md
│   └── security.md
└── workflows/
    ├── issue-workflow.md    # Issue 操作流程
    ├── pr-workflow.md       # PR 操作流程
    ├── review-workflow.md   # 评审流程
    └── test-workflow.md     # 测试流程
```

### 核心规范要点

| 规范类别 | 要求 |
|---------|------|
| 命名规范 | 包名小写简短，导出名称大驼峰 |
| 提交限制 | 单次提交不超过 800 行 |
| 测试覆盖率 | 新功能 >= 70%，核心模块 >= 80% |
| 安全要求 | Token 禁止硬编码，仅使用环境变量 |

---

## 3. Skills 体系（AI 自动化）

GitCode CLI Skills 已迁移到独立仓库：

- Web：[gitcode-cli/skills](https://gitcode.com/gitcode-cli/skills)
- SSH：`git@gitcode.com:gitcode-cli/skills.git`

该仓库提供核心命令 Skill 和端到端工作流 Skill，包括认证、仓库、Issue、PR、Review、
Release、回归测试、安全检查和流水线分析等场景。完整清单和安装方式以其 `README.md` 为准。

本仓库不再追踪 `.ai/skills/`、`.claude/skills/` 或 `.codex/skills/` 副本。项目内部开发规则
仍以 `spec/` 为准，Skill 只提供执行方法。

```bash
# 克隆到 CLI 工作树之外，避免把外部仓库误加入当前提交
git clone git@gitcode.com:gitcode-cli/skills.git ../gitcode-cli-skills
```

---

## 4. CLI 命令开发流程

### 完整开发流程

```
提交 Issue → 打标签 → 创建分支 → 分支开发 → 编写测试 →
本地测试 → 实际命令测试 → 安全审查 → 提交 PR →
Issue 评论 → PR 审查评论 → 关闭 Issue → 合并 PR
```

### 分支命名规范

| 类型 | 格式 | 示例 |
|------|------|------|
| BUG 修复 | `bugfix/issue-<number>` | `bugfix/issue-33` |
| 新特性 | `feature/issue-<number>` | `feature/issue-23` |
| 文档更新 | `docs/issue-<number>` | `docs/issue-5` |

### 命令目录结构

```
pkg/cmd/<category>/
├── <category>.go          # Parent command
├── <action>/
│   ├── <action>.go        # Subcommand implementation
│   └── <action>_test.go   # Subcommand tests
```

### 命令开发模板

```go
// Package <action> implements the <category> <action> command
package <action>

import (
    "fmt"
    "net/http"
    "os"

    "github.com/MakeNowJust/heredoc/v2"
    "github.com/spf13/cobra"

    "gitcode.com/gitcode-cli/cli/api"
    cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
    "gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type <Action>Options struct {
    IO         *iostreams.IOStreams
    HttpClient func() (*http.Client, error)
    Repository string
}

func NewCmd<Action>(f *cmdutil.Factory, runF func(*<Action>Options) error) *cobra.Command {
    opts := &<Action>Options{
        IO:         f.IOStreams,
        HttpClient: f.HttpClient,
    }

    cmd := &cobra.Command{
        Use:   "<action> [<number>]",
        Short: "<Short description>",
        Args:  cobra.MaximumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            if runF != nil {
                return runF(opts)
            }
            return <action>Run(opts)
        },
    }

    cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
    return cmd
}

func <action>Run(opts *<Action>Options) error {
    // Implementation
    return nil
}
```

---

## 5. 测试流程与限制

### 测试仓库限制（重要！）

**允许使用的测试仓库**:
- `infra-test/gctest1` ← **首选测试仓库**
- `infra-test` 组织下的其他仓库

**禁止行为**:
- 使用个人仓库测试
- 使用其他组织或用户的仓库测试
- 使用 `gitcode-cli/cli` 测试

### 测试检查清单

```bash
# 1. 构建
go build -o ./gc ./cmd/gc

# 2. 单元测试
go test ./...

# 3. 实际命令测试（使用 infra-test/gctest1）
./gc issue list -R infra-test/gctest1 --state open
./gc issue view 1 -R infra-test/gctest1
./gc pr list -R infra-test/gctest1 --state all
./gc repo view infra-test/gctest1
```

---

## 6. CI/CD 配置

### GitHub Actions 工作流

| Workflow | 文件 | 功能 |
|----------|------|------|
| CI | `.github/workflows/ci.yml` | Lint + Test + Build + Docker |
| Release | `.github/workflows/release.yml` | GoReleaser + RPM/DEB + PyPI |

### CI 流程

```
Push/PR → Lint → Test (多平台矩阵) → Build → Docker
```

### Release 流程

```
推送 v* 标签 → GitHub Actions 触发 → GoReleaser 构建 → 上传资产 → PyPI 发布
```

---

## 7. 项目管理

### issues-plan 目录

```
issues-plan/
├── README.md                    # 需求管理说明
├── 01-requirements-overview.md  # 需求总清单
├── PROGRESS.md                  # 进度跟踪表
└── milestones/                  # 里程碑追踪
    ├── m1-foundation.md
    ├── m2-auth.md
    ├── m3-repo.md
    ├── m4-issue.md
    ├── m5-pr.md
    └── m6-release.md
```

### 里程碑状态

| 里程碑 | 状态 | 进度 |
|--------|------|------|
| M1 基础架构 | ✅ 完成 | 7/7 |
| M2 认证功能 | ✅ 完成 | 8/8 |
| M3 仓库功能 | ✅ 完成 | 6/6 |
| M4 Issue功能 | ✅ 完成 | 8/8 |
| M5 PR功能 | ✅ 完成 | 9/9 |
| M6 Release功能 | ✅ 完成 | 6/6 |
| M7 文档与基础设施 | ✅ 完成 | 5/5 |

---

## 8. Agent Teams 模式

### 协作架构

```
Team Lead (orchestrator)
    ├── Issue-Reviewer Agent    → 评审 Issue、打标签
    ├── PR-Reviewer Agent       → 代码审查、安全检查
    ├── Developer Agent         → 编写代码、运行测试
    └── Tester Agent            → 执行实际命令测试
```

### 自动化触发路径

```
用户请求 → 读取 AGENTS.md / CLAUDE.md 与相关 spec
    → gitcode-issue-review → Issue 分析
    → gitcode-pr-create → PR 创建
    → gitcode-pr-review → 工程评审
    → gitcode-security-check → 安全检查
```

---

## 9. 全流程 AI 自动化模式

### 核心原则

**先验证，后修复。禁止看到 Issue 就直接写代码。**

### Issue 修复验证流程

```
1. 用当前版本验证问题是否存在
2. 检查时间线（Issue 创建时间、代码提交时间）
3. 判断是否需要修复
   → 问题不存在：评论说明并关闭
   → 问题存在：继续修复流程
```

### 自动化开发流程

```
用户提出需求 → 创建 Issue → 按 spec 完成验证
    ↓
打标签 → 创建分支 → 按命令模板开发与测试
    ↓
本地开发 → 运行测试 → 使用 gitcode-pr-review 辅助独立评审
    ↓
提交 PR → 合并 → 调用打包脚本发布 Release
```

---

## 10. 关键文件路径索引

| 类别 | 文件路径 |
|------|---------|
| AI 开发入口 | `CLAUDE.md` |
| 规范索引 | `spec/README.md` |
| 开发流程 | `spec/workflows/development-workflow.md` |
| 命令模板 | `spec/foundations/command-template.md` |
| GitCode CLI Skills | `https://gitcode.com/gitcode-cli/skills` |
| 需求总清单 | `issues-plan/01-requirements-overview.md` |
| 进度跟踪 | `issues-plan/PROGRESS.md` |
| CI 配置 | `.github/workflows/ci.yml` |
| Release 配置 | `.github/workflows/release.yml` |
| 打包指南 | `docs/PACKAGING.md` |
| 命令参考 | `docs/COMMANDS.md` |

---

**最后更新**: 2026-03-28
