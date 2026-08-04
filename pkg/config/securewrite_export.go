package config

import "os"

// SecureWriteFile is the exported hardened-write helper for sibling packages
// (e.g. pkg/larkcli) that persist credential-adjacent state under the gc config
// directory. It delegates to the platform-specific secureWriteFile so that
// symlink redirection and permission hardening stay consistent across the
// codebase.
func SecureWriteFile(path string, data []byte, perm os.FileMode) error {
	return secureWriteFile(path, data, perm)
}
