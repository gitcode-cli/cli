package contracts

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

const npmSecurityFlags = "--yes --ignore-scripts --registry=https://registry.npmjs.org --@gitcode-cli:registry=https://registry.npmjs.org"
const canonicalNPMBootstrap = "npx " + npmSecurityFlags + " @gitcode-cli/cli@latest install"
const canonicalNPMGlobal = "npm install -g --ignore-scripts --registry=https://registry.npmjs.org --@gitcode-cli:registry=https://registry.npmjs.org @gitcode-cli/cli@latest"

var npxInstallPattern = regexp.MustCompile(`npx[^` + "`\"\r\n" + `]*@gitcode-cli/cli[^` + "`\"\r\n" + `]*install`)
var npmGlobalPattern = regexp.MustCompile(`npm install -g[^` + "`\"\r\n" + `]*@gitcode-cli/cli[^` + "`\"\r\n" + `]*`)

func TestInstallDocsUseCanonicalNPMBootstrap(t *testing.T) {
	version := strings.TrimSpace(readRepositoryFile(t, "VERSION"))
	currentReleasePath := filepath.Join("docs", "releases", "v"+version+".md")
	paths := []string{
		"README.md",
		filepath.Join("npm", "README.md"),
		filepath.Join("docs", "INTRODUCTION.md"),
		filepath.Join("docs", "PACKAGING.md"),
		filepath.Join("docs", "AI-GUIDE.md"),
		currentReleasePath,
		filepath.Join("spec", "delivery", "release-process.md"),
	}
	for _, relativePath := range paths {
		content := readRepositoryFile(t, relativePath)
		if !strings.Contains(content, canonicalNPMBootstrap) {
			t.Errorf("%s does not contain canonical npm bootstrap %q", relativePath, canonicalNPMBootstrap)
		}
		assertNPMInstallLinesAreCanonical(t, relativePath, content)
	}
	assertNPMInstallLinesAreCanonical(t, filepath.Join("npm", "bin", "gc.js"), readRepositoryFile(t, filepath.Join("npm", "bin", "gc.js")))
}

func TestREADMEPrioritizesBootstrapAndExplainsLocalInstall(t *testing.T) {
	content := readRepositoryFile(t, "README.md")
	bootstrap := strings.Index(content, canonicalNPMBootstrap)
	sourceBuild := strings.Index(content, "### 从源码构建")
	if bootstrap < 0 || sourceBuild < 0 || bootstrap > sourceBuild {
		t.Fatal("README must show the canonical npm bootstrap before source installation")
	}
	if !strings.Contains(content, "`npm i @gitcode-cli/cli` 或 `npm install @gitcode-cli/cli`") {
		t.Fatal("README must distinguish project-local npm dependencies from CLI installation")
	}
}

func TestCurrentReleaseNotesProvidePinnedBootstrap(t *testing.T) {
	version := strings.TrimSpace(readRepositoryFile(t, "VERSION"))
	releasePath := filepath.Join("docs", "releases", "v"+version+".md")
	content := readRepositoryFile(t, releasePath)
	pinned := "npx " + npmSecurityFlags + " @gitcode-cli/cli@" + version + " install"
	if !strings.Contains(content, pinned) {
		t.Fatalf("%s must provide pinned bootstrap command %q", releasePath, pinned)
	}
}

func assertNPMInstallLinesAreCanonical(t *testing.T, path, content string) {
	t.Helper()
	version := strings.TrimSpace(readRepositoryFile(t, "VERSION"))
	pinned := "npx " + npmSecurityFlags + " @gitcode-cli/cli@" + version + " install"
	for _, command := range npxInstallPattern.FindAllString(content, -1) {
		if command == canonicalNPMBootstrap || command == pinned || strings.HasPrefix(command, canonicalNPMBootstrap+" --target-dir ") {
			continue
		}
		t.Errorf("%s contains non-canonical npx install command %q", path, command)
	}
	for _, command := range npmGlobalPattern.FindAllString(content, -1) {
		if command == canonicalNPMGlobal {
			continue
		}
		t.Errorf("%s contains non-canonical global npm install command %q", path, command)
	}
}

func readRepositoryFile(t *testing.T, relativePath string) string {
	t.Helper()
	path := filepath.Join("..", "..", relativePath)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
