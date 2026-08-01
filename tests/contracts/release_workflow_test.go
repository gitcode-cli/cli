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
	assertOrdered(t, workflow, "- name: Create and push tag", "- name: Run GoReleaser")
}

func TestReleaseWorkflowUsesTrackedNotesAndExactTag(t *testing.T) {
	workflow := readReleaseWorkflow(t)
	required := []string{
		`test "${GITHUB_REF_NAME}" = "main"`,
		`test -f "docs/releases/${VERSION_TAG}.md"`,
		`git diff --exit-code`,
		`TAG_SHA="$(git rev-list -n 1 "${VERSION_TAG}")"`,
		`test "${TAG_SHA}" = "${HEAD_SHA}"`,
		`--release-notes=docs/releases/${{ steps.version.outputs.VERSION_TAG }}.md`,
		`dist/gc_linux_amd64`,
		`dist/gc_linux_arm64`,
		`contents: read`,
	}
	for _, value := range required {
		if !strings.Contains(workflow, value) {
			t.Errorf("release workflow missing %q", value)
		}
	}
}

func TestReleaseWorkflowIsValidYAML(t *testing.T) {
	workflow := readReleaseWorkflow(t)
	var document any
	if err := yaml.Unmarshal([]byte(workflow), &document); err != nil {
		t.Fatalf("parse release workflow: %v", err)
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

func assertOrdered(t *testing.T, content, first, second string) {
	t.Helper()
	firstIndex := strings.Index(content, first)
	secondIndex := strings.Index(content, second)
	if firstIndex < 0 || secondIndex < 0 || firstIndex >= secondIndex {
		t.Fatalf("expected %q before %q", first, second)
	}
}
