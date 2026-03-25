# 标签和里程碑命令

## 标签命令 (label)

### label list - 列出标签

```bash
gc label list -R infra-test/gctest1
```

### label create - 创建标签

```bash
gc label create "bug" -R infra-test/gctest1 --color "#ff0000" --description "Bug report"
```

### label delete - 删除标签

```bash
gc label delete bug -R infra-test/gctest1
```

---

## 里程碑命令 (milestone)

### milestone list - 列出里程碑

```bash
gc milestone list -R infra-test/gctest1
```

### milestone create - 创建里程碑

```bash
gc milestone create "v1.0" -R infra-test/gctest1 --description "First release"
```

### milestone view - 查看里程碑

```bash
gc milestone view 1 -R infra-test/gctest1
```

---

## 其他命令

### version - 显示版本

```bash
gc version
```

### help - 帮助

```bash
# 显示帮助
gc help

# 显示命令帮助
gc help issue
gc help issue create
```