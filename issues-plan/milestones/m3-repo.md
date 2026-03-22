# 里程碑 3: 仓库功能

## 概述

实现 GitCode 仓库管理功能，包括克隆、创建、查看、Fork 和删除仓库。

**预计工期**: 1 周

**依赖**: 里程碑 2 (认证功能)

**目标**: 用户能够通过命令行管理 GitCode 仓库

---

## 任务清单

### REPO-001: 仓库克隆

**优先级**: P0

**任务描述**:

实现 `gc repo clone` 命令，克隆远程仓库。

**文件**:

```
pkg/cmd/repo/clone/clone.go
pkg/cmd/repo/clone/http.go
git/clone.go
```

**功能**:

- 支持 `OWNER/REPO` 格式
- 支持 URL 格式
- 支持 HTTPS/SSH 协议选择
- 支持 Git 参数传递

**验收标准**:

- [ ] `gc repo clone owner/repo` 克隆成功
- [ ] `gc repo clone https://gitcode.com/owner/repo` URL 克隆
- [ ] `gc repo clone owner/repo --depth 1` 浅克隆
- [ ] `gc repo clone owner/repo --branch main` 指定分支
- [ ] 自动检测并使用配置的协议

**示例**:

```bash
# 基本克隆
$ gc repo clone owner/repo
Cloning into 'repo'...
✓ Cloned repository owner/repo

# 指定目录
$ gc repo clone owner/repo my-dir
✓ Cloned repository owner/repo to my-dir

# 浅克隆
$ gc repo clone owner/repo --depth 1
✓ Cloned repository owner/repo (shallow)

# 使用 SSH
$ gc repo clone owner/repo --git-protocol ssh
✓ Cloned repository owner/repo via SSH
```

---

### REPO-002: 仓库创建

**优先级**: P0

**任务描述**:

实现 `gc repo create` 命令，创建新仓库。

**文件**:

```
pkg/cmd/repo/create/create.go
api/queries_repo.go
```

**功能**:

- 创建公开/私有仓库
- 设置描述和主页
- 支持模板仓库
- 创建后自动克隆

**验收标准**:

- [ ] 创建公开仓库
- [ ] 创建私有仓库
- [ ] 设置描述
- [ ] 创建后自动克隆
- [ ] 显示仓库 URL

**示例**:

```bash
# 交互式创建
$ gc repo create
? Repository name: my-repo
? Description: My repository
? Visibility: Public
? Add a README? Yes
✓ Created repository owner/my-repo
? Clone the repository? Yes
✓ Cloned repository owner/my-repo

# 命令行创建
$ gc repo create my-repo --public --description "My repo"
✓ Created repository owner/my-repo
https://gitcode.com/owner/my-repo

# 创建并克隆
$ gc repo create my-repo --private --clone
✓ Created and cloned repository owner/my-repo
```

---

### REPO-003: 仓库 Fork

**优先级**: P1

**任务描述**:

实现 `gc repo fork` 命令，Fork 远程仓库。

**文件**:

```
pkg/cmd/repo/fork/fork.go
api/queries_repo.go
```

**功能**:

- Fork 到个人命名空间
- Fork 到组织
- 添加为远程

**验收标准**:

- [ ] Fork 到个人账户
- [ ] Fork 到指定组织
- [ ] 克隆 Fork 后的仓库
- [ ] 添加为远程仓库

**示例**:

```bash
# Fork 当前仓库
$ gc repo fork
✓ Forked repository owner/repo to myuser/repo

# Fork 并克隆
$ gc repo fork owner/repo --clone
✓ Forked and cloned repository owner/repo to myuser/repo

# Fork 到组织
$ gc repo fork owner/repo --org myorg
✓ Forked repository owner/repo to myorg/repo

# 添加为远程
$ gc repo fork --remote
✓ Forked repository and added as remote 'fork'
```

---

### REPO-004: 仓库查看

**优先级**: P1

**任务描述**:

实现 `gc repo view` 命令，查看仓库详情。

**文件**:

```
pkg/cmd/repo/view/view.go
api/queries_repo.go
```

**验收标准**:

