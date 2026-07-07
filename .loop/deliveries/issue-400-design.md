## 方案设计: Issue #400

### 方案
新增 `secureWriteFile` 包级私有函数，替代 `writeAuthState`/`writeConfigState` 中的裸 `os.WriteFile`，提供三点加固：

| 加固点 | 机制 | 防护威胁 |
|--------|------|----------|
| 原子拒绝符号链接 | Unix `os.OpenFile(O_NOFOLLOW)` 打开时原子拒绝 symlink，消除 Lstat→WriteFile TOCTOU；Windows 保留 `os.Lstat` 检测（无 `O_NOFOLLOW` 等价物，依赖 ACL 缓解） | 凭证重定向攻击（auth.json 被替换为软链，写入覆盖目标/泄露内容）+ TOCTOU 竞态绕过 |
| fd 权限硬化 | 对打开的 fd `f.Chmod(perm)`，作用在写入的 inode 而非 path（无 path 竞态） | umask 宽松导致新建文件权限过宽 + 既有文件权限未收紧 + Chmod 跟随 symlink 的 TOCTOU |
| 目录权限硬化 | `writeAuthState`/`writeConfigState` 在 `MkdirAll(configDir, 0o700)` 后补 `os.Chmod(configDir, 0o700)` | 配置目录权限宽松（与 `Write()` 的 `config.go:212` 风格一致） |

### 实现位置
跨平台分离，按 build tag 两个文件：

`pkg/config/securewrite_unix.go`（`//go:build !windows`）：
```go
func secureWriteFile(path string, data []byte, perm os.FileMode) error {
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC | syscall.O_NOFOLLOW
	f, err := os.OpenFile(path, flags, perm)
	if err != nil {
		if errors.Is(err, syscall.ELOOP) {
			return fmt.Errorf("refusing to write to symlink: %s", path)
		}
		return fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer f.Close()
	if err := f.Chmod(perm); err != nil {
		return fmt.Errorf("failed to chmod %s: %w", path, err)
	}
	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	return nil
}
```

`pkg/config/securewrite_windows.go`（`//go:build windows`）：保留原 `Lstat`→`WriteFile`→`Chmod` 三步（非原子，但 Windows reparse-point 语义与 ACL 模型提供独立缓解）。

### 调用点
| 位置 | 原 | 修复后 |
|------|----|----|
| auth_config.go (writeAuthState) | `os.MkdirAll` + `os.WriteFile(c.authStatePath(), data, 0o600)` | `os.MkdirAll` + `os.Chmod(configDir, 0o700)` + `secureWriteFile(c.authStatePath(), data, 0o600)` |
| config.go (writeConfigState) | `os.MkdirAll` + `os.WriteFile(c.configStatePath(), data, 0o600)` | `os.MkdirAll` + `os.Chmod(configDir, 0o700)` + `secureWriteFile(c.configStatePath(), data, 0o600)` |

### 与已合入安全改动的关系
| Issue/PR | 防护方向 | 关系 |
|----------|----------|------|
| #317 (PR !336) warn before printing auth token | 读路径 — token 输出 stdout | 正交 |
| #358 (PR !335) validate URL scheme in browser.Open | URL scheme 注入 | 正交 |
| !343 require interactive token disclosure (75acec0) | 读路径 — token 交互式确认 | 正交 |
| #400 (本 issue) secureWriteFile | 写路径 — 凭证文件落盘 | 与上述无重叠 |

### 风险
- **symlink 误伤**：若用户合法地用 symlink 管理 auth.json（如 dotfiles 管理），新行为会拒绝写入。但 auth.json 含敏感凭证，不应作为 symlink 暴露，拒绝是安全正确的。错误信息明确提示 "refusing to write to symlink"，用户可删除 symlink 后重试。
- **Windows 兼容**：symlink 语义与 Unix 不同，权限模型不同。测试用 `runtime.GOOS == "windows"` t.Skip；生产代码 `os.Lstat` 在 Windows 上对 symlink 检测行为一致（reparse point），`os.Chmod` 在 Windows 上基本 no-op，不影响。
- **既有文件覆盖**：Unix 实现用 `O_TRUNC` 覆盖写，symlink 在打开时被拒，不会泄露旧内容。

### 威胁模型假设
- **防护范围**：末段 path 为符号链接的静态替换 + TOCTOU 竞态（主要威胁，已用 `O_NOFOLLOW` 原子消除）
- **依赖目录 0700**：`configDir` 已 `Chmod 0o700`，跨用户攻击者无法改父目录，无法竞态；同用户攻击者本可读 `auth.json`，竞态无增益
- **不防护**：
  - 父目录组件为符号链接（需 `openat2(RESOLVE_NO_SYMLINKS)`，Linux 5.6+，非跨平台，过度工程化）
  - hard link 攻击（`O_NOFOLLOW` 不防，但被 `configDir 0700` 缓解，同用户场景无意义）
  - Windows 的 TOCTOU（`O_NOFOLLOW` 不可移植，Windows 保留 `Lstat`，依赖 ACL 模型）
- **结论**：当前方案消除主要 TOCTOU，残留威胁被独立的目录权限缓解，是跨平台条件下的合理最优解；100% 彻底需 `openat2`，成本与可利用性不匹配

### 系统测试限制
自动化系统测试在当前架构下不可行：
- `writeConfigState` 无 CLI 命令入口（`pkg/cmd/` 下无 `config` 命令，无人调用 `cfg.Set()`/`cfg.Write()`）
- `writeAuthState` 经 `gc auth login --with-token` 触发，但 `login.go:137` 在写前强制 `api.VerifyToken`（需真实 token + 网络），项目规则禁止 AI/脚本向 stdin 提供真实 token

因此采用「UT 充分覆盖 + 人工真实命令验证」策略。

## 人工真实命令验证清单（TTY 执行，AI 不参与 token 输入）
```bash
# 前置：./gc auth status 确认认证可用；token 须与 infra-test 关联，不得用个人/生产 token

# 场景1：symlink 拒绝
mkdir -p ~/.config/gc
target=$(mktemp)
ln -sf "$target" ~/.config/gc/auth.json
echo "<真实token>" | ./gc auth login --with-token --hostname gitcode.com
# 期望：失败，stderr 含 "refusing to write to symlink"
cat "$target"   # 应为空（未被写入）

# 场景2：权限硬化
rm -f ~/.config/gc/auth.json
touch ~/.config/gc/auth.json && chmod 644 ~/.config/gc/auth.json
echo "<真实token>" | ./gc auth login --with-token --hostname gitcode.com
# 期望：成功（✓ Logged in as ...）
stat -c '%a' ~/.config/gc/auth.json   # Linux: 600；macOS: stat -f '%Lp'

# 清理：用正常方式重新 login 恢复认证状态
```
