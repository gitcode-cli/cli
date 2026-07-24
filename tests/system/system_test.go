//go:build system

package system_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestReadScripts(t *testing.T) {
	root := repoRoot(t)
	gcBin := buildOrUseBinary(t, root)
	readRepo := getenvDefault("GC_SYSTEM_REPO", "infra-test/gctest1")
	writeRepo := getenvDefault("GC_SYSTEM_WRITE_REPO", readRepo)

	requireInfraRepo(t, "GC_SYSTEM_REPO", readRepo)
	requireInfraRepo(t, "GC_SYSTEM_WRITE_REPO", writeRepo)

	testscript.Run(t, testscript.Params{
		Dir:       "testdata/read",
		Condition: systemCondition,
		Setup: func(env *testscript.Env) error {
			setupEnv(t, env, root, gcBin, readRepo, writeRepo)
			return nil
		},
		Cmds: systemCmds(),
	})
}

func TestWriteScripts(t *testing.T) {
	if os.Getenv("GC_SYSTEM_WRITE") != "1" {
		t.Skip("set GC_SYSTEM_WRITE=1 to run write-path system scripts")
	}

	root := repoRoot(t)
	gcBin := buildOrUseBinary(t, root)
	readRepo := getenvDefault("GC_SYSTEM_REPO", "infra-test/gctest1")
	writeRepo := getenvDefault("GC_SYSTEM_WRITE_REPO", readRepo)

	requireInfraRepo(t, "GC_SYSTEM_REPO", readRepo)
	requireInfraRepo(t, "GC_SYSTEM_WRITE_REPO", writeRepo)

	testscript.Run(t, testscript.Params{
		Dir:       "testdata/write",
		Condition: systemCondition,
		Setup: func(env *testscript.Env) error {
			setupEnv(t, env, root, gcBin, readRepo, writeRepo)
			return nil
		},
		Cmds: systemCmds(),
	})
}

func systemCmds() map[string]func(ts *testscript.TestScript, neg bool, args []string) {
	return map[string]func(ts *testscript.TestScript, neg bool, args []string){
		"defer-delete-label": cmdDeferDeleteLabel,
		"defer-close-issue":  cmdDeferCloseIssue,
		"json-assert":        cmdJSONAssert,
		"json-value":         cmdJSONValue,
		"json-ok":            cmdJSONOK,
		"require-infra":      cmdRequireInfra,
		"stdout2env":         cmdStdout2Env,
		"unique-name":        cmdUniqueName,
	}
}

func setupEnv(t *testing.T, env *testscript.Env, root, gcBin, readRepo, writeRepo string) {
	t.Helper()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("resolve user home: %v", err)
	}

	env.Setenv("ROOT", root)
	env.Setenv("GC_BIN", gcBin)
	env.Setenv("SYSTEM_REPO", readRepo)
	env.Setenv("WRITE_REPO", writeRepo)
	env.Setenv(homeEnvName(), home)
	copyOptionalEnv(env, "APPDATA")
	copyOptionalEnv(env, "GC_CONFIG_DIR")
	copyOptionalEnv(env, "XDG_CONFIG_HOME")
	copyOptionalEnv(env, "GC_TOKEN")
	copyOptionalEnv(env, "GITCODE_TOKEN")
	copyOptionalEnv(env, "GC_SYSTEM_ASSIGNEE")
}

func buildOrUseBinary(t *testing.T, root string) string {
	t.Helper()

	if gcBin := os.Getenv("GC_BIN"); gcBin != "" {
		return gcBin
	}

	name := "gc"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	out := filepath.Join(t.TempDir(), name)
	cmd := exec.Command("go", "build", "-o", out, "./cmd/gc")
	cmd.Dir = root
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build gc: %v\n%s", err, output)
	}
	return out
}

func repoRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	root, err := filepath.Abs(filepath.Join(wd, "..", ".."))
	if err != nil {
		t.Fatalf("resolve repository root: %v", err)
	}
	return root
}

func getenvDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func copyOptionalEnv(env *testscript.Env, name string) {
	if value := os.Getenv(name); value != "" {
		env.Setenv(name, value)
	}
}

func systemCondition(condition string) (bool, error) {
	name, ok := strings.CutPrefix(condition, "env:")
	if !ok || name == "" {
		return false, fmt.Errorf("unknown condition %q", condition)
	}
	return os.Getenv(name) != "", nil
}

