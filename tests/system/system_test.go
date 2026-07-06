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
		Dir: "testdata/read",
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
		Dir: "testdata/write",
		Setup: func(env *testscript.Env) error {
			setupEnv(t, env, root, gcBin, readRepo, writeRepo)
			return nil
		},
		Cmds: systemCmds(),
	})
}

func systemCmds() map[string]func(ts *testscript.TestScript, neg bool, args []string) {
	return map[string]func(ts *testscript.TestScript, neg bool, args []string){
		"defer-close-issue": cmdDeferCloseIssue,
		"json-ok":           cmdJSONOK,
		"require-infra":     cmdRequireInfra,
		"stdout2env":        cmdStdout2Env,
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
	if neg {
		ts.Fatalf("unsupported: ! require-infra")
	}
	if len(args) != 1 {
		ts.Fatalf("usage: require-infra repo")
	}
	if err := validateInfraRepo("repo", args[0]); err != nil {
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
