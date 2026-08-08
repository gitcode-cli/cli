package api

import (
	"net/url"
	"testing"
)

func TestQueryBuilderEmpty(t *testing.T) {
	q := newQueryBuilder()
	if got := q.String(); got != "" {
		t.Fatalf("empty String = %q, want %q", got, "")
	}
}

func TestQueryBuilderSkipsEmptyAndZero(t *testing.T) {
	q := newQueryBuilder().
		Set("event", "").
		Set("status", "").
		SetInt("page", 0).
		SetInt("per_page", -1).
		SetInt64("startTime", 0).
		SetInt64("endTime", -5)
	if got := q.String(); got != "" {
		t.Fatalf("all-empty String = %q, want %q (no trailing ?)", got, "")
	}
}

func TestQueryBuilderSetsNonEmpty(t *testing.T) {
	q := newQueryBuilder().
		Set("event", "push").
		Set("status", "success").
		SetInt("page", 2).
		SetInt("per_page", 50).
		SetInt64("startTime", 1700000000).
		SetInt64("endTime", 1800000000)
	got := q.String()
	if got == "" {
		t.Fatal("String = empty, want query")
	}
	parsed, err := url.Parse(got)
	if err != nil {
		t.Fatalf("url.Parse(%q) err = %v", got, err)
	}
	want := map[string]string{
		"event":     "push",
		"status":    "success",
		"page":      "2",
		"per_page":  "50",
		"startTime": "1700000000",
		"endTime":   "1800000000",
	}
	for k, w := range want {
		if g := parsed.Query().Get(k); g != w {
			t.Errorf("param %q = %q, want %q", k, g, w)
		}
	}
}

// TestQueryBuilderBitIdenticalLegacy confirms the builder produces the same
// output as the legacy url.Values pattern it replaces.
func TestQueryBuilderBitIdenticalLegacy(t *testing.T) {
	// Legacy (with the same skip-empty/skip-zero guards the api code uses)
	values := url.Values{}
	if e := "push"; e != "" {
		values.Set("event", e)
	}
	if s := ""; s != "" {
		values.Set("status", s)
	}
	if pp := 50; pp > 0 {
		values.Set("per_page", itoa(pp))
	}
	if p := 0; p > 0 {
		values.Set("page", itoa(p))
	}
	legacy := ""
	if len(values) > 0 {
		legacy = "?" + values.Encode()
	}
	// Builder mirrors the same skip rules.
	b := newQueryBuilder().
		Set("event", "push").
		Set("status", "").
		SetInt("per_page", 50).
		SetInt("page", 0)
	got := b.String()
	if got != legacy {
		t.Fatalf("builder = %q, legacy = %q, want bit-identical", got, legacy)
	}
}
