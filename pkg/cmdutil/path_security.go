package cmdutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// MaxInputFileSize is the maximum size allowed for user-provided input files
// read via flags such as --input, --body-file, and --custom-fields-file.
// 10 MiB. The limit is a secondary defense that bounds the amount of data a
// single flag invocation can exfiltrate; the primary defense is the
// current-directory containment check in ReadInputFile (see issue #397).
const MaxInputFileSize int64 = 10 * 1024 * 1024

// ReadInputFile reads a user-provided input file for flags whose content is
// sent to the GitCode API (e.g. --input, --body-file, --custom-fields-file).
//
// It applies two defenses against path-traversal-assisted data exfiltration:
//
//  1. The path is cleaned and resolved to an absolute path, then verified to
//     be inside the current working directory. This rejects traversal such as
//     "../../../etc/passwd" and absolute paths outside the cwd. Users may
//     still read any file inside the cwd (including subdirectories).
//  2. The file size is capped at MaxInputFileSize to bound exfiltration volume.
//
// Unlike safeOutputPath (release/download), this does not restrict to a
// caller-supplied directory: it uses the process cwd, because these flags are
// invoked interactively from the user's working directory.
func ReadInputFile(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("input file path must not be empty")
	}
	cleaned := filepath.Clean(path)

	absPath, err := filepath.Abs(cleaned)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve input path: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	rel, err := filepath.Rel(cwd, absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to validate input path: %w", err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return nil, fmt.Errorf("input file path must be within the current directory: %s", path)
	}

	f, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("input path is a directory, not a file: %s", path)
	}

	data, err := io.ReadAll(io.LimitReader(f, MaxInputFileSize+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > MaxInputFileSize {
		return nil, fmt.Errorf("input file %s exceeds size limit of %d bytes", path, MaxInputFileSize)
	}
	return data, nil
}
