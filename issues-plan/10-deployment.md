# 部署发布需求

本文档详细描述 gitcode-cli 的构建、发布和部署流程。

## 发布策略

### 版本命名规范

遵循语义化版本规范 (Semantic Versioning)：

```
vMAJOR.MINOR.PATCH[-PRERELEASE]

示例：
v1.0.0        # 正式版本
v1.0.0-beta.1 # 预发布版本
v1.0.1        # 补丁版本
v1.1.0        # 次要版本
v2.0.0        # 主要版本
```

### 发布周期

| 版本类型 | 发布频率 | 说明 |
|----------|----------|------|
| PATCH | 按需 | Bug 修复 |
| MINOR | 每月 | 新功能 |
| MAJOR | 按需 | 重大更新 |

---

## DEPLOY-001: 构建流程

### 功能描述

定义多平台构建流程，生成可执行文件。

### 构建目标

| 平台 | 架构 | 输出文件 |
|------|------|----------|
| Linux | amd64 | gc-linux-amd64 |
| Linux | arm64 | gc-linux-arm64 |
| macOS | amd64 | gc-darwin-amd64 |
| macOS | arm64 | gc-darwin-arm64 |
| Windows | amd64 | gc-windows-amd64.exe |
| Windows | arm64 | gc-windows-arm64.exe |

### Makefile 定义

```makefile
# Makefile

VERSION ?= $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(shell git rev-parse HEAD) -X main.date=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)"

.PHONY: all build test clean

all: build

build:
	go build $(LDFLAGS) -o bin/gc ./cmd/gc

build-all:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/gc-linux-amd64 ./cmd/gc
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/gc-linux-arm64 ./cmd/gc
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/gc-darwin-amd64 ./cmd/gc
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/gc-darwin-arm64 ./cmd/gc
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/gc-windows-amd64.exe ./cmd/gc
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o dist/gc-windows-arm64.exe ./cmd/gc

test:
	go test -v -race -coverprofile=coverage.out ./...

clean:
	rm -rf bin/ dist/ coverage.out

install: build
	cp bin/gc /usr/local/bin/gc

.PHONY: release
release: clean build-all
	cd dist && sha256sum * > sha256sums.txt
```

### 构建信息嵌入

```go
// cmd/gc/main.go
package main

var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)

func main() {
    // 构建信息可通过 gc version 查看
}
```

### 验收标准

- [ ] 支持跨平台构建
- [ ] 版本信息正确嵌入
- [ ] 构建产物可执行
- [ ] 支持 Makefile 构建

---

## DEPLOY-002: GitHub Actions 发布

### 功能描述

使用 GitHub Actions 自动化构建和发布流程。

