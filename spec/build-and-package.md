# 本地构建与打包规范

本文件定义 gitcode-cli 的本地构建、打包和产物校验规则。

本规范面向开发者和 AI 协作者，目标是让本地交付流程可重复、可验证、可审查。

## 1. 适用范围

本规范适用于以下场景：

- 本地构建 `gc` 二进制
- 本地构建多平台二进制
- 本地构建 DEB / RPM / PyPI 包
- 本地验证构建产物
- 提交前检查本地构建目录和产物边界

本规范不定义 CI 自动化规则。CI 相关内容在后续 `spec/ci-workflows.md` 中单独维护。

## 2. 权威边界

- 构建和打包的正式规范以本文件为准
- 用户如何使用打包产物，以 `docs/PACKAGING.md` 为准
- 当前仓库的脚本实现以 `Makefile` 和 `scripts/package.sh`、`scripts/build.sh`、`scripts/build-release.sh` 为准
- 若脚本行为与本规范不一致，应优先修正脚本或同步更新本规范，不允许长期漂移

## 3. 标准构建命令

### 3.1 本地开发构建

开发阶段的标准本地构建命令为：

```bash
go build -o ./gc ./cmd/gc
```

该命令用于：

- 日常开发验证
- 实际命令测试
- issue 修复后的快速回归

### 3.2 Makefile 构建

仓库已提供以下标准构建入口：

```bash
make build
make build-all
make clean
```

使用约束如下：

- `make build` 用于当前平台构建
- `make build-all` 用于多平台二进制构建
- `make clean` 用于清理 `bin/`、`dist/` 和覆盖率产物

### 3.3 快照式发布构建

需要验证发布产物时，可使用：

```bash
make release-local
make release-snapshot
```

这些命令用于本地快照验证，不等同于正式 release。

## 4. 标准打包方式

### 4.1 推荐脚本

本地打包优先使用：

```bash
./scripts/package.sh <version> [target]
```

例如：

```bash
./scripts/package.sh v0.3.4 release
./scripts/package.sh v0.3.4 linux
./scripts/package.sh v0.3.4 deb
./scripts/package.sh v0.3.4 rpm
./scripts/package.sh v0.3.4 pypi
```

原因：

- 该脚本已经承担版本同步、二进制构建和打包收口
- 比手工调用零散工具更接近当前仓库真实发布流程

### 4.2 允许的产物类型

当前仓库允许生成以下产物：

- 当前平台本地二进制：`./gc`
- 多平台二进制：`bin/` 或 `dist/` 中的对应产物
- DEB 包：`gc_*.deb`
- RPM 包：`gc-*.rpm`
- PyPI 包：`gitcode_cli-*.whl`、`gitcode_cli-*.tar.gz`

## 5. 构建前置条件

执行本地构建或打包前，必须确认：

- Go 环境可用
- 依赖已安装
- 需要 DEB / RPM 打包时，`nfpm` 已安装
- 需要 PyPI 打包时，Python build 工具已安装
- 需要真实命令验证时，已设置 `GC_TOKEN`

常用准备命令：

```bash
go mod tidy
go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest
pip install --upgrade build wheel setuptools
```

## 6. 本地构建验证

### 6.1 最小验证

任何涉及命令行为的改动，在提交前至少完成：

```bash
go build -o ./gc ./cmd/gc
./gc version
```

### 6.2 功能验证

如果改动影响具体命令，必须继续做：

- 相关包级单元测试
- `go test ./...`
- 至少一个真实命令验证

真实命令验证只能使用 `infra-test/*` 仓库。

### 6.3 打包验证

如果改动影响构建、打包、版本信息或发布产物，必须继续验证：

- 构建脚本可执行
- 目标产物实际生成
- 产物命名与文档描述一致
- 至少抽样验证一个安装或运行路径

## 7. 产物边界

构建和打包产物默认属于本地产物，不应直接提交到仓库。

常见本地产物包括：

- `gc`
- `bin/`
- `dist/`
- `build/`
- `*.deb`
- `*.rpm`
- `*.tar.gz`
- `*.zip`
- `*.egg-info/`

提交前必须确认：

- 本地产物未被误提交
- `.gitignore` 已覆盖新增产物类型
- 本地临时目录未进入版本控制

## 8. 文档同步要求

当构建、打包流程或产物命名发生变化时，必须同步检查：

- `docs/PACKAGING.md`
- `README.md`
- `spec/build-and-package.md`
- `spec/release-process.md`
- `AGENTS.md`
- `CLAUDE.md`
- 相关 AI skills

## 9. 禁止事项

以下行为不允许出现：

- 在未验证产物实际可用前更新发布说明
- 在主工作区堆积大量构建产物后直接提交
- 将本地测试凭证写入构建脚本或文档
- 在未同步文档的情况下修改打包产物命名
- 使用未纳入仓库规范的个人脚本替代标准流程

## 10. 当前执行基线

当前仓库的构建与打包基线如下：

1. 开发验证优先使用 `go build -o ./gc ./cmd/gc`
2. 多平台构建优先使用 `make build-all`
3. 本地打包优先使用 `scripts/package.sh`
4. 发布前的构建检查必须结合 `go test ./...` 和真实命令验证
5. 构建与打包规则变更后必须同步更新用户文档和 AI 协作文档
