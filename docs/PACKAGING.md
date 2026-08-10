# 本地打包和发布指南

> 项目概述和功能介绍请参阅 [README.md](../README.md)，命令使用请参阅 [COMMANDS.md](./COMMANDS.md)，开发与规范入口请参阅 [spec/README.md](../spec/README.md)。

本文档说明如何在本地构建验证 DEB/RPM/PyPI 包，以及如何通过标准 workflow 发布并同步 Release。

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
# 发布前全量打包验证（构建 DEB + RPM + PyPI，推荐）
./scripts/package.sh v0.10.3 release

# 构建所有包（DEB + RPM + PyPI）
./scripts/package.sh v0.10.3

# 仅构建 Linux 包（DEB + RPM）
./scripts/package.sh v0.10.3 linux

# 仅构建 DEB 包
./scripts/package.sh v0.10.3 deb

# 仅构建 PyPI 包
./scripts/package.sh v0.10.3 pypi
```

### 构建目标

| 目标 | 说明 | 输出文件 |
|------|------|----------|
| `release` | DEB + RPM + PyPI（发布前全量验证，推荐） | 全部包 |
| `all` | 构建所有包（默认） | 全部包 |
| `linux` | DEB + RPM 包 | `gc_*.deb`, `gc-*.rpm` |
| `deb` | 仅 DEB 包 | `gc_*.deb` |
| `rpm` | 仅 RPM 包 | `gc-*.rpm` |
| `pypi` | 仅 PyPI 包 | `gitcode_cli-*.whl` |

### 脚本功能

`package.sh` 自动完成：

| 步骤 | 说明 |
|------|------|
| 1. 版本同步 | 通过 `VERSION` + `scripts/sync-package-version.sh` 对齐 nFPM、Python 与 npm metadata，再同步版本化文档 |
| 2. 构建 Linux 二进制 | amd64 + arm64 |
| 3. 构建 DEB 包 | amd64 + arm64 |
| 4. 构建 RPM 包 | x86_64 + aarch64 |
| 5. 构建 PyPI 包 | 多平台二进制 + wheel + sdist |

### 构建产物

```bash
dist/
├── gc_0.10.3_amd64.deb              # DEB amd64
├── gc_0.10.3_arm64.deb              # DEB arm64
├── gc-0.10.3-1.x86_64.rpm           # RPM x86_64
├── gc-0.10.3-1.aarch64.rpm          # RPM aarch64
├── gc_linux_amd64                  # Linux 二进制 amd64
├── gc_linux_arm64                  # Linux 二进制 arm64
├── gitcode_cli-0.10.3-py3-none-any.whl  # PyPI wheel
├── gitcode_cli-0.10.3.tar.gz        # PyPI sdist
└── gitcode-cli-cli-0.10.3.tgz       # npm tarball（正式 release workflow）
```

---

## 前置要求

### 安装 nfpm

[nfpm](https://github.com/goreleaser/nfpm) 是 DEB/RPM 包构建工具。

```bash
go install github.com/goreleaser/nfpm/v2/cmd/nfpm@v2.41.3
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
# 1. 本地验证所有包；这些制品不得用于正式发布
./scripts/package.sh v0.10.3 release

# 2. 发布准备 PR 合入两个远端 main 后，触发正式 workflow
gh workflow run release.yml -R gitcode-cli/cli -f version=v0.10.3
gh run watch <run-id> -R gitcode-cli/cli

# 3. workflow 全部成功后，同步同一 tag 到 GitCode（SSH）
git fetch github tag v0.10.3
git push origin refs/tags/v0.10.3

# 4. 下载 GitHub workflow 生成的正式制品
mkdir -p dist/github-release
gh release download v0.10.3 -R gitcode-cli/cli --dir dist/github-release

# 5. 验证覆盖全部正式资产的 SHA-256 清单
cd dist/github-release
sha256sum -c gc_0.10.3_checksums.txt
cd ../..

# 6. 使用受跟踪的同一份说明创建 GitCode Release
gc release create v0.10.3 -R gitcode-cli/cli \
  --title "GitCode CLI v0.10.3" \
  --notes-file docs/releases/v0.10.3.md \
  --target main \
  --json

