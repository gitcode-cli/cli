# gc-regression

使用 `gc` 完成最小真实命令回归验证。

## 触发场景

- 修改 `gc` 命令行为后做回归
- 升级或变更认证逻辑后做快速验证
- 发布前做最小真实命令检查

## 最小回归建议

```bash
# 构建
go build -o ./gc ./cmd/gc

# 查看版本
./gc version

# 查看认证状态
./gc auth status

# 查看仓库
./gc repo view infra-test/gctest1

# 列出 issue
./gc issue list -R infra-test/gctest1 --state open
```

## 如有现成脚本

如果目标项目本身提供了回归脚本，优先复用项目已有脚本。

在 gitcode-cli 仓库内，当前可直接运行：

```bash
./scripts/regression-core.sh
```

## 使用约束

- 真实命令回归应使用安全测试仓库
- 写操作命令应明确区分是否允许污染测试数据
- 回归结果要记录通过项和未覆盖项
