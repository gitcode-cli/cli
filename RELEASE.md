# GitCode CLI 版本发布指南

本文档定义 GitCode CLI 的发布流程和发布产物。

## 发布方式

**推荐方式**：通过 GitHub Actions 自动发布。

### 发布触发条件

| 触发方式 | 说明 |
|---------|------|
| 推送 `v*` 标签 | 自动触发完整发布流程 |
| 手动触发 | 通过 GitHub Actions 界面手动触发 |

### 发布命令

```bash
# 1. 确保代码已合并到 main 分支
git checkout main
git pull

# 2. 创建并推送标签
git tag v0.3.0
git push origin v0.3.0

# 3. GitHub Actions 自动执行发布
```

---

## 发布产物

每次发布生成以下三类产物：

### 1. RPM 包（RedHat 系列）

适用于 RHEL、CentOS、Fedora、openSUSE 等。

| 架构 | 文件名 | 适用系统 |
|------|--------|----------|
| x86_64 | `gc-{version}-1.x86_64.rpm` | Intel/AMD 64 位系统 |
| aarch64 | `gc-{version}-1.aarch64.rpm` | ARM 64 位系统 |

**安装方式**：
```bash
# 下载并安装
sudo rpm -i gc-0.3.0-1.x86_64.rpm

# 或使用 yum/dnf
sudo yum install gc-0.3.0-1.x86_64.rpm
```

### 2. DEB 包（Debian 系列）

适用于 Debian、Ubuntu、Linux Mint 等。

| 架构 | 文件名 | 适用系统 |
|------|--------|----------|
| amd64 | `gc_{version}_amd64.deb` | Intel/AMD 64 位系统 |
| arm64 | `gc_{version}_arm64.deb` | ARM 64 位系统 |

**安装方式**：
```bash
# 下载并安装
sudo dpkg -i gc_0.3.0_amd64.deb

# 或使用 apt
sudo apt install ./gc_0.3.0_amd64.deb
```

### 3. PyPI 包（Python）

跨平台 Python 包，内置各平台二进制文件。

| 包名 | 安装命令 |
|------|----------|
| `gitcode-cli` | `pip install gitcode-cli` |

**支持平台**：
- Linux (x86_64, arm64)
- macOS (x86_64, arm64)
- Windows (x86_64)

**安装方式**：
```bash
pip install gitcode-cli

# 或指定版本
pip install gitcode-cli==0.3.0
```

---

## 发布产物汇总表

| 产物类型 | 格式 | 架构 | 目标系统 | 发布位置 |
|---------|------|------|---------|---------|
| RPM | `.rpm` | x86_64 | RedHat/CentOS/Fedora | GitHub Release |
| RPM | `.rpm` | aarch64 | RedHat/CentOS/Fedora (ARM) | GitHub Release |
| DEB | `.deb` | amd64 | Debian/Ubuntu | GitHub Release |
| DEB | `.deb` | arm64 | Debian/Ubuntu (ARM) | GitHub Release |
| PyPI | wheel | all | 跨平台 | PyPI |

---

## 前置配置

### GitHub Secrets 配置

发布前需在 GitHub 仓库设置以下 Secrets：

| Secret 名称 | 必需 | 用途 |
|------------|------|------|
| `GITHUB_TOKEN` | ✅ 自动 | GitHub Release 创建 |
| `PYPI_API_TOKEN` | ✅ 必须 | PyPI 包发布 |
| `HOMEBREW_TAP_GITHUB_TOKEN` | ⚠️ 可选 | Homebrew tap 更新 |
| `SCOOP_BUCKET_GITHUB_TOKEN` | ⚠️ 可选 | Scoop bucket 更新 |

### PyPI Token 获取

