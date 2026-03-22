# 里程碑 1: 基础架构

## 概述

建立 gitcode-cli 的基础架构和核心模块，为后续功能开发奠定基础。

**预计工期**: 1 周

**目标**: 完成项目骨架、基础模块和开发环境搭建

---

## 任务清单

### INFRA-001: 项目初始化

**优先级**: P0

**任务描述**:

- 初始化 Go 模块
- 创建目录结构
- 配置依赖管理

**验收标准**:

- [ ] `go mod init github.com/gitcode-com/gitcode-cli`
- [ ] 目录结构符合设计规范
- [ ] 依赖配置完整

**命令**:

```bash
mkdir -p gitcode-cli/{cmd/gc,pkg/cmd,pkg/cmdutil,pkg/iostreams,internal,api,git,context}
cd gitcode-cli
go mod init github.com/gitcode-com/gitcode-cli
```

---

### INFRA-002: Root 命令实现

**优先级**: P0

**任务描述**:

- 创建根命令框架
- 集成 Cobra
- 实现版本命令

**文件**:

```
cmd/gc/main.go
pkg/cmd/root/root.go
pkg/cmd/version/version.go
```

**验收标准**:

- [ ] `gc` 命令可执行
- [ ] `gc version` 显示版本信息
- [ ] `gc help` 显示帮助信息

**示例输出**:

```bash
$ gc version
gc version 0.1.0-dev (2026-03-22)
https://gitcode.com/gitcode-cli

$ gc help
gc is a command line tool for GitCode.

Usage:
  gc [command]

Available Commands:
  help        Help about any command
  version     Print gc version

Use "gc [command] --help" for more information about a command.
```

---

### INFRA-003: Factory 模式实现

**优先级**: P0

**任务描述**:

- 实现 Factory 接口
- 定义依赖注入模式
- 创建 Mock Factory

**文件**:

```
pkg/cmdutil/factory.go
pkg/cmdutil/factory_mock.go
```

**设计**:

```go
// pkg/cmdutil/factory.go
type Factory struct {
    IOStreams  *iostreams.IOStreams
    HttpClient func() (*http.Client, error)
    Config     func() (gc.Config, error)
    BaseRepo   func() (gtrepo.Interface, error)
    Branch     func() (string, error)
    Prompter   prompter.Prompter
    Browser    browser.Browser
}
```

**验收标准**:

- [ ] Factory 接口定义完整
- [ ] 支持 Mock 测试
- [ ] 依赖可配置

---

### INFRA-004: IOStreams 模块

**优先级**: P0

**任务描述**:

- 实现输入输出流管理
- 支持颜色输出
- 支持分页器

**文件**:

```
pkg/iostreams/iostreams.go
pkg/iostreams/color.go
```

**功能**:

- 标准输入/输出/错误流管理
- 颜色主题支持
- 分页器集成
- 终端检测

**验收标准**:

- [ ] 支持 In/Out/Err 流
- [ ] 支持颜色输出
- [ ] 支持 NO_COLOR 环境变量
- [ ] 支持分页器

---

### INFRA-005: 配置基础结构

**优先级**: P0

**任务描述**:

- 定义配置接口
- 实现配置文件解析
- 支持环境变量

**文件**:

```
internal/config/config.go
internal/config/config_file.go
internal/config/hosts_config.go
```

**验收标准**:

- [ ] 支持 YAML 配置
- [ ] 支持 `~/.config/gc/` 目录
- [ ] 支持环境变量覆盖
- [ ] 自动创建配置目录

---

### INFRA-006: Git 操作封装

**优先级**: P0

**任务描述**:

- 封装 Git 命令
- 获取仓库信息
- 分支操作

**文件**:

```
git/git.go
git/branch.go
git/remote.go
git/inspect.go
```

**验收标准**:

- [ ] 获取当前分支名
- [ ] 获取远程仓库信息
- [ ] 获取仓库根目录
- [ ] 支持 Git 操作

---

### INFRA-007: Makefile 和 CI/CD

**优先级**: P1

**任务描述**:

- 创建 Makefile
- 配置 GitHub Actions
- 设置测试流程

**文件**:

```
Makefile
.github/workflows/ci.yml
```

**验收标准**:

- [ ] `make build` 构建成功
- [ ] `make test` 运行测试
- [ ] CI 自动运行

---

## 依赖关系

```
INFRA-001 (项目初始化)
    ↓
INFRA-003 (Factory) ← INFRA-004 (IOStreams)
    ↓
INFRA-002 (Root命令)
    ↓
INFRA-005 (配置) + INFRA-006 (Git)
    ↓
INFRA-007 (CI/CD)
```

---

## 完成标准

里程碑 M1 完成需满足：

1. ✅ 项目可编译运行
2. ✅ `gc version` 命令可用
3. ✅ 基础模块测试通过
4. ✅ CI/CD 流程配置完成

---

## 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| Go 版本兼容性 | 中 | 明确 Go 1.21+ 要求 |
| 跨平台问题 | 低 | 优先开发 Linux/macOS |
| 依赖冲突 | 低 | 使用 go.mod 管理 |

---

## 相关文档

- [issues-plan/02-architecture.md](../02-architecture.md)
- [CLAUDE.md](../../CLAUDE.md)

---

**最后更新**: 2026-03-22