// Package api provides HTTP client configuration
package api

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Default timeout values
const (
	DefaultTimeout  = 30 * time.Second
	UploadTimeout   = 10 * time.Minute
	DownloadTimeout = 10 * time.Minute
	LongPollTimeout = 5 * time.Minute
)

// DefaultHTTPClient returns an HTTP client with default configuration
func DefaultHTTPClient() *http.Client {
	return NewHTTPClient(DefaultTimeout)
}

// NewHTTPClient creates an HTTP client with specified timeout
func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout:   timeout,
		Transport: newDefaultTransport(),
	}
}

// NewHTTPClientWithRetry creates an HTTP client with retry middleware
func NewHTTPClientWithRetry(timeout time.Duration, cfg RetryConfig) *http.Client {
	transport := newDefaultTransport()
	return &http.Client{
		Timeout:   timeout,
		Transport: RetryMiddleware(transport, cfg),
	}
}

// NewHTTPClientWithRetryAndLogger creates an HTTP client with retry middleware and debug logging
func NewHTTPClientWithRetryAndLogger(timeout time.Duration, cfg RetryConfig, logger func(string)) *http.Client {
	transport := newDefaultTransport()
	return &http.Client{
		Timeout:   timeout,
		Transport: RetryMiddlewareWithLogger(transport, cfg, logger),
	}
}

// NewUploadHTTPClient returns an HTTP client suitable for file uploads
func NewUploadHTTPClient() *http.Client {
	return NewHTTPClientWithRetry(UploadTimeout, DefaultRetryConfig())
}

// NewDownloadHTTPClient returns an HTTP client suitable for file downloads
func NewDownloadHTTPClient() *http.Client {
	return NewHTTPClientWithRetry(DownloadTimeout, DefaultRetryConfig())
}

// newDefaultTransport returns the default HTTP transport
func newDefaultTransport() http.RoundTripper {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// ParseTimeoutFromEnv parses timeout from GC_TIMEOUT environment variable
func ParseTimeoutFromEnv() time.Duration {
	if v := os.Getenv("GC_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil && d > 0 {
			return d
		}
		// Try parsing as seconds if no unit specified
		seconds, err := strconv.Atoi(v)
		if err == nil && seconds > 0 {
			return time.Duration(seconds) * time.Second
		}
	}
	return DefaultTimeout
}

// IsDebugEnabled returns true if GC_DEBUG is set
func IsDebugEnabled() bool {
	return os.Getenv("GC_DEBUG") != "" || os.Getenv("GC_API_DEBUG") != ""
}
