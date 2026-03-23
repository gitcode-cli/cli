# 本地打包和发布指南

本文档说明如何在本地构建 DEB/RPM 包并使用 `gc` 命令发布 Release。

## 前置要求

### 安装 nfpm

[nfpm](https://github.com/goreleaser/nfpm) 是一个通用的包构建工具。

```bash
# 方式一：使用 go install
go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest

# 方式二：使用安装脚本
curl -sfL https://install.goreleaser.com/github.com/goreleaser/nfpm.sh | sh -s -- -b ~/bin
```

### 设置认证

```bash
# 设置 Token 环境变量
export GC_TOKEN="your_gitcode_token"

# 或添加到 shell 配置永久生效
echo 'export GC_TOKEN="your_gitcode_token"' >> ~/.bashrc
source ~/.bashrc
```

---

## 构建 DEB/RPM 包

### 1. 构建二进制文件

```bash
# 创建输出目录
mkdir -p dist

# 构建 Linux amd64
GOOS=linux GOARCH=amd64 go build -o dist/gc_linux_amd64 ./cmd/gc

# 构建 Linux arm64
GOOS=linux GOARCH=arm64 go build -o dist/gc_linux_arm64 ./cmd/gc
```

### 2. 配置 nfpm

项目已包含 nfpm 配置文件：

- `nfpm-amd64.yaml` - amd64 架构配置
- `nfpm-arm64.yaml` - arm64 架构配置

配置文件示例 (`nfpm-amd64.yaml`)：

```yaml
name: "gc"
arch: "amd64"
platform: "linux"
version: "0.2.0"
section: "default"
priority: "extra"
maintainer: "gitcode-cli contributors"
description: |
  GitCode CLI - Command line tool for GitCode
  Provides convenient access to GitCode features including:
  - Authentication management
  - Repository operations
  - Issue and PR management
vendor: "gitcode-cli"
homepage: "https://gitcode.com/gitcode-cli/cli"
license: "MIT"
contents:
  - src: ./dist/gc_linux_amd64
    dst: /usr/bin/gc
  - src: ./completions/gc.bash
    dst: /usr/share/bash-completion/completions/gc
  - src: ./completions/gc.zsh
    dst: /usr/share/zsh/vendor-completions/_gc
  - src: ./completions/gc.fish
    dst: /usr/share/fish/completions/gc.fish
```

### 3. 更新版本号

```bash
VERSION="0.2.0"

# 更新 nfpm 配置文件中的版本号
sed -i "s/version: .*/version: \"$VERSION\"/" nfpm-amd64.yaml
sed -i "s/version: .*/version: \"$VERSION\"/" nfpm-arm64.yaml
```

### 4. 构建 DEB 包

```bash
# 构建 amd64 DEB
nfpm package -f nfpm-amd64.yaml -p deb -t dist/

# 构建 arm64 DEB
nfpm package -f nfpm-arm64.yaml -p deb -t dist/
```

### 5. 构建 RPM 包

```bash
# 构建 amd64 RPM
nfpm package -f nfpm-amd64.yaml -p rpm -t dist/

# 构建 arm64 RPM
nfpm package -f nfpm-arm64.yaml -p rpm -t dist/
```

### 6. 查看构建结果

```bash
ls -la dist/*.deb dist/*.rpm
```

输出示例：
```
dist/gc_0.2.0_amd64.deb
dist/gc_0.2.0_arm64.deb
dist/gc-0.2.0-1.x86_64.rpm
dist/gc-0.2.0-1.aarch64.rpm
```

---

## 发布 Release

### 1. 创建 Release

```bash
gc release create v0.2.0 -R owner/repo \
  --title "gc v0.2.0" \
  --notes "Release notes here"
```

> **注意**：`--notes` 参数是必需的，不带此参数可能返回 400 错误。

### 2. 上传资产

```bash
# 上传单个文件
gc release upload v0.2.0 dist/gc_0.2.0_amd64.deb -R owner/repo

# 上传多个文件
gc release upload v0.2.0 \
  dist/gc_0.2.0_amd64.deb \
  dist/gc_0.2.0_arm64.deb \
  dist/gc-0.2.0-1.x86_64.rpm \
  dist/gc-0.2.0-1.aarch64.rpm \
  -R owner/repo
```

### 3. 查看 Release

```bash
gc release view v0.2.0 -R owner/repo
```

### 4. 列出所有 Releases

```bash
gc release list -R owner/repo
```

### 5. 下载资产

```bash
# 下载所有资产
gc release download v0.2.0 -R owner/repo

# 下载到指定目录
gc release download v0.2.0 -R owner/repo -o ./downloads/

# 下载指定文件
gc release download v0.2.0 gc_0.2.0_amd64.deb -R owner/repo
```

---

## 完整示例

以下是一个完整的打包和发布流程：

```bash
#!/bin/bash
set -e

# 配置
VERSION="0.2.0"
REPO="owner/repo"
TOKEN="your_token"

# 设置 Token
export GC_TOKEN="$TOKEN"

# 创建输出目录
mkdir -p dist

# 构建二进制文件
echo "Building binaries..."
GOOS=linux GOARCH=amd64 go build -o dist/gc_linux_amd64 ./cmd/gc
GOOS=linux GOARCH=arm64 go build -o dist/gc_linux_arm64 ./cmd/gc

# 更新版本号
echo "Updating version..."
sed -i "s/version: .*/version: \"$VERSION\"/" nfpm-amd64.yaml
sed -i "s/version: .*/version: \"$VERSION\"/" nfpm-arm64.yaml

# 构建 DEB 包
echo "Building DEB packages..."
nfpm package -f nfpm-amd64.yaml -p deb -t dist/
nfpm package -f nfpm-arm64.yaml -p deb -t dist/

# 构建 RPM 包
echo "Building RPM packages..."
nfpm package -f nfpm-amd64.yaml -p rpm -t dist/
nfpm package -f nfpm-arm64.yaml -p rpm -t dist/

# 创建 Release
echo "Creating release..."
gc release create v$VERSION -R $REPO \
  --title "gc v$VERSION" \
  --notes "GitCode CLI v$VERSION release"

# 上传资产
echo "Uploading assets..."
gc release upload v$VERSION \
  dist/gc_${VERSION}_amd64.deb \
  dist/gc_${VERSION}_arm64.deb \
  dist/gc-${VERSION}-1.x86_64.rpm \
  dist/gc-${VERSION}-1.aarch64.rpm \
  -R $REPO

echo "Done! Release v$VERSION published."
```

---

## 安装指南

### DEB (Debian/Ubuntu)

```bash
# 下载并安装
wget https://gitcode.com/owner/repo/releases/download/v0.2.0/gc_0.2.0_amd64.deb
sudo dpkg -i gc_0.2.0_amd64.deb

# 或使用 gc 命令下载
gc release download v0.2.0 gc_0.2.0_amd64.deb -R owner/repo
sudo dpkg -i gc_0.2.0_amd64.deb
```

### RPM (RHEL/CentOS/Fedora)

```bash
# 下载并安装
wget https://gitcode.com/owner/repo/releases/download/v0.2.0/gc-0.2.0-1.x86_64.rpm
sudo rpm -i gc-0.2.0-1.x86_64.rpm

# 或使用 gc 命令下载
gc release download v0.2.0 gc-0.2.0-1.x86_64.rpm -R owner/repo
sudo rpm -i gc-0.2.0-1.x86_64.rpm
```

### 验证安装

```bash
gc version
```

---

## 常见问题

### Q: 创建 Release 时返回 400 错误

确保包含 `--notes` 参数：
```bash
gc release create v1.0.0 -R owner/repo --title "v1.0.0" --notes "Release notes"
```

### Q: 上传失败

1. 检查 Token 是否有效：`gc auth status`
2. 检查仓库是否存在且有写入权限
3. 检查 Release 是否已创建

### Q: 版本号已存在

每个 tag 只能创建一个 Release，使用新的版本号：
```bash
gc release create v1.0.1 -R owner/repo --title "v1.0.1" --notes "..."
```

### Q: nfpm 找不到命令

确保 nfpm 在 PATH 中：
```bash
export PATH="$HOME/go/bin:$HOME/bin:$PATH"
# 或添加到 shell 配置
echo 'export PATH="$HOME/go/bin:$HOME/bin:$PATH"' >> ~/.bashrc
```

---

**最后更新**: 2026-03-23