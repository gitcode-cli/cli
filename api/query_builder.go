package api

import "net/url"

// queryBuilder assembles a URL query string while skipping empty/zero values,
// eliminating the repeated if-set boilerplate across the api package query
// functions. Its String method produces a bit-identical result to the legacy
// `url.Values` + `if len(values) > 0 { endpoint += "?" + values.Encode() }`
// pattern.
type queryBuilder struct{ v url.Values }

// newQueryBuilder returns a ready-to-use queryBuilder.
func newQueryBuilder() *queryBuilder { return &queryBuilder{v: url.Values{}} }

// Set adds key=value, skipping empty values (mirrors
// `if v != "" { values.Set(key, v) }`).
func (q *queryBuilder) Set(key, value string) *queryBuilder {
	if value != "" {
		q.v.Set(key, value)
	}
	return q
}

// SetInt adds key=itoa(n), skipping n <= 0 (mirrors
// `if n > 0 { values.Set(key, itoa(n)) }`).
func (q *queryBuilder) SetInt(key string, n int) *queryBuilder {
	if n > 0 {
		q.v.Set(key, itoa(n))
	}
	return q
}

// SetInt64 adds key=itoa64(n), skipping n <= 0 (mirrors
// `if n > 0 { values.Set(key, itoa64(n)) }`).
func (q *queryBuilder) SetInt64(key string, n int64) *queryBuilder {
	if n > 0 {
		q.v.Set(key, itoa64(n))
	}
	return q
}

// String returns the encoded query with a leading "?" or "" when no parameters
// were set (mirrors `if len(values) > 0 { endpoint += "?" + values.Encode() }`).
func (q *queryBuilder) String() string {
	if len(q.v) == 0 {
		return ""
	}
	return "?" + q.v.Encode()
}
