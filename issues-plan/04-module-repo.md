# 仓库模块需求 (repo)

本文档详细描述 gitcode-cli 仓库模块的功能需求、验收标准和 API 映射。

## 模块概述

仓库模块提供 GitCode 仓库的管理功能，包括克隆、创建、Fork、查看、列出和删除仓库。

### 命令结构

```
gc repo <command>

Commands:
  clone    Clone a repository locally
  create   Create a new repository
  fork     Fork a repository
  view     View a repository
  list     List repositories
  delete   Delete a repository
```

### 仓库标识格式

| 格式 | 示例 | 描述 |
|------|------|------|
| OWNER/REPO | `owner/repo` | 简短格式 |
| URL | `https://gitcode.com/owner/repo` | 完整 URL |
| 当前目录 | - | 自动检测 |

---

## REPO-001: repo clone - 克隆仓库

### 功能描述

克隆 GitCode 仓库到本地。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --git-protocol | | string | | Git 协议 (https/ssh) |
| --depth | -d | int | 0 | 浅克隆深度 |
| --branch | -b | string | | 克隆指定分支 |
| --single-branch | | bool | false | 只克隆指定分支 |
| --recursive | -r | bool | false | 递归克隆子模块 |
| --quiet | -q | bool | false | 静默模式 |

### 使用示例

```bash
# 克隆仓库
gc repo clone owner/repo

# 使用 SSH 协议
gc repo clone owner/repo --git-protocol ssh

# 克隆到指定目录
gc repo clone owner/repo my-project

# 浅克隆
gc repo clone owner/repo --depth 1

# 克隆指定分支
gc repo clone owner/repo --branch develop

# 克隆子模块
gc repo clone owner/repo --recursive
```

### 验收标准

- [ ] 支持 OWNER/REPO 格式
- [ ] 支持完整 URL 格式
- [ ] 支持 HTTPS 和 SSH 协议
- [ ] 正确处理浅克隆参数
- [ ] 正确处理分支参数
- [ ] 显示克隆进度

### API 端点

此命令主要使用 Git 协议，不需要直接调用 GitCode API。

---

## REPO-002: repo create - 创建仓库

### 功能描述

在 GitCode 上创建新仓库。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --description | -d | string | | 仓库描述 |
| --homepage | | string | | 主页 URL |
| --private | | bool | false | 创建私有仓库 |
| --public | | bool | false | 创建公开仓库 |
| --clone | | bool | false | 创建后克隆 |
| --gitignore | | string | | .gitignore 模板 |
| --license | | string | | 开源许可证模板 |
| --template | | string | | 模板仓库 |

### 使用示例

```bash
# 创建公开仓库
gc repo create my-repo --public

# 创建私有仓库
gc repo create my-repo --private

# 创建并克隆
gc repo create my-repo --public --clone

# 使用模板创建
gc repo create my-repo --template owner/template-repo

# 添加 .gitignore
gc repo create my-repo --gitignore Go
```

### 验收标准

- [ ] 正确创建公开/私有仓库
- [ ] 正确设置描述和主页
- [ ] 支持模板仓库
- [ ] 支持自动克隆
- [ ] 显示创建成功的仓库 URL

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/user/repos` | POST | 创建个人仓库 |
| `/api/v5/orgs/{org}/repos` | POST | 创建组织仓库 |
| `/api/v5/repos/{owner}/{repo}/generate` | POST | 从模板创建 |

### 测试用例映射

- 参考 `gc-api-doc/test/test_repositories.py`

---

## REPO-003: repo fork - Fork 仓库

### 功能描述

Fork 其他用户的仓库到自己账户或组织。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --clone | -c | bool | false | Fork 后克隆 |
| --remote | -r | bool | false | 添加为远程 |
| --remote-name | | string | origin | 远程名称 |
| --organization | -o | string | | Fork 到组织 |
| --fork-name | | string | | Fork 仓库名称 |

### 使用示例

```bash
# Fork 仓库
gc repo fork owner/repo

# Fork 并克隆
gc repo fork owner/repo --clone