func requireInfraRepo(t *testing.T, name, repo string) {
	t.Helper()

	if err := validateInfraRepo(name, repo); err != nil {
		t.Fatal(err)
	}
}

func validateInfraRepo(name, repo string) error {
	if repo == "" {
		return fmt.Errorf("%s is required", name)
	}
	if !strings.HasPrefix(repo, "infra-test/") {
		return fmt.Errorf("%s must be infra-test/*, got %q", name, repo)
	}
	if repo == "infra-test/" || strings.Count(repo, "/") != 1 {
		return fmt.Errorf("%s must be an owner/repo path under infra-test, got %q", name, repo)
	}
	return nil
}

func homeEnvName() string {
	if runtime.GOOS == "windows" {
		return "USERPROFILE"
	}
	return "HOME"
}

func cmdRequireInfra(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 1 {
		ts.Fatalf("usage: require-infra repo")
	}
	err := validateInfraRepo("repo", args[0])
	if neg {
		if err == nil {
			ts.Fatalf("%q unexpectedly passed infra-test validation", args[0])
		}
		return
	}
	if err != nil {
		ts.Fatalf("%v", err)
	}
}

func cmdJSONOK(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 1 {
		ts.Fatalf("usage: json-ok file")
	}
	var value any
	err := json.Unmarshal([]byte(ts.ReadFile(args[0])), &value)
	if neg {
		if err == nil {
			ts.Fatalf("%s is valid JSON", args[0])
		}
		return
	}
	if err != nil {
		ts.Fatalf("%s is not valid JSON: %v", args[0], err)
	}
}

func cmdJSONAssert(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 3 {
		ts.Fatalf("usage: json-assert file path type")
	}

	value, err := parseJSONFile(ts, args[0])
	if err != nil {
		ts.Fatalf("%s is not valid JSON: %v", args[0], err)
	}
	actual, ok, err := lookupJSONPath(value, args[1])
	if err != nil {
		ts.Fatalf("invalid JSON path %q: %v", args[1], err)
	}
	matches := ok && jsonTypeMatches(actual, args[2])
	if neg {
		if matches {
			ts.Fatalf("%s %s unexpectedly matched type %s", args[0], args[1], args[2])
		}
		return
	}
	if !ok {
		ts.Fatalf("%s %s is missing", args[0], args[1])
	}
	if !matches {
		ts.Fatalf("%s %s has type %s, want %s", args[0], args[1], jsonType(actual), args[2])
	}
}

func cmdJSONValue(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 3 {
		ts.Fatalf("usage: json-value file path expected")
	}
	value, err := parseJSONFile(ts, args[0])
	if err != nil {
		ts.Fatalf("%s is not valid JSON: %v", args[0], err)
	}
	actual, ok, err := lookupJSONPath(value, args[1])
	if err != nil {
		ts.Fatalf("invalid JSON path %q: %v", args[1], err)
	}
	matches := ok && fmt.Sprint(actual) == args[2]
	if neg {
		if matches {
			ts.Fatalf("%s %s unexpectedly equals %q", args[0], args[1], args[2])
		}
		return
	}
	if !ok {
		ts.Fatalf("%s %s is missing", args[0], args[1])
	}
	if !matches {
		ts.Fatalf("%s %s = %q, want %q", args[0], args[1], fmt.Sprint(actual), args[2])
	}
}

func parseJSONFile(ts *testscript.TestScript, file string) (any, error) {
	var value any
	err := json.Unmarshal([]byte(ts.ReadFile(file)), &value)
	return value, err
}

