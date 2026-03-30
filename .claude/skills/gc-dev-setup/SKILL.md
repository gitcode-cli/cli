---
name: gc-dev-setup
description: 初始化 gitcode-cli 项目本地开发环境。当用户说"初始化本地开发环境"、"搭建本地开发环境"、"init dev environment"、"setup local dev"时触发，或用户有开发需求时首先检查本地开发环境是否就绪。检查标准：本地代码编译无错误。
---

# GC 开发环境初始化

为 gitcode-cli (gc) 项目初始化本地开发环境。

## 触发条件

主动触发此 skill 的情况：
- 用户说"初始化本地开发环境"或"搭建本地开发环境"
- 用户说"init dev environment"、"setup local dev"、"初始化开发环境"
- 用户有开发需求（创建功能、修复 Bug）- 先验证环境是否就绪
- 用户刚拉取代码准备开始工作

## 环境就绪标准

本地开发环境被认为"就绪"的条件：
- Go 已安装且可用
- 项目构建成功（`./gc` 二进制文件存在且可运行）
- 无编译错误

## 工作流程

### 步骤 1：拉取最新代码

```bash
git pull origin main
# 或当前分支
git pull
```

如果不在 git 仓库中或需要先克隆，引导用户克隆项目。

### 步骤 2：检查 Go 环境

验证 Go 是否已安装：

```bash
go version
```

如果 Go 未安装：
1. Ubuntu/Debian: `sudo apt update && sudo apt install -y golang-go`
2. 或从 https://go.dev/dl/ 下载安装
3. 国内用户可能需要设置 GOPROXY: `export GOPROXY=https://goproxy.cn,direct`

### 步骤 3：构建项目

构建 gc 二进制文件：

```bash
# 国内用户设置 GOPROXY
export GOPROXY=https://goproxy.cn,direct

# 构建
go build -o ./gc ./cmd/gc
```

### 步骤 4：验证构建

运行版本命令验证：

```bash
./gc version
```

如果命令执行无错误，则环境已就绪。无需检查认证状态。

如果需要继续开发或验证核心链路，优先执行：

```bash
./scripts/regression-core.sh
```

说明：
- 该脚本会用临时 `GC_CONFIG_DIR` 验证 auth、repo view、issue list/view 和非 Git 目录错误路径。
- 默认不执行 `pr create` 这类写路径；需要时按 `docs/REGRESSION.md` 显式开启。

预期输出：
```
gc version dev
  commit: none
  built:  unknown
https://gitcode.com/gitcode-cli/cli
```

## 快速环境检查

快速验证环境是否就绪：

```bash
# 重新构建并验证
GOPROXY=https://goproxy.cn,direct go build -o ./gc ./cmd/gc && ./gc version && ./scripts/regression-core.sh
```

如果无错误，则环境已就绪；如果只是做最小环境探测，`./gc version` 通过即可。

## 常见问题

### Go 未找到
- 安装 Go: `sudo apt install -y golang-go` (Ubuntu)
- 或从 https://go.dev/dl/ 下载

### 构建失败（网络错误）
- 设置 GOPROXY: `export GOPROXY=https://goproxy.cn,direct`

### 权限被拒绝
- 确保不使用 `sudo` 执行 go build
- 检查文件权限

### PATH 中存在旧版本
- 不要将 `./gc` 复制到 `~/bin/` 或其他 PATH 目录
- 始终在项目目录中直接使用 `./gc`

## 输出格式

初始化完成后，提供简洁的摘要：

```
✓ 本地开发环境已就绪
- Go 版本: go1.22.x
- gc 二进制: ./gc
- 状态: 构建成功，命令正常
```

无需检查认证状态。只要 `./gc version` 执行无错误，环境即被认为就绪。
