package contracts

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestReleaseWorkflowOrdersTagBeforeGoReleaser(t *testing.T) {
	workflow := readReleaseWorkflow(t)
	assertOrdered(t, workflow, "- name: Validate release notes", "- name: Create and push tag")
	assertOrdered(t, workflow, "- name: Validate release tree", "- name: Create and push tag")
	assertOrdered(t, workflow, "- name: Preflight packages and wheel", "- name: Create and push tag")
	assertOrdered(t, workflow, "- name: Create and push tag", "- name: Build GoReleaser assets")
}

func TestReleaseWorkflowUsesTrackedNotesAndExactTag(t *testing.T) {
	workflow := readReleaseWorkflow(t)
	required := []string{
		`test "${GITHUB_REF_NAME}" = "main"`,
		`test -f "docs/releases/${VERSION_TAG}.md"`,
		`git diff --exit-code`,
		`TAG_SHA="$(git rev-list -n 1 "${VERSION_TAG}")"`,
		`test "${TAG_SHA}" = "${HEAD_SHA}"`,
		`args: release --clean --skip=publish`,
		`dist/gc_linux_amd64`,
		`dist/gc_linux_arm64`,
		`contents: read`,
		`persist-credentials: false`,
		`bash scripts/prepare-package-assets.sh`,
		`name: release-assets`,
		`xargs -0 sha256sum`,
		`sha256sum -c "gc_${VERSION_NUM}_checksums.txt"`,
	}
	for _, value := range required {
		if !strings.Contains(workflow, value) {
			t.Errorf("release workflow missing %q", value)
		}
	}
}

func TestReleaseWorkflowIsValidYAML(t *testing.T) {
	parseReleaseWorkflow(t)
}

func TestReleaseWorkflowConfiguresPyPIPublishing(t *testing.T) {
	document := parseReleaseWorkflow(t)
	jobs := requireMapping(t, document, "jobs")
	pypi := requireMapping(t, jobs, "pypi")
	for _, dependency := range []string{"preflight", "publish"} {
		if !containsString(requireStringList(t, pypi, "needs"), dependency) {
			t.Fatalf("pypi job must depend on %q", dependency)
		}
	}

	permissions := requireMapping(t, pypi, "permissions")
	for key, want := range map[string]string{"contents": "read", "id-token": "write"} {
		if got := requireString(t, permissions, key); got != want {
			t.Fatalf("pypi permission %s = %q, want %q", key, got, want)
		}
	}

	environment := requireMapping(t, pypi, "environment")
	if got := requireString(t, environment, "name"); got != "pypi" {
		t.Fatalf("pypi environment = %q, want pypi", got)
	}
	pypiYAML, err := yaml.Marshal(pypi)
	if err != nil {
		t.Fatalf("marshal pypi job: %v", err)
	}
	for _, forbidden := range []string{"gh release upload", "pip install", "python -m build"} {
		if strings.Contains(string(pypiYAML), forbidden) {
			t.Fatalf("pypi job must not run %q while holding OIDC permission", forbidden)
		}
	}
	if !strings.Contains(readReleaseWorkflow(t), "pip install build==1.4.1") {
		t.Fatal("pypi build tool must be pinned to build 1.4.1")
	}
}

func TestReleaseWorkflowSeparatesWritePermissions(t *testing.T) {
	document := parseReleaseWorkflow(t)
	jobs := requireMapping(t, document, "jobs")
	wants := map[string]string{
		"preflight": "read",
		"tag":       "write",
		"artifacts": "read",
		"publish":   "write",
		"pypi":      "read",
	}
	for jobName, want := range wants {
		job := requireMapping(t, jobs, jobName)
		permissions := requireMapping(t, job, "permissions")
		if got := requireString(t, permissions, "contents"); got != want {
			t.Fatalf("job %s contents permission = %q, want %q", jobName, got, want)
		}
	}
}

func TestPythonBuildBackendIsPinned(t *testing.T) {
	path := filepath.Join("..", "..", "pyproject.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	content := string(data)
	for _, dependency := range []string{"setuptools==80.9.0", "wheel==0.45.1"} {
		if !strings.Contains(content, dependency) {
			t.Fatalf("Python build dependency %q is not pinned", dependency)
		}
	}
}

func TestReleaseWorkflowPinsPublishingDependencies(t *testing.T) {
	workflow := readReleaseWorkflow(t)
	if strings.Contains(workflow, "@latest") {
		t.Fatal("release workflow must not install publishing tools from @latest")
	}
	for _, line := range strings.Split(workflow, "\n") {
		if strings.Contains(line, "uses:") && !pinnedActionPattern.MatchString(line) {
			t.Errorf("release action is not pinned to a commit: %s", strings.TrimSpace(line))
		}
	}
}

func TestPackageScriptPinsNfpm(t *testing.T) {
	path := filepath.Join("..", "..", "scripts", "package.sh")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	script := string(data)
	if !strings.Contains(script, "nfpm/v2/cmd/nfpm@v2.41.3") || strings.Contains(script, "nfpm@latest") {
		t.Fatal("package script must recommend the Go 1.22-compatible pinned nfpm version")
	}
}

var pinnedActionPattern = regexp.MustCompile(`uses:\s+[^@\s]+@[0-9a-f]{40}(?:\s+#.*)?$`)

func readReleaseWorkflow(t *testing.T) string {
	t.Helper()
	path := filepath.Join("..", "..", ".github", "workflows", "release.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func parseReleaseWorkflow(t *testing.T) map[string]any {
	t.Helper()
	var document map[string]any
	if err := yaml.Unmarshal([]byte(readReleaseWorkflow(t)), &document); err != nil {
		t.Fatalf("parse release workflow: %v", err)
	}
	return document
}

func requireMapping(t *testing.T, parent map[string]any, key string) map[string]any {
	t.Helper()
	value, ok := parent[key].(map[string]any)
	if !ok {
		t.Fatalf("workflow key %q is not a mapping", key)
	}
	return value
}

func requireString(t *testing.T, parent map[string]any, key string) string {
	t.Helper()
	value, ok := parent[key].(string)
	if !ok {
		t.Fatalf("workflow key %q is not a string", key)
	}
	return value
}

func requireStringList(t *testing.T, parent map[string]any, key string) []string {
	t.Helper()
	values, ok := parent[key].([]any)
	if !ok {
		t.Fatalf("workflow key %q is not a list", key)
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		item, ok := value.(string)
		if !ok {
			t.Fatalf("workflow key %q contains a non-string value", key)
		}
		result = append(result, item)
	}
	return result
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func assertOrdered(t *testing.T, content, first, second string) {
	t.Helper()
	firstIndex := strings.Index(content, first)
	secondIndex := strings.Index(content, second)
	if firstIndex < 0 || secondIndex < 0 || firstIndex >= secondIndex {
		t.Fatalf("expected %q before %q", first, second)
	}
}
