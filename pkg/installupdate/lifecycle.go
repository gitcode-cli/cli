// Package installupdate manages the npm-bootstrap updater lifecycle.
package installupdate

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gitcode.com/gitcode-cli/cli/pkg/config"
)

const updateTTL = 24 * time.Hour

// Manifest describes an npm-bootstrap installation.
type Manifest struct {
	Distribution string `json:"distribution"`
	Version      string `json:"version"`
	TargetDir    string `json:"targetDir"`
	Node         string `json:"node"`
	NPM          string `json:"npm"`
	Helper       string `json:"helper"`
	path         string
}

type updateState struct {
	NextCheck   int64         `json:"nextCheck,omitempty"`
	NoticeShown bool          `json:"noticeShown,omitempty"`
	Summary     *stateSummary `json:"summary,omitempty"`
}

type stateSummary struct {
	Message string `json:"message"`
	Shown   bool   `json:"shown"`
}

// LoadBootstrapManifest returns the manifest adjacent to the running binary.
func LoadBootstrapManifest() (*Manifest, error) {
	binary := os.Getenv("GITCODE_CLI_BINARY")
	if binary == "" {
		var err error
		binary, err = os.Executable()
		if err != nil {
			return nil, err
		}
	}
	path := filepath.Join(filepath.Dir(binary), ".gitcode-install.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	if manifest.Distribution != "npm-bootstrap" || manifest.TargetDir == "" || manifest.Helper == "" {
		return nil, fmt.Errorf("invalid npm-bootstrap manifest")
	}
	manifest.path = path
	return &manifest, nil
}

// ManifestPath returns the source path for the loaded manifest.
func (m *Manifest) ManifestPath() string {
	return m.path
}

// AfterCommand shows pending status and schedules a due bootstrap update.
// Errors are intentionally reported to stderr and never change command status.
func AfterCommand(cfg config.Config, errOut io.Writer, noUpdate, noInteractive bool) {
	manifest, err := LoadBootstrapManifest()
	if err != nil {
		return
	}
	statePath := StatePath()
	state := readState(statePath)
	if state.Summary != nil && !state.Summary.Shown {
		fmt.Fprintln(errOut, state.Summary.Message)
		state.Summary.Shown = true
		writeState(statePath, state)
	}
	if disabled(cfg, noUpdate, noInteractive) {
		return
	}
	if !state.NoticeShown {
		fmt.Fprintln(errOut, "GitCode CLI installed by npm bootstrap checks daily and automatically applies stable updates after commands.")
		fmt.Fprintln(errOut, `Set "gitcode config set update.mode notify|off" or GC_NO_UPDATE_CHECK=1 to change this behavior.`)
		state.NoticeShown = true
		writeState(statePath, state)
	}
	if state.NextCheck > time.Now().UnixMilli() {
		return
	}
	if err := StartDetached(manifest, false); err != nil {
		fmt.Fprintf(errOut, "update check skipped: %v\n", err)
	}
}

// StatePath returns the cross-platform npm update state file.
func StatePath() string {
	if dir := os.Getenv("GC_STATE_DIR"); dir != "" {
		return filepath.Join(dir, "update-state.json")
	}
	if local := os.Getenv("LOCALAPPDATA"); runtime.GOOS == "windows" && local != "" {
		return filepath.Join(local, "gitcode-cli", "update-state.json")
	}
	home, _ := os.UserHomeDir()
	root := os.Getenv("XDG_STATE_HOME")
	if root == "" {
		root = filepath.Join(home, ".local", "state")
	}
	return filepath.Join(root, "gitcode-cli", "update-state.json")
}

// StartDetached runs the copied bootstrap helper after this process exits.
func StartDetached(manifest *Manifest, force bool) error {
	node := manifest.Node
	if node == "" {
		node = "node"
	}
	args := []string{
		manifest.Helper,
		"--background",
		"--parent-pid", strconv.Itoa(os.Getpid()),
		"--manifest", manifest.ManifestPath(),
	}
	if force {
		args = append(args, "--force")
	}
	cmd := exec.Command(node, args...)
	cmd.Env = updaterEnvironment()
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Process.Release()
}

// RunCheck executes an explicit foreground update check.
func RunCheck(manifest *Manifest, jsonOutput bool, out, errOut io.Writer) error {
	node := manifest.Node
	if node == "" {
		node = "node"
	}
	args := []string{manifest.Helper, "--check", "--manifest", manifest.ManifestPath()}
	if jsonOutput {
		args = append(args, "--json")
	}
	cmd := exec.Command(node, args...)
	cmd.Env = updaterEnvironment()
	cmd.Stdout = out
	cmd.Stderr = errOut
	return cmd.Run()
}

func disabled(cfg config.Config, noUpdate, noInteractive bool) bool {
	if noUpdate || noInteractive || truthy(os.Getenv("GC_NO_UPDATE_CHECK")) || truthy(os.Getenv("CI")) {
		return true
	}
	mode := strings.ToLower(os.Getenv("GC_UPDATE_MODE"))
	if mode == "" && cfg != nil {
		mode, _ = cfg.Get("gitcode.com", "update.mode")
	}
	return mode == "off"
}

func truthy(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes":
		return true
	default:
		return false
	}
}

func updaterEnvironment() []string {
	blocked := map[string]struct{}{
		"GC_TOKEN":        {},
		"GITCODE_TOKEN":   {},
		"NPM_TOKEN":       {},
		"NODE_AUTH_TOKEN": {},
	}
	clean := make([]string, 0, len(os.Environ())+1)
	for _, item := range os.Environ() {
		key, _, _ := strings.Cut(item, "=")
		if _, ok := blocked[strings.ToUpper(key)]; !ok {
			clean = append(clean, item)
		}
	}
	return append(clean, "GC_NO_UPDATE_CHECK=1")
}

func readState(path string) updateState {
	data, err := os.ReadFile(path)
	if err != nil {
		return updateState{}
	}
	var state updateState
	if json.Unmarshal(data, &state) != nil {
		return updateState{}
	}
	return state
}

func writeState(path string, state updateState) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return
	}
	temp := path + ".tmp"
	if os.WriteFile(temp, append(data, '\n'), 0o600) != nil {
		return
	}
	_ = os.Remove(path)
	_ = os.Rename(temp, path)
}

// DueAt returns the next automatic-check time used by tests and diagnostics.
func DueAt(now time.Time) int64 {
	return now.Add(updateTTL).UnixMilli()
}
