# 本地打包和发布指南

> 项目概述和功能介绍请参阅 [README.md](../README.md)，命令使用请参阅 [COMMANDS.md](./COMMANDS.md)，开发指南请参阅 [CLAUDE.md](../CLAUDE.md)。

本文档说明如何在本地构建 DEB/RPM/PyPI 包并使用 `gc` 命令发布 Release。

---

## 目录

- [快速开始](#快速开始)
- [前置要求](#前置要求)
- [发布 Release](#发布-release)
- [安装指南](#安装指南)
- [常见问题](#常见问题)
- [附录：手动构建步骤](#附录手动构建步骤)

---

## 重要流程

**每次打包发布后，必须同步更新以下文档中的版本信息：**

| 文档 | 需要更新的内容 |
|------|---------------|
| `README.md` | Release badge、下载链接、版本号 |
| `docs/AI-GUIDE.md` | 安装命令中的版本号 |

> **注意**：使用 `package.sh` 脚本会自动同步版本号，无需手动更新。

---

## 快速开始

### 方式一：使用打包脚本（推荐）

使用 `scripts/package.sh` 一键完成版本同步、构建和打包：

```bash
# 发布用（构建 DEB + RPM + PyPI，推荐）
./scripts/package.sh v0.3.4 release

# 构建所有包（DEB + RPM + PyPI）
./scripts/package.sh v0.3.4

# 仅构建 Linux 包（DEB + RPM）
./scripts/package.sh v0.3.4 linux

# 仅构建 DEB 包
./scripts/package.sh v0.3.4 deb

# 仅构建 PyPI 包
./scripts/package.sh v0.3.4 pypi
```

### 构建目标

| 目标 | 说明 | 输出文件 |
|------|------|----------|
| `release` | DEB + RPM + PyPI（发布用，推荐） | 全部包 |
| `all` | 构建所有包（默认） | 全部包 |
| `linux` | DEB + RPM 包 | `gc_*.deb`, `gc-*.rpm` |
| `deb` | 仅 DEB 包 | `gc_*.deb` |
| `rpm` | 仅 RPM 包 | `gc-*.rpm` |
| `pypi` | 仅 PyPI 包 | `gitcode_cli-*.whl` |

### 脚本功能

`package.sh` 自动完成：

| 步骤 | 说明 |
|------|------|
| 1. 版本同步 | 自动更新 6 个配置文件 + 3 个文档 |
| 2. 构建 Linux 二进制 | amd64 + arm64 |
| 3. 构建 DEB 包 | amd64 + arm64 |
| 4. 构建 RPM 包 | x86_64 + aarch64 |
| 5. 构建 PyPI 包 | 多平台二进制 + wheel |

### 构建产物

```bash
dist/
├── gc_0.3.6_amd64.deb              # DEB amd64
├── gc_0.3.6_arm64.deb              # DEB arm64
├── gc-0.3.6-1.x86_64.rpm           # RPM x86_64
├── gc-0.3.6-1.aarch64.rpm          # RPM aarch64
├── gc_linux_amd64                  # Linux 二进制 amd64
├── gc_linux_arm64                  # Linux 二进制 arm64
├── gitcode_cli-0.3.6-py3-none-any.whl  # PyPI wheel
└── gitcode_cli-0.3.6.tar.gz        # PyPI sdist
```

---

## 前置要求

### 安装 nfpm

[nfpm](https://github.com/goreleaser/nfpm) 是 DEB/RPM 包构建工具。

```bash
go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest
```

> **注意**：`go install` 安装到 `~/go/bin/`，`package.sh` 会自动查找此路径。

### 安装 Python build 工具（PyPI 包需要）

```bash
pip install --upgrade build wheel setuptools
```

### 设置认证

```bash
export GC_TOKEN="your_gitcode_token"
```

---

## 发布 Release

### 完整发布流程

```bash
# 1. 构建所有包
./scripts/package.sh v0.3.4 release

# 2. 创建 Release
gc release create v0.3.4 -R gitcode-cli/cli \
  --title "GitCode CLI v0.3.4" \
  --notes "$(cat <<'EOF'
## 更新内容

### 新功能
- 功能描述

### Bug 修复
- 修复描述

### 修复的 Issue
- Fixes Issue XX

## 安装方式

### Wheel 包（推荐，跨平台）

创建虚拟环境并安装：

    python3 -m venv .venv
    source .venv/bin/activate
    pip install https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gitcode_cli-0.3.6-py3-none-any.whl

Windows 用户激活虚拟环境：

    .venv\Scripts\activate

### DEB (Debian/Ubuntu)

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gc_0.3.6_amd64.deb
    sudo dpkg -i gc_0.3.6_amd64.deb

ARM64 设备：

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gc_0.3.6_arm64.deb
    sudo dpkg -i gc_0.3.6_arm64.deb

### RPM (RHEL/CentOS/Fedora)

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gc-0.3.6-1.x86_64.rpm
    sudo rpm -i gc-0.3.6-1.x86_64.rpm

ARM64 设备：

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gc-0.3.6-1.aarch64.rpm
    sudo rpm -i gc-0.3.6-1.aarch64.rpm

### Linux 二进制

AMD64：

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gc_linux_amd64
    chmod +x gc_linux_amd64
    sudo mv gc_linux_amd64 /usr/local/bin/gc

ARM64：

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gc_linux_arm64
    chmod +x gc_linux_arm64
    sudo mv gc_linux_arm64 /usr/local/bin/gc

## 验证安装

    gc version
EOF
)"

# 3. 上传所有包
gc release upload v0.3.4 \
  dist/gc_linux_amd64 \
  dist/gc_linux_arm64 \
  dist/gc_0.3.6_amd64.deb \
  dist/gc_0.3.6_arm64.deb \
  dist/gc-0.3.6-1.x86_64.rpm \
  dist/gc-0.3.6-1.aarch64.rpm \
  dist/gitcode_cli-0.3.6-py3-none-any.whl \
  -R gitcode-cli/cli
```

> **注意**：将示例中的版本号 `0.3.4` 替换为实际版本号。

### Release Notes 要求

**重要**：所有下载链接必须使用完整路径：

```
https://gitcode.com/gitcode-cli/cli/releases/download/v{VERSION}/{FILENAME}
```

**禁止**只写 `pip install xxx.whl` 不提供下载地址！

**代码块格式警告**：

GitCode 会错误渲染代码块内的 `#` 开头行为标题！

错误格式（会导致格式混乱）：

    ```bash
    # 创建虚拟环境
    python3 -m venv .venv
    ```

正确格式（注释放在代码块外）：

    创建虚拟环境：

    ```bash
    python3 -m venv .venv
    ```

正确格式（使用其他注释符号）：

    ```bash
    :: 创建虚拟环境
    python3 -m venv .venv
    ```

#### Release Notes 模板

```markdown
## 更新内容

### 新功能
- 功能描述

### Bug 修复
- 修复描述

### 修复的 Issue
- Fixes Issue XX

## 安装方式

### Wheel 包（推荐，跨平台）

创建虚拟环境并安装：

    python3 -m venv .venv
    source .venv/bin/activate
    pip install https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gitcode_cli-0.3.6-py3-none-any.whl

Windows 用户激活虚拟环境：

    .venv\Scripts\activate

### DEB (Debian/Ubuntu)

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gc_0.3.6_amd64.deb
    sudo dpkg -i gc_0.3.6_amd64.deb

ARM64 设备：

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gc_0.3.6_arm64.deb
    sudo dpkg -i gc_0.3.6_arm64.deb

### RPM (RHEL/CentOS/Fedora)

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gc-0.3.6-1.x86_64.rpm
    sudo rpm -i gc-0.3.6-1.x86_64.rpm

ARM64 设备：

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gc-0.3.6-1.aarch64.rpm
    sudo rpm -i gc-0.3.6-1.aarch64.rpm

### Linux 二进制

AMD64：

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gc_linux_amd64
    chmod +x gc_linux_amd64
    sudo mv gc_linux_amd64 /usr/local/bin/gc

ARM64：

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gc_linux_arm64
    chmod +x gc_linux_arm64
    sudo mv gc_linux_arm64 /usr/local/bin/gc

## 验证安装

    gc version
```

#### 注意事项

1. **版本号替换**：将模板中的 `0.3.4` 替换为实际版本号
2. **避免 `#` 字符问题**：GitCode 会错误渲染代码块内的 `#`
   - Issue 引用使用 `Issue XX` 格式，不使用 `#XX`
   - PR 引用使用 `PR XX` 格式，不使用 `#XX`
   - 代码块注释单独成行，避免行内注释
3. **完整下载路径**：所有安装命令必须包含完整下载 URL

### 发布命令参考

```bash
# 查看 Release
gc release view v0.3.4 -R gitcode-cli/cli

# 列出所有 Releases
gc release list -R gitcode-cli/cli

# 下载资产
gc release download v0.3.4 -R gitcode-cli/cli
```

---

## 安装指南

### Wheel 包（跨平台，推荐）

```bash
python3 -m venv .venv
source .venv/bin/activate

pip install https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gitcode_cli-0.3.6-py3-none-any.whl
```

### DEB (Debian/Ubuntu)

```bash
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gc_0.3.6_amd64.deb
sudo dpkg -i gc_0.3.6_amd64.deb
```

### RPM (RHEL/CentOS/Fedora)

```bash
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gc-0.3.6-1.x86_64.rpm
sudo rpm -i gc-0.3.6-1.x86_64.rpm
```

### Linux 二进制

```bash
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.6/gc_linux_amd64
chmod +x gc_linux_amd64
sudo mv gc_linux_amd64 /usr/local/bin/gc
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
gc release create v1.0.0 -R gitcode-cli/cli --title "v1.0.0" --notes "Release notes"
```

### Q: nfpm 找不到命令

`package.sh` 会自动查找 `~/go/bin/nfpm`。如果提示找不到：

```bash
# 确认已安装
ls ~/go/bin/nfpm

# 如果没有，安装
go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest
```

### Q: 上传失败

1. 检查 Token：`gc auth status`
2. 检查仓库权限
3. 确认 Release 已创建

### Q: 版本号已存在

每个 tag 只能创建一个 Release，使用新版本号。

---

## 附录：手动构建步骤

> **注意**：推荐使用 `package.sh` 脚本，以下手动步骤仅供参考。

### 手动构建 DEB/RPM 包

```bash
# 1. 构建二进制
mkdir -p dist
GOOS=linux GOARCH=amd64 go build -o dist/gc_linux_amd64 ./cmd/gc
GOOS=linux GOARCH=arm64 go build -o dist/gc_linux_arm64 ./cmd/gc

# 2. 更新版本号
VERSION="0.3.4"
sed -i "s/version: .*/version: \"$VERSION\"/" nfpm-amd64.yaml
sed -i "s/version: .*/version: \"$VERSION\"/" nfpm-arm64.yaml

# 3. 构建 DEB
~/go/bin/nfpm package -f nfpm-amd64.yaml -p deb -t dist/
~/go/bin/nfpm package -f nfpm-arm64.yaml -p deb -t dist/

# 4. 构建 RPM
~/go/bin/nfpm package -f nfpm-amd64.yaml -p rpm -t dist/
~/go/bin/nfpm package -f nfpm-arm64.yaml -p rpm -t dist/
```

### 手动构建 PyPI 包

```bash
# 1. 构建多平台二进制
mkdir -p gc_cli/bin
GOOS=linux GOARCH=amd64 go build -o gc_cli/bin/gc-linux-amd64 ./cmd/gc
GOOS=linux GOARCH=arm64 go build -o gc_cli/bin/gc-linux-arm64 ./cmd/gc
GOOS=darwin GOARCH=amd64 go build -o gc_cli/bin/gc-darwin-amd64 ./cmd/gc
GOOS=darwin GOARCH=arm64 go build -o gc_cli/bin/gc-darwin-arm64 ./cmd/gc
GOOS=windows GOARCH=amd64 go build -o gc_cli/bin/gc-windows-amd64.exe ./cmd/gc

# 2. 更新版本号
VERSION="0.3.4"
sed -i "s/version = \".*/version = \"$VERSION\"/" pyproject.toml
sed -i "s/__version__ = \".*/__version__ = \"$VERSION\"/" gc_cli/__init__.py

# 3. 构建 wheel
python3 -m build --wheel --sdist --outdir dist/
```

### nfpm 配置示例

项目包含两个 nfpm 配置文件：

- `nfpm-amd64.yaml` - amd64/x86_64 架构
- `nfpm-arm64.yaml` - arm64/aarch64 架构

配置文件结构：

```yaml
name: "gc"
arch: "amd64"
platform: "linux"
version: "0.3.4"
maintainer: "gitcode-cli contributors"
description: "GitCode CLI - Command line tool for GitCode"
homepage: "https://gitcode.com/gitcode-cli/cli"
license: "MIT"
contents:
  - src: ./dist/gc_linux_amd64
    dst: /usr/bin/gc
  - src: ./completions/gc.bash
    dst: /usr/share/bash-completion/completions/gc
```

---

## Release 说明编写规范

### GitCode Markdown 渲染问题

**已知问题**：GitCode 的 Markdown 渲染器会错误地将代码块内的 `#` 渲染成标题。

**错误示例**：

```bash
source .venv/bin/activate  # Linux/macOS
# .venv\Scripts\activate   # Windows
```

渲染后 `# Windows` 会显示为一级标题，导致格式混乱。

**正确做法**：

1. 使用普通代码块，不指定语法高亮
2. 注释单独成行
3. 避免行内注释

**推荐格式**：

```
# 创建虚拟环境
python3 -m venv .venv
source .venv/bin/activate

# Windows 用户使用
.venv\Scripts\activate
```

---

**最后更新**: 2026-03-27