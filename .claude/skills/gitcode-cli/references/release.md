# Release 命令 (release)

> 本文档是 Claude 参考层，不是命令行为真相源。
> Release 命令行为以 `docs/COMMANDS.md`、`docs/PACKAGING.md` 和 `spec/release-process.md` 为准。

## release create - 创建 Release

```bash
# 创建 Release（建议包含 --notes 参数）
gc release create v1.0.0 -R infra-test/gctest1 --title "Version 1.0.0" --notes "Release notes"

# 创建预发布 Release
gc release create v1.0.0-beta -R infra-test/gctest1 --title "v1.0.0 Beta" --notes "Beta release" --prerelease

# 创建草稿 Release
gc release create v1.0.0 -R infra-test/gctest1 --title "v1.0.0" --notes "Draft" --draft

# 指定目标分支
gc release create v1.0.0 -R infra-test/gctest1 --title "v1.0.0" --notes "Release" --target main
```

> **注意**: `--notes` 参数是必需的，不带此参数可能返回 400 错误。

## release list - 列出 Releases

```bash
gc release list -R infra-test/gctest1
```

## release view - 查看 Release

```bash
# 查看 Release 详情
gc release view v1.0.0 -R infra-test/gctest1

# 在浏览器中打开
gc release view v1.0.0 -R infra-test/gctest1 --web
```

## release upload - 上传资产

```bash
# 上传单个文件
gc release upload v1.0.0 app.zip -R infra-test/gctest1

# 上传多个文件
gc release upload v1.0.0 app.zip checksum.txt -R infra-test/gctest1
```

## release download - 下载资产

```bash
# 下载所有资产到当前目录
gc release download v1.0.0 -R infra-test/gctest1

# 下载到指定目录
gc release download v1.0.0 -R infra-test/gctest1 -o ./downloads/

# 下载指定文件
gc release download v1.0.0 app.zip -R infra-test/gctest1
```

## release edit - 编辑 Release

```bash
# 修改标题
gc release edit v1.0.0 --title "New title" -R infra-test/gctest1

# 修改说明
gc release edit v1.0.0 --notes "New release notes" -R infra-test/gctest1
```

## release delete - 删除 Release

```bash
gc release delete v1.0.0 -R infra-test/gctest1
```
