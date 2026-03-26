# 构建打包流程

本文档定义项目的构建和打包流程。

## 流程概览

```
本地构建 → 运行测试 → 打包 DEB/RPM → 验证安装
```

## 1. 本地构建

### 开发构建

```bash
# 构建本地版本
go build -o ./gc ./cmd/gc

# 使用本地版本测试
./gc issue list -R owner/repo
```

### 发布构建

```bash
# 使用构建脚本
./scripts/build-release.sh v0.2.10

# 手动构建（带版本号注入）
VERSION=v0.2.10
go build -ldflags "-X main.version=${VERSION#v}" -o ./gc ./cmd/gc
```

## 2. 版本号管理

### 版本号位置

发布新版本时，以下文件版本号**必须保持一致**：

| 文件 | 版本字段位置 |
|------|-------------|
| `nfpm-amd64.yaml` | `version: "x.x.x"` |
| `nfpm-arm64.yaml` | `version: "x.x.x"` |
| `pyproject.toml` | `version = "x.x.x"` |
| `gc_cli/__init__.py` | `__version__ = "x.x.x"` |

### 版本号同步

```bash
# 更新版本号（示例：v0.2.10）
VERSION=0.2.10

# nfpm 配置
sed -i "s/version: \".*\"/version: \"$VERSION\"/" nfpm-amd64.yaml
sed -i "s/version: \".*\"/version: \"$VERSION\"/" nfpm-arm64.yaml

# Python 包
sed -i "s/version = \".*\"/version = \"$VERSION\"/" pyproject.toml
sed -i "s/__version__ = \".*\"/__version__ = \"$VERSION\"/" gc_cli/__init__.py
```

## 3. 打包 DEB/RPM

### 前置条件

```bash
# 安装 nfpm
go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest
```

### 打包命令

```bash
# AMD64 DEB
~/go/bin/nfpm pkg -f nfpm-amd64.yaml -p deb -t dist/gc_0.2.10_amd64.deb

# AMD64 RPM
~/go/bin/nfpm pkg -f nfpm-amd64.yaml -p rpm -t dist/gc-0.2.10-1.x86_64.rpm

# ARM64 DEB
~/go/bin/nfpm pkg -f nfpm-arm64.yaml -p deb -t dist/gc_0.2.10_arm64.deb

# ARM64 RPM
~/go/bin/nfpm pkg -f nfpm-arm64.yaml -p rpm -t dist/gc-0.2.10-1.aarch64.rpm
```

### 使用打包脚本

```bash
# 推荐：一键打包所有格式
./scripts/package.sh v0.2.10
```

### 打包产物

```
dist/
├── gc_0.2.10_amd64.deb
├── gc_0.2.10_arm64.deb
├── gc-0.2.10-1.x86_64.rpm
└── gc-0.2.10-1.aarch64.rpm
```

## 4. 验证安装

### DEB 安装验证

```bash
# 安装
sudo dpkg -i dist/gc_0.2.10_amd64.deb

# 验证
gc version
gc auth status

# 卸载
sudo dpkg -r gc
```

### RPM 安装验证

```bash
# 安装
sudo rpm -i dist/gc-0.2.10-1.x86_64.rpm

# 验证
gc version
gc auth status

# 卸载
sudo rpm -e gc
```

## 5. PyPI 包打包

### 构建 Wheel

```bash
# 安装构建工具
pip install build

# 构建
python -m build

# 产物
# dist/gitcode_cli-0.2.10-py3-none-any.whl
```

### 上传到 PyPI

```bash
# 安装 twine
pip install twine

# 上传
twine upload dist/gitcode_cli-0.2.10-py3-none-any.whl
```

## 6. 常见问题

### 构建失败

```bash
# 清理缓存
go clean -cache

# 重新下载依赖
go mod download

# 重新构建
go build -o ./gc ./cmd/gc
```

### nfpm 找不到

```bash
# 检查安装位置
ls ~/go/bin/nfpm

# 添加到 PATH（临时）
export PATH=$PATH:~/go/bin
```

### 版本号不一致

```bash
# 检查所有版本号
grep -E "version|__version__" nfpm-*.yaml pyproject.toml gc_cli/__init__.py
```

## 检查清单

- [ ] 版本号已同步更新
- [ ] 本地构建成功
- [ ] 单元测试通过
- [ ] DEB 包生成成功
- [ ] RPM 包生成成功
- [ ] 安装验证通过

---

**最后更新**: 2026-03-26