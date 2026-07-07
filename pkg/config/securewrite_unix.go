//go:build !windows

package config

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

// secureWriteFile writes data to path with hardened permission semantics.
// On Unix, opens with O_NOFOLLOW to atomically reject a symlink at the final
// path component, eliminating the Lstat→WriteFile TOCTOU window that would
// allow credential redirection. The returned fd is used for both Write and
// Chmod so the permission hardening targets the exact inode that received
// the data. Threat model assumptions documented in
// .loop/deliveries/issue-400-design.md.
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