# Fork 到组织
gc repo fork owner/repo --organization my-org

# Fork 并添加远程
gc repo fork owner/repo --remote
```

### 验收标准

- [ ] 正确 Fork 仓库
- [ ] 支持 Fork 到组织
- [ ] 支持自动克隆
- [ ] 支持添加远程
- [ ] 显示 Fork 后的仓库 URL

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}/forks` | POST | Fork 仓库 |

### 测试用例映射

- 参考 `gc-api-doc/test/test_repositories.py`

---

## REPO-004: repo view - 查看仓库

### 功能描述

查看仓库的详细信息。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --web | -w | bool | false | 在浏览器中打开 |
| --json | | []string | | 输出指定 JSON 字段 |

### 使用示例

```bash
# 查看当前仓库
gc repo view

# 查看指定仓库
gc repo view owner/repo

# 在浏览器中打开
gc repo view --web

# JSON 输出
gc repo view --json name,description,url
```

### 输出示例

```
owner/repo
Repository description

  URL: https://gitcode.com/owner/repo
  Visibility: Public
  Stars: 100
  Forks: 10
  Issues: 5
  Default branch: main

  Clone URLs:
    HTTPS: https://gitcode.com/owner/repo.git
    SSH: git@gitcode.com:owner/repo.git
```

### 验收标准

- [ ] 正确显示仓库名称和描述
- [ ] 显示可见性（公开/私有）
- [ ] 显示统计信息
- [ ] 显示克隆 URL
- [ ] 支持 JSON 输出

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}` | GET | 获取仓库信息 |

### 测试用例映射

- 参考 `gc-api-doc/test/test_repositories.py`

---

## REPO-005: repo list - 列出仓库

### 功能描述

列出用户或组织的仓库。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --owner | -o | string | | 指定所有者 |
| --limit | -L | int | 30 | 最大数量 |
| --visibility | | string | | 可见性过滤 (public/private/all) |
| --sort | | string | updated | 排序方式 |
| --json | | []string | | JSON 输出 |

### 使用示例

```bash
# 列出自己的仓库
gc repo list

# 列出指定用户的仓库
gc repo list --owner username

# 列出组织的仓库
gc repo list --owner my-org

# 限制数量
gc repo list --limit 10
```

### 输出示例

```
NAME                DESCRIPTION              VISIBILITY  UPDATED
owner/repo-1        First repository         public      2 days ago
owner/repo-2        Second repository        private     1 week ago
```

### 验收标准

- [ ] 正确列出仓库
- [ ] 支持用户/组织过滤
- [ ] 支持可见性过滤
- [ ] 支持分页
- [ ] 支持 JSON 输出

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/user/repos` | GET | 列出自己的仓库 |
| `/api/v5/users/{username}/repos` | GET | 列出用户的仓库 |
| `/api/v5/orgs/{org}/repos` | GET | 列出组织的仓库 |

### 测试用例映射

- 参考 `gc-api-doc/test/test_repositories.py`

---

## REPO-006: repo delete - 删除仓库

### 功能描述

删除指定的仓库（危险操作）。

### 命令参数

| 参数 | 短参数 | 类型 | 默认值 | 说明 |
|------|--------|------|--------|------|
| --yes | -y | bool | false | 跳过确认 |

### 使用示例

```bash
# 删除仓库（需要确认）
gc repo delete owner/repo

# 跳过确认
gc repo delete owner/repo --yes
```

### 验收标准

- [ ] 要求输入仓库名确认
- [ ] 正确删除仓库
- [ ] 显示删除成功的确认信息
- [ ] 支持 --yes 跳过确认

### API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}` | DELETE | 删除仓库 |

### 测试用例映射

- 参考 `gc-api-doc/test/test_repositories.py`

---

## 相关文档

- [gc-api-doc/doc/03-repositories.md](../../gc-api-doc/doc/03-repositories.md)
- [gc-api-doc/test/test_repositories.py](../../gc-api-doc/test/test_repositories.py)

---

**最后更新**: 2026-03-22