### 发布工作流

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Run tests
        run: go test -v -race ./...

      - name: Build
        run: make build-all

      - name: Create checksums
        run: cd dist && sha256sum * > sha256sums.txt

      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: dist/*
          generate_release_notes: true
```

### 持续集成工作流

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go: ['1.21', '1.22']
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}
      - name: Build
        run: go build ./...
      - name: Test
        run: go test -v -race ./...

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - uses: golangci/golangci-lint-action@v3
```

### 验收标准

- [ ] Tag 触发自动发布
- [ ] 多平台构建产物
- [ ] 自动生成 Release Notes
- [ ] 包含校验和文件

---

## DEPLOY-003: 包管理器分发

### 功能描述

通过各平台包管理器分发 gitcode-cli。

### Homebrew (macOS/Linux)

```ruby
# Formula/gc.rb
class Gc < Formula
  desc "GitCode CLI - Command line tool for GitCode"
  homepage "https://gitcode.com/gitcode-cli"
  version "1.0.0"

  on_macos do
    on_intel do
      url "https://gitcode.com/gitcode-cli/cli/releases/download/v#{version}/gc-darwin-amd64"
      sha256 ""
    end
    on_arm do
      url "https://gitcode.com/gitcode-cli/cli/releases/download/v#{version}/gc-darwin-arm64"
      sha256 ""
    end
  end

  on_linux do
    on_intel do
      url "https://gitcode.com/gitcode-cli/cli/releases/download/v#{version}/gc-linux-amd64"
      sha256 ""
    end
    on_arm do
      url "https://gitcode.com/gitcode-cli/cli/releases/download/v#{version}/gc-linux-arm64"
      sha256 ""
    end
  end

  def install
    bin.install "gc"
  end

  test do
    assert_match "gc version", shell_output("#{bin}/gc version")
  end
end
```

### 安装方式

```bash
# Homebrew
brew install gitcode-com/tap/gc

# 或通过自定义 Tap
brew tap gitcode-com/tap
brew install gc

# Shell 脚本安装
curl -sSL https://gitcode.com/install.sh | sh

# 手动安装
# Linux/macOS
curl -L https://gitcode.com/gitcode-cli/cli/releases/latest/download/gc-$(uname -s)-$(uname -m) -o gc
chmod +x gc
sudo mv gc /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri https://gitcode.com/gitcode-cli/cli/releases/latest/download/gc-windows-amd64.exe -OutFile gc.exe
```

### 安装脚本

```bash
#!/bin/bash
# install.sh

set -e

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
VERSION=${1:-"latest"}

# 架构映射
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
esac

echo "Installing gc for $OS/$ARCH..."

# 下载
if [ "$VERSION" = "latest" ]; then
    URL="https://gitcode.com/gitcode-cli/cli/releases/latest/download/gc-$OS-$ARCH"
else
    URL="https://gitcode.com/gitcode-cli/cli/releases/download/$VERSION/gc-$OS-$ARCH"
fi

curl -L "$URL" -o /tmp/gc
chmod +x /tmp/gc

# 安装
if [ -w /usr/local/bin ]; then
    mv /tmp/gc /usr/local/bin/gc
else
    sudo mv /tmp/gc /usr/local/bin/gc
fi

echo "gc installed successfully!"
gc version
```

### 验收标准

- [ ] 支持 Homebrew 安装
- [ ] 支持脚本安装
- [ ] 支持手动安装
- [ ] 安装后可直接运行

---

## DEPLOY-004: 版本检查和更新

### 功能描述

支持自动检查更新和提示用户升级。

### 版本命令

```bash
gc version
# gc version 1.0.0 (2026-03-22)
# https://gitcode.com/gitcode-cli

gc version --check
# A new version of gc is available: 1.1.0
# Upgrade with: brew upgrade gc

gc upgrade
# Upgrading gc...
# Successfully upgraded to v1.1.0
```

### 实现设计

```go
// pkg/cmd/version/version.go

type VersionInfo struct {
    Version    string
    Commit     string
    Date       string
    Latest     string
    UpdateAvailable bool
}

func CheckForUpdate(currentVersion string) (*VersionInfo, error) {
    // 检查 GitHub API 获取最新版本
    // 比较版本号
    // 返回更新信息
}
```

### 更新检查频率

- 每天最多检查一次
- 结果缓存到 `state.yml`
- 可通过配置禁用

### 验收标准

- [ ] `gc version` 显示版本信息
- [ ] `gc version --check` 检查更新
- [ ] 显示升级提示
- [ ] 支持禁用更新检查

---

## 发布检查清单

发布新版本前必须检查：

### 代码质量

- [ ] 所有测试通过
- [ ] 代码覆盖率达标
- [ ] Lint 检查通过
- [ ] 无安全漏洞

### 文档更新

- [ ] CHANGELOG.md 更新
- [ ] README.md 更新（如有必要）
- [ ] 文档站点更新（如有）

### 版本管理

- [ ] 版本号符合语义化规范
- [ ] Git Tag 正确创建
- [ ] Release Notes 完整

### 构建验证

- [ ] 所有平台构建成功
- [ ] 可执行文件测试通过
- [ ] 安装脚本测试通过

---

## CHANGELOG 规范

```markdown
# CHANGELOG

All notable changes to this project will be documented in this file.

## [Unreleased]

## [1.0.0] - 2026-03-22

### Added
- Initial release
- `gc auth login` - OAuth Device Flow authentication
- `gc auth status` - View authentication status
- `gc repo clone` - Clone repositories
- `gc repo create` - Create repositories
- `gc issue create/list/view` - Issue management
- `gc pr create/list/view/review` - PR management

### Changed
-

### Fixed
-

### Security
-
```

---

## 相关文档

- [gc-design/docs/deployment/release.md](https://gitcode.com/afly-infra/gc-design/blob/main/docs/deployment/release.md)
- [gc-design/docs/deployment/ci-cd.md](https://gitcode.com/afly-infra/gc-design/blob/main/docs/deployment/ci-cd.md)

---

**最后更新**: 2026-03-22