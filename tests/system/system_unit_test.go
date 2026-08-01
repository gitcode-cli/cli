//go:build system

package system_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestLookupJSONPath(t *testing.T) {
	value := map[string]any{
		"repo": map[string]any{
			"name": "gctest1",
		},
		"items": []any{
			map[string]any{"number": "1"},
		},
	}

	tests := []struct {
		name string
		path string
		want any
		ok   bool
	}{
		{name: "root", path: ".", want: value, ok: true},
		{name: "object key", path: "repo.name", want: "gctest1", ok: true},
		{name: "leading dot", path: ".repo.name", want: "gctest1", ok: true},
		{name: "array index", path: "items[0].number", want: "1", ok: true},
		{name: "missing key", path: "repo.missing", ok: false},
		{name: "missing index", path: "items[1].number", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := lookupJSONPath(value, tt.path)
			if err != nil {
				t.Fatalf("lookupJSONPath returned error: %v", err)
			}
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if ok && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateInfraRepo(t *testing.T) {
	valid := []string{"infra-test/gctest1", "infra-test/another-repo"}
	for _, repo := range valid {
		if err := validateInfraRepo("repo", repo); err != nil {
			t.Fatalf("validateInfraRepo(%q) returned error: %v", repo, err)
		}
	}

	invalid := []string{"", "gitcode-cli/cli", "personal/repo", "infra-test/", "infra-test/nested/repo"}
	for _, repo := range invalid {
		if err := validateInfraRepo("repo", repo); err == nil {
			t.Fatalf("validateInfraRepo(%q) unexpectedly succeeded", repo)
		}
	}
}

func TestSystemCondition(t *testing.T) {
	const name = "GC_SYSTEM_CONDITION_TEST"
	t.Setenv(name, "set")
	ok, err := systemCondition("env:" + name)
	if err != nil || !ok {
		t.Fatalf("systemCondition() = %v, %v, want true, nil", ok, err)
	}
	t.Setenv(name, "")
	ok, err = systemCondition("env:" + name)
	if err != nil || ok {
		t.Fatalf("systemCondition() = %v, %v, want false, nil", ok, err)
	}
	if _, err := systemCondition("unsupported"); err == nil {
		t.Fatal("systemCondition() error = nil, want unknown condition error")
	}
}

func TestJSONTypeMatches(t *testing.T) {
	tests := []struct {
		value any
		want  string
		ok    bool
	}{
		{value: "text", want: "string", ok: true},
		{value: "text", want: "nonempty-string", ok: true},
		{value: "", want: "nonempty-string", ok: false},
		{value: float64(1), want: "number", ok: true},
		{value: []any{}, want: "array", ok: true},
		{value: map[string]any{}, want: "object", ok: true},
		{value: nil, want: "null", ok: true},
	}

	for _, tt := range tests {
		if got := jsonTypeMatches(tt.value, tt.want); got != tt.ok {
			t.Fatalf("jsonTypeMatches(%v, %q) = %v, want %v", tt.value, tt.want, got, tt.ok)
		}
	}
}

func TestLookupJSONPathSupportsAssigneeLogin(t *testing.T) {
	value := map[string]any{
		"assignees": []any{map[string]any{"login": "alice"}},
	}
	got, ok, err := lookupJSONPath(value, "assignees[0].login")
	if err != nil {
		t.Fatalf("lookupJSONPath returned error: %v", err)
	}
	if !ok || got != "alice" {
		t.Fatalf("lookupJSONPath got %v, %v, want alice, true", got, ok)
	}
}

func TestUniqueNameShape(t *testing.T) {
	name, err := uniqueName("system-test-label", "label-lifecycle", 1234, 0)
	if err != nil {
		t.Fatalf("uniqueName returned error: %v", err)
	}
	matched, err := regexp.MatchString(`^system-test-label-label-lifecycle-1234-[0-9a-f]{32}$`, name)
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Fatalf("uniqueName returned %q with unexpected shape", name)
	}
}

func TestUniqueNameHonorsMaxLength(t *testing.T) {
	name, err := uniqueName("system-test-label", "label-lifecycle", 1234, 50)
	if err != nil {
		t.Fatalf("uniqueName returned error: %v", err)
	}
	if got := len([]rune(name)); got > 50 {
		t.Fatalf("uniqueName length = %d, want at most 50: %q", got, name)
	}
	matched, err := regexp.MatchString(`^system-test[^-]*-1234-[0-9a-f]{32}$`, name)
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Fatalf("uniqueName returned %q without the expected prefix and suffix", name)
	}
}

func TestUniqueNameRejectsTooSmallMaxLength(t *testing.T) {
	if _, err := uniqueName("label", "test", 1234, 38); err == nil {
		t.Fatal("uniqueName error = nil, want max length error")
	}
}

func TestParsePositiveInt(t *testing.T) {
	for _, value := range []string{"", "invalid", "0", "-1"} {
		if _, err := parsePositiveInt(value); err == nil {
			t.Fatalf("parsePositiveInt(%q) error = nil", value)
		}
	}
	if got, err := parsePositiveInt("50"); err != nil || got != 50 {
		t.Fatalf("parsePositiveInt(50) = %d, %v, want 50, nil", got, err)
	}
}

func TestDeleteLabelIfExists(t *testing.T) {
	tests := []struct {
		name       string
		listOutput string
		listErr    error
		deleteErr  error
		wantCalls  int
		wantError  string
	}{
		{name: "absent", listOutput: `[{"name":"other"}]`, wantCalls: 1},
		{name: "present", listOutput: `[{"name":"target"}]`, wantCalls: 2},
		{name: "list failure", listErr: errors.New("network"), wantCalls: 1, wantError: "list labels"},
		{name: "malformed response", listOutput: `{`, wantCalls: 1, wantError: "parse label list"},
		{name: "delete failure", listOutput: `[{"name":"target"}]`, deleteErr: errors.New("denied"), wantCalls: 2, wantError: "delete label"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			run := func(args ...string) ([]byte, error) {
				calls++
				if calls == 1 {
					return []byte(tt.listOutput), tt.listErr
				}
				return nil, tt.deleteErr
			}
			err := deleteLabelIfExists(run, "infra-test/gctest1", "target")
			if calls != tt.wantCalls {
				t.Fatalf("calls = %d, want %d", calls, tt.wantCalls)
			}
			if tt.wantError == "" && err != nil {
				t.Fatalf("deleteLabelIfExists returned error: %v", err)
			}
			if tt.wantError != "" && (err == nil || !strings.Contains(err.Error(), tt.wantError)) {
				t.Fatalf("error = %v, want containing %q", err, tt.wantError)
			}
		})
	}
}