1. 登录 [PyPI](https://pypi.org/)
2. 进入 Account settings → API tokens
3. 创建新 token，选择 "Entire account" 或指定项目
4. 将 token 添加到 GitHub Secrets 的 `PYPI_API_TOKEN`

---

## 发布流程

### 完整发布流程

```
┌─────────────────────────────────────────────────────────────┐
│                      推送 v* 标签                            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    GitHub Actions 触发                       │
└─────────────────────────────────────────────────────────────┘
                              │
                    ┌─────────┴─────────┐
                    ▼                   ▼
           ┌─────────────┐       ┌─────────────┐
           │  GoReleaser │       │    PyPI     │
           │   构建      │       │   发布      │
           └─────────────┘       └─────────────┘
                    │                   │
                    ▼                   ▼
           ┌─────────────┐       ┌─────────────┐
           │ GitHub      │       │ PyPI        │
           │ Release     │       │ gitcode-cli │
           │ - RPM x2    │       │             │
           │ - DEB x2    │       │             │
           │ - 源码包    │       │             │
           └─────────────┘       └─────────────┘
```

> **说明**：GoReleaser 负责构建 RPM/DEB 包并发布到 GitHub Release。PyPI 发布在 GoReleaser 完成后单独执行。

### 发布步骤详解

```bash
# 步骤 1: 确保在 main 分支
git checkout main
git pull origin main

# 步骤 2: 运行测试确保代码质量
make test
make lint

# 步骤 3: 更新版本号（如果有版本文件）
# - pyproject.toml
# - gc_cli/__init__.py

# 步骤 4: 提交版本更新
git add .
git commit -m "chore: prepare for release v0.3.0"

# 步骤 5: 创建标签
git tag -a v0.3.0 -m "Release v0.3.0"

# 步骤 6: 推送标签触发发布
git push origin main --tags

# 步骤 7: 监控 GitHub Actions 执行状态
# https://github.com/your-org/your-repo/actions

# 步骤 8: 发布完成后验证
# - 检查 GitHub Release 页面
# - 检查 PyPI: https://pypi.org/project/gitcode-cli/
```

---

## 发布后验证

### 验证 RPM 包

```bash
# 在 RedHat/CentOS/Fedora 系统上
wget https://github.com/your-org/your-repo/releases/download/v0.3.0/gc-0.3.0-1.x86_64.rpm
sudo rpm -i gc-0.3.0-1.x86_64.rpm
gc version
```

### 验证 DEB 包

```bash
# 在 Debian/Ubuntu 系统上
wget https://github.com/your-org/your-repo/releases/download/v0.3.0/gc_0.3.0_amd64.deb
sudo dpkg -i gc_0.3.0_amd64.deb
gc version
```

### 验证 PyPI 包

```bash
# 创建虚拟环境测试
python -m venv test-env
source test-env/bin/activate
pip install gitcode-cli==0.3.0
gc version
```

---

## 版本命名规范

遵循 [语义化版本](https://semver.org/) 规范：

```
vMAJOR.MINOR.PATCH[-PRERELEASE]

示例：
v1.0.0         # 正式版本
v1.0.0-beta.1  # 预发布版本
v1.0.1         # Bug 修复版本
v1.1.0         # 新功能版本
v2.0.0         # 重大更新版本
```

### 版本递增规则

| 版本类型 | 递增条件 | 示例 |
|---------|---------|------|
| MAJOR | 不兼容的 API 修改 | v1.0.0 → v2.0.0 |
| MINOR | 向后兼容的功能新增 | v1.0.0 → v1.1.0 |
| PATCH | 向后兼容的问题修复 | v1.0.0 → v1.0.1 |

---

## 发布检查清单

发布前必须确认以下事项：

### 代码质量
- [ ] 所有测试通过 (`make test`)
- [ ] Lint 检查通过 (`make lint`)
- [ ] 代码覆盖率达标
- [ ] 无已知严重 Bug

### 文档更新
- [ ] CHANGELOG.md 已更新
- [ ] README.md 版本信息已更新（如需要）
- [ ] 升级指南已更新（如需要）

### 配置检查
- [ ] GitHub Secrets 已正确配置
- [ ] PyPI Token 有效
- [ ] 版本号已更新

### 发布后验证
- [ ] GitHub Release 页面正确
- [ ] RPM 包可安装运行
- [ ] DEB 包可安装运行
- [ ] PyPI 包可安装运行

---

## 故障排除

### 问题：PyPI 发布失败

**原因**：Token 无效或版本号已存在

**解决方案**：
1. 检查 `PYPI_API_TOKEN` 是否有效
2. 确认版本号未被使用过
3. PyPI 不允许重复发布同一版本

### 问题：RPM/DEB 包安装失败

**原因**：依赖缺失或架构不匹配

**解决方案**：
1. 确认系统架构 (`uname -m`)
2. 检查是否缺少依赖包
3. 使用正确的包管理器安装

---

## 相关文档

- [CONTRIBUTING.md](./CONTRIBUTING.md) - 贡献指南
- [docs/PACKAGING.md](./docs/PACKAGING.md) - 本地打包指南
- [issues-plan/10-deployment.md](./issues-plan/10-deployment.md) - 部署需求

---

**最后更新**: 2026-03-24