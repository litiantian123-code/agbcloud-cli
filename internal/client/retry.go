// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// RetryConfig defines retry behavior
type RetryConfig struct {
	MaxRetries    int
	InitialDelay  time.Duration
	MaxDelay      time.Duration
	BackoffFactor float64
}

// DefaultRetryConfig returns a sensible default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:    3,
		InitialDelay:  500 * time.Millisecond,
		MaxDelay:      5 * time.Second,
		BackoffFactor: 2.0,
	}
}

// IsRetryableError determines if an error should trigger a retry
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Network connection errors that are typically transient
	errorStr := err.Error()

	// First check for non-retryable errors
	nonRetryableErrors := []string{
		"invalid uri for request",
		"parse",
		"malformed",
		"unsupported protocol",
	}

	for _, nonRetryableErr := range nonRetryableErrors {
		if strings.Contains(strings.ToLower(errorStr), nonRetryableErr) {
			return false
		}
	}

	// Common transient network errors
	retryableErrors := []string{
		"bad file descriptor",
		"connection refused",
		"connection reset by peer",
		"no such host",
		"network is unreachable",
		"timeout",
		"deadline exceeded",
		"temporary failure",
		"i/o timeout",
		"broken pipe",
	}

	for _, retryableErr := range retryableErrors {
		if strings.Contains(strings.ToLower(errorStr), retryableErr) {
			return true
		}
	}

	// Check for specific error types
	// Note: netErr.Temporary() is deprecated since Go 1.18
	// We rely on string matching and timeout checks instead
	if netErr, ok := err.(net.Error); ok {
		return netErr.Timeout()
	}

	// Check for syscall errors
	if opErr, ok := err.(*net.OpError); ok {
		if syscallErr, ok := opErr.Err.(*syscall.Errno); ok {
			switch *syscallErr {
			case syscall.ECONNREFUSED, syscall.ECONNRESET, syscall.ETIMEDOUT:
				return true
			}
		}
	}

	return false
}

// IsRetryableHTTPStatus determines if an HTTP status code should trigger a retry
func IsRetryableHTTPStatus(statusCode int) bool {
	// Retry on server errors and some client errors
	switch statusCode {
	case http.StatusRequestTimeout, // 408
		http.StatusTooManyRequests,     // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout:      // 504
		return true
	}
	return false
}

// RetryableHTTPClient wraps an HTTP client with retry functionality
type RetryableHTTPClient struct {
	client      *http.Client
	retryConfig *RetryConfig
}

// NewRetryableHTTPClient creates a new HTTP client with retry capability
func NewRetryableHTTPClient(client *http.Client, config *RetryConfig) *RetryableHTTPClient {
	if config == nil {
		config = DefaultRetryConfig()
	}
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	return &RetryableHTTPClient{
		client:      client,
		retryConfig: config,
	}
}

// Do executes an HTTP request with retry logic
func (r *RetryableHTTPClient) Do(req *http.Request) (*http.Response, error) {
	var lastErr error
	delay := r.retryConfig.InitialDelay

	for attempt := 0; attempt <= r.retryConfig.MaxRetries; attempt++ {
		// Clone the request for each attempt (in case body needs to be re-read)
		reqClone := req.Clone(req.Context())

		log.Debugf("[RETRY] Attempt %d/%d for %s %s",
			attempt+1, r.retryConfig.MaxRetries+1, req.Method, req.URL.String())

		resp, err := r.client.Do(reqClone)

		// Success case
		if err == nil && !IsRetryableHTTPStatus(resp.StatusCode) {
			if attempt > 0 {
				log.Infof("[RETRY] Request succeeded on attempt %d", attempt+1)
			}
			return resp, nil
		}

		// Store the error for potential retry
		if err != nil {
			lastErr = err
			log.Debugf("[RETRY] Attempt %d failed with error: %v", attempt+1, err)
		} else {
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
			log.Debugf("[RETRY] Attempt %d failed with HTTP status: %d", attempt+1, resp.StatusCode)
			// Close the response body to avoid resource leak
			resp.Body.Close()
		}

		// Don't retry if this is the last attempt
		if attempt == r.retryConfig.MaxRetries {
			break
		}

		// Check if the error is retryable
		shouldRetry := false
		if err != nil {
			shouldRetry = IsRetryableError(err)
			if !shouldRetry {
				log.Debugf("[RETRY] Error is not retryable, stopping attempts")
				break
			}
		} else if resp != nil {
			shouldRetry = IsRetryableHTTPStatus(resp.StatusCode)
			if !shouldRetry {
				log.Debugf("[RETRY] HTTP status is not retryable, stopping attempts")
				break
			}
		}

		// Wait before retrying
		log.Infof("[RETRY] Request failed (attempt %d/%d), retrying in %v...",
			attempt+1, r.retryConfig.MaxRetries+1, delay)

		select {
		case <-time.After(delay):
			// Delay elapsed, continue to next attempt
		case <-req.Context().Done():
			return nil, req.Context().Err()
		}

		// Calculate next delay with exponential backoff
		delay = time.Duration(float64(delay) * r.retryConfig.BackoffFactor)
		if delay > r.retryConfig.MaxDelay {
			delay = r.retryConfig.MaxDelay
		}
	}

	log.Warnf("[RETRY] All %d attempts failed, giving up", r.retryConfig.MaxRetries+1)
	return nil, fmt.Errorf("request failed after %d attempts, last error: %w",
		r.retryConfig.MaxRetries+1, lastErr)
}