- [ ] 显示仓库基本信息
- [ ] 显示 Star/Fork 数量
- [ ] 显示语言统计
- [ ] 在浏览器中打开

**示例**:

```bash
$ gc repo view owner/repo
owner/repo
   My awesome repository

   Language: Go
   Stars: 100
   Forks: 20
   Issues: 5
   Watchers: 10

   Default branch: main
   License: MIT

   https://gitcode.com/owner/repo
```

---

### REPO-005: 仓库列表

**优先级**: P1

**任务描述**:

实现 `gc repo list` 命令，列出仓库。

**文件**:

```
pkg/cmd/repo/list/list.go
api/queries_repo.go
```

**功能**:

- 列出个人仓库
- 列出组织仓库
- 过滤和排序

**验收标准**:

- [ ] 列出个人仓库
- [ ] 按类型过滤（owner/member）
- [ ] 按可见性过滤（public/private）
- [ ] 支持分页
- [ ] 支持 JSON 输出

**示例**:

```bash
# 列出个人仓库
$ gc repo list
NAME            VISIBILITY  UPDATED
my-repo         public      2 days ago
another-repo    private     1 week ago

# 列出组织仓库
$ gc repo list --org myorg

# 过滤
$ gc repo list --limit 10 --public
```

---

### REPO-006: 仓库删除

**优先级**: P2

**任务描述**:

实现 `gc repo delete` 命令，删除仓库。

**文件**:

```
pkg/cmd/repo/delete/delete.go
api/queries_repo.go
```

**验收标准**:

- [ ] 需要确认才能删除
- [ ] 输入仓库名确认
- [ ] 显示删除成功

**示例**:

```bash
$ gc repo delete owner/repo
? This will delete repository owner/repo. Continue? Yes
? Type owner/repo to confirm: owner/repo
✓ Deleted repository owner/repo
```

---

## API 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/api/v5/repos/{owner}/{repo}` | GET | 获取仓库 |
| `/api/v5/repos/{owner}/{repo}` | POST | 创建仓库 |
| `/api/v5/repos/{owner}/{repo}` | DELETE | 删除仓库 |
| `/api/v5/repos/{owner}/{repo}/forks` | POST | Fork 仓库 |
| `/api/v5/user/repos` | GET | 列出用户仓库 |
| `/api/v5/orgs/{org}/repos` | GET | 列出组织仓库 |

---

## 依赖关系

```
REPO-001 (Clone) ─┐
                  │
REPO-002 (Create) ├─→ API Client
                  │
REPO-003 (Fork) ──┤
                  │
REPO-004 (View) ──┤
                  │
REPO-005 (List) ──┤
                  │
REPO-006 (Delete)─┘
```

---

## 完成标准

里程碑 M3 完成需满足：

1. ✅ `gc repo clone` 正确克隆仓库
2. ✅ `gc repo create` 创建仓库成功
3. ✅ `gc repo fork` Fork 仓库成功
4. ✅ `gc repo view` 显示仓库信息
5. ✅ `gc repo list` 列出仓库
6. ✅ 单元测试覆盖率 ≥ 80%

---

## 测试用例

### 单元测试

```bash
go test ./pkg/cmd/repo/... -v
```

### 集成测试

```bash
go test -tags=integration ./pkg/cmd/repo/... -v
```

### 手动测试清单

- [ ] 克隆公开仓库
- [ ] 克隆私有仓库
- [ ] 创建公开仓库
- [ ] 创建私有仓库
- [ ] Fork 仓库
- [ ] 查看仓库信息
- [ ] 列出仓库
- [ ] 删除仓库（需确认）

---

## 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| Git 命令失败 | 中 | 完善错误处理 |
| 大仓库克隆慢 | 低 | 支持 --depth |
| API 限流 | 中 | 实现重试机制 |

---

## 相关文档

- [issues-plan/04-module-repo.md](../04-module-repo.md)
- [gc-api-doc/doc/03-repositories.md](../../../gc-api-doc/doc/03-repositories.md)
- [gc-api-doc/test/test_repositories.py](../../../gc-api-doc/test/test_repositories.py)

---

**最后更新**: 2026-03-22