package installupdate

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadBootstrapManifestAdjacentToBinary(t *testing.T) {
	dir := t.TempDir()
	binary := filepath.Join(dir, "gitcode")
	t.Setenv("GITCODE_CLI_BINARY", binary)
	manifestPath := filepath.Join(dir, ".gitcode-install.json")
	data := []byte(`{"distribution":"npm-bootstrap","version":"1.2.3","targetDir":"` + escapedJSON(dir) + `","helper":"` + escapedJSON(filepath.Join(dir, "helper.js")) + `"}`)
	if err := os.WriteFile(manifestPath, data, 0o600); err != nil {
		t.Fatal(err)
	}
	manifest, err := LoadBootstrapManifest()
	if err != nil {
		t.Fatal(err)
	}
	if manifest.Version != "1.2.3" || manifest.ManifestPath() != manifestPath {
		t.Fatalf("unexpected manifest: %#v", manifest)
	}
}

func TestAfterCommandShowsSummaryButDoesNotScheduleWhenDisabled(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("GITCODE_CLI_BINARY", filepath.Join(dir, "gitcode"))
	t.Setenv("GC_STATE_DIR", filepath.Join(dir, "state"))
	manifest := Manifest{Distribution: "npm-bootstrap", Version: "1.2.3", TargetDir: dir, Helper: filepath.Join(dir, "missing.js")}
	data, _ := json.Marshal(manifest)
	if err := os.WriteFile(filepath.Join(dir, ".gitcode-install.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	writeState(StatePath(), updateState{Summary: &stateSummary{Message: "updated", Shown: false}})
	out := &bytes.Buffer{}
	AfterCommand(nil, out, true, false)
	if out.String() != "updated\n" {
		t.Fatalf("output = %q", out.String())
	}
	if !readState(StatePath()).Summary.Shown {
		t.Fatal("summary should be marked shown")
	}
}

func TestDueAtUsesTwentyFourHourTTL(t *testing.T) {
	now := time.Unix(100, 0)
	if got := DueAt(now); got != now.Add(24*time.Hour).UnixMilli() {
		t.Fatalf("DueAt() = %d", got)
	}
}

func TestUpdaterEnvironmentStripsCredentials(t *testing.T) {
	t.Setenv("GC_TOKEN", "gitcode-secret")
	t.Setenv("GITCODE_TOKEN", "legacy-secret")
	t.Setenv("NPM_TOKEN", "npm-secret")
	t.Setenv("NODE_AUTH_TOKEN", "node-secret")
	for _, item := range updaterEnvironment() {
		upper := strings.ToUpper(item)
		for _, key := range []string{"GC_TOKEN=", "GITCODE_TOKEN=", "NPM_TOKEN=", "NODE_AUTH_TOKEN="} {
			if strings.HasPrefix(upper, key) {
				t.Fatalf("credential leaked to updater environment: %s", key)
			}
		}
	}
}

func escapedJSON(value string) string {
	data, _ := json.Marshal(value)
	return string(data[1 : len(data)-1])
}
