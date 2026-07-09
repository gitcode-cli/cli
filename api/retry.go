// Package api provides retry middleware for HTTP client
package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// RetryConfig configures retry behavior
type RetryConfig struct {
	MaxRetries  int           // Maximum number of retries (default: 3)
	InitialWait time.Duration // Initial wait time (default: 1s)
	MaxWait     time.Duration // Maximum wait time (default: 30s)
	Multiplier  float64       // Backoff multiplier (default: 2.0)
}

// DefaultRetryConfig returns the default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:  3,
		InitialWait: 1 * time.Second,
		MaxWait:     30 * time.Second,
		Multiplier:  2.0,
	}
}

// retryTransport wraps a RoundTripper with retry logic
type retryTransport struct {
	base   http.RoundTripper
	cfg    RetryConfig
	logger func(string)
}

// RetryMiddleware wraps a RoundTripper with retry logic
func RetryMiddleware(rt http.RoundTripper, cfg RetryConfig) http.RoundTripper {
	return &retryTransport{base: rt, cfg: cfg}
}

// RetryMiddlewareWithLogger wraps a RoundTripper with retry logic and debug logging
func RetryMiddlewareWithLogger(rt http.RoundTripper, cfg RetryConfig, logger func(string)) http.RoundTripper {
	return &retryTransport{base: rt, cfg: cfg, logger: logger}
}

// sanitizeError removes URL and host details from network errors to avoid
// leaking endpoint information into CI logs. url.Error.Error() formats as
// `Get "https://host/path": <underlying>`, so we unwrap to the inner error.
func sanitizeError(err error) string {
	var ue *url.Error
	if errors.As(err, &ue) {
		if ue.Err != nil {
			return ue.Err.Error()
		}
		return ue.Op
	}
	return err.Error()
}

// RoundTrip executes the request with retry logic
func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// For rewindable bodies (GetBody != nil), retries can restore the body.
	// For non-rewindable bodies (GetBody == nil, e.g. *os.File or stdin),
	// we skip retry buffering to avoid loading the entire body into memory.
	// The first attempt proceeds without retry capability.
	canRetry := req.Body == nil || req.GetBody != nil
	if !canRetry && t.cfg.MaxRetries > 0 {
		if t.logger != nil {
			t.logger("retry: body is non-rewindable (GetBody==nil), skipping retry to avoid memory buffering")
		}
		return t.base.RoundTrip(req)
	}

	var lastErr error
	var lastResp *http.Response

	for attempt := 0; attempt <= t.cfg.MaxRetries; attempt++ {
		// Restore body for each attempt
		if req.Body != nil {
			body, err := req.GetBody()
			if err != nil {
				return nil, fmt.Errorf("failed to restore request body: %w", err)
			}
			req.Body = body
		}

		resp, err := t.base.RoundTrip(req)
		if err != nil {
			lastErr = err
			if t.logger != nil {
				t.logger(fmt.Sprintf("retry: attempt %d failed with error: %s", attempt+1, sanitizeError(err)))
			}
			if t.shouldRetryOnError(err) {
				wait := t.calculateWait(attempt)
				if t.logger != nil {
					t.logger(fmt.Sprintf("retry: waiting %v before retry %d", wait, attempt+2))
				}
				time.Sleep(wait)
				continue
			}
			return nil, err
		}

		// Check status code for retry
		if t.shouldRetryOnStatus(resp.StatusCode) {
			lastResp = resp

			// Handle Rate Limit (429)
			if resp.StatusCode == 429 {
				wait := t.handleRateLimit(resp)
				if wait > 0 {
					if t.logger != nil {
						t.logger(fmt.Sprintf("retry: rate limited, waiting %v (Retry-After)", wait))
					}
					time.Sleep(wait)
					resp.Body.Close()
					continue
				}
			}

			// Exponential backoff for server errors
			wait := t.calculateWait(attempt)
			if t.logger != nil {
				t.logger(fmt.Sprintf("retry: server error %d, waiting %v before retry %d", resp.StatusCode, wait, attempt+2))
			}
			resp.Body.Close()
			time.Sleep(wait)
			continue
		}

		// Success or non-retryable status
		return resp, nil
	}

	// All retries exhausted
	if lastResp != nil {
		return lastResp, nil
	}
	return nil, fmt.Errorf("request failed after %d retries: %w", t.cfg.MaxRetries, lastErr)
}

// shouldRetryOnError determines if we should retry based on error
func (t *retryTransport) shouldRetryOnError(err error) bool {
	// Retry on network errors (connection refused, timeout, etc.)
	if err == nil {
		return false
	}

	// Don't retry on context cancellation
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	// Retry on temporary network errors
	return true
}

// shouldRetryOnStatus determines if we should retry based on HTTP status code
func (t *retryTransport) shouldRetryOnStatus(statusCode int) bool {
	switch statusCode {
	case 401: // Unauthorized - don't retry
		return false
	case 429: // Rate Limited - retry with Retry-After
		return true
	case 500, 502, 503, 504: // Server errors - retry with backoff
		return true
	default:
		return false
	}
}

// handleRateLimit extracts Retry-After header value
func (t *retryTransport) handleRateLimit(resp *http.Response) time.Duration {
	retryAfter := resp.Header.Get("Retry-After")
	if retryAfter == "" {
		return 0
	}

	// Try parsing as seconds
	seconds, err := strconv.Atoi(retryAfter)
	if err == nil && seconds > 0 {
		wait := time.Duration(seconds) * time.Second
		if wait > t.cfg.MaxWait {
			wait = t.cfg.MaxWait
		}
		return wait
	}

	// Try parsing as date (RFC 850)
	// For simplicity, we just use default wait if date format is used
	return t.cfg.InitialWait
}

// calculateWait calculates exponential backoff wait time
func (t *retryTransport) calculateWait(attempt int) time.Duration {
	wait := t.cfg.InitialWait
	for i := 0; i < attempt; i++ {
		wait = time.Duration(float64(wait) * t.cfg.Multiplier)
		if wait > t.cfg.MaxWait {
			wait = t.cfg.MaxWait
			break
		}
	}
	return wait
}