func TestDeleteLabelIfExistsPaginates(t *testing.T) {
	firstPage := make([]map[string]string, 100)
	for i := range firstPage {
		firstPage[i] = map[string]string{"name": "other"}
	}
	firstOutput, err := json.Marshal(firstPage)
	if err != nil {
		t.Fatal(err)
	}

	calls := 0
	wantArgs := [][]string{
		{"label", "list", "-R", "infra-test/gctest1", "--limit", "100", "--page", "1", "--json"},
		{"label", "list", "-R", "infra-test/gctest1", "--limit", "100", "--page", "2", "--json"},
		{"label", "delete", "target", "-R", "infra-test/gctest1", "--yes"},
	}
	run := func(args ...string) ([]byte, error) {
		calls++
		if calls > len(wantArgs) || !reflect.DeepEqual(args, wantArgs[calls-1]) {
			t.Fatalf("call %d args = %v", calls, args)
		}
		switch calls {
		case 1:
			return firstOutput, nil
		case 2:
			return []byte(`[{"name":"target"}]`), nil
		case 3:
			return nil, nil
		}
		return nil, nil
	}
	if err := deleteLabelIfExists(run, "infra-test/gctest1", "target"); err != nil {
		t.Fatalf("deleteLabelIfExists returned error: %v", err)
	}
	if calls != 3 {
		t.Fatalf("calls = %d, want 3", calls)
	}
}