func lookupJSONPath(value any, path string) (any, bool, error) {
	if path == "." || path == "" {
		return value, true, nil
	}
	for path != "" {
		switch {
		case strings.HasPrefix(path, "."):
			path = path[1:]
			key, rest := nextJSONKey(path)
			if key == "" {
				return nil, false, fmt.Errorf("empty object key")
			}
			object, ok := value.(map[string]any)
			if !ok {
				return nil, false, nil
			}
			value, ok = object[key]
			if !ok {
				return nil, false, nil
			}
			path = rest
		case strings.HasPrefix(path, "["):
			end := strings.Index(path, "]")
			if end < 0 {
				return nil, false, fmt.Errorf("missing closing ]")
			}
			index, err := parseJSONIndex(path[1:end])
			if err != nil {
				return nil, false, err
			}
			array, ok := value.([]any)
			if !ok || index < 0 || index >= len(array) {
				return nil, false, nil
			}
			value = array[index]
			path = path[end+1:]
		default:
			key, rest := nextJSONKey(path)
			if key == "" {
				return nil, false, fmt.Errorf("empty object key")
			}
			object, ok := value.(map[string]any)
			if !ok {
				return nil, false, nil
			}
			value, ok = object[key]
			if !ok {
				return nil, false, nil
			}
			path = rest
		}
	}
	return value, true, nil
}

func nextJSONKey(path string) (string, string) {
	nextDot := strings.Index(path, ".")
	nextBracket := strings.Index(path, "[")
	next := -1
	switch {
	case nextDot >= 0 && nextBracket >= 0:
		next = minInt(nextDot, nextBracket)
	case nextDot >= 0:
		next = nextDot
	case nextBracket >= 0:
		next = nextBracket
	}
	if next < 0 {
		return path, ""
	}
	return path[:next], path[next:]
}

func parseJSONIndex(value string) (int, error) {
	index, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid array index %q", value)
	}
	return index, nil
}

func jsonTypeMatches(value any, want string) bool {
	switch want {
	case "present":
		return true
	case "string":
		_, ok := value.(string)
		return ok
	case "nonempty-string":
		text, ok := value.(string)
		return ok && text != ""
	case "number":
		_, ok := value.(float64)
		return ok
	case "bool":
		_, ok := value.(bool)
		return ok
	case "object":
		_, ok := value.(map[string]any)
		return ok
	case "array":
		_, ok := value.([]any)
		return ok
	case "null":
		return value == nil
	default:
		return false
	}
}

func jsonType(value any) string {
	switch value.(type) {
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "bool"
	case map[string]any:
		return "object"
	case []any:
		return "array"
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%T", value)
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func cmdStdout2Env(ts *testscript.TestScript, neg bool, args []string) {
	if neg {
		ts.Fatalf("unsupported: ! stdout2env")
	}
	if len(args) != 2 {
		ts.Fatalf("usage: stdout2env VAR regexp-with-one-capture")
	}
	re, err := regexp.Compile(args[1])
	ts.Check(err)
	matches := re.FindStringSubmatch(ts.ReadFile("stdout"))
	if len(matches) < 2 {
		ts.Fatalf("stdout did not match %q", args[1])
	}
	ts.Setenv(args[0], matches[1])
}

func cmdDeferCloseIssue(ts *testscript.TestScript, neg bool, args []string) {
	if neg {
		ts.Fatalf("unsupported: ! defer-close-issue")
	}
	if len(args) != 1 {
		ts.Fatalf("usage: defer-close-issue issue-number")
	}
	issueNumber := args[0]
	gcBin := ts.Getenv("GC_BIN")
	writeRepo := ts.Getenv("WRITE_REPO")
	ts.Defer(func() {
		_ = exec.Command(gcBin, "issue", "close", issueNumber, "-R", writeRepo, "--yes").Run()
	})
}

func cmdDeferDeleteLabel(ts *testscript.TestScript, neg bool, args []string) {
	if neg {
		ts.Fatalf("unsupported: ! defer-delete-label")
	}
	if len(args) != 1 {
		ts.Fatalf("usage: defer-delete-label label-name")
	}
	labelName := args[0]
	gcBin := ts.Getenv("GC_BIN")
	writeRepo := ts.Getenv("WRITE_REPO")
	ts.Defer(func() {
		_ = exec.Command(gcBin, "label", "delete", labelName, "-R", writeRepo, "--yes").Run()
	})
}

func cmdUniqueName(ts *testscript.TestScript, neg bool, args []string) {
	if neg {
		ts.Fatalf("unsupported: ! unique-name")
	}
	if len(args) != 2 {
		ts.Fatalf("usage: unique-name VAR prefix")
	}
	ts.Setenv(args[0], uniqueName(args[1], ts.Name(), os.Getpid()))
}

func uniqueName(prefix, testName string, pid int) string {
	return fmt.Sprintf("%s-%s-%d", prefix, testName, pid)
}
