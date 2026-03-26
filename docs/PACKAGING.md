# 本地打包和发布指南

> 项目概述和功能介绍请参阅 [README.md](../README.md)，命令使用请参阅 [COMMANDS.md](./COMMANDS.md)，开发指南请参阅 [CLAUDE.md](../CLAUDE.md)。

本文档说明如何在本地构建 DEB/RPM/PyPI 包并使用 `gc` 命令发布 Release。

## 重要流程

**每次打包发布后，必须同步更新 README.md 中的下载版本信息！**

## 快速开始

### 一键打包（推荐）

使用 `scripts/package.sh` 脚本一键完成版本同步和打包：

```bash
# 构建所有包（DEB + RPM + PyPI）
./scripts/package.sh v0.3.0

# 仅构建 Linux 包（DEB + RPM）
./scripts/package.sh v0.3.0 linux

# 仅构建 DEB 包
./scripts/package.sh v0.3.0 deb

# 仅构建 PyPI 包
./scripts/package.sh v0.3.0 pypi

# 发布用（DEB + RPM + PyPI）
./scripts/package.sh v0.3.0 release
```

### 构建目标

| 目标 | 说明 |
|------|------|
| `all` | 构建所有包（默认） |
| `deb` | 仅构建 DEB 包 |
| `rpm` | 仅构建 RPM 包 |
| `linux` | 构建 DEB + RPM 包 |
| `pypi` | 仅构建 PyPI 包 |
| `release` | 构建 DEB + RPM + PyPI（发布用） |

### 版本自动同步

脚本会自动同步版本号到以下文件：

**配置文件：**
- `nfpm-amd64.yaml`
- `nfpm-arm64.yaml`
- `pyproject.toml`
- `gc_cli/__init__.py`

**文档文件：**
- `README.md` - Release badge 和下载链接
- `docs/AI-GUIDE.md` - 安装命令中的版本号
- `docs/PACKAGING.md` - 示例命令中的版本号

## 前置要求

### 安装 nfpm

