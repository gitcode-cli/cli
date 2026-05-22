package cmdutil

import (
	"bytes"
	"os"
)

var utf8BOM = []byte{0xef, 0xbb, 0xbf}

// ReadTextFile reads a user-provided text file and strips a UTF-8 BOM when present.
func ReadTextFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes.TrimPrefix(content, utf8BOM)), nil
}
