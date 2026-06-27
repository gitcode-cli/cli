package testutil

import "net/http"

// RoundTripFunc is an http.RoundTripper implemented by a function.
// It allows tests to mock HTTP transport without a real server.
type RoundTripFunc func(*http.Request) (*http.Response, error)

// RoundTrip implements http.RoundTripper by calling the underlying function.
func (fn RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

// NewRoundTripFunc wraps fn in a RoundTripFunc for use as an http.RoundTripper.
func NewRoundTripFunc(fn func(*http.Request) (*http.Response, error)) RoundTripFunc {
	return RoundTripFunc(fn)
}