[nfpm](https://github.com/goreleaser/nfpm) 是一个通用的包构建工具。

```bash
# 使用 go install（安装到 ~/go/bin/nfpm）
go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest
```

> **重要**：`go install` 会将 nfpm 安装到 `~/go/bin/` 目录。如果该目录不在 PATH 中，有以下解决方案：

**方案一：添加到 PATH（推荐）**

```bash
# 临时添加（当前终端有效）
export PATH="$HOME/go/bin:$PATH"

# 永久添加（写入 shell 配置）
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

**方案二：使用完整路径**

```bash
# 直接使用完整路径调用
~/go/bin/nfpm --version
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
version: "0.3.0"
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
VERSION="0.3.0"

# 更新 nfpm 配置文件中的版本号
sed -i "s/version: .*/version: \"$VERSION\"/" nfpm-amd64.yaml
sed -i "s/version: .*/version: \"$VERSION\"/" nfpm-arm64.yaml
```

### 4. 构建 DEB 包

```bash
# 如果 nfpm 在 PATH 中
nfpm package -f nfpm-amd64.yaml -p deb -t dist/
nfpm package -f nfpm-arm64.yaml -p deb -t dist/

# 如果 nfpm 不在 PATH 中，使用完整路径
~/go/bin/nfpm package -f nfpm-amd64.yaml -p deb -t dist/
~/go/bin/nfpm package -f nfpm-arm64.yaml -p deb -t dist/
```

### 5. 构建 RPM 包

```bash
# 如果 nfpm 在 PATH 中
nfpm package -f nfpm-amd64.yaml -p rpm -t dist/
nfpm package -f nfpm-arm64.yaml -p rpm -t dist/

# 如果 nfpm 不在 PATH 中，使用完整路径
~/go/bin/nfpm package -f nfpm-amd64.yaml -p rpm -t dist/
~/go/bin/nfpm package -f nfpm-arm64.yaml -p rpm -t dist/
```

### 6. 查看构建结果

```bash
ls -la dist/*.deb dist/*.rpm
```

输出示例：
```
dist/gc_0.3.1_amd64.deb
dist/gc_0.3.1_arm64.deb
dist/gc-0.3.1-1.x86_64.rpm
dist/gc-0.3.1-1.aarch64.rpm
```

---

## 发布 Release

### 1. 创建 Release

```bash
gc release create v0.3.0 -R owner/repo \
  --title "gc v0.3.0" \
  --notes "Release notes here"
```

> **注意**：`--notes` 参数是必需的，不带此参数可能返回 400 错误。

### 2. 上传资产

```bash
# 上传单个文件
gc release upload v0.3.0 dist/gc_0.3.1_amd64.deb -R owner/repo

# 上传所有包（包括 wheel）
gc release upload v0.3.0 \
  dist/gc_0.3.1_amd64.deb \
  dist/gc_0.3.1_arm64.deb \
  dist/gc-0.3.1-1.x86_64.rpm \
  dist/gc-0.3.1-1.aarch64.rpm \
  dist/gitcode_cli-0.3.1-py3-none-any.whl \
  -R owner/repo
```

### 3. 查看 Release

```bash
gc release view v0.3.0 -R owner/repo
```

### 4. 列出所有 Releases

```bash
gc release list -R owner/repo
```

### 5. 下载资产

```bash
# 下载所有资产
gc release download v0.3.0 -R owner/repo

# 下载到指定目录
gc release download v0.3.0 -R owner/repo -o ./downloads/

# 下载指定文件
gc release download v0.3.0 gc_0.3.1_amd64.deb -R owner/repo
```

---

## 完整示例

以下是一个完整的打包和发布流程：

```bash
#!/bin/bash
set -e

# 配置
VERSION="0.3.0"
REPO="owner/repo"
TOKEN="your_token"

# nfpm 路径（根据实际情况修改）
# 如果 nfpm 在 PATH 中，使用: NFPAM="nfpm"
# 如果 nfpm 不在 PATH 中，使用完整路径: NFPAM="$HOME/go/bin/nfpm"
NFPAM="$HOME/go/bin/nfpm"

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
$NFPAM package -f nfpm-amd64.yaml -p deb -t dist/
$NFPAM package -f nfpm-arm64.yaml -p deb -t dist/

# 构建 RPM 包
echo "Building RPM packages..."
$NFPAM package -f nfpm-amd64.yaml -p rpm -t dist/
$NFPAM package -f nfpm-arm64.yaml -p rpm -t dist/

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

## 构建 PyPI 包

### 1. 前置要求

确保系统安装了 Python 3.8+ 和 pip。

```bash
# 检查 Python 版本
python3 --version

# 安装 build 工具（推荐使用虚拟环境）
python3 -m venv .venv
source .venv/bin/activate
pip install --upgrade build wheel setuptools
```

### 2. 准备二进制文件

PyPI 包依赖预编译的二进制文件，需要先构建对应平台的二进制：

```bash
# 创建输出目录
mkdir -p gc_cli/bin

# 构建 Linux amd64
GOOS=linux GOARCH=amd64 go build -o gc_cli/bin/gc-linux-amd64 ./cmd/gc

# 构建 Linux arm64
GOOS=linux GOARCH=arm64 go build -o gc_cli/bin/gc-linux-arm64 ./cmd/gc

# 构建 macOS amd64
GOOS=darwin GOARCH=amd64 go build -o gc_cli/bin/gc-darwin-amd64 ./cmd/gc

# 构建 macOS arm64
GOOS=darwin GOARCH=arm64 go build -o gc_cli/bin/gc-darwin-arm64 ./cmd/gc

# 构建 Windows amd64
GOOS=windows GOARCH=amd64 go build -o gc_cli/bin/gc-windows-amd64.exe ./cmd/gc
```

### 3. 更新版本号

确保 `pyproject.toml` 和 `gc_cli/__init__.py` 中的版本号一致：

```bash
VERSION="0.3.0"

# 更新 pyproject.toml
sed -i "s/version = \".*/version = \"$VERSION\"/" pyproject.toml

# 更新 gc_cli/__init__.py
sed -i "s/__version__ = \".*/__version__ = \"$VERSION\"/" gc_cli/__init__.py
```

### 4. 构建包

```bash
# 激活虚拟环境（如果使用）
source .venv/bin/activate

# 构建 wheel 和 sdist
python -m build --wheel --sdist
```

构建产物位于 `dist/` 目录：

```
dist/
├── gitcode_cli-0.3.1-py3-none-any.whl
└── gitcode_cli-0.3.1.tar.gz
```

### 5. 本地测试安装

```bash
# 创建测试虚拟环境
python -m venv /tmp/gc-test-env
source /tmp/gc-test-env/bin/activate

# 安装 wheel
pip install dist/gitcode_cli-0.3.1-py3-none-any.whl

# 测试命令
gc version

# 清理测试环境
deactivate
rm -rf /tmp/gc-test-env
```

### 6. 上传到 PyPI（CI 自动化）

发布到 PyPI 由 GitHub Actions 自动完成，无需手动上传。详见 [RELEASE.md](../RELEASE.md)。

如需手动上传到 TestPyPI 进行测试：

```bash
# 安装 twine
pip install twine

# 上传到 TestPyPI
twine upload --repository testpypi dist/*

# 从 TestPyPI 安装测试
pip install --index-url https://test.pypi.org/simple/ gitcode-cli
```

---

## 安装指南

### Wheel 包（跨平台，推荐）

> **推荐**: 从 Release 归档下载 wheel 包，版本与发布一致。

```bash
# 创建虚拟环境
python3 -m venv .venv
source .venv/bin/activate  # Linux/macOS
# .venv\Scripts\activate   # Windows

# 安装（一行命令）
pip install https://gitcode.com/gitcode-cli/cli/releases/download/v0.3.1/gitcode_cli-0.3.1-py3-none-any.whl
```

### PyPI（备选）

> ⚠️ **注意**: PyPI 官方源可能有同步延迟，推荐使用上方 wheel 包下载

```bash
# 创建虚拟环境
python3 -m venv .venv
source .venv/bin/activate  # Linux/macOS
# .venv\Scripts\activate   # Windows

# 使用官方 PyPI 源安装
pip install -i https://pypi.org/simple/ gitcode-cli
```

### DEB (Debian/Ubuntu)

```bash
# 下载并安装
wget https://gitcode.com/owner/repo/releases/download/v0.3.1/gc_0.3.1_amd64.deb
sudo dpkg -i gc_0.3.1_amd64.deb

# 或使用 gc 命令下载
gc release download v0.3.0 gc_0.3.1_amd64.deb -R owner/repo
sudo dpkg -i gc_0.3.1_amd64.deb
```

### RPM (RHEL/CentOS/Fedora)

```bash
# 下载并安装
wget https://gitcode.com/owner/repo/releases/download/v0.3.1/gc-0.3.1-1.x86_64.rpm
sudo rpm -i gc-0.3.1-1.x86_64.rpm

# 或使用 gc 命令下载
gc release download v0.3.0 gc-0.3.1-1.x86_64.rpm -R owner/repo
sudo rpm -i gc-0.3.1-1.x86_64.rpm
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

**原因**：`go install` 安装的程序默认放在 `~/go/bin/`，该目录可能不在 PATH 中。

**解决方案一：添加到 PATH（推荐）**

```bash
# 检查 nfpm 是否存在
ls ~/go/bin/nfpm

# 添加到 PATH
export PATH="$HOME/go/bin:$PATH"

# 永久生效，添加到 shell 配置
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc

# 验证
nfpm --version
```

**解决方案二：使用完整路径**

```bash
# 直接使用完整路径
~/go/bin/nfpm --version

# 或在脚本开头定义变量
NFPAM="$HOME/go/bin/nfpm"
$NFPAM package -f nfpm-amd64.yaml -p deb -t dist/
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

1. **移除语法高亮**：使用普通代码块，不指定 `bash`
2. **注释单独成行**：将 `#` 注释放在单独一行
3. **避免行内注释**：不要在命令后使用 `# 注释` 格式

**推荐格式**：

```
# 创建虚拟环境
python3 -m venv .venv
source .venv/bin/activate

# Windows 用户使用
.venv\Scripts\activate
```

或使用纯文本说明：

```
python3 -m venv .venv
source .venv/bin/activate       Linux/macOS
.venv\Scripts\activate          Windows
```

---

**最后更新**: 2026-03-26