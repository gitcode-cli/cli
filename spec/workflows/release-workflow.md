# Release 流程

本文档定义版本发布的完整流程。

## 流程概览

```
更新版本号 → 构建 → 打包 → 创建 Release → 上传包 → 更新文档 → 提交推送
```

## 1. 发布前准备

### 检查清单

- [ ] 所有 Issue 已关闭或移至下一版本
- [ ] 单元测试全部通过
- [ ] 功能已在测试仓库验证
- [ ] CHANGELOG 已更新（如有）

### 确认分支

```bash
# 确保在 main 分支
git checkout main
git pull
```

## 2. 更新版本号

### 需要更新的文件

| 文件 | 更新内容 |
|------|----------|
| `nfpm-amd64.yaml` | `version: "x.x.x"` |
| `nfpm-arm64.yaml` | `version: "x.x.x"` |
| `pyproject.toml` | `version = "x.x.x"` |
| `gc_cli/__init__.py` | `__version__ = "x.x.x"` |
| `README.md` | 下载链接版本号、Release badge |

### 更新命令

```bash
VERSION=0.2.10

# 更新 nfpm 配置
sed -i "s/version: \".*\"/version: \"$VERSION\"/" nfpm-amd64.yaml nfpm-arm64.yaml

# 更新 Python 包
sed -i "s/version = \".*\"/version = \"$VERSION\"/" pyproject.toml
sed -i "s/__version__ = \".*\"/__version__ = \"$VERSION\"/" gc_cli/__init__.py
```

## 3. 构建和打包

### 构建二进制

```bash
# 使用构建脚本
./scripts/build-release.sh v0.2.10
```

### 打包 DEB/RPM

```bash
# 使用打包脚本
./scripts/package.sh v0.2.10
```

### 预期产物

```
dist/
├── gc_0.2.10_amd64.deb
├── gc_0.2.10_arm64.deb
├── gc-0.2.10-1.x86_64.rpm
├── gc-0.2.10-1.aarch64.rpm
└── gitcode_cli-0.2.10-py3-none-any.whl
```

## 4. 创建 Release

### 命令行创建

```bash
./gc release create v0.2.10 \
  --title "v0.2.10" \
  --notes "## 变更内容

### 新增
- xxx 命令

### 修复
- 修复 xxx 问题

### 安装

\`\`\`bash
# DEB
sudo dpkg -i gc_0.2.10_amd64.deb

# RPM
sudo rpm -i gc-0.2.10-1.x86_64.rpm

# PyPI（使用虚拟环境）
python -m venv .venv
source .venv/bin/activate
pip install gitcode-cli
\`\`\`
" \
  -R gitcode-cli/cli
```

### Web 界面创建

1. 访问 https://gitcode.com/gitcode-cli/cli/releases/new
2. 填写 Tag、Title、Notes
3. 点击 Create Release

## 5. 上传包

### 命令行上传

```bash
./gc release upload v0.2.10 \
  dist/gc_0.2.10_amd64.deb \
  dist/gc_0.2.10_arm64.deb \
  dist/gc-0.2.10-1.x86_64.rpm \
  dist/gc-0.2.10-1.aarch64.rpm \
  dist/gitcode_cli-0.2.10-py3-none-any.whl \
  -R gitcode-cli/cli
```

### 验证上传

1. 访问 Release 页面
2. 确认所有包都已上传
3. 测试下载链接

## 6. 发布 PyPI 包

```bash
# 安装 twine
pip install twine

# 上传
twine upload dist/gitcode_cli-0.2.10-py3-none-any.whl
```

## 7. 更新文档

### README.md 更新

```markdown
## 安装

### 下载安装包

[![Release](https://img.shields.io/badge/Release-v0.2.10-blue)](https://gitcode.com/gitcode-cli/cli/releases/tag/v0.2.10)

| 平台 | 下载 |
|------|------|
| Linux AMD64 (DEB) | [gc_0.2.10_amd64.deb](https://gitcode.com/gitcode-cli/cli/releases/download/v0.2.10/gc_0.2.10_amd64.deb) |
| Linux ARM64 (DEB) | [gc_0.2.10_arm64.deb](https://gitcode.com/gitcode-cli/cli/releases/download/v0.2.10/gc_0.2.10_arm64.deb) |
| Linux AMD64 (RPM) | [gc-0.2.10-1.x86_64.rpm](https://gitcode.com/gitcode-cli/cli/releases/download/v0.2.10/gc-0.2.10-1.x86_64.rpm) |
| Linux ARM64 (RPM) | [gc-0.2.10-1.aarch64.rpm](https://gitcode.com/gitcode-cli/cli/releases/download/v0.2.10/gc-0.2.10-1.aarch64.rpm) |
```

## 8. 提交版本更新

```bash
# 暂存更改
git add nfpm-*.yaml pyproject.toml gc_cli/__init__.py README.md

# 提交
git commit -m "chore: release v0.2.10

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"

# 推送
git push
```

## 完整发布流程

```bash
# 1. 确保在 main 分支
git checkout main && git pull

# 2. 更新版本号
VERSION=0.2.10
sed -i "s/version: \".*\"/version: \"$VERSION\"/" nfpm-amd64.yaml nfpm-arm64.yaml
sed -i "s/version = \".*\"/version = \"$VERSION\"/" pyproject.toml
sed -i "s/__version__ = \".*\"/__version__ = \"$VERSION\"/" gc_cli/__init__.py

# 3. 构建和打包
./scripts/package.sh v$VERSION

# 4. 创建 Release
./gc release create v$VERSION --title "v$VERSION" --notes "Release notes" -R gitcode-cli/cli

# 5. 上传包
./gc release upload v$VERSION \
  dist/gc_${VERSION}_amd64.deb \
  dist/gc_${VERSION}_arm64.deb \
  dist/gc-${VERSION}-1.x86_64.rpm \
  dist/gc-${VERSION}-1.aarch64.rpm \
  -R gitcode-cli/cli

# 6. 发布 PyPI
twine upload dist/gitcode_cli-${VERSION}-py3-none-any.whl

# 7. 更新 README.md（手动）

# 8. 提交版本更新
git add . && git commit -m "chore: release v$VERSION" && git push
```

## 发布后验证

```bash
# 1. 验证 Release 页面
# 访问 https://gitcode.com/gitcode-cli/cli/releases

# 2. 验证下载
wget https://gitcode.com/gitcode-cli/cli/releases/download/v0.2.10/gc_0.2.10_amd64.deb

# 3. 验证安装
sudo dpkg -i gc_0.2.10_amd64.deb
gc version

# 4. 验证 PyPI
pip install gitcode-cli==0.2.10
```

## 检查清单

- [ ] 版本号已更新（所有文件一致）
- [ ] 构建成功
- [ ] 打包成功（DEB/RPM/Whl）
- [ ] Release 已创建
- [ ] 所有包已上传
- [ ] PyPI 已发布
- [ ] README.md 已更新
- [ ] 版本更新已提交推送

---

**最后更新**: 2026-03-26