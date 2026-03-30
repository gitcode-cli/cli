# gc-release

使用 `gc` 完成 GitCode release 相关操作。

## 触发场景

- 创建 release
- 查看或列出 release
- 上传 release 资产
- 下载 release 资产

## 常用命令

```bash
# 创建 release
gc release create v1.0.0 -R owner/repo --title "Version 1.0.0" --notes "Release notes"

# 列出 release
gc release list -R owner/repo

# 查看 release
gc release view v1.0.0 -R owner/repo

# 上传资产
gc release upload v1.0.0 file.tar.gz -R owner/repo

# 下载资产
gc release download v1.0.0 -R owner/repo
gc release download v1.0.0 app.zip -R owner/repo -o ./downloads/
```

## 使用约束

- `release create` 应显式提供 `--notes`
- 上传前应确认 release 已存在且 tag 正确
- 下载或上传前应确认文件名和版本一致
