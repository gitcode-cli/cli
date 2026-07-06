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