# 7. 将同一批正式制品上传到 GitCode，不得重新构建
gc release upload v0.10.3 dist/github-release/* -R gitcode-cli/cli --json
```

> **注意**：将示例中的版本号 `0.10.3` 替换为实际版本号。正式制品必须携带准确 commit SHA；本地验证包不得上传。

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

### Wheel 包（跨平台、隔离安装）

内置全平台二进制（Linux x64/ARM、macOS Intel/Apple Silicon、Windows x64），创建虚拟环境并安装：

    python3 -m venv .venv
    source .venv/bin/activate
    pip install https://gitcode.com/gitcode-cli/cli/releases/download/v0.10.3/gitcode_cli-0.10.3-py3-none-any.whl

Windows 用户激活虚拟环境：

    .\.venv\Scripts\Activate.ps1

Windows PowerShell 用户建议运行：

    gitcode version

说明：wheel 会同时安装 `gc` 和 `gitcode` 两个命令入口，功能相同。Windows 使用 `py -m pip install --user ...` 时，脚本会安装到 Python user scheme 的 `Scripts` 目录；请运行 `py -c "import os, sysconfig; print(sysconfig.get_path('scripts', os.name + '_user'))"` 获取准确路径，将其加入用户 `PATH` 后重新打开终端，配置前可直接运行 `py -m gc_cli version`。PowerShell 预置 `gc` 作为 `Get-Content` 别名；如果 `gc version` 被解析为读取文件，请改用 `gitcode version`、`gc.exe version` 或 `py -m gc_cli version`。

### DEB (Debian/Ubuntu)

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.10.3/gc_0.10.3_amd64.deb
    sudo dpkg -i gc_0.10.3_amd64.deb

ARM64 设备：

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.10.3/gc_0.10.3_arm64.deb
    sudo dpkg -i gc_0.10.3_arm64.deb

### RPM (RHEL/CentOS/Fedora)

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.10.3/gc-0.10.3-1.x86_64.rpm
    sudo rpm -i gc-0.10.3-1.x86_64.rpm

ARM64 设备：

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.10.3/gc-0.10.3-1.aarch64.rpm
    sudo rpm -i gc-0.10.3-1.aarch64.rpm

### Linux 二进制

AMD64：

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.10.3/gc_linux_amd64
    chmod +x gc_linux_amd64
    sudo mv gc_linux_amd64 /usr/local/bin/gc

ARM64：

    wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.10.3/gc_linux_arm64
    chmod +x gc_linux_arm64
    sudo mv gc_linux_arm64 /usr/local/bin/gc

### Homebrew (macOS/Linux)

    brew install gitcode-cli/homebrew-tap/gc

更新到最新版本：

    brew upgrade gc

shell 补全（bash/zsh/fish）随安装自动配置。formula 由 GoReleaser 在发布流程中生成并推送到 [gitcode-cli/homebrew-tap](https://github.com/gitcode-cli/homebrew-tap)（见 `.goreleaser.yaml` `brews:` 与 release workflow `brew` job）。

### npm (跨平台)

    npx @gitcode-cli/cli@latest install     # 一行 bootstrap：装二进制 + 补全
    npm install -g @gitcode-cli/cli          # 或全局安装

npm 包 `@gitcode-cli/cli` 内置 Linux/macOS/Windows 多平台二进制（`npm/bin/platforms/`），Node wrapper 按平台选择并 exec。Windows bootstrap 同时安装 `gc.exe` 与 `gitcode.exe`。

正式 npm tarball 不再独立编译二进制：release workflow 的 `artifacts` job 调用 `scripts/prepare-npm-package.sh`，从同一批 GoReleaser 归档/裸二进制组装 npm 包，并将 `.tgz` 纳入 Release SHA-256 清单。`npm` job 只下载已验证的 Release artifact 并执行 OIDC Trusted Publishing；目标版本已存在时，必须下载 registry tarball 与 Release tarball 比对 SHA-256，内容一致才允许幂等跳过。

npm global 与 bootstrap 默认 `auto` 更新 stable 版本，提供 `notify` / `off`、24 小时 TTL、跨进程锁、健康检查和回滚。安装/升级验证必须覆盖：

    gitcode doctor install --json
    gitcode update --check --json
    gitcode config set update.mode off

更新器只允许操作 `@gitcode-cli/cli`；不得自动卸载 pip/Homebrew/DEB/RPM，不得提权或重写 PATH。鉴权使用 **OIDC Trusted Publishing**（`id-token: write`，无 `NPM_TOKEN`）；`npm/package.json` 的 `repository.url` 须保持为 `https://github.com/gitcode-cli/cli.git`。

## 验证安装

    gc version
```

#### 注意事项

1. **版本号替换**：将模板中的 `0.10.3` 替换为实际版本号
2. **避免 `#` 字符问题**：GitCode 会错误渲染代码块内的 `#`
   - Issue 引用使用 `Issue XX` 格式，不使用 `#XX`
   - PR 引用使用 `PR XX` 格式，不使用 `#XX`
   - 代码块注释单独成行，避免行内注释
3. **完整下载路径**：所有安装命令必须包含完整下载 URL

### 发布命令参考

```bash
# 查看 Release
gc release view v0.10.3 -R gitcode-cli/cli

# 列出所有 Releases
gc release list -R gitcode-cli/cli

# 下载资产
gc release download v0.10.3 -R gitcode-cli/cli
```

---

## 安装指南

### Wheel 包（跨平台、隔离安装）

内置全平台二进制（Linux x64/ARM、macOS Intel/Apple Silicon、Windows x64）：

```bash
python3 -m venv .venv
source .venv/bin/activate

pip install https://gitcode.com/gitcode-cli/cli/releases/download/v0.10.3/gitcode_cli-0.10.3-py3-none-any.whl

# Windows PowerShell 中推荐使用 gitcode
gitcode version
```

说明：wheel 会同时安装 `gc` 和 `gitcode` 两个命令入口，功能相同。Windows 使用 `py -m pip install --user ...` 时，脚本会安装到 Python user scheme 的 `Scripts` 目录；请运行 `py -c "import os, sysconfig; print(sysconfig.get_path('scripts', os.name + '_user'))"` 获取准确路径，将其加入用户 `PATH` 后重新打开终端，配置前可直接运行 `py -m gc_cli version`。PowerShell 预置 `gc` 作为 `Get-Content` 别名；如果 `gc version` 被解析为读取文件，请改用 `gitcode version`、`gc.exe version` 或 `py -m gc_cli version`。

### DEB (Debian/Ubuntu)

```bash
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.10.3/gc_0.10.3_amd64.deb
sudo dpkg -i gc_0.10.3_amd64.deb
```

DEB/RPM packages install both `gc` and `gitcode`; on Linux they are equivalent.

### RPM (RHEL/CentOS/Fedora)

```bash
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.10.3/gc-0.10.3-1.x86_64.rpm
sudo rpm -i gc-0.10.3-1.x86_64.rpm
```

DEB/RPM packages install both `gc` and `gitcode`; on Linux they are equivalent.

### 多渠道 PATH 冲突

pip、npm、Homebrew、DEB/RPM 和手工 archive 都可能提供同名命令。安装器只诊断、不替用户卸载其他渠道或修改 PATH：

```bash
gitcode doctor install
gitcode doctor install --json
```

先根据 `selected`、`candidates`、`distribution` 和 `recommendations` 确认实际命令，再由用户选择调整 PATH、升级原渠道或显式卸载旧渠道。Bash 更改 PATH 后运行 `hash -r`，Zsh 运行 `rehash`，PowerShell 重新打开终端；Windows 不应全局删除内置 `gc`/`Get-Content` alias，使用 `gitcode`。

### Linux 二进制

```bash
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.10.3/gc_linux_amd64
chmod +x gc_linux_amd64
sudo mv gc_linux_amd64 /usr/local/bin/gc
```

### 验证安装

```bash
gc version
gitcode version  # DEB/RPM packages
gitcode doctor install
```

---

## 常见问题

### Q: 创建 Release 时返回 400 错误

确保包含 `--notes` 参数：

```bash
gc release create vX.Y.Z -R gitcode-cli/cli --title "vX.Y.Z" --notes "Release notes"
```

### Q: nfpm 找不到命令

`package.sh` 会自动查找 `~/go/bin/nfpm`。如果提示找不到：

```bash
# 确认已安装
ls ~/go/bin/nfpm

# 如果没有，安装
go install github.com/goreleaser/nfpm/v2/cmd/nfpm@v2.41.3
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
VERSION="0.10.3"
bash scripts/sync-package-version.sh "$VERSION"

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
VERSION="0.10.3"
bash scripts/sync-package-version.sh "$VERSION"

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
version: "0.10.3"
maintainer: "gitcode-cli contributors"
description: "GitCode CLI - Command line tool for GitCode"
homepage: "https://gitcode.com/gitcode-cli/cli"
license: "MIT"
scripts:
  postinstall: ./build/scripts/postinstall.sh
contents:
  - src: ./dist/gc_linux_amd64
    dst: /usr/bin/gc
  - src: /usr/bin/gc
    dst: /usr/bin/gitcode
    type: symlink
  - src: ./completions/gc.bash
    dst: /usr/share/bash-completion/completions/gc
  - src: ./build/completions/gitcode.bash
    dst: /usr/share/bash-completion/completions/gitcode
```

---

## Release 说明编写规范

### GitCode Markdown 渲染问题

**已知问题**：GitCode 的 Markdown 渲染器会错误地将代码块内的 `#` 渲染成标题。

**错误示例**：

```bash
source .venv/bin/activate  # Linux/macOS
# .\.venv\Scripts\Activate.ps1  # Windows PowerShell
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
.\.venv\Scripts\Activate.ps1
```

---

**最后更新**: 2026-03-27
