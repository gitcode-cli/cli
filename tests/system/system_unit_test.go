//go:build system

package system_test

import (
	"reflect"
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
	name := uniqueName("system-test-label", "label-lifecycle", 1234)
	if name != "system-test-label-label-lifecycle-1234" {
		t.Fatalf("uniqueName returned %q", name)
	}
}
