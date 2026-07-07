//go:build windows

package config

import (
	"fmt"
	"os"
)

// secureWriteFile writes data to path with hardened permission semantics.
// On Windows, O_NOFOLLOW is not portable; Lstat rejects symlinks (non-atomic,
// but Windows reparse-point semantics differ and the configDir ACL model
// provides independent mitigation). Chmod is largely a no-op on Windows but
// kept for parity with the Unix path. Threat model assumptions documented in
// .loop/deliveries/issue-400-design.md.
func secureWriteFile(path string, data []byte, perm os.FileMode) error {
	if li, err := os.Lstat(path); err == nil && li.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("refusing to write to symlink: %s", path)
	}
	if err := os.WriteFile(path, data, perm); err != nil {
		return err
	}
	return os.Chmod(path, perm)
}